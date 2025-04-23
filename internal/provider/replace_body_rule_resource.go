// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var _ resource.Resource = &ReplaceBodyRuleResource{}
var _ resource.ResourceWithImportState = &ReplaceBodyRuleResource{}

func NewReplaceBodyRuleResource() resource.Resource {
	return &ReplaceBodyRuleResource{}
}

type ReplaceBodyRuleResource struct {
	client *api.Client
}

type ReplaceBodyRuleResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Pattern      types.String `tfsdk:"pattern"`
	Replacement  types.String `tfsdk:"replacement"`
	NoCase       types.Bool   `tfsdk:"no_case"`
	OnlyOnFlag   types.Int32  `tfsdk:"only_on_flag"`
	OnlyOnNoFlag types.Int32  `tfsdk:"only_on_no_flag"`
}

func (r *ReplaceBodyRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replace_body_rule"
}

func (r *ReplaceBodyRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a `ReplaceBodyRule`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the rule, should be unique for all different rules.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pattern": schema.StringAttribute{
				MarkdownDescription: "The pattern to be matched.",
				Computed:            true,
				Optional:            true,
			},
			"replacement": schema.StringAttribute{
				MarkdownDescription: "The replacement string.",
				Required:            true,
			},
			"no_case": schema.BoolAttribute{
				MarkdownDescription: "Ignore case when comparing the strings.",
				Computed:            true,
				Optional:            true,
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

func (r *ReplaceBodyRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ReplaceBodyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ReplaceBodyRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating a resource")

	response, err := r.client.AddRule("5", data.Id.ValueString(), api.GeneralRule{
		Pattern:      data.Pattern.ValueStringPointer(),
		Replacement:  data.Replacement.ValueStringPointer(),
		NoCase:       data.NoCase.ValueBoolPointer(),
		OnlyOnFlag:   data.OnlyOnFlag.ValueInt32Pointer(),
		OnlyOnNoFlag: data.OnlyOnNoFlag.ValueInt32Pointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create replace body rule, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.ReplaceBodyRules[len(response.ReplaceBodyRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Pattern = types.StringValue(rule.Pattern)
	data.Replacement = types.StringValue(rule.Replacement)
	data.NoCase = types.BoolPointerValue(rule.CaseIndependent)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	tflog.Trace(ctx, "created a resource replace body rule")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplaceBodyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ReplaceBodyRuleResourceModel

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read replace body rule, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.ReplaceBodyRules[len(response.ReplaceBodyRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Pattern = types.StringValue(rule.Pattern)
	data.Replacement = types.StringValue(rule.Replacement)
	data.NoCase = types.BoolPointerValue(rule.CaseIndependent)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplaceBodyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ReplaceBodyRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	response, err := r.client.ModifyRule(data.Id.ValueString(), api.GeneralRule{
		Pattern:      data.Pattern.ValueStringPointer(),
		Replacement:  data.Replacement.ValueStringPointer(),
		NoCase:       data.NoCase.ValueBoolPointer(),
		OnlyOnFlag:   data.OnlyOnFlag.ValueInt32Pointer(),
		OnlyOnNoFlag: data.OnlyOnNoFlag.ValueInt32Pointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update replace body rule, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.ReplaceBodyRules[len(response.ReplaceBodyRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Pattern = types.StringValue(rule.Pattern)
	data.Replacement = types.StringValue(rule.Replacement)
	data.NoCase = types.BoolPointerValue(rule.CaseIndependent)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplaceBodyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ReplaceBodyRuleResourceModel

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

func (r *ReplaceBodyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data ReplaceBodyRuleResourceModel

	response, err := r.client.ShowRule(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read replace body rule for import, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.ReplaceBodyRules[len(response.ReplaceBodyRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Pattern = types.StringValue(rule.Pattern)
	data.Replacement = types.StringValue(rule.Replacement)
	data.NoCase = types.BoolPointerValue(rule.CaseIndependent)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
