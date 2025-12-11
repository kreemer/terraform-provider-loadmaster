// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var (
	_ datasource.DataSource              = &RealServerDataSource{}
	_ datasource.DataSourceWithConfigure = &RealServerDataSource{}
)

func NewRealServerDataSource() datasource.DataSource {
	return &RealServerDataSource{}
}

type RealServerDataSource struct {
	client *api.Client
}

type RealServerDataSourceModel struct {
	Id               types.Int32  `tfsdk:"id"`
	VirtualServiceId types.String `tfsdk:"virtual_service_id"`
	Address          types.String `tfsdk:"address"`
	Port             types.Int32  `tfsdk:"port"`
	Weight           types.Int32  `tfsdk:"weight"`
	Forward          types.String `tfsdk:"forward"`
	Enable           types.Bool   `tfsdk:"enable"`
	Limit            types.Int32  `tfsdk:"limit"`
	Critical         types.Bool   `tfsdk:"critical"`
	Follow           types.Int32  `tfsdk:"follow"`
	DnsName          types.String `tfsdk:"dns_name"`
}

func (d *RealServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_real_server"
}

func (d *RealServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about a real server.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The real server id. This is also called `RIndex` in the LoadMaster API.",
				Required:            true,
			},
			"virtual_service_id": schema.StringAttribute{
				MarkdownDescription: "The id of the virtual service. This is also called `Index` in the LoadMaster API.",
				Required:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The address of the real server. Should be an IP address.",
				Computed:            true,
			},
			"port": schema.Int32Attribute{
				MarkdownDescription: "The port of the real server.",
				Computed:            true,
			},
			"weight": schema.Int32Attribute{
				MarkdownDescription: "The weight of the real server.",
				Computed:            true,
			},
			"forward": schema.StringAttribute{
				MarkdownDescription: "The forward of the real server.",
				Optional:            true,
				Computed:            true,
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "The enable of the real server.",
				Computed:            true,
			},
			"limit": schema.Int32Attribute{
				MarkdownDescription: "The limit of the real server.",
				Computed:            true,
			},
			"critical": schema.BoolAttribute{
				MarkdownDescription: "The critical of the real server.",
				Computed:            true,
			},
			"follow": schema.Int32Attribute{
				MarkdownDescription: "The follow of the real server.",
				Computed:            true,
			},
			"dns_name": schema.StringAttribute{
				MarkdownDescription: "The dns name of the real server.",
				Computed:            true,
			},
		},
	}
}

func (d *RealServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RealServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RealServerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operation := ClientBackoff(func() (*api.ListRealServerResponse, error) {
		return d.client.ShowRealServer(data.VirtualServiceId.ValueString(), "!"+strconv.Itoa(int(data.Id.ValueInt32())))
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read real server, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	real_server_response := response.Rs[len(response.Rs)-1]
	data.VirtualServiceId = types.StringValue(strconv.Itoa(int(real_server_response.VSIndex)))
	data.Address = types.StringValue(real_server_response.Address)
	data.Port = types.Int32Value(real_server_response.Port)
	data.Weight = types.Int32Value(real_server_response.Weight)
	data.Forward = types.StringValue(real_server_response.Forward)
	data.Enable = types.BoolPointerValue(real_server_response.Enable)
	data.Limit = types.Int32Value(real_server_response.Limit)
	data.Critical = types.BoolValue(*real_server_response.Critical)
	data.Follow = types.Int32Value(real_server_response.Follow)
	data.DnsName = types.StringValue(real_server_response.DnsName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
