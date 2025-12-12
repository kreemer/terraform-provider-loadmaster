// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var _ resource.Resource = &OwaspCustomDataResource{}
var _ resource.ResourceWithImportState = &OwaspCustomDataResource{}

func NewOwaspCustomDataResource() resource.Resource {
	return &OwaspCustomDataResource{}
}

type OwaspCustomDataResource struct {
	client *api.Client
}

type OwaspCustomDataResourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Data     types.String `tfsdk:"data"`
}

func (r OwaspCustomDataResource) getMarker() string {
	return "# LoadMaster API MÄrker\n"
}

func (r *OwaspCustomDataResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_owasp_custom_data"
}

func (r *OwaspCustomDataResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a `OwaspCustomData`.\n\nBeware: The LoadMaster API base64 encodes the data and returns this format only if there exists a multibyte character. This resource places a marker line in every resource to ensure consistent behavior.",

		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data, should be unique for all different data.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "The content of the custom data.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *OwaspCustomDataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OwaspCustomDataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OwaspCustomDataResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating a resource")

	// Decoding shenanigans
	content := base64.StdEncoding.EncodeToString([]byte(r.getMarker() + data.Data.ValueString()))

	operation := ClientBackoff(func() (*api.LoadMasterResponse, error) {
		return r.client.AddOwaspCustomData(data.Filename.ValueString(), content)
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create owasp custom data, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Filename = types.StringValue(data.Filename.ValueString())
	data.Data = types.StringValue(data.Data.ValueString())

	tflog.Trace(ctx, "created a resource owasp custom rule")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OwaspCustomDataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OwaspCustomDataResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operation := ClientBackoff(func() (*api.LoadMasterDataResponse, error) {
		return r.client.ShowOwaspCustomData(data.Filename.ValueString())
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		if serr, ok := err.(*api.LoadMasterError); ok && serr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read owasp custom data, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	// Decoding shenanigans
	content_bytes, err := base64.StdEncoding.DecodeString(response.Data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to decode owasp custom data, got error: %s", err))
		return
	}
	content := strings.TrimSuffix(strings.TrimPrefix(string(content_bytes), r.getMarker()), "\r\n")

	data.Filename = types.StringValue(data.Filename.ValueString())
	data.Data = types.StringValue(content)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OwaspCustomDataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OwaspCustomDataResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.AddError("Client Error", "Unable to update owasp custom data, got error")
}

func (r *OwaspCustomDataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OwaspCustomDataResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filename := strings.TrimSuffix(data.Filename.ValueString(), filepath.Ext(data.Filename.ValueString()))

	operation := ClientBackoff(func() (*api.LoadMasterResponse, error) {
		return r.client.DeleteOwaspCustomData(filename)
	})
	_, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete owasp custom data, got error: %s", err))
		return
	}
}

func (r *OwaspCustomDataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data OwaspCustomDataResourceModel

	operation := ClientBackoff(func() (*api.LoadMasterDataResponse, error) {
		return r.client.ShowOwaspCustomData(req.ID)
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read owasp custom data for import, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Filename = types.StringValue(req.ID)
	data.Data = types.StringValue(response.Data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
