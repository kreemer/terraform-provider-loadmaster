// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var _ resource.Resource = &AddHeaderRuleResource{}
var _ resource.ResourceWithImportState = &AddHeaderRuleResource{}

func NewAddHeaderRuleResource() resource.Resource {
	return &AddHeaderRuleResource{}
}

type AddHeaderRuleResource struct {
	client *api.Client
}

type AddHeaderRuleResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Header       types.String `tfsdk:"header"`
	OnlyOnFlag   types.Int32  `tfsdk:"only_on_flag"`
	OnlyOnNoFlag types.Int32  `tfsdk:"only_on_no_flag"`
	Replacement  types.String `tfsdk:"replacement"`
}

func (r *AddHeaderRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_add_header_rule"
}

func (r *AddHeaderRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a `AddHeaderRule`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the rule, should be unique for all different rules.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"header": schema.StringAttribute{
				MarkdownDescription: "Name of the header field to be added.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"replacement": schema.StringAttribute{
				MarkdownDescription: "The replacement string. You can enter a maximum of 255 characters in this parameter.",
				Required:            true,
			},
			"only_on_flag": schema.Int32Attribute{
				MarkdownDescription: "Only try to execute this rule if the specified flag is set.",
				Optional:            true,
				Computed:            true,
			},
			"only_on_no_flag": schema.Int32Attribute{
				MarkdownDescription: "Only try to execute this rule if the specified flag is not set.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *AddHeaderRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AddHeaderRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AddHeaderRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating a resource")

	response, err := r.client.AddRule("1", data.Id.ValueString(), api.GeneralRule{
		Header:       data.Header.ValueStringPointer(),
		Replacement:  data.Replacement.ValueString(),
		Onlyonflag:   int(data.OnlyOnFlag.ValueInt32()),
		Onlyonnoflag: int(data.OnlyOnNoFlag.ValueInt32()),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create add header rule, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.AddHeaderRules[len(response.AddHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Replacement = types.StringValue(rule.HeaderValue)
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))

	tflog.Trace(ctx, "created a resource add header rule")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddHeaderRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AddHeaderRuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.ShowRule(data.Id.ValueString())
	if err != nil {
		if serr, ok := err.(*api.LoadMasterError); ok && serr.Message == "Rule not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read add header rule, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.AddHeaderRules[len(response.AddHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Replacement = types.StringValue(rule.HeaderValue)
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddHeaderRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AddHeaderRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	response, err := r.client.ModifyRule(data.Id.ValueString(), api.GeneralRule{
		Header:       data.Header.ValueStringPointer(),
		Replacement:  data.Replacement.ValueString(),
		Onlyonflag:   int(data.OnlyOnFlag.ValueInt32()),
		Onlyonnoflag: int(data.OnlyOnNoFlag.ValueInt32()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update add header rule, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.AddHeaderRules[len(response.AddHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Replacement = types.StringValue(rule.HeaderValue)
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddHeaderRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AddHeaderRuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteRule(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rule, got error: %s", err))
		return
	}
}

func (r *AddHeaderRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data AddHeaderRuleResourceModel

	response, err := r.client.ShowRule(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read add header rule for import, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.AddHeaderRules[len(response.AddHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Replacement = types.StringValue(rule.HeaderValue)
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
