package cloudcasa

import (
	"context"
	"errors"

	cloudcasa "terraform-provider-cloudcasa/cloudcasa/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceKubebackup{}
	_ resource.ResourceWithConfigure   = &resourceKubebackup{}
	_ resource.ResourceWithImportState = &resourceKubebackup{}
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
	Select_namespaces []types.String        `tfsdk:"select_namespaces"`
	Snapshot_pvs      types.Bool            `tfsdk:"snapshot_persistent_volumes"`
	Copy_pvs          types.Bool            `tfsdk:"copy_persistent_volumes"`
	Delete_snapshots  types.Bool            `tfsdk:"delete_snapshot_after_copy"`
	Kubeoffload_id    types.String          `tfsdk:"kubeoffload_id"`
	Updated           types.String          `tfsdk:"updated"`
	Created           types.String          `tfsdk:"created"`
	Etag              types.String          `tfsdk:"etag"`
	Offload_etag      types.String          `tfsdk:"offload_etag"`
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

// TODO: RequiredWith https://developer.hashicorp.com/terraform/plugin/framework/migrating/attributes-blocks/fields
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
			// TODO: ListNested to SetNested? Sets are UNORDERED pairs, lists are ORDERED
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
			"offload_etag": schema.StringAttribute{
				Computed: true,
			},
			"copy_persistent_volumes": schema.BoolAttribute{
				Optional: true,
			},
			"delete_snapshot_after_copy": schema.BoolAttribute{
				Optional: true,
			},
			"kubeoffload_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
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

func ValidateKubebackupSelections(plan kubebackupResourceModel) error {
	// Validate namespace option selections
	if plan.All_namespaces.ValueBool() && plan.Select_namespaces != nil {
		return errors.New("Set all_namespaces to true to snapshot every namespace OR define a set of namespaces to snapshot with the select_namespaces attribute.")
	}

	if !plan.All_namespaces.ValueBool() && plan.Select_namespaces == nil {
		return errors.New("Define a set of namespaces to snapshot with the select_namespaces attribute, or set all_namespaces to true.")
	}

	// Validate Copy options
	if !plan.Copy_pvs.ValueBool() && plan.Delete_snapshots.ValueBool() {
		return errors.New("delete_snapshot_after_copy requires copy_persistent_volumes to be true.")
	}

	return nil
}

// CreateKubebackupFromPlan initializes a request body from TF values
func CreateKubebackupFromPlan(plan kubebackupResourceModel) (cloudcasa.CreateKubebackupReq, error) {

	// Build 'source' dict of kubebackup body from plan
	reqBodySource := cloudcasa.KubebackupSource{
		All_namespaces:            plan.All_namespaces.ValueBool(),
		SnapshotPersistentVolumes: plan.Snapshot_pvs.ValueBool(),
	}
	if plan.Select_namespaces != nil {
		reqBodySource.Namespaces = ConvertTfStringList(plan.Select_namespaces)
	}

	// Build main kubebackup body from plan
	reqBody := cloudcasa.CreateKubebackupReq{
		Name:    plan.Name.ValueString(),
		Cluster: plan.Kubecluster_id.ValueString(),
		Source:  reqBodySource,
	}

	// Validate namespace option selections
	if plan.All_namespaces.ValueBool() && plan.Select_namespaces != nil {
		return reqBody, errors.New("Set all_namespaces to true to snapshot every namespace OR define a set of namespaces to snapshot with the select_namespaces attribute.")
	}

	if !plan.All_namespaces.ValueBool() && plan.Select_namespaces == nil {
		return reqBody, errors.New("Define a set of namespaces to snapshot with the select_namespaces attribute, or set all_namespaces to true.")
	}

	// Validate Copy options
	if !plan.Copy_pvs.ValueBool() && plan.Delete_snapshots.ValueBool() {
		return reqBody, errors.New("delete_snapshot_after_copy requires copy_persistent_volumes to be true.")
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
			return reqBody, errors.New("Retention is set but backup job will not run. run_on_apply must be true to run the job without selecting a policy.")
		}
	}

	// If run_on_apply, set trigger_type to ADHOC
	if plan.Run.ValueBool() {
		reqBody.Trigger_type = "ADHOC"
	} else {
		reqBody.Trigger_type = "SCHEDULED"

		// Exit if no policy is defined for scheduled backup
		if plan.Policy_id.IsNull() {
			return reqBody, errors.New("Kubebackups run on a schedule by default and require a policy. To run an adhoc backup, set run_on_apply.")
		}
	}

	return reqBody, nil
}

