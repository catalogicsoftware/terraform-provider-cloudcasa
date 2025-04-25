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
	Copy_policy       types.String          `tfsdk:"copy_policy"`
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

// Schema defines the schema for the resource.
func (r *resourceKubebackup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CloudCasa kubebackup configuration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "CloudCasa resource ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "CloudCasa resource name",
			},
			"kubecluster_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the kubecluster to back up",
			},
			"policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of a policy for scheduling this backup",
			},
			"pre_hooks": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Pre-backup app hooks to execute. See https://docs.cloudcasa.io/help/configuration-apphook.html for details",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"template": schema.BoolAttribute{
							Required:    true,
							Description: "Set to use a predefined hook template",
						},
						"namespaces": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "List of namespaces to run the selected hook in",
						},
						"hooks": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "ID of a hook created in CloudCasa",
						},
					},
				},
			},
			"post_hooks": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Post-backup app hooks to execute. See https://docs.cloudcasa.io/help/configuration-apphook.html for details",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"template": schema.BoolAttribute{
							Required:    true,
							Description: "Set to use a predefined hook template",
						},
						"namespaces": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "List of namespaces to run the selected hook in",
						},
						"hooks": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "ID of a hook created in CloudCasa",
						},
					},
				},
			},
			// run_on_apply will determine trigger_type
			"run_on_apply": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to run the backup immediately after creation or update. If enabled, this will also cause the backup to run on each terraform apply",
			},
			"retention": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of days to retain backup data for",
			},
			"all_namespaces": schema.BoolAttribute{
				Required:    true,
				Description: "Set to backup all namespaces, otherwise use the select_namespaces attribute to list namespaces",
			},
			"select_namespaces": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of namespaces to include in the backup",
			},
			"snapshot_persistent_volumes": schema.BoolAttribute{
				Required:    true,
				Description: "Set to snapshot persistent volumes. If false, PVs will be ignored",
			},
			"copy_policy": schema.StringAttribute{
				Computed:    true,
				Description: "ID of a policy used for scheduling copy backups, used internally by CloudCasa.",
			},
			"updated": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the CloudCasa resource",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the CloudCasa resource",
			},
			"etag": schema.StringAttribute{
				Computed:    true,
				Description: "Etag generated by CloudCasa, used for updating resources in place",
			},
			"offload_etag": schema.StringAttribute{
				Computed:    true,
				Description: "Etag of the associated offload resource generated by CloudCasa, used for updating resources in place",
			},
			"copy_persistent_volumes": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, persistent volume data will be copied and offloaded to S3 storage. This will create and manage an associated kubeoffload resource in CloudCasa",
			},
			"delete_snapshot_after_copy": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to delete resource snapshots after performing data offload",
			},
			"kubeoffload_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "ID of the associated kubeoffload resource created for Copy backups",
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

