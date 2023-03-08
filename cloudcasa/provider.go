package cloudcasa

import (
	"context"
	"os"

	cloudcasa "terraform-provider-cloudcasa/cloudcasa/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &cloudcasaProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &cloudcasaProvider{}
}

// cloudcasaProvider is the provider implementation.
type cloudcasaProvider struct{}

// cloudcasaProviderModel maps provider schema data to a Go type.
type cloudcasaProviderModel struct {
	Apikey types.String `tfsdk:"apikey"`
	Email  types.String `tfsdk:"email"`
}

// Metadata returns the provider type name.
func (p *cloudcasaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudcasa"
}

// Schema defines the provider-level schema for configuration data.
func (p *cloudcasaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Optional: true,
			},
			"apikey": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Configure prepares a CloudCasa API client for data sources and resources.
func (p *cloudcasaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config cloudcasaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Apikey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown CloudCasa API Key",
			"The provider cannot create the CloudCasa API client as there is an unknown configuration value for the CloudCasa API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CLOUDCASA_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	apikey := os.Getenv("CLOUDCASA_KEY")

	if !config.Apikey.IsNull() {
		apikey = config.Apikey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing CloudCasa API Key",
			"The provider cannot create the CloudCasa API client as there is a missing or empty value for the CloudCasa API key. "+
				"Set the host value in the configuration or use the CLOUDCASA_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new CloudCasa client using the configuration values
	client, err := cloudcasa.NewClient(&apikey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create CloudCasa API Client",
			"An unexpected error occurred when creating the CloudCasa API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"CloudCasa Client Error: "+err.Error(),
		)
		return
	}

	// Make the CloudCasa client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *cloudcasaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *cloudcasaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceKubecluster,
	}
}
