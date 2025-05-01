package cloudcasa

import (
	"context"
	"fmt"

	cloudcasa "terraform-provider-cloudcasa/cloudcasa/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceObjectstore{}
	_ resource.ResourceWithConfigure   = &resourceObjectstore{}
	_ resource.ResourceWithImportState = &resourceObjectstore{}
)

func NewResourceObjectstore() resource.Resource {
	return &resourceObjectstore{}
}

type resourceObjectstore struct {
	Client *cloudcasa.Client
}

// objectstoreResourceModel maps the resource schema data.
type objectstoreResourceModel struct {
	Name              types.String `tfsdk:"name"`
	Id                types.String `tfsdk:"id"`
	Private           types.Bool   `tfsdk:"private"`
	ProxyCluster      types.String `tfsdk:"proxy_cluster"`
	ProviderType      types.String `tfsdk:"provider_type"`
	// Common fields
	BucketName        types.String `tfsdk:"bucket_name"`
	Region            types.String `tfsdk:"region"`
	SkipTlsValidation types.Bool   `tfsdk:"skip_tls_validation"`
	// S3-specific fields
	EndpointUrl       types.String `tfsdk:"endpoint_url"`
	AccessKey         types.String `tfsdk:"access_key"`
	SecretKey         types.String `tfsdk:"secret_key"`
	// Azure-specific fields
	SubscriptionId    types.String `tfsdk:"subscription_id"`
	TenantId          types.String `tfsdk:"tenant_id"`
	ClientId          types.String `tfsdk:"client_id"`
	ClientSecret      types.String `tfsdk:"client_secret"`
	Cloud             types.String `tfsdk:"cloud"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	StorageAccountName types.String `tfsdk:"storage_account_name"`
	// Common fields
	Updated           types.String `tfsdk:"updated"`
	Created           types.String `tfsdk:"created"`
	Etag              types.String `tfsdk:"etag"`
}

// Metadata returns the data source type name.
func (r *resourceObjectstore) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objectstore"
}

// Schema defines the schema for the resource.
func (r *resourceObjectstore) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CloudCasa objectstore configuration",
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
			"private": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable for storage location which is isolated from CloudCasa server. Defaults to false.",
			},
			"proxy_cluster": schema.StringAttribute{
				Optional:    true,
				Description: "The proxy cluster ID to use for connecting to the object store. Required if private is true.",
			},
			"provider_type": schema.StringAttribute{
				Required:    true,
				Description: "Object store provider type. Supported values: 's3' and 'azure'.",
				Validators: []validator.String{
					stringvalidator.OneOf("s3", "azure"),
				},
			},
			"bucket_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the S3 bucket or Azure storage container. Required for 's3' and 'azure' provider types.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The region for the storage provider. Required for 's3' and 'azure' provider types.",
			},
			"skip_tls_validation": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to skip TLS validation for the object store. Defaults to false.",
			},
			"endpoint_url": schema.StringAttribute{
				Optional:    true,
				Description: "The endpoint URL for the S3 provider. Required if provider_type is 's3'.",
			},
			"subscription_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Azure subscription ID. Required if provider_type is 'azure'.",
				Sensitive:   true,
			},
			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Azure tenant ID. Required if provider_type is 'azure'.",
				Sensitive:   true,
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Azure client ID. Required if provider_type is 'azure'.",
				Sensitive:   true,
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Description: "The Azure client secret. Required if provider_type is 'azure'.",
				Sensitive:   true,
			},
			"cloud": schema.StringAttribute{
				Optional:    true,
				Description: "The Azure cloud type. Required if provider_type is 'azure'. Default: 'Public'.",
			},
			"resource_group_name": schema.StringAttribute{
				Optional:    true,
				Description: "The Azure resource group name. Required if provider_type is 'azure'.",
			},
			"storage_account_name": schema.StringAttribute{
				Optional:    true,
				Description: "The Azure storage account name. Required if provider_type is 'azure'.",
			},
			"access_key": schema.StringAttribute{
				Optional:    true,
				Description: "The access key for authenticating with the S3 provider.",
				Sensitive:   true,
			},
			"secret_key": schema.StringAttribute{
				Optional:    true,
				Description: "The secret key for authenticating with the S3 provider.",
				Sensitive:   true,
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *resourceObjectstore) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.Client = req.ProviderData.(*cloudcasa.Client)
}

// CreateObjectstoreFromPlan creates an Objectstore struct with the required fields from the plan
func CreateObjectstoreFromPlan(plan objectstoreResourceModel) (*cloudcasa.Objectstore, error) {
	// Create objectstore object
	objectstore := &cloudcasa.Objectstore{
		Name:                      plan.Name.ValueString(),
		Private:                   plan.Private.ValueBool(),
		SkipTlsCertificateValidation: plan.SkipTlsValidation.ValueBool(),
	}

	// Translate from user-facing "s3" to CloudCasa API's "aws"
	if plan.ProviderType.ValueString() == "s3" {
		objectstore.ProviderType = "aws"
	} else {
		objectstore.ProviderType = plan.ProviderType.ValueString()
	}

	// Validate provider_type - current supported values are 's3' and 'azure'
	if plan.ProviderType.ValueString() != "s3" && plan.ProviderType.ValueString() != "azure" {
		return nil, fmt.Errorf("provider_type must be 's3' or 'azure', got: %s", plan.ProviderType.ValueString())
	}

	// Handle proxy cluster for private objectstores
	if !plan.ProxyCluster.IsNull() {
		// Add single cluster to cluster list
		objectstore.ProxyClusterList = []string{plan.ProxyCluster.ValueString()}
	}

	// Common fields for all provider types
	if plan.BucketName.IsNull() {
		return nil, fmt.Errorf("bucket_name is required for all provider types")
	}
	if plan.Region.IsNull() {
		return nil, fmt.Errorf("region is required for all provider types")
	}
	objectstore.BucketName = plan.BucketName.ValueString()
	objectstore.Region = plan.Region.ValueString()

	// Add S3-specific fields if provider_type is 's3'
	if plan.ProviderType.ValueString() == "s3" {
		// Validate all required fields for S3
		if plan.EndpointUrl.IsNull() {
			return nil, fmt.Errorf("endpoint_url is required when provider_type is 's3'")
		}
		
		// Set up the S3Provider with endpoint
		objectstore.S3Provider.Endpoint = plan.EndpointUrl.ValueString()

		// Add credentials if provided - access_key and secret_key should be provided together
		if !plan.AccessKey.IsNull() && !plan.SecretKey.IsNull() {
			objectstore.S3Provider.Credentials.AccessKey = plan.AccessKey.ValueString()
			objectstore.S3Provider.Credentials.SecretKey = plan.SecretKey.ValueString()
		} else if !plan.AccessKey.IsNull() || !plan.SecretKey.IsNull() {
			// Only one of access_key or secret_key was provided
			return nil, fmt.Errorf("both access_key and secret_key must be provided together for authentication")
		}
	}
	
	// Add Azure-specific fields if provider_type is 'azure'
	if plan.ProviderType.ValueString() == "azure" {
		// Validate all required fields for Azure
		if plan.SubscriptionId.IsNull() {
			return nil, fmt.Errorf("subscription_id is required when provider_type is 'azure'")
		}
		if plan.TenantId.IsNull() {
			return nil, fmt.Errorf("tenant_id is required when provider_type is 'azure'")
		}
		if plan.ClientId.IsNull() {
			return nil, fmt.Errorf("client_id is required when provider_type is 'azure'")
		}
		if plan.ClientSecret.IsNull() {
			return nil, fmt.Errorf("client_secret is required when provider_type is 'azure'")
		}
		if plan.ResourceGroupName.IsNull() {
			return nil, fmt.Errorf("resource_group_name is required when provider_type is 'azure'")
		}
		if plan.StorageAccountName.IsNull() {
			return nil, fmt.Errorf("storage_account_name is required when provider_type is 'azure'")
		}
		
		// Set up Azure credentials
		objectstore.S3Provider.Credentials.SubscriptionId = plan.SubscriptionId.ValueString()
		objectstore.S3Provider.Credentials.TenantId = plan.TenantId.ValueString()
		objectstore.S3Provider.Credentials.ClientId = plan.ClientId.ValueString()
		objectstore.S3Provider.Credentials.ClientSecret = plan.ClientSecret.ValueString()
		
		// Set up Azure-specific provider fields
		objectstore.S3Provider.ResourceGroupName = plan.ResourceGroupName.ValueString()
		objectstore.S3Provider.StorageAccountName = plan.StorageAccountName.ValueString()
		
		// Set cloud if provided, otherwise use default
		if !plan.Cloud.IsNull() {
			objectstore.S3Provider.Cloud = plan.Cloud.ValueString()
		} else {
			objectstore.S3Provider.Cloud = "Public" // Default cloud type
		}
	}

	// Validate that proxy_cluster is provided if private is true
	if plan.Private.ValueBool() {
		if plan.ProxyCluster.IsNull() {
			return nil, fmt.Errorf("proxy_cluster is required when private is true")
		}
		// Ensure we have a valid non-empty cluster ID
		if plan.ProxyCluster.ValueString() == "" {
			return nil, fmt.Errorf("proxy_cluster cannot be empty when private is true")
		}
	}

	return objectstore, nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *resourceObjectstore) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan objectstoreResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create request object from plan
	objectstore, err := CreateObjectstoreFromPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Objectstore",
			err.Error(),
		)
		return
	}

	// Create objectstore in cloudcasa
	createResp, err := r.Client.CreateObjectstore(*objectstore)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Objectstore",
			err.Error(),
		)
		return
	}

	// Set fields in plan
	plan.Id = types.StringValue(createResp.Id)
	plan.Name = types.StringValue(createResp.Name)
	plan.Created = types.StringValue(createResp.Created)
	plan.Updated = types.StringValue(createResp.Updated)
	plan.Etag = types.StringValue(createResp.Etag)

	// Save state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *resourceObjectstore) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state objectstoreResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed objectstore value from CloudCasa
	objectstoreId := state.Id.ValueString()
	objectstoreResp, err := r.Client.GetObjectstore(objectstoreId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading CloudCasa Objectstore",
			fmt.Sprintf("Could not read CloudCasa objectstore ID %s: %s", objectstoreId, err),
		)
		return
	}

	// Update state with refreshed values
	state.Id = types.StringValue(objectstoreResp.Id)
	state.Name = types.StringValue(objectstoreResp.Name)
	state.Private = types.BoolValue(objectstoreResp.Private)
	
	// Translate provider_type from API to user-facing value
	if objectstoreResp.ProviderType == "aws" {
		state.ProviderType = types.StringValue("s3")
	} else {
		state.ProviderType = types.StringValue(objectstoreResp.ProviderType)
	}
	
	state.Created = types.StringValue(objectstoreResp.Created)
	state.Updated = types.StringValue(objectstoreResp.Updated)
	state.Etag = types.StringValue(objectstoreResp.Etag)
	
	// Handle proxy clusters (take first one if available)
	if len(objectstoreResp.ProxyClusterList) > 0 {
		state.ProxyCluster = types.StringValue(objectstoreResp.ProxyClusterList[0])
	} else {
		state.ProxyCluster = types.StringNull()
	}
	
	// Set S3-specific values if present
	if objectstoreResp.ProviderType == "aws" {
		if objectstoreResp.BucketName != "" {
			state.BucketName = types.StringValue(objectstoreResp.BucketName)
		}
		
		if objectstoreResp.S3Provider.Endpoint != "" {
			state.EndpointUrl = types.StringValue(objectstoreResp.S3Provider.Endpoint)
		}
		
		if objectstoreResp.Region != "" {
			state.Region = types.StringValue(objectstoreResp.Region)
		}
		
		state.SkipTlsValidation = types.BoolValue(objectstoreResp.SkipTlsCertificateValidation)
		
		// We only set the access key if it's returned, secret key is almost always masked
		if objectstoreResp.S3Provider.Credentials.AccessKey != "" {
			state.AccessKey = types.StringValue(objectstoreResp.S3Provider.Credentials.AccessKey)
		}
		
		// We don't set the secret key from the response as it's often masked
		// The API might not return the actual secret value for security reasons
	} else if objectstoreResp.ProviderType == "azure" {
		// Set Azure-specific values
		if objectstoreResp.BucketName != "" {
			state.BucketName = types.StringValue(objectstoreResp.BucketName)
		}
		
		if objectstoreResp.Region != "" {
			state.Region = types.StringValue(objectstoreResp.Region)
		}
		
		// Set Azure provider details
		if objectstoreResp.S3Provider.Cloud != "" {
			state.Cloud = types.StringValue(objectstoreResp.S3Provider.Cloud)
		}
		
		if objectstoreResp.S3Provider.ResourceGroupName != "" {
			state.ResourceGroupName = types.StringValue(objectstoreResp.S3Provider.ResourceGroupName)
		}
		
		if objectstoreResp.S3Provider.StorageAccountName != "" {
			state.StorageAccountName = types.StringValue(objectstoreResp.S3Provider.StorageAccountName)
		}
		
		// Set Azure credentials
		if objectstoreResp.S3Provider.Credentials.SubscriptionId != "" {
			state.SubscriptionId = types.StringValue(objectstoreResp.S3Provider.Credentials.SubscriptionId)
		}
		
		if objectstoreResp.S3Provider.Credentials.TenantId != "" {
			state.TenantId = types.StringValue(objectstoreResp.S3Provider.Credentials.TenantId)
		}
		
		if objectstoreResp.S3Provider.Credentials.ClientId != "" {
			state.ClientId = types.StringValue(objectstoreResp.S3Provider.Credentials.ClientId)
		}
		
		state.SkipTlsValidation = types.BoolValue(objectstoreResp.SkipTlsCertificateValidation)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourceObjectstore) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan objectstoreResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state objectstoreResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create request object from plan
	objectstore, err := CreateObjectstoreFromPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Objectstore",
			err.Error(),
		)
		return
	}

	// Update objectstore in cloudcasa
	objectstoreId := state.Id.ValueString()
	objectstoreResp, err := r.Client.UpdateObjectstore(objectstoreId, *objectstore, state.Etag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating CloudCasa Objectstore",
			fmt.Sprintf("Could not update CloudCasa objectstore ID %s: %s", objectstoreId, err),
		)
		return
	}

	// Update state with refreshed values
	plan.Id = types.StringValue(objectstoreResp.Id)
	plan.Name = types.StringValue(objectstoreResp.Name)
	plan.Updated = types.StringValue(objectstoreResp.Updated)
	plan.Created = types.StringValue(objectstoreResp.Created)
	plan.Etag = types.StringValue(objectstoreResp.Etag)
	
	// The API might return other updated fields, so let's use the Read function to get a full refresh
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceObjectstore) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state objectstoreResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete objectstore
	objectstoreId := state.Id.ValueString()
	err := r.Client.DeleteObjectstore(objectstoreId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting CloudCasa Objectstore",
			fmt.Sprintf("Could not delete CloudCasa objectstore ID %s: %s", objectstoreId, err),
		)
		return
	}
}

// ImportState imports an existing objectstore resource for management by Terraform.
func (r *resourceObjectstore) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}