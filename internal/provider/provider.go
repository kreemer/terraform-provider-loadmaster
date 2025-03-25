// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kreemer/loadmaster-go-client/api"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &LoadMasterProvider{}
var _ provider.ProviderWithFunctions = &LoadMasterProvider{}
var _ provider.ProviderWithEphemeralResources = &LoadMasterProvider{}

// ScaffoldingProvider defines the provider implementation.
type LoadMasterProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type LoadMasterProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	ApiKey   types.String `tfsdk:"api_key"`
}

func (p *LoadMasterProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "loadmaster"
	resp.Version = p.version
}

func (p *LoadMasterProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "IP address and port of the loadmaster instance",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for the loadmaster instance",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for the loadmaster instance",
				Optional:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for the loadmaster instance",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *LoadMasterProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LoadMasterProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown LoadMaster API Host",
			"The provider cannot create the LoadMaster API client as there is an unknown configuration value for the LoadMaster host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOADMASTER_HOST environment variable.",
		)
	}

	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown LoadMaster API Username",
			"The provider cannot create the LoadMaster API client as there is an unknown configuration value for the LoadMaster username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOADMASTER_USERNAME environment variable.",
		)
	}

	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown LoadMaster API Password",
			"The provider cannot create the LoadMaster API client as there is an unknown configuration value for the LoadMaster password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOADMASTER_PASSWORD environment variable.",
		)
	}

	if data.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown LoadMaster API Key",
			"The provider cannot create the LoadMaster API client as there is an unknown configuration value for the LoadMaster API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOADMASTER_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("LOADMASTER_HOST")
	username := os.Getenv("LOADMASTER_USERNAME")
	password := os.Getenv("LOADMASTER_PASSWORD")
	apiKey := os.Getenv("LOADMASTER_API_KEY")

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing HashiCups API Host",
			"The provider cannot create the HashiCups API client as there is a missing or empty value for the HashiCups API host. "+
				"Set the host value in the configuration or use the HASHICUPS_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" && (username == "" || password == "") {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing HashiCups API Username",
			"The provider cannot create the HashiCups API client as there is a missing or empty value for the HashiCups API username. "+
				"Set the username value in the configuration or use the HASHICUPS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var client *api.Client
	if apiKey != "" {
		client = api.NewClientWithApiKey(apiKey, host)
	} else {
		client = api.NewClientWithUsernamePassword(username, password, host)
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }
}

func (p *LoadMasterProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVirtualServiceResource,
	}
}

func (p *LoadMasterProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return nil
}

func (p *LoadMasterProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVirtualServiceDataSource,
	}
}

func (p *LoadMasterProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LoadMasterProvider{
			version: version,
		}
	}
}
