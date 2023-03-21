package cloudcasa

import (
	"context"
	"os/exec"
	"time"

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
	_ resource.Resource                = &resourceKubecluster{}
	_ resource.ResourceWithConfigure   = &resourceKubecluster{}
	_ resource.ResourceWithImportState = &resourceKubecluster{}
)

func NewResourceKubecluster() resource.Resource {
	return &resourceKubecluster{}
}

type resourceKubecluster struct {
	Client *cloudcasa.Client
}

// kubeclusterResourceModel maps the resource schema data.
type kubeclusterResourceModel struct {
	Name          types.String `tfsdk:"name"`
	Id            types.String `tfsdk:"id"`
	Auto_install  types.Bool   `tfsdk:"auto_install"`
	Cc_user_email types.String `tfsdk:"cc_user_email"`
	Updated       types.String `tfsdk:"updated"`
	Created       types.String `tfsdk:"created"`
	Org_id        types.String `tfsdk:"org_id"`
	Etag          types.String `tfsdk:"etag"`
	Status        types.Map    `tfsdk:"status"`
	Links         types.Map    `tfsdk:"links"`
	Agent_url     types.String `tfsdk:"agent_url"`
}

// API Response Objects
type CreateKubeclusterResp struct {
	Id            string   `json:"_id"`
	Name          string   `json:"name"`
	Cc_user_email string   `json:"cc_user_email"`
	Updated       string   `json:"_updated"`
	Created       string   `json:"_created"`
	Etag          string   `json:"_etag"`
	Org_id        string   `json:"org_id"`
	Status        string   `json:"_status"`
	Links         struct{} `json:"_links"`
}

// Metadata returns the data source type name.
func (r *resourceKubecluster) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubecluster"
}

// Schema defines the schema for the resource.
func (r *resourceKubecluster) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"auto_install": schema.BoolAttribute{
				Optional: true,
			},
			"cc_user_email": schema.StringAttribute{
				Computed: true,
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
			"org_id": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"links": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"agent_url": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *resourceKubecluster) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.Client = req.ProviderData.(*cloudcasa.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *resourceKubecluster) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan kubeclusterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Define kubecluster object
	reqBody := map[string]string{
		"name": plan.Name.ValueString(),
	}

	// Create kubecluster in cloudcasa
	createResp, err := r.Client.CreateKubecluster(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Kubecluster",
			err.Error(),
		)
		return
	}

	// Set fields in plan
	plan.Id = types.StringValue(createResp.Id)
	plan.Name = types.StringValue(createResp.Name)
	plan.Cc_user_email = types.StringValue(createResp.Cc_user_email)
	plan.Created = types.StringValue(createResp.Created)
	plan.Updated = types.StringValue(createResp.Updated)
	plan.Etag = types.StringValue(createResp.Etag)
	plan.Org_id = types.StringValue(createResp.Org_id)

	// if auto_install is false return now. Otherwise proceed with agent installation
	if !plan.Auto_install.ValueBool() {
		plan.Agent_url = types.StringNull()
		plan.Links = types.MapNull(types.StringType)
		plan.Status = types.MapNull(types.StringType)
		diags = resp.State.Set(ctx, plan)
		return
	}

	var kubeclusterStatus cloudcasa.KubeclusterStatus

	// wait 1m for agent URL
	for i := 1; i < 12; i++ {
		getKubeclusterResp, err := r.Client.GetKubecluster(createResp.Id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error checking Kubecluster status after creation",
				err.Error(),
			)
			return
		}
		kubeclusterStatus = getKubeclusterResp.Status
		if len(kubeclusterStatus.Agent_url) > 0 {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Check that Agent URL was fetched successfully
	if len(kubeclusterStatus.Agent_url) == 0 {
		resp.Diagnostics.AddError(
			"Failed to get Agent URL for kubecluster",
			"Timed out waiting for Agent URL: "+err.Error(),
		)
		return
	}

	// Set agent url from response
	plan.Agent_url = types.StringValue(kubeclusterStatus.Agent_url)

	// TODO: add tip to make sure kubeconfig env var is set?
	// OR we can accept kubeconfig as an input option?
	kubectlCmd := exec.Command("kubectl", "apply", "-f", kubeclusterStatus.Agent_url)
	_, err = kubectlCmd.Output()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to apply kubeagent manifest",
			err.Error(),
		)
		return
	}

	// Now wait for cluster to be ACTIVE
	// Wait 5min?
	for i := 1; i < 60; i++ {
		getKubeclusterResp, err := r.Client.GetKubecluster(createResp.Id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error checking Kubecluster status after applying Agent manifest",
				err.Error(),
			)
			return
		}
		kubeclusterStatus = getKubeclusterResp.Status
		if kubeclusterStatus.State == "ACTIVE" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Check if state was set to ACTIVE
	if kubeclusterStatus.State != "ACTIVE" {
		resp.Diagnostics.AddError(
			"CloudCasa Agent installation failed",
			"Timed out waiting for cluster to reach ACTIVE state: "+err.Error(),
		)
		return
	}

	// Save state before returning
	plan.Links = types.MapNull(types.StringType)
	plan.Status = types.MapNull(types.StringType)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// TODO: if auto-install is enabled, we should check the status on each refresh
