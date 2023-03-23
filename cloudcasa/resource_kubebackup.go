package cloudcasa

import (
	"context"
	"fmt"

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
	Id                types.String          `tfsdk:"id"`
	Name              types.String          `tfsdk:"name"`
	Kubecluster_id    types.String          `tfsdk:"kubecluster_id"`
	Policy_id         types.String          `tfsdk:"policy_id"`
	Pre_hooks         []kubebackupHookModel `tfsdk:"pre_hooks"`
	Post_hooks        []kubebackupHookModel `tfsdk:"post_hooks"`
	Run               types.Bool            `tfsdk:"run_on_apply"`
	Retention         types.Int64           `tfsdk:"retention"`
	All_namespaces    types.Bool            `tfsdk:"all_namespaces"`
	Select_namespaces types.Set             `tfsdk:"select_namespaces"`
	Snapshot_pvs      types.Bool            `tfsdk:"snapshot_persistent_volumes"`
	Updated           types.String          `tfsdk:"updated"`
	Created           types.String          `tfsdk:"created"`
	Etag              types.String          `tfsdk:"etag"`
	// Pause             types.Bool   `tfsdk:"pause"`
}

type kubebackupHookModel struct {
	Template   types.Bool     `tfsdk:"template"`
	Namespaces []types.String `tfsdk:"namespaces"`
	Hooks      []types.String `tfsdk:"hooks"`
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
			"pre_hooks": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"template": schema.BoolAttribute{
							Required: true,
						},
						"namespaces": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
						"hooks": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"post_hooks": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"template": schema.BoolAttribute{
							Required: true,
						},
						"namespaces": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
						"hooks": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			// run_on_apply will determine trigger_type
			// TODO: implement /run API on every apply by forcing GET
			// like we do for kubeclusters
			"run_on_apply": schema.BoolAttribute{
				Optional: true,
			},
			"retention": schema.Int64Attribute{
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
	reqBodySource := cloudcasa.KubebackupSource{
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

	// For each Hook in pre_hooks, convert string values and append
	if plan.Pre_hooks != nil {
		for _, v := range plan.Pre_hooks {
			thisHook := cloudcasa.KubebackupHook{
				Template:   v.Template.ValueBool(),
				Namespaces: ConvertTfStringList(v.Namespaces),
				Hooks:      ConvertTfStringList(v.Hooks),
			}
			reqBody.Pre_hooks = append(reqBody.Pre_hooks, thisHook)
		}
	}
	if plan.Post_hooks != nil {
		for _, v := range plan.Post_hooks {
			thisHook := cloudcasa.KubebackupHook{
				Template:   v.Template.ValueBool(),
				Namespaces: ConvertTfStringList(v.Namespaces),
				Hooks:      ConvertTfStringList(v.Hooks),
			}
			reqBody.Post_hooks = append(reqBody.Post_hooks, thisHook)
		}
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

	// DEBUG
	resp.Diagnostics.AddWarning(
		"pv option before create",
		fmt.Sprint(reqBody.Source.SnapshotPersistentVolumes),
	)

	// Create resource in CloudCasa
	createResp, err := r.Client.CreateKubebackup(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Kubebackup",
			err.Error(),
		)
		return
	}

	// DEBUG
	resp.Diagnostics.AddWarning(
		"pv option after create",
		fmt.Sprint(createResp.Source.SnapshotPersistentVolumes),
	)

	// Set fields in plan
	plan.Id = types.StringValue(createResp.Id)
	plan.Created = types.StringValue(createResp.Created)
	plan.Updated = types.StringValue(createResp.Updated)
	plan.Etag = types.StringValue(createResp.Etag)

	diags = resp.State.Set(ctx, plan)

	// If run_on_apply is false return now. Otherwise continue and run the job
	if !plan.Run.ValueBool() {
		return
	}

	// Select options for backup
	// tODO: Handle offloads, not supported rn
	backupType := "kubebackups"
	if createResp.Source.SnapshotPersistentVolumes {
		backupType = "kubeoffloads"
	}

	var retentionDays int
	if plan.Retention.IsNull() {
		retentionDays = 7
	} else {
		retentionDays = int(plan.Retention.ValueInt64())
	}

	// Run backup
	runResp, err := r.Client.RunKubebackup(createResp.Id, backupType, retentionDays)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Running Kubebackup",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.AddWarning("Received runResp. Job should be running...", runResp.Id)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if backup.status.lastjobruntime is 0 (nil) : we just created the kubebackup,
	// so first job matching the filter is the correct job.
	// if timestamp is not 0, we have a last run timestamp so job has ran before
	// first job since that timestamp is the correct one.

	// DEBUG
	resp.Diagnostics.AddWarning("lastJobRunTime", fmt.Sprint(runResp.Status.LastJobRunTime))

	// Get Job ID
	jobResp, err := r.Client.GetJobFromBackupdef(runResp.Id, runResp.Status.LastJobRunTime)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for backup job to start",
			err.Error(),
		)
		return
	}

	// DEBUG
	resp.Diagnostics.AddWarning("Found Job with ID", jobResp.Id)

	// watch job
	jobStatusResp, err := r.Client.WatchJobUntilComplete(jobResp.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching running job status",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.AddWarning("Job state", jobStatusResp.State)

	if jobStatusResp.State != "COMPLETED" {
		resp.Diagnostics.AddWarning("Job finished in an incomplete state", fmt.Sprint("Job %w finished in state %s. This means the job completed successfully, but some resources might have been missed. Check logs in the CloudCasa UI for more information.", jobResp.Id, jobStatusResp.State))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *resourceKubebackup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state kubebackupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Kubecluster from CloudCasa
	kubebackup, err := r.Client.GetKubebackup(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Kuebbackup from CloudCasa",
			"Could not read Kubebackup with ID "+state.Id.ValueString()+" :"+err.Error(),
		)
		return
	}

	// Overwrite values with refreshed state
	// Set fields in plan
	state.Id = types.StringValue(kubebackup.Id)
	state.Name = types.StringValue(kubebackup.Name)
	state.Kubecluster_id = types.StringValue(kubebackup.Cluster)

	// Check if Optional values are Null
	if kubebackup.Policy == "" {
		state.Policy_id = types.StringNull()
	} else {
		state.Policy_id = types.StringValue(kubebackup.Policy)
	}

	// convert list values
	// preHooksList, diags := types.ListValueFrom(ctx, types.StringType, kubebackup.Pre_hooks)
	// resp.Diagnostics.Append(diags...)
	// state.Pre_hooks = basetypes.SetValue(preHooksList)

	// postHooksList, diags := types.ListValueFrom(ctx, types.StringType, kubebackup.Post_hooks)
	// resp.Diagnostics.Append(diags...)
	// state.Post_hooks = basetypes.SetValue(postHooksList)

	// check for errors from list conversion
	if resp.Diagnostics.HasError() {
		return
	}

	// Job runtime options are not read from API:
	//state.Run
	//state.Retention
	//state.All_namespaces
	//state.Select_namespaces
	//state.Snapshot_pvs

	state.Updated = types.StringValue(kubebackup.Updated)
	state.Created = types.StringValue(kubebackup.Created)
	state.Etag = types.StringValue(kubebackup.Etag)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourceKubebackup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceKubebackup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}
