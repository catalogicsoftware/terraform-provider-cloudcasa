// Copyright 2025 Catalogic Software, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Apikey         types.String `tfsdk:"apikey"`
	CloudcasaUrl   types.String `tfsdk:"cloudcasa_url"`
	AllowInsecureTLS types.Bool   `tfsdk:"insecure_tls"`
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
			"cloudcasa_url": schema.StringAttribute{
				Description: "CloudCasa URL (for Selfhosted CloudCasa users). Defaults to https://home.cloudcasa.io",
				Optional:    true,
			},
			"insecure_tls": schema.BoolAttribute{
				Description: "Allow insecure TLS connections to CloudCasa. Defaults to false. Intended for Selfhosted CloudCasa servers with self-signed certificates.",
				Optional:    true,
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

	// Get CloudCasa URL if provided or use default
	var cloudcasaUrl *string
	if !config.CloudcasaUrl.IsNull() {
		value := config.CloudcasaUrl.ValueString()
		cloudcasaUrl = &value
	}

	// Get insecure TLS setting
	allowInsecureTLS := false
	if !config.AllowInsecureTLS.IsNull() {
		allowInsecureTLS = config.AllowInsecureTLS.ValueBool()
	}

	// Create a new CloudCasa client using the configuration values
	client, err := cloudcasa.NewClient(&apikey, cloudcasaUrl, allowInsecureTLS)
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
		NewResourceObjectstore,
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