// createKubebackupFromPlan initializes a request body from TF values
func createKubebackupFromPlan(plan kubebackupResourceModel) (cloudcasa.Kubebackup, error) {
	// Build 'source' dict of kubebackup body from plan
	kubebackupSource := cloudcasa.KubebackupSource{
		All_namespaces:            plan.All_namespaces.ValueBool(),
		SnapshotPersistentVolumes: plan.Snapshot_pvs.ValueBool(),
	}
	if plan.Select_namespaces != nil {
		kubebackupSource.Namespaces = ConvertTfStringList(plan.Select_namespaces)
	}

	// Build main kubebackup body from plan
	kubebackup := cloudcasa.Kubebackup{
		Name:    plan.Name.ValueString(),
		Cluster: plan.Kubecluster_id.ValueString(),
		Source:  kubebackupSource,
	}

	// Validate namespace option selections
	if plan.All_namespaces.ValueBool() && plan.Select_namespaces != nil {
		return kubebackup, errors.New("set all_namespaces to true to snapshot every namespace OR define a list of namespaces to snapshot with the select_namespaces attribute.")
	}

	if !plan.All_namespaces.ValueBool() && plan.Select_namespaces == nil {
		return kubebackup, errors.New("define a list of namespaces to snapshot with the select_namespaces attribute, or set all_namespaces to true.")
	}

	// Validate Copy options
	if !plan.Copy_pvs.ValueBool() && plan.Delete_snapshots.ValueBool() {
		return kubebackup, errors.New("delete_snapshot_after_copy requires copy_persistent_volumes to be true.")
	}

	// Check optional fields
	if !plan.Policy_id.IsNull() {
		// For scheduled COPY jobs, we need to pass the policy to the API initially
		// But we'll clear it later after getting the offload ID
		kubebackup.Policy = plan.Policy_id.ValueString()
	}

	// For each Hook in pre_hooks, convert string values and append
	if plan.Pre_hooks != nil {
		for _, v := range plan.Pre_hooks {
			thisHook := cloudcasa.KubebackupHook{
				Template:   v.Template.ValueBool(),
				Namespaces: ConvertTfStringList(v.Namespaces),
				Hooks:      ConvertTfStringList(v.Hooks),
			}
			kubebackup.Pre_hooks = append(kubebackup.Pre_hooks, thisHook)
		}
	}
	if plan.Post_hooks != nil {
		for _, v := range plan.Post_hooks {
			thisHook := cloudcasa.KubebackupHook{
				Template:   v.Template.ValueBool(),
				Namespaces: ConvertTfStringList(v.Namespaces),
				Hooks:      ConvertTfStringList(v.Hooks),
			}
			kubebackup.Post_hooks = append(kubebackup.Post_hooks, thisHook)
		}
	}

	// If retention is set, check that run_on_apply is true
	if !plan.Retention.IsNull() {
		if !plan.Run.ValueBool() {
			return kubebackup, errors.New("retention is set but backup job will not run. run_on_apply must be true to run the job without selecting a policy.")
		}
	}

	// If run_on_apply, set trigger_type to ADHOC
	if plan.Run.ValueBool() {
		kubebackup.Trigger_type = "ADHOC"
	} else {
		kubebackup.Trigger_type = "SCHEDULED"

		// Exit if no policy is defined for scheduled backup
		if plan.Policy_id.IsNull() {
			return kubebackup, errors.New("Kubebackups run on a schedule by default and require a policy. To run an adhoc backup, set run_on_apply.")
		}
	}

	return kubebackup, nil
}

// Create a Kubeoffload request object from the Terraform plan
func createKubeoffloadFromPlan(plan kubebackupResourceModel) (cloudcasa.Kubeoffload, error) {

	var req cloudcasa.Kubeoffload
	// Populate required fields
	req.Run_backup = true
	req.Name = plan.Name.ValueString() // Add name field for the API
	
	// Set trigger_type based on run_on_apply
	if plan.Run.ValueBool() {
		req.Trigger_type = "ADHOC"
	} else {
		req.Trigger_type = "SCHEDULED"
	}
	
	// Set delete_snapshots field if specified
	if !plan.Delete_snapshots.IsNull() {
		req.Delete_snapshots = plan.Delete_snapshots.ValueBool()
	}

	// Set the policy ID - for COPY jobs we use policy_id if set, otherwise copy_policy
	if !plan.Policy_id.IsNull() {
		req.Policy = plan.Policy_id.ValueString()
	} else if !plan.Copy_policy.IsNull() && plan.Copy_policy.ValueString() != "" {
		req.Policy = plan.Copy_policy.ValueString()
	} else {
		// If policy_id is null and copy_pvs is true, ensure we set an empty policy
		req.Policy = ""
	}

	// Set the cluster ID
	req.Cluster = plan.Kubecluster_id.ValueString()

	return req, nil
}

