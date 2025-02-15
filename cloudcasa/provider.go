package cloudcasa

import (
	"context"

	cloudcasa "terraform-provider-cloudcasa/cloudcasa/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
}

// Metadata returns the provider type name.
func (p *cloudcasaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudcasa"
}

// Schema defines the provider-level schema for configuration data.
func (p *cloudcasaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with CloudCasa",
		Attributes: map[string]schema.Attribute{
			"apikey": schema.StringAttribute{
				Description: "CloudCasa API Key for authentication. Visit https://docs.cloudcasa.io/help/apikeys.html for more details",
				Required:    true,
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

	// Validate and set config values
	apikey := config.Apikey.ValueString()
	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"missing CloudCasa API Key",
			"the provider recieved an empty value for apikey."+
				"Set a valid API key - refer to https://docs.cloudcasa.io/help/apikeys.html for details",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new CloudCasa client using the configuration values
	client, err := cloudcasa.NewClient(&apikey)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to create CloudCasa API client",
			"an unexpected error occurred when creating the CloudCasa API client. "+
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
		NewResourceKubebackup,
		NewResourcePolicy,
	}
}

// ConvertTfStringList converts a list of TF StringValues to a list of Go string
func ConvertTfStringList(tfList []basetypes.StringValue) []string {
	var stringList []string
	for _, v := range tfList {
		stringList = append(stringList, v.ValueString())
	}
	return stringList
}

// ConvertStringListTf converts a list strings to list of TF StringValues
func ConvertStringListTf(stringList []string) []basetypes.StringValue {
	var tfList []basetypes.StringValue
	for _, v := range stringList {
		tfList = append(tfList, types.StringValue(v))
	}
	return tfList
}
