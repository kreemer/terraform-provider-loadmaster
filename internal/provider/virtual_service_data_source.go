// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var (
	_ datasource.DataSource              = &VirtualServiceDataSource{}
	_ datasource.DataSourceWithConfigure = &VirtualServiceDataSource{}
)

func NewVirtualServiceDataSource() datasource.DataSource {
	return &VirtualServiceDataSource{}
}

// ExampleDataSource defines the data source implementation.
type VirtualServiceDataSource struct {
	client *api.Client
}

// ExampleDataSourceModel describes the data source data model.
type VirtualServiceDataSourceModel struct {
	Id       types.Int32  `tfsdk:"id"`
	Address  types.String `tfsdk:"address"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

func (d *VirtualServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_service"
}

func (d *VirtualServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "Virtual service index",
				Required:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "Virtual service index",
				Computed:            true,
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "Virtual service index",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Virtual service index",
				Computed:            true,
			},
		},
	}
}

func (d *VirtualServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *VirtualServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VirtualServiceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	id := int(data.Id.ValueInt32())

	response, err := d.client.ShowVirtualService(id)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.Int32Value(int32(response.Index))
	data.Address = types.StringValue(response.Address)
	data.Port = types.StringValue(response.Port)
	data.Protocol = types.StringValue(response.Protocol)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
