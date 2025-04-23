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

var _ resource.Resource = &DeleteHeaderRuleResource{}
var _ resource.ResourceWithImportState = &DeleteHeaderRuleResource{}

func NewDeleteHeaderRuleResource() resource.Resource {
	return &DeleteHeaderRuleResource{}
}

type DeleteHeaderRuleResource struct {
	client *api.Client
}

type DeleteHeaderRuleResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Header       types.String `tfsdk:"header"`
	OnlyOnFlag   types.Int32  `tfsdk:"only_on_flag"`
	OnlyOnNoFlag types.Int32  `tfsdk:"only_on_no_flag"`
}

func (r *DeleteHeaderRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delete_header_rule"
}

func (r *DeleteHeaderRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a `DeleteHeaderRule`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the rule, should be unique for all different rules.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"header": schema.StringAttribute{
				MarkdownDescription: "Name of the header field to be removed.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
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

func (r *DeleteHeaderRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DeleteHeaderRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeleteHeaderRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating a resource")

	response, err := r.client.AddRule("2", data.Id.ValueString(), api.GeneralRule{
		Pattern:      data.Header.ValueStringPointer(),
		OnlyOnFlag:   data.OnlyOnFlag.ValueInt32Pointer(),
		OnlyOnNoFlag: data.OnlyOnNoFlag.ValueInt32Pointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create delete header rule, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.DeleteHeaderRules[len(response.DeleteHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringValue(rule.Pattern)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	tflog.Trace(ctx, "created a resource delete header rule")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeleteHeaderRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeleteHeaderRuleResourceModel

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read delete header rule, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.DeleteHeaderRules[len(response.DeleteHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringValue(rule.Pattern)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeleteHeaderRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeleteHeaderRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	response, err := r.client.ModifyRule(data.Id.ValueString(), api.GeneralRule{
		Pattern:      data.Header.ValueStringPointer(),
		OnlyOnFlag:   data.OnlyOnFlag.ValueInt32Pointer(),
		OnlyOnNoFlag: data.OnlyOnNoFlag.ValueInt32Pointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update delete header rule, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.DeleteHeaderRules[len(response.DeleteHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringValue(rule.Pattern)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeleteHeaderRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeleteHeaderRuleResourceModel

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

func (r *DeleteHeaderRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data DeleteHeaderRuleResourceModel

	response, err := r.client.ShowRule(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read delete header rule for import, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.DeleteHeaderRules[len(response.DeleteHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringValue(rule.Pattern)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
