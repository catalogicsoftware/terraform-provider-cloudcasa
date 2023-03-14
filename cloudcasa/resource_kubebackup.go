package cloudcasa

import (
	"context"

	cloudcasa "terraform-provider-cloudcasa/cloudcasa/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &resourceKubebackup{}
	_ resource.ResourceWithConfigure = &resourceKubebackup{}
)

func NewResourceKubebackup() resource.Resource {
	return &resourceKubebackup{}
}

type resourceKubebackup struct {
	Client *cloudcasa.Client
}

// kubebackupResourceModel maps the resource schema data.
type kubebackupResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Kubecluster_id    types.String `tfsdk:"kubecluster_id"`
	Policy_id         types.String `tfsdk:"policy_id"`
	Pre_hooks         types.Set    `tfsdk:"pre_hooks"`
	Post_hooks        types.Set    `tfsdk:"post_hooks"`
	Run               types.Bool   `tfsdk:"run_on_apply"`
	Retention         types.Int64  `tfsdk:"retention"`
	All_namespaces    types.Bool   `tfsdk:"all_namespaces"`
	Select_namespaces types.Set    `tfsdk:"select_namespaces"`
	Snapshot_pvs      types.Bool   `tfsdk:"snapshot_persistent_volumes"`
	Updated           types.String `tfsdk:"updated"`
	Created           types.String `tfsdk:"created"`
	Etag              types.String `tfsdk:"etag"`
	// Pause             types.Bool   `tfsdk:"pause"`
}

// API Response Objects
type CreateKubebackupResp struct {
	Id   string `json:"_id"`
	Name string `json:"name"`
}

// Metadata returns the data source type name.
func (r *resourceKubebackup) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubebackup"
}

// Schema defines the schema for the resource.
func (r *resourceKubebackup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"kubecluster_id": schema.StringAttribute{
				Required: true,
			},
			"policy_id": schema.StringAttribute{
				Optional: true,
			},
			"pre_hooks": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"post_hooks": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			// run_on_apply will determine trigger_type
			// TODO: implement /run API on every apply by forcing GET
			// like we do for kubeclusters
			"run_on_apply": schema.BoolAttribute{
				Optional: true,
			},
			"retention": schema.NumberAttribute{
				Optional: true,
			},
			"all_namespaces": schema.BoolAttribute{
				// TODO: Better validation between these two namespace attrs
				Required: true,
			},
			"select_namespaces": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"snapshot_persistent_volumes": schema.BoolAttribute{
				Required: true,
			},
			// "pause": schema.BoolAttribute{
			// 	Optional: true,
			// },
			"updated": schema.StringAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *resourceKubebackup) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.Client = req.ProviderData.(*cloudcasa.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *resourceKubebackup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan kubebackupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate namespace option selections
	if plan.All_namespaces.ValueBool() && !plan.Select_namespaces.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Kubebackup Definition",
			"Set all_namespaces to true to copy every namespace OR define a set of namespaces to copy with the select_namespace attribute.",
		)
		return
	}

	// Build 'source' dict of request body from plan
	reqBodySource := cloudcasa.CreateKubebackupReqSource{
		All_namespaces:            plan.All_namespaces.ValueBool(),
		SnapshotPersistentVolumes: plan.Snapshot_pvs.ValueBool(),
	}
	if !plan.Select_namespaces.IsNull() {
		plan.Select_namespaces.ElementsAs(ctx, reqBodySource.Namespaces, false)
	}

	// Build main request body from plan
	reqBody := cloudcasa.CreateKubebackupReq{
		Name:    plan.Name.ValueString(),
		Cluster: plan.Kubecluster_id.ValueString(),
		Source:  reqBodySource,
	}

	// Check optional fields
	if !plan.Policy_id.IsNull() {
		reqBody.Policy = plan.Policy_id.ValueString()
	}
	if !plan.Pre_hooks.IsNull() {
		plan.Pre_hooks.ElementsAs(ctx, reqBody.Pre_hooks, false)
	}
	if !plan.Post_hooks.IsNull() {
		plan.Post_hooks.ElementsAs(ctx, reqBody.Post_hooks, false)
	}

	// If retention is set, check that run_on_apply is true
	if !plan.Retention.IsNull() {
		if !plan.Run.ValueBool() {
			resp.Diagnostics.AddError(
				"Invalid Kubebackup Definition",
				"Retention is set but backup job will not run. run_on_apply must be true to specify retention outside of a policy.",
			)
		}
	}

	// If run_on_apply, set trigger_type to ADHOC
	if plan.Run.ValueBool() {
		reqBody.Trigger_type = "ADHOC"
	} else {
		reqBody.Trigger_type = "SCHEDULED"

		// WARN user if no policy is defined
		if plan.Policy_id.IsNull() {
			resp.Diagnostics.AddError(
				"No policy defined for kubebackup",
				"Kubebackups run on a schedule by default and require a policy definition. To run an Adhoc backup, set run_on_apply.",
			)
		}
	}

	// Create resource in CloudCasa
	createResp, err := r.Client.CreateKubebackup(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Kubebackup",
			err.Error(),
		)
		return
	}

	// Set fields in plan
	plan.Id = types.StringValue(createResp.Id)
	plan.Created = types.StringValue(createResp.Created)
	plan.Updated = types.StringValue(createResp.Updated)
	plan.Etag = types.StringValue(createResp.Etag)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: run backup

}

// Read refreshes the Terraform state with the latest data.
func (r *resourceKubebackup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourceKubebackup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceKubebackup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}
