// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var _ resource.Resource = &RealServerResource{}
var _ resource.ResourceWithImportState = &RealServerResource{}

func NewRealServerResource() resource.Resource {
	return &RealServerResource{}
}

type RealServerResource struct {
	client *api.Client
}

type RealServerResourceModel struct {
	Id               types.Int32  `tfsdk:"id"`
	VirtualServiceId types.Int32  `tfsdk:"virtual_service_id"`
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

func (r *RealServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_real_server"
}

func (r *RealServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages real server.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the real server. This is also called `RIndex` in the LoadMaster API.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"virtual_service_id": schema.Int32Attribute{
				MarkdownDescription: "The id of the virtual service. This is also called `VIndex` in the LoadMaster API.",
				Required:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The address of the real server. Should be an IP address.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int32Attribute{
				MarkdownDescription: "The port of the real server.",
				Required:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"weight": schema.Int32Attribute{
				MarkdownDescription: "The weight of the real server.",
				Optional:            true,
				Computed:            true,
			},
			"forward": schema.StringAttribute{
				MarkdownDescription: "The forward of the real server.",
				Optional:            true,
				Computed:            true,
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "The enable of the real server.",
				Optional:            true,
				Computed:            true,
			},
			"limit": schema.Int32Attribute{
				MarkdownDescription: "The limit of the real server.",
				Optional:            true,
				Computed:            true,
			},
			"critical": schema.BoolAttribute{
				MarkdownDescription: "The critical of the real server.",
				Optional:            true,
				Computed:            true,
			},
			"follow": schema.Int32Attribute{
				MarkdownDescription: "The follow of the real server.",
				Optional:            true,
				Computed:            true,
			},
			"dns_name": schema.StringAttribute{
				MarkdownDescription: "The dns name of the real server.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *RealServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RealServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RealServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "virtual_service_id", data.VirtualServiceId)
	tflog.Debug(ctx, "creating a resource")

	response, err := r.client.AddRealServer(data.VirtualServiceId.String(), data.Address.ValueString(), data.Port.String(), api.RealServerParameters{
		Weight:   int(data.Weight.ValueInt32()),
		Forward:  data.Forward.ValueString(),
		Enable:   bool2ptr(data.Enable.ValueBool()),
		Limit:    int(data.Limit.ValueInt32()),
		Critical: bool2ptr(data.Critical.ValueBool()),
		Follow:   int(data.Follow.ValueInt32()),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create real server, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	real_server_response := response.Rs[len(response.Rs)-1]
	data.Id = types.Int32Value(int32(real_server_response.RsIndex))
	data.VirtualServiceId = types.Int32Value(int32(real_server_response.VSIndex))
	data.Address = types.StringValue(real_server_response.Address)
	data.Port = types.Int32Value(int32(real_server_response.Port))
	data.Weight = types.Int32Value(int32(real_server_response.Weight))
	data.Forward = types.StringValue(real_server_response.Forward)
	data.Enable = types.BoolValue(*real_server_response.Enable)
	data.Limit = types.Int32Value(int32(real_server_response.Limit))
	data.Critical = types.BoolValue(*real_server_response.Critical)
	data.Follow = types.Int32Value(int32(real_server_response.Follow))
	data.DnsName = types.StringValue(real_server_response.DnsName)

	tflog.Trace(ctx, "created a resource real server")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RealServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RealServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.ShowRealServer(data.VirtualServiceId.String(), "!"+data.Id.String())
	if err != nil {
		if serr, ok := err.(*api.LoadMasterError); ok && serr.Message == "Unknown VS" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read real server, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	real_server_response := response.Rs[len(response.Rs)-1]
	data.Id = types.Int32Value(int32(real_server_response.RsIndex))
	data.VirtualServiceId = types.Int32Value(int32(real_server_response.VSIndex))
	data.Address = types.StringValue(real_server_response.Address)
	data.Port = types.Int32Value(int32(real_server_response.Port))
	data.Weight = types.Int32Value(int32(real_server_response.Weight))
	data.Forward = types.StringValue(real_server_response.Forward)
	data.Enable = types.BoolValue(*real_server_response.Enable)
	data.Limit = types.Int32Value(int32(real_server_response.Limit))
	data.Critical = types.BoolValue(*real_server_response.Critical)
	data.Follow = types.Int32Value(int32(real_server_response.Follow))
	data.DnsName = types.StringValue(real_server_response.DnsName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RealServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RealServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	response, err := r.client.ModifyRealServer(data.VirtualServiceId.String(), "!"+data.Id.String(), api.RealServerParameters{
		Weight:   int(data.Weight.ValueInt32()),
		Forward:  data.Forward.ValueString(),
		Enable:   bool2ptr(data.Enable.ValueBool()),
		Limit:    int(data.Limit.ValueInt32()),
		Critical: bool2ptr(data.Critical.ValueBool()),
		Follow:   int(data.Follow.ValueInt32()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update real server, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	real_server_response := response.Rs[len(response.Rs)-1]
	data.Id = types.Int32Value(int32(real_server_response.RsIndex))
	data.VirtualServiceId = types.Int32Value(int32(real_server_response.VSIndex))
	data.Address = types.StringValue(real_server_response.Address)
	data.Port = types.Int32Value(int32(real_server_response.Port))
	data.Weight = types.Int32Value(int32(real_server_response.Weight))
	data.Forward = types.StringValue(real_server_response.Forward)
	data.Enable = types.BoolValue(*real_server_response.Enable)
	data.Limit = types.Int32Value(int32(real_server_response.Limit))
	data.Critical = types.BoolValue(*real_server_response.Critical)
	data.Follow = types.Int32Value(int32(real_server_response.Follow))
	data.DnsName = types.StringValue(real_server_response.DnsName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RealServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RealServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteRealServer(data.VirtualServiceId.String(), "!"+data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete real server, got error: %s", err))
		return
	}
}

func (r *RealServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data RealServerResourceModel

	id_list := strings.Split(req.ID, "/")

	if len(id_list) != 2 {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to parse ID: %s", req.ID))
		return
	}

	response, err := r.client.ShowRealServer(id_list[0], "!"+id_list[1])
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read real server for import, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")
	real_server_response := response.Rs[len(response.Rs)-1]
	data.Id = types.Int32Value(int32(real_server_response.RsIndex))
	data.VirtualServiceId = types.Int32Value(int32(real_server_response.VSIndex))
	data.Address = types.StringValue(real_server_response.Address)
	data.Port = types.Int32Value(int32(real_server_response.Port))
	data.Weight = types.Int32Value(int32(real_server_response.Weight))
	data.Forward = types.StringValue(real_server_response.Forward)
	data.Enable = types.BoolValue(*real_server_response.Enable)
	data.Limit = types.Int32Value(int32(real_server_response.Limit))
	data.Critical = types.BoolValue(*real_server_response.Critical)
	data.Follow = types.Int32Value(int32(real_server_response.Follow))
	data.DnsName = types.StringValue(real_server_response.DnsName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
