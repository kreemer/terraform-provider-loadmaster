// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var _ resource.Resource = &VirtualServiceResource{}
var _ resource.ResourceWithImportState = &VirtualServiceResource{}

func NewVirtualServiceResource() resource.Resource {
	return &VirtualServiceResource{}
}

type VirtualServiceResource struct {
	client *api.Client
}

type VirtualServiceResourceModel struct {
	Id       types.Int32  `tfsdk:"id"`
	Address  types.String `tfsdk:"address"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Nickname types.String `tfsdk:"nickname"`
	Enabled  types.Bool   `tfsdk:"enabled"`
}

func (r *VirtualServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_service"
}

func (r *VirtualServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual service.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the virtual service. This is also called `Index` in the LoadMaster API.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The address of the virtual service. Should be an IP address of an interface attached to the LoadMaster.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "The port of the virtual service.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol of the virtual service, either `tcp` or `udp`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nickname": schema.StringAttribute{
				MarkdownDescription: "The nickname of the virtual service.",
				Computed:            true,
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "If the virtual service is enabled.",
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (r *VirtualServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VirtualServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtualServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "address", data.Address)
	ctx = tflog.SetField(ctx, "port", data.Port)
	ctx = tflog.SetField(ctx, "protocol", data.Protocol)
	tflog.Debug(ctx, "creating a resource")

	response, err := r.client.AddVirtualService(data.Address.ValueString(), data.Port.ValueString(), data.Protocol.ValueString(), api.VirtualServiceParameters{
		VirtualServiceParametersBasicProperties: &api.VirtualServiceParametersBasicProperties{
			NickName: data.Nickname.ValueString(),
			Enable:   bool2ptr(data.Enabled.ValueBool()),
		},
	})

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read virtual service, got error: %s", err))
		return
	}

	data.Id = types.Int32Value(int32(response.Index))
	data.Address = types.StringValue(response.Address)
	data.Port = types.StringValue(response.Port)
	data.Protocol = types.StringValue(response.Protocol)
	data.Nickname = types.StringValue(response.NickName)
	data.Enabled = types.BoolValue(*response.Enable)

	tflog.Trace(ctx, "created a resource virtual service")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtualServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt32())
	response, err := r.client.ShowVirtualService(id)
	if err != nil {
		if response != nil {
			if response.Code == 422 && response.Message == "Unknown VS" {
				resp.State.RemoveResource(ctx)
				return
			}
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read virtual service, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Id = types.Int32Value(int32(response.Index))
	data.Address = types.StringValue(response.Address)
	data.Port = types.StringValue(response.Port)
	data.Protocol = types.StringValue(response.Protocol)
	data.Nickname = types.StringValue(response.NickName)
	data.Enabled = types.BoolValue(*response.Enable)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VirtualServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	id := int(data.Id.ValueInt32())
	response, err := r.client.ModifyVirtualService(id, api.VirtualServiceParameters{
		VirtualServiceParametersBasicProperties: &api.VirtualServiceParametersBasicProperties{
			NickName: data.Nickname.ValueString(),
			Enable:   bool2ptr(data.Enabled.ValueBool()),
		},
	})
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
	data.Nickname = types.StringValue(response.NickName)
	data.Enabled = types.BoolValue(*response.Enable)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtualServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt32())
	_, err := r.client.DeleteVirtualService(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read virtual service, got error: %s", err))
		return
	}
}

func (r *VirtualServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data VirtualServiceResourceModel

	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to convert ID to integer: %s", err))
		return
	}

	response, err := r.client.ShowVirtualService(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read virtual service, got error: %s", err))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Id = types.Int32Value(int32(response.Index))
	data.Address = types.StringValue(response.Address)
	data.Port = types.StringValue(response.Port)
	data.Protocol = types.StringValue(response.Protocol)
	data.Nickname = types.StringValue(response.NickName)
	data.Enabled = types.BoolValue(*response.Enable)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