func (r *resourceKubebackup) runBackupJob(plan kubebackupResourceModel, backupId string) error {
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

	// Initialize computed fields with default values for all cases
	plan.Copy_policy = types.StringValue("")
	plan.Offload_etag = types.StringValue("")
	plan.Kubeoffload_id = types.StringValue("")

	// Create Kubebackup object from plan values
	reqBody, err := createKubebackupFromPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"invalid Kubebackup definition",
			err.Error(),
		)
		return
	}

	// If Copy options are present, initialize kubeoffload body
	var copyReqBody cloudcasa.Kubeoffload
	if plan.Copy_pvs.ValueBool() {
		copyReqBody, err = createKubeoffloadFromPlan(plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"invalid Kubeoffload definition",
				err.Error(),
			)
			return
		}
	}

	// Create kubebackup resource in CloudCasa
	createResp, err := r.Client.CreateKubebackup(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"error creating Kubebackup",
			err.Error(),
		)
		return
	}

	// Update basic fields in plan
	plan.Id = types.StringValue(createResp.Id)
	plan.Created = types.StringValue(createResp.Created)
	plan.Updated = types.StringValue(createResp.Updated)
	plan.Etag = types.StringValue(createResp.Etag)

	// Save the state here before setting copy options
	diags = resp.State.Set(ctx, plan)

	var createKubeoffloadResp *cloudcasa.Kubeoffload
	backupId := createResp.Id

	// Set backupdef ID for kubeoffload
	if plan.Copy_pvs.ValueBool() {
		copyReqBody.Backupdef = createResp.Id
		createKubeoffloadResp, err = r.Client.CreateKubeoffload(copyReqBody)
		if err != nil {
			// If creating the kubeoffload fails, delete the kubebackup to avoid leaving it in an inconsistent state
			deleteErr := r.Client.DeleteKubebackup(createResp.Id)
			if deleteErr != nil {
				resp.Diagnostics.AddError(
					"error deleting kubebackup after failed kubeoffload creation",
					deleteErr.Error(),
				)
			}
			
			resp.Diagnostics.AddError(
				"error creating Kubeoffload",
				err.Error(),
			)
			
			// Remove the resource from state
			resp.State.RemoveResource(ctx)
			return
		}

		// Set offload ID, use Offload ID as backupId for copy jobs
		plan.Kubeoffload_id = types.StringValue(createKubeoffloadResp.Id)
		backupId = createKubeoffloadResp.Id

		// Append copydef ID to original kubebackup request
		reqBody.Copydef = createKubeoffloadResp.Id
		
		// For COPY jobs, clear the policy from the kubebackup (it should only be on the kubeoffload)
		if !plan.Policy_id.IsNull() && plan.Copy_pvs.ValueBool() {
			reqBody.Policy = ""
		}

		// Update kubebackup resource in CloudCasa
		putResp, err := r.Client.UpdateKubebackup(createResp.Id, reqBody, createResp.Etag)
		if err != nil {
			resp.Diagnostics.AddError(
				"error updating Kubebackup",
				err.Error(),
			)
			return
		}

		// Update fields in plan
		plan.Updated = types.StringValue(putResp.Updated)
		plan.Etag = types.StringValue(putResp.Etag)
		plan.Offload_etag = types.StringValue(createKubeoffloadResp.Etag)
		
		// For COPY jobs, store the policy in copy_policy as well
		if !plan.Policy_id.IsNull() {
			plan.Copy_policy = types.StringValue(plan.Policy_id.ValueString())
		} else {
			plan.Copy_policy = types.StringValue("")
		}
	} else {
		// For non-COPY jobs, ensure these fields have known values
		plan.Copy_policy = types.StringValue("")
		plan.Offload_etag = types.StringValue("")
		plan.Kubeoffload_id = types.StringValue("")
	}

	// If run_on_apply is false return now. Otherwise continue and run the job
	if !plan.Run.ValueBool() {
		diags = resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		return
	}

	err = r.runBackupJob(plan, backupId)
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
func (r *resourceKubebackup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state kubebackupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Kubebackup from CloudCasa
	kubebackup, err := r.Client.GetKubebackup(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading Kubebackup with ID "+state.Id.ValueString(),
			err.Error(),
		)
		return
	}

	// Overwrite state values with refreshed cloudcasa info
	state.Name = types.StringValue(kubebackup.Name)
	state.Kubecluster_id = types.StringValue(kubebackup.Cluster)
	state.Kubeoffload_id = types.StringValue(kubebackup.Copydef)

	state.Snapshot_pvs = types.BoolValue(kubebackup.Source.SnapshotPersistentVolumes)
	state.All_namespaces = types.BoolValue(kubebackup.Source.All_namespaces)
	if kubebackup.Source.All_namespaces != true {
		state.Select_namespaces = ConvertStringListTf(kubebackup.Source.Namespaces)
	}

	// check hooks fields and convert
	state.Pre_hooks, state.Post_hooks = nil, nil
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

	// For non-COPY jobs, get policy from kubebackup
	if kubebackup.Copydef == "" {
		state.Copy_pvs = types.BoolValue(false)
		state.Kubeoffload_id = types.StringValue("")
		state.Copy_policy = types.StringValue("")
		state.Offload_etag = types.StringValue("")
		state.Delete_snapshots = types.BoolValue(false)
		
		if kubebackup.Policy != "" {
			state.Policy_id = types.StringValue(kubebackup.Policy)
		} else {
			state.Policy_id = types.StringNull()
		}
	} else {
		// This is a COPY job, get the kubeoffload to get its policy
		state.Copy_pvs = types.BoolValue(true)
		kubeoffload, err := r.Client.GetKubeoffload(kubebackup.Copydef)
		if err != nil {
			resp.Diagnostics.AddError(
				"error reading Kubeoffload with ID "+kubebackup.Copydef,
				err.Error(),
			)
			return
		}
		
		state.Offload_etag = types.StringValue(kubeoffload.Etag)
		state.Delete_snapshots = types.BoolValue(kubeoffload.Delete_snapshots)
		
		// For COPY jobs, the policy_id from Terraform is stored in the kubeoffload's policy field
		// so we need to keep it in the policy_id field in the state
		if kubeoffload.Policy != "" {
			state.Policy_id = types.StringValue(kubeoffload.Policy)
			state.Copy_policy = types.StringValue(kubeoffload.Policy)
		} else {
			state.Policy_id = types.StringNull()
			state.Copy_policy = types.StringValue("")
		}
	}

	// Set metadata
	state.Created = types.StringValue(kubebackup.Created)
	state.Updated = types.StringValue(kubebackup.Updated)
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

	// Initialize computed fields with default values for all cases
	plan.Copy_policy = types.StringValue("")
	plan.Offload_etag = types.StringValue("")
	plan.Kubeoffload_id = types.StringValue("")

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
	reqBody, err := createKubebackupFromPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"invalid kubebackup definition",
			err.Error(),
		)
		return
	}

	// If Copy options are present, initialize kubeoffload body
	var copyReqBody cloudcasa.Kubeoffload
	if plan.Copy_pvs.ValueBool() {
		copyReqBody, err = createKubeoffloadFromPlan(plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"invalid Kubeoffload definition",
				err.Error(),
			)
			return
		}
	}

	// Update kubebackup resource in CloudCasa
	updateResp, err := r.Client.UpdateKubebackup(plan.Id.ValueString(), reqBody, state.Etag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating Kubebackup",
			err.Error(),
		)
		return
	}

	// Update basic fields in plan
	plan.Id = types.StringValue(updateResp.Id)
	plan.Created = types.StringValue(updateResp.Created)
	plan.Updated = types.StringValue(updateResp.Updated)
	plan.Etag = types.StringValue(updateResp.Etag)
	
	// Initialize computed fields with empty strings for non-copy jobs
	// or if we're switching from copy to non-copy
	if !plan.Copy_pvs.ValueBool() {
		plan.Copy_policy = types.StringValue("")
		plan.Offload_etag = types.StringValue("")
		// Ensure kubeoffload_id is set to a known empty value
		plan.Kubeoffload_id = types.StringValue("")
	}

	diags = resp.State.Set(ctx, plan)

	// Update kubeoffload if copy is selected
	backupId := updateResp.Id

	// Set backupdef ID for kubeoffload
	if plan.Copy_pvs.ValueBool() {
		copyReqBody.Backupdef = updateResp.Id

		// Determine whether to create or update the kubeoffload
		var updateKubeoffloadResp *cloudcasa.Kubeoffload
		
		// If kubeoffload_id is empty or null in the state, we need to create a new kubeoffload
		if state.Kubeoffload_id.IsNull() || state.Kubeoffload_id.ValueString() == "" {
			// Create a new kubeoffload
			updateKubeoffloadResp, err = r.Client.CreateKubeoffload(copyReqBody)
			if err != nil {
				// If creating the kubeoffload fails, we need to roll back to the original state
				// or delete the kubebackup if it was already modified
				
				// First try to restore the kubebackup without copy settings
				reqBody.Copydef = ""
				restoreResp, restoreErr := r.Client.UpdateKubebackup(updateResp.Id, reqBody, updateResp.Etag)
				if restoreErr != nil {
					// If we can't restore the kubebackup, try to delete it to avoid leaving it in an inconsistent state
					deleteErr := r.Client.DeleteKubebackup(updateResp.Id)
					if deleteErr != nil {
						resp.Diagnostics.AddError(
							"error deleting kubebackup after failed kubeoffload creation",
							deleteErr.Error(),
						)
					}
					resp.Diagnostics.AddError(
						"error creating kubeoffload and failed to restore original kubebackup",
						err.Error() + " and " + restoreErr.Error(),
					)
					// Force a refresh to get the updated state from the server
					resp.State.RemoveResource(ctx)
					return
				}
				
				// Update the state to reflect the restored kubebackup
				plan.Updated = types.StringValue(restoreResp.Updated)
				plan.Etag = types.StringValue(restoreResp.Etag)
				plan.Copy_pvs = types.BoolValue(false)
				
				resp.Diagnostics.AddError(
					"error creating kubeoffload",
					err.Error(),
				)
				return
			}
		} else {
			// GET kubeoffload to get current ETAG
			getKubeoffload, err := r.Client.GetKubeoffload(state.Kubeoffload_id.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"error fetching kubeoffload",
					err.Error(),
				)
				return
			}
			plan.Offload_etag = types.StringValue(getKubeoffload.Etag)

			// Update the existing kubeoffload
			updateKubeoffloadResp, err = r.Client.UpdateKubeoffload(state.Kubeoffload_id.ValueString(), copyReqBody, getKubeoffload.Etag)
			if err != nil {
				resp.Diagnostics.AddError(
					"error updating kubeoffload",
					err.Error(),
				)
				return
			}
		}

		// Set offload ID, use Offload ID as backupId for copy jobs
		plan.Kubeoffload_id = types.StringValue(updateKubeoffloadResp.Id)
		backupId = updateKubeoffloadResp.Id

		// Append copydef ID to original kubebackup request
		reqBody.Copydef = updateKubeoffloadResp.Id
		
		// For COPY jobs, clear the policy from the kubebackup (it should only be on the kubeoffload)
		if !plan.Policy_id.IsNull() && plan.Copy_pvs.ValueBool() {
			reqBody.Policy = ""
		}

		// Update kubebackup resource in CloudCasa
		putResp, err := r.Client.UpdateKubebackup(updateResp.Id, reqBody, updateResp.Etag)
		if err != nil {
			resp.Diagnostics.AddError(
				"error updating Kubebackup",
				err.Error(),
			)
			return
		}

		// Update fields in plan
		plan.Updated = types.StringValue(putResp.Updated)
		plan.Etag = types.StringValue(putResp.Etag)
		plan.Offload_etag = types.StringValue(updateKubeoffloadResp.Etag)
		
		// For COPY jobs, store the policy in copy_policy as well
		if !plan.Policy_id.IsNull() {
			plan.Copy_policy = types.StringValue(plan.Policy_id.ValueString())
		} else {
			// If policy_id is null but copy_pvs is true, ensure copy_policy is empty
			plan.Copy_policy = types.StringValue("")
		}
	} else {
		// For non-COPY jobs, ensure these fields have known values
		plan.Copy_policy = types.StringValue("")
		plan.Offload_etag = types.StringValue("")
		plan.Kubeoffload_id = types.StringValue("")
	}

	// If run_on_apply is false return now. Otherwise continue and run the job
	if !plan.Run.ValueBool() {
		diags = resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		return
	}

	err = r.runBackupJob(plan, backupId)
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
				"error deleting Kubeoffload resource",
				err.Error(),
			)
			return
		}
	}

	err := r.Client.DeleteKubebackup(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error deleting Kubebackup resource",
			err.Error(),
		)
		return
	}
}

func (r *resourceKubebackup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
