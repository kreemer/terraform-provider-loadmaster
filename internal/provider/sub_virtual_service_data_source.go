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
	_ datasource.DataSource              = &SubVirtualServiceDataSource{}
	_ datasource.DataSourceWithConfigure = &SubVirtualServiceDataSource{}
)

func NewSubVirtualServiceDataSource() datasource.DataSource {
	return &SubVirtualServiceDataSource{}
}

type SubVirtualServiceDataSource struct {
	client *api.Client
}

type SubVirtualServiceDataSourceModel struct {
	Id               types.Int32  `tfsdk:"id"`
	VirtualServiceId types.Int32  `tfsdk:"virtual_service_id"`
	Nickname         types.String `tfsdk:"nickname"`
}

func (d *SubVirtualServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sub_virtual_service"
}

func (d *SubVirtualServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about a sub virtual service.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The sub virtual service id. This is also called `Index` in the LoadMaster API.",
				Required:            true,
			},
			"virtual_service_id": schema.Int32Attribute{
				MarkdownDescription: "The id of the virtual service. This is also called `Index` in the LoadMaster API.",
				Computed:            true,
			},
			"nickname": schema.StringAttribute{
				MarkdownDescription: "The nickname of the sub virtual service.",
				Computed:            true,
			},
		},
	}
}

func (d *SubVirtualServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SubVirtualServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SubVirtualServiceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt32())
	response, err := d.client.ShowVirtualService(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read virtual service, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Id = types.Int32Value(int32(response.Index))
	data.VirtualServiceId = types.Int32Value(int32(response.MasterVSID))
	data.Nickname = types.StringValue(response.NickName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
