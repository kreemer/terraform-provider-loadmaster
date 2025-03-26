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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var _ resource.Resource = &SubVirtualServiceResource{}
var _ resource.ResourceWithImportState = &SubVirtualServiceResource{}

func NewSubVirtualServiceResource() resource.Resource {
	return &SubVirtualServiceResource{}
}

type SubVirtualServiceResource struct {
	client *api.Client
}

type SubVirtualServiceResourceModel struct {
	Id               types.Int32  `tfsdk:"id"`
	VirtualServiceId types.Int32  `tfsdk:"virtual_service_id"`
	Nickname         types.String `tfsdk:"nickname"`
}

func (r *SubVirtualServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sub_virtual_service"
}

func (r *SubVirtualServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a sub virtual service.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the sub virtual service. This is also called `Index` in the LoadMaster API.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"virtual_service_id": schema.Int32Attribute{
				MarkdownDescription: "The id of the virtual service. This is also called `Index` in the LoadMaster API.",
				Required:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"nickname": schema.StringAttribute{
				MarkdownDescription: "The nickname of the sub virtual service.",
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (r *SubVirtualServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SubVirtualServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SubVirtualServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "virtual_service_id", data.VirtualServiceId)
	tflog.Debug(ctx, "creating a resource")

	response, err := r.client.AddSubVirtualService(int(data.VirtualServiceId.ValueInt32()), api.VirtualServiceParameters{})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create sub virtual service, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Id = types.Int32Value(int32(response.SubVS[len(response.SubVS)-1].VSIndex))
	response, err = r.client.ModifySubVirtualService(int(data.Id.ValueInt32()), api.VirtualServiceParameters{
		VirtualServiceParametersBasicProperties: &api.VirtualServiceParametersBasicProperties{
			NickName: data.Nickname.ValueString(),
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to modify sub virtual service, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Nickname = types.StringValue(response.NickName)

	tflog.Trace(ctx, "created a resource sub virtual service")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubVirtualServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SubVirtualServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt32())
	response, err := r.client.ShowSubVirtualService(id)
	if err != nil {
		if serr, ok := err.(*api.LoadMasterError); ok && serr.Message == "Unknown VS" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read sub virtual service, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Id = types.Int32Value(int32(response.Index))

	data.Nickname = types.StringValue(response.NickName)
	data.VirtualServiceId = types.Int32Value(int32(response.MasterVSID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubVirtualServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SubVirtualServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	id := int(data.Id.ValueInt32())
	response, err := r.client.ModifySubVirtualService(id, api.VirtualServiceParameters{
		VirtualServiceParametersBasicProperties: &api.VirtualServiceParametersBasicProperties{
			NickName: data.Nickname.ValueString(),
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read virtual service, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Id = types.Int32Value(int32(response.Index))
	data.Nickname = types.StringValue(response.NickName)
	data.VirtualServiceId = types.Int32Value(int32(response.MasterVSID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubVirtualServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SubVirtualServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt32())
	_, err := r.client.DeleteSubVirtualService(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read sub virtual service, got error: %s", err))
		return
	}
}

func (r *SubVirtualServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data SubVirtualServiceResourceModel

	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to convert ID to integer: %s", err))
		return
	}

	response, err := r.client.ShowSubVirtualService(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read virtual service, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Id = types.Int32Value(int32(response.Index))
	data.Nickname = types.StringValue(response.NickName)
	data.VirtualServiceId = types.Int32Value(int32(response.MasterVSID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
