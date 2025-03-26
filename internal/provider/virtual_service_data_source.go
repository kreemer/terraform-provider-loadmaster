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

type VirtualServiceDataSource struct {
	client *api.Client
}

type VirtualServiceDataSourceModel struct {
	Id       types.Int32  `tfsdk:"id"`
	Address  types.String `tfsdk:"address"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Type     types.String `tfsdk:"type"`
	Nickname types.String `tfsdk:"nickname"`
	Enabled  types.Bool   `tfsdk:"enabled"`
}

func (d *VirtualServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_service"
}

func (d *VirtualServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Use this data source to retrieve information about a virtual service.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The virtual service id. This is also called `Index` in the LoadMaster API.",
				Required:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The address of the virtual service. Should be an IP address of an interface attached to the LoadMaster.",
				Computed:            true,
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "The port of the virtual service.",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol of the virtual service, either `tcp` or `udp`.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the virtual service, either `gen`, `http`, `http2`, `ts`, `tls` or `log`.",
				Computed:            true,
			},
			"nickname": schema.StringAttribute{
				MarkdownDescription: "The nickname of the virtual service.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "If the virtual service is enabled.",
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
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *VirtualServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VirtualServiceDataSourceModel

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
	data.Address = types.StringValue(response.Address)
	data.Port = types.StringValue(response.Port)
	data.Protocol = types.StringValue(response.Protocol)
	data.Type = types.StringValue(response.VSType)
	data.Nickname = types.StringValue(response.NickName)
	data.Enabled = types.BoolValue(*response.Enable)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
