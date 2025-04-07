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

var _ resource.Resource = &ModifyUrlRuleResource{}
var _ resource.ResourceWithImportState = &ModifyUrlRuleResource{}

func NewModifyUrlRuleResource() resource.Resource {
	return &ModifyUrlRuleResource{}
}

type ModifyUrlRuleResource struct {
	client *api.Client
}

type ModifyUrlRuleResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Pattern      types.String `tfsdk:"pattern"`
	Replacement  types.String `tfsdk:"replacement"`
	OnlyOnFlag   types.Int32  `tfsdk:"only_on_flag"`
	OnlyOnNoFlag types.Int32  `tfsdk:"only_on_no_flag"`
}

func (r *ModifyUrlRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_modify_url_rule"
}

func (r *ModifyUrlRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a `ModifyUrlRule`.",

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
				Required:            true,
			},
			"replacement": schema.StringAttribute{
				MarkdownDescription: "The replacement string.",
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

func (r *ModifyUrlRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ModifyUrlRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ModifyUrlRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating a resource")

	response, err := r.client.AddRule("4", data.Id.ValueString(), api.GeneralRule{
		Pattern:      data.Pattern.ValueString(),
		Replacement:  data.Replacement.ValueString(),
		Onlyonflag:   int(data.OnlyOnFlag.ValueInt32()),
		Onlyonnoflag: int(data.OnlyOnNoFlag.ValueInt32()),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create modify url rule, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.ModifyURLRules[len(response.ModifyURLRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Pattern = types.StringValue(rule.Pattern)
	data.Replacement = types.StringValue(rule.Replacement)
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))

	tflog.Trace(ctx, "created a resource modify url rule")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ModifyUrlRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ModifyUrlRuleResourceModel

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read modify url rule, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.ModifyURLRules[len(response.ModifyURLRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Pattern = types.StringValue(rule.Pattern)
	data.Replacement = types.StringValue(rule.Replacement)
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ModifyUrlRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ModifyUrlRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	response, err := r.client.ModifyRule(data.Id.ValueString(), api.GeneralRule{
		Pattern:      data.Pattern.ValueString(),
		Replacement:  data.Replacement.ValueString(),
		Onlyonflag:   int(data.OnlyOnFlag.ValueInt32()),
		Onlyonnoflag: int(data.OnlyOnNoFlag.ValueInt32()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update modify url rule, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.ModifyURLRules[len(response.ModifyURLRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Pattern = types.StringValue(rule.Pattern)
	data.Replacement = types.StringValue(rule.Replacement)
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ModifyUrlRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ModifyUrlRuleResourceModel

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

func (r *ModifyUrlRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data ModifyUrlRuleResourceModel

	response, err := r.client.ShowRule(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read modify url rule for import, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.ModifyURLRules[len(response.ModifyURLRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Pattern = types.StringValue(rule.Pattern)
	data.Replacement = types.StringValue(rule.Replacement)
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
