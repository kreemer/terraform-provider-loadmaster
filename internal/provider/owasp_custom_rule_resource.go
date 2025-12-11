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

var _ resource.Resource = &OwaspCustomRuleResource{}
var _ resource.ResourceWithImportState = &OwaspCustomRuleResource{}

func NewOwaspCustomRuleResource() resource.Resource {
	return &OwaspCustomRuleResource{}
}

type OwaspCustomRuleResource struct {
	client *api.Client
}

type OwaspCustomRuleResourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Data     types.String `tfsdk:"data"`
}

func (r OwaspCustomRuleResource) getMarker() string {
	return "# LoadMaster API MÄrker\n"
}

func (r *OwaspCustomRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_owasp_custom_rule"
}

func (r *OwaspCustomRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a `OwaspCustomRule`.\n\nBeware: The LoadMaster API base64 encodes the data and returns this format only if there exists a multibyte character. This resource places a marker line in every resource to ensure consistent behavior.",

		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				MarkdownDescription: "Identifier of the rule, should be unique for all different rules.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "The content of the custom rule.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *OwaspCustomRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OwaspCustomRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OwaspCustomRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating a resource")

	content := base64.StdEncoding.EncodeToString([]byte(r.getMarker() + data.Data.ValueString()))

	operation := ClientBackoff(func() (*api.LoadMasterResponse, error) {
		return r.client.AddOwaspCustomRule(data.Filename.ValueString(), content)
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create owasp custom rule, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Filename = types.StringValue(data.Filename.ValueString())
	data.Data = types.StringValue(data.Data.ValueString())

	tflog.Trace(ctx, "created a resource owasp custom rule")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OwaspCustomRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OwaspCustomRuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filename := strings.TrimSuffix(data.Filename.ValueString(), filepath.Ext(data.Filename.ValueString()))

	operation := ClientBackoff(func() (*api.LoadMasterDataResponse, error) {
		return r.client.ShowOwaspCustomRule(filename)
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		if serr, ok := err.(*api.LoadMasterError); ok && serr.Message == "Unknown Rule" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read owasp custom rule, got error: %s", err))
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

func (r *OwaspCustomRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OwaspCustomRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.AddError("Client Error", "Unable to update owasp custom rule, got error")

}

func (r *OwaspCustomRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OwaspCustomRuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filename := strings.TrimSuffix(data.Filename.ValueString(), filepath.Ext(data.Filename.ValueString()))

	operation := ClientBackoff(func() (*api.LoadMasterResponse, error) {
		return r.client.DeleteOwaspCustomRule(filename)
	})
	_, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rule, got error: %s", err))
		return
	}
}

func (r *OwaspCustomRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data OwaspCustomRuleResourceModel

	filename := strings.TrimSuffix(req.ID, filepath.Ext(req.ID))

	operation := ClientBackoff(func() (*api.LoadMasterDataResponse, error) {
		return r.client.ShowOwaspCustomRule(filename)
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read owasp custom rule for import, got error: %s", err))
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