func (r *resourceKubebackup) RunKubebackup(plan kubebackupResourceModel, backupId string) error {
	// Select options for backup
	var retentionDays int
	if plan.Retention.IsNull() {
		retentionDays = 7
	} else {
		retentionDays = int(plan.Retention.ValueInt64())
	}

	var backupJobId string
	var lastJobRunTime int64

	// Run Offload or Backup
	if plan.Copy_pvs.ValueBool() {
		runResp, err := r.Client.RunKubeoffload(backupId, retentionDays)
		if err != nil {
			return errors.New("error running kubeoffload: " + err.Error())
		}
		backupJobId = runResp.Backupdef
		lastJobRunTime = runResp.Status.LastJobRunTime
	} else {
		runResp, err := r.Client.RunKubebackup(backupId, retentionDays)
		if err != nil {
			return errors.New("error running kubebackup: " + err.Error())
		}
		backupJobId = runResp.Id
		lastJobRunTime = runResp.Status.LastJobRunTime
	}

	// if backup.status.lastjobruntime is 0 (nil) : we just created the kubebackup,
	// so first job matching the filter is the correct job.
	// if timestamp is not 0, we have a last run timestamp so job has ran before
	// first job since that timestamp is the correct one.

	// Get Job ID
	jobResp, err := r.Client.GetJobFromBackupdef(backupJobId, lastJobRunTime)
	if err != nil {
		if err != nil {
			return errors.New("error waiting for job to start: " + err.Error())
		}
	}

	// watch job
	_, err = r.Client.WatchJobUntilComplete(jobResp.Id)
	if err != nil {
		if err != nil {
			return errors.New("error fetching job status: " + err.Error())
		}
	}

	return nil
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

	reqBody, err := CreateKubebackupFromPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Kubebackup Definition",
			err.Error(),
		)
		return
	}

	// If Copy options are present, initialize kubeoffload body
	var copyReqBody cloudcasa.CreateKubeoffloadReq
	if plan.Copy_pvs.ValueBool() {
		// Build kubeoffload body
		copyReqBody.Name = plan.Name.ValueString()
		copyReqBody.Cluster = plan.Kubecluster_id.ValueString()
		copyReqBody.Delete_snapshots = plan.Delete_snapshots.ValueBool()

		// Check optional kubeoffload fields
		if !plan.Policy_id.IsNull() {
			copyReqBody.Policy = plan.Policy_id.ValueString()
		}

		if plan.Run.ValueBool() {
			copyReqBody.Trigger_type = "ADHOC"
			copyReqBody.Run_backup = true
		} else {
			copyReqBody.Trigger_type = "SCHEDULED"
			copyReqBody.Run_backup = false
		}

	}

	// Create kubebackup resource in CloudCasa
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
	plan.Kubeoffload_id = types.StringNull()
	plan.Offload_etag = types.StringNull()

	diags = resp.State.Set(ctx, plan)

	var createKubeoffloadResp *cloudcasa.Kubeoffload
	backupId := createResp.Id

	// Set backupdef ID for kubeoffload
	if plan.Copy_pvs.ValueBool() {
		copyReqBody.Backupdef = createResp.Id
		createKubeoffloadResp, err = r.Client.CreateKubeoffload(copyReqBody)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Kubeoffload",
				err.Error(),
			)
			return
		}

		// Set offload ID, use Offload ID as backupId for copy jobs
		plan.Kubeoffload_id = types.StringValue(createKubeoffloadResp.Id)
		backupId = createKubeoffloadResp.Id

		// Append copydef ID to original kubebackup request
		reqBody.Copydef = createKubeoffloadResp.Id

		// Update kubebackup resource in CloudCasa
		putResp, err := r.Client.UpdateKubebackup(createResp.Id, reqBody, createResp.Etag)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Kubebackup",
				err.Error(),
			)
			return
		}

		// Update fields in plan
		plan.Updated = types.StringValue(putResp.Updated)
		plan.Etag = types.StringValue(putResp.Etag)
		plan.Offload_etag = types.StringValue(createKubeoffloadResp.Etag)
	}

	// If run_on_apply is false return now. Otherwise continue and run the job
	if !plan.Run.ValueBool() {
		return
	}

	err = r.RunKubebackup(plan, backupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"error running backup job",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data from CloudCasa
// TODO: set from CC resource like in resource_policy.go
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
			"Error getting Kubebackup from CloudCasa",
			"Could not read Kubebackup with ID "+state.Id.ValueString()+" :"+err.Error(),
		)
		return
	}

	// Overwrite state values with refreshed cloudcasa info
	state.Id = types.StringValue(kubebackup.Id)
	state.Name = types.StringValue(kubebackup.Name)
	state.Kubecluster_id = types.StringValue(kubebackup.Cluster)
	state.Kubeoffload_id = types.StringValue(kubebackup.Copydef)

	state.Snapshot_pvs = types.BoolValue(kubebackup.Source.SnapshotPersistentVolumes)
	state.All_namespaces = types.BoolValue(kubebackup.Source.All_namespaces)
	if kubebackup.Source.All_namespaces != true {
		state.Select_namespaces = ConvertStringListTf(kubebackup.Source.Namespaces)
	}

	// check hooks fields and convert
	if kubebackup.Pre_hooks != nil {
		for _, v := range kubebackup.Pre_hooks {
			thisHook := kubebackupHookModel{
				Template:   types.BoolValue(v.Template),
				Namespaces: ConvertStringListTf(v.Namespaces),
				Hooks:      ConvertStringListTf(v.Hooks),
			}
			state.Pre_hooks = append(state.Pre_hooks, thisHook)
		}
	}
	if kubebackup.Post_hooks != nil {
		for _, v := range kubebackup.Post_hooks {
			thisHook := kubebackupHookModel{
				Template:   types.BoolValue(v.Template),
				Namespaces: ConvertStringListTf(v.Namespaces),
				Hooks:      ConvertStringListTf(v.Hooks),
			}
			state.Post_hooks = append(state.Post_hooks, thisHook)
		}
	}

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
	// Retrieve values from plan
	var plan kubebackupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from TF state
	// need etag value to edit the existing object
	var state kubebackupResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// CLUSTER ID is immutable
	if plan.Kubecluster_id.ValueString() != state.Kubecluster_id.ValueString() {
		resp.Diagnostics.AddError(
			"invalid kubebackup definition",
			"cluster id cannot be changed for an existing kubebackup",
		)
		return
	}

	//update other fields
	reqBody, err := CreateKubebackupFromPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"invalid kubebackup definition",
			err.Error(),
		)
		return
	}

	// If Copy options are present, initialize kubeoffload body
	var copyReqBody cloudcasa.CreateKubeoffloadReq
	if plan.Copy_pvs.ValueBool() {
		// Add offload ID to kubebackup request
		reqBody.Copydef = state.Kubeoffload_id.ValueString()

		// Build kubeoffload body
		copyReqBody.Name = plan.Name.ValueString()
		copyReqBody.Cluster = plan.Kubecluster_id.ValueString()
		copyReqBody.Delete_snapshots = plan.Delete_snapshots.ValueBool()

		// Check optional kubeoffload fields
		if !plan.Policy_id.IsNull() {
			copyReqBody.Policy = plan.Policy_id.ValueString()
		}

		if plan.Run.ValueBool() {
			copyReqBody.Trigger_type = "ADHOC"
			copyReqBody.Run_backup = true
		} else {
			copyReqBody.Trigger_type = "SCHEDULED"
			copyReqBody.Run_backup = false
		}

	}

	// Update kubebackup resource in CloudCasa
	updateResp, err := r.Client.UpdateKubebackup(plan.Id.ValueString(), reqBody, state.Etag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kubebackup",
			err.Error(),
		)
		return
	}

	// Set fields in plan
	plan.Id = types.StringValue(updateResp.Id)
	plan.Created = types.StringValue(updateResp.Created)
	plan.Updated = types.StringValue(updateResp.Updated)
	plan.Etag = types.StringValue(updateResp.Etag)
	plan.Kubeoffload_id = types.StringValue(updateResp.Copydef)
	plan.Offload_etag = state.Offload_etag

	diags = resp.State.Set(ctx, plan)

	// Update kubeoffload if copy is selected
	var updateKubeoffloadResp *cloudcasa.Kubeoffload
	backupId := updateResp.Id

	// Set backupdef ID for kubeoffload
	if plan.Copy_pvs.ValueBool() {
		copyReqBody.Backupdef = updateResp.Id

		// GET kubeoffload to get current ETAG
		// TODO: do same for kubebackup/kubecluster?
		getKubeoffload, err := r.Client.GetKubeoffload(state.Kubeoffload_id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"error fetching kubeoffload",
				err.Error(),
			)
			return
		}
		plan.Offload_etag = types.StringValue(getKubeoffload.Etag)

		updateKubeoffloadResp, err = r.Client.UpdateKubeoffload(state.Kubeoffload_id.ValueString(), copyReqBody, getKubeoffload.Etag)
		if err != nil {
			resp.Diagnostics.AddError(
				"error updating kubeoffload",
				err.Error(),
			)
			return
		}

		// Set offload ID, use Offload ID as backupId for copy jobs
		plan.Kubeoffload_id = types.StringValue(updateKubeoffloadResp.Id)
		backupId = updateKubeoffloadResp.Id

		// Append copydef ID to original kubebackup request
		reqBody.Copydef = updateKubeoffloadResp.Id

		// Update kubebackup resource in CloudCasa
		putResp, err := r.Client.UpdateKubebackup(updateResp.Id, reqBody, updateResp.Etag)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Kubebackup",
				err.Error(),
			)
			return
		}

		// Update fields in plan
		plan.Updated = types.StringValue(putResp.Updated)
		plan.Etag = types.StringValue(putResp.Etag)
		plan.Offload_etag = types.StringValue(updateKubeoffloadResp.Etag)
	}

	// If run_on_apply is false return now. Otherwise continue and run the job
	if !plan.Run.ValueBool() {
		return
	}

	err = r.RunKubebackup(plan, backupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"error running backup job",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceKubebackup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state kubebackupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete kubeoffload first if valid
	if state.Copy_pvs.ValueBool() {
		err := r.Client.DeleteKubeoffload(state.Kubeoffload_id.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Error Deleting Kubeoffload resource",
				err.Error(),
			)
			return
		}
	}

	err := r.Client.DeleteKubebackup(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kubebackup resource",
			err.Error(),
		)
		return
	}
}

func (r *resourceKubebackup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
