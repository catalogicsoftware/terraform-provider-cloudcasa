package cloudcasa

import (
	"context"

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
	_ resource.Resource                = &resourcePolicy{}
	_ resource.ResourceWithConfigure   = &resourcePolicy{}
	_ resource.ResourceWithImportState = &resourcePolicy{}
)

func NewResourcePolicy() resource.Resource {
	return &resourcePolicy{}
}

type resourcePolicy struct {
	Client *cloudcasa.Client
}

type policyResourceModel struct {
	Id        types.String          `tfsdk:"id"`
	Name      types.String          `tfsdk:"name"`
	Timezone  types.String          `tfsdk:"timezone"`
	Schedules []policyScheduleModel `tfsdk:"schedules"`
	Updated   types.String          `tfsdk:"updated"`
	Created   types.String          `tfsdk:"created"`
	Etag      types.String          `tfsdk:"etag"`
}

type policyScheduleModel struct {
	Retention types.Int64  `tfsdk:"retention"`
	Locked    types.Bool   `tfsdk:"locked"`
	Cron_spec types.String `tfsdk:"cron_spec"`
}

// Metadata returns the data source type name.
func (r *resourcePolicy) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

// Schema defines the schema for the resource.
func (r *resourcePolicy) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
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
			"timezone": schema.StringAttribute{
				Required:    true,
				Description: "TZ string for the defined Cronjob. Ex: America/New_York",
			},
			"schedules": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"retention": schema.Int64Attribute{
							Required:    true,
							Description: "Number of days to retain backup data for",
						},
						"locked": schema.BoolAttribute{
							Required:    true,
							Description: "Enable SafeLock for backups (CloudCasa Premium only)",
						},
						"cron_spec": schema.StringAttribute{
							Required:    true,
							Description: "Cron expression for backup schedule. Ex: '0 4 * * sun'",
						},
					},
				},
			},
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
func (r *resourcePolicy) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.Client = req.ProviderData.(*cloudcasa.Client)
}

// setPlanFromPolicy sets Terraform state values from a given CloudCasa Policy resource
// TODO: implement this function for other resources
func (plan *policyResourceModel) setPlanFromPolicy(policy *cloudcasa.Policy) error {
	// Set fields in plan
	plan.Id = types.StringValue(policy.Id)
	plan.Timezone = types.StringValue(policy.Timezone)
	plan.Created = types.StringValue(policy.Created)
	plan.Updated = types.StringValue(policy.Updated)
	plan.Etag = types.StringValue(policy.Etag)

	// Remove existing Schedules from the plan so we can overwrite
	plan.Schedules = []policyScheduleModel{}

	// TODO: do this for backups pre_hooks
	// Convert Schedules body from CC to TF
	for _, v := range policy.Schedules {
		thisSchedule := policyScheduleModel{
			Retention: types.Int64Value(v.RetainDays),
			Locked:    types.BoolValue(v.Locked),
			Cron_spec: types.StringValue(v.Schedule.CronSpec),
		}
		plan.Schedules = append(plan.Schedules, thisSchedule)
	}
	return nil
}

// createPolicyFromPlan initializes a cloudcasa.Policy from TF values
func createPolicyFromPlan(plan policyResourceModel) (cloudcasa.Policy, error) {
	// Initialize CC policy body
	policy := cloudcasa.Policy{
		Name:     plan.Name.ValueString(),
		Timezone: plan.Timezone.ValueString(),
	}

	// Create body for Schedules
	for _, v := range plan.Schedules {
		thisSchedule := cloudcasa.PolicySchedule{
			RetainDays: v.Retention.ValueInt64(),
			Locked:     v.Locked.ValueBool(),
			Schedule: cloudcasa.ScheduleStruct{
				CronSpec: v.Cron_spec.ValueString(),
			},
		}
		policy.Schedules = append(policy.Schedules, thisSchedule)
	}
	return policy, nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *resourcePolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize CC policy body
	reqBody, err := createPolicyFromPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"error initalizing TF plan for Policy",
			err.Error(),
		)
		return
	}

	// Create Policy resource in CloudCasa
	createResp, err := r.Client.CreatePolicy(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"error creating Policy",
			err.Error(),
		)
		return
	}

	err = plan.setPlanFromPolicy(createResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating TF state for created Policy",
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
func (r *resourcePolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Kubecluster from CloudCasa
	policy, err := r.Client.GetPolicy(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading Policy with ID "+state.Id.ValueString(),
			err.Error(),
		)
		return
	}

	// Update state with values from CloudCasa
	err = state.setPlanFromPolicy(policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating TF state for policy with ID "+state.Id.ValueString(),
			err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourcePolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from TF state
	// need etag value to edit the existing object
	var state policyResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize CC policy body
	reqBody, err := createPolicyFromPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"error initalizing TF plan for Policy",
			err.Error(),
		)
		return
	}

	// Update Policy resource in CloudCasa
	updateResp, err := r.Client.UpdatePolicy(plan.Id.ValueString(), reqBody, state.Etag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating Policy",
			err.Error(),
		)
		return
	}

	err = plan.setPlanFromPolicy(updateResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating TF state for updated policy",
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
func (r *resourcePolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeletePolicy(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error deleting Policy resource",
			err.Error(),
		)
		return
	}
}

func (r *resourcePolicy) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