// and apply the agent with current kubeconfig if not active.
// TODO: CHECK STATUS!
func (r *resourceKubecluster) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state kubeclusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Kubecluster from CloudCasa
	getKubeclusterResp, err := r.Client.GetKubecluster(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Kubecluster from CloudCasa",
			"Could not read Kubecluster with ID "+state.Id.ValueString()+" :"+err.Error(),
		)
		return
	}

	// Overwrite values with refreshed state
	// Set fields in plan
	state.Id = types.StringValue(getKubeclusterResp.Id)
	state.Name = types.StringValue(getKubeclusterResp.Name)
	state.Cc_user_email = types.StringValue(getKubeclusterResp.Cc_user_email)
	state.Created = types.StringValue(getKubeclusterResp.Created)
	state.Updated = types.StringValue(getKubeclusterResp.Updated)
	state.Etag = types.StringValue(getKubeclusterResp.Etag)
	state.Org_id = types.StringValue(getKubeclusterResp.Org_id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourceKubecluster) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan kubeclusterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from TF state
	// (need etag value to edit the existing object)
	// TODO: Load etag from TF state OR GET from CC API?
	var state kubeclusterResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	// NAME is editable
	reqBody := map[string]string{
		"name": plan.Name.ValueString(),
	}

	// Update kubecluster in CloudCasa
	updateResp, err := r.Client.UpdateKubecluster(plan.Id.ValueString(), reqBody, state.Etag.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kubecluster",
			err.Error(),
		)
		return
	}

	// Overwrite values with refreshed state
	// Set fields in plan
	plan.Id = types.StringValue(updateResp.Id)
	plan.Name = types.StringValue(updateResp.Name)
	plan.Cc_user_email = types.StringValue(updateResp.Cc_user_email)
	plan.Created = types.StringValue(updateResp.Created)
	plan.Updated = types.StringValue(updateResp.Updated)
	plan.Etag = types.StringValue(updateResp.Etag)
	plan.Org_id = types.StringValue(updateResp.Org_id)

	// Check that Agent URL was fetched successfully
	kubeclusterStatus := updateResp.Status
	if len(kubeclusterStatus.Agent_url) == 0 {
		plan.Agent_url = types.StringValue("")
	} else {
		plan.Agent_url = types.StringValue(kubeclusterStatus.Agent_url)
	}

	// Save state before returning
	plan.Links = types.MapNull(types.StringType)
	plan.Status = types.MapNull(types.StringType)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceKubecluster) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state kubeclusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteKubecluster(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kubecluster resource",
			"Could not delete Kubecluster, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *resourceKubecluster) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
