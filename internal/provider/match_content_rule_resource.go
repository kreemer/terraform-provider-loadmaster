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

var _ resource.Resource = &MatchContentRuleResource{}
var _ resource.ResourceWithImportState = &MatchContentRuleResource{}

func NewMatchContentRuleResource() resource.Resource {
	return &MatchContentRuleResource{}
}

type MatchContentRuleResource struct {
	client *api.Client
}

type MatchContentRuleResourceModel struct {
	Id           types.String `tfsdk:"id"`
	MatchType    types.String `tfsdk:"match_type"`
	IncHost      types.Bool   `tfsdk:"inc_host"`
	NoCase       types.Bool   `tfsdk:"no_case"`
	Negate       types.Bool   `tfsdk:"negate"`
	IncQuery     types.Bool   `tfsdk:"inc_query"`
	Header       types.String `tfsdk:"header"`
	Pattern      types.String `tfsdk:"pattern"`
	SetOnMatch   types.Int32  `tfsdk:"set_on_match"`
	OnlyOnFlag   types.Int32  `tfsdk:"only_on_flag"`
	OnlyOnNoFlag types.Int32  `tfsdk:"only_on_no_flag"`
	MustFail     types.Bool   `tfsdk:"must_fail"`
}

func (r *MatchContentRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_match_content_rule"
}

func (r *MatchContentRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a `MatchContentRule`.",

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
			"match_type": schema.StringAttribute{
				MarkdownDescription: "The type of matching to be performed by the rule. Should be either `regex`, `prefix` or `postfix`.",
				Computed:            true,
				Optional:            true,
			},
			"inc_host": schema.BoolAttribute{
				MarkdownDescription: "Prepend the hostname to request URI before performing the match.",
				Computed:            true,
				Optional:            true,
			},
			"no_case": schema.BoolAttribute{
				MarkdownDescription: "Ignore case when comparing the strings.",
				Computed:            true,
				Optional:            true,
			},
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Invert the sense of the match.",
				Computed:            true,
				Optional:            true,
			},
			"inc_query": schema.BoolAttribute{
				MarkdownDescription: "Append the query string to the URI before performing a match.",
				Computed:            true,
				Optional:            true,
			},
			"header": schema.StringAttribute{
				MarkdownDescription: "The header field name that should be matched. If no header field is set, the default is to match in the URL. Set this to `body` to match on the body of a request.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"set_on_match": schema.Int32Attribute{
				MarkdownDescription: "If the rule is successfully matched, set the specified flag.",
				Optional:            true,
				Computed:            true,
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
			"must_fail": schema.BoolAttribute{
				MarkdownDescription: "If this rule is matched, then always fail to connect.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MatchContentRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MatchContentRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MatchContentRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating a resource")

	response, err := r.client.AddRule("0", data.Id.ValueString(), api.GeneralRule{
		Header:       data.Header.ValueStringPointer(),
		Pattern:      data.Pattern.ValueString(),
		Matchtype:    data.MatchType.ValueString(),
		Inchost:      data.IncHost.ValueBoolPointer(),
		Nocase:       data.NoCase.ValueBoolPointer(),
		Negate:       data.Negate.ValueBoolPointer(),
		Incquery:     data.IncQuery.ValueBoolPointer(),
		Setonmatch:   int(data.SetOnMatch.ValueInt32()),
		Onlyonflag:   int(data.OnlyOnFlag.ValueInt32()),
		Onlyonnoflag: int(data.OnlyOnNoFlag.ValueInt32()),
		Mustfail:     data.MustFail.ValueBoolPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create match content rule, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.MatchContentRules[len(response.MatchContentRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Pattern = types.StringValue(rule.Pattern)
	data.MatchType = types.StringValue(rule.Matchtype)
	data.IncHost = types.BoolPointerValue(rule.Inchost)
	data.NoCase = types.BoolPointerValue(rule.Nocase)
	data.Negate = types.BoolPointerValue(rule.Negate)
	data.IncQuery = types.BoolPointerValue(rule.Inchost)
	data.SetOnMatch = types.Int32Value(int32(rule.Setonmatch))
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))
	data.MustFail = types.BoolPointerValue(rule.Mustfail)

	tflog.Trace(ctx, "created a resource match content rule")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MatchContentRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MatchContentRuleResourceModel

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read match content rule, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.MatchContentRules[len(response.MatchContentRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Pattern = types.StringValue(rule.Pattern)
	data.MatchType = types.StringValue(rule.Matchtype)
	data.IncHost = types.BoolPointerValue(rule.Inchost)
	data.NoCase = types.BoolPointerValue(rule.Nocase)
	data.Negate = types.BoolPointerValue(rule.Negate)
	data.IncQuery = types.BoolPointerValue(rule.Inchost)
	data.SetOnMatch = types.Int32Value(int32(rule.Setonmatch))
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))
	data.MustFail = types.BoolPointerValue(rule.Mustfail)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MatchContentRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MatchContentRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	response, err := r.client.ModifyRule(data.Id.ValueString(), api.GeneralRule{
		Header:       data.Header.ValueStringPointer(),
		Pattern:      data.Pattern.ValueString(),
		Matchtype:    data.MatchType.ValueString(),
		Inchost:      data.IncHost.ValueBoolPointer(),
		Nocase:       data.IncHost.ValueBoolPointer(),
		Negate:       data.IncHost.ValueBoolPointer(),
		Incquery:     data.IncHost.ValueBoolPointer(),
		Setonmatch:   int(data.SetOnMatch.ValueInt32()),
		Onlyonflag:   int(data.OnlyOnFlag.ValueInt32()),
		Onlyonnoflag: int(data.OnlyOnNoFlag.ValueInt32()),
		Mustfail:     data.IncHost.ValueBoolPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update match content rule, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.MatchContentRules[len(response.MatchContentRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Pattern = types.StringValue(rule.Pattern)
	data.MatchType = types.StringValue(rule.Matchtype)
	data.IncHost = types.BoolPointerValue(rule.Inchost)
	data.NoCase = types.BoolPointerValue(rule.Nocase)
	data.Negate = types.BoolPointerValue(rule.Negate)
	data.IncQuery = types.BoolPointerValue(rule.Inchost)
	data.SetOnMatch = types.Int32Value(int32(rule.Setonmatch))
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))
	data.MustFail = types.BoolPointerValue(rule.Mustfail)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MatchContentRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MatchContentRuleResourceModel

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

func (r *MatchContentRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data MatchContentRuleResourceModel

	response, err := r.client.ShowRule(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read match content rule for import, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.MatchContentRules[len(response.MatchContentRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Pattern = types.StringValue(rule.Pattern)
	data.MatchType = types.StringValue(rule.Matchtype)
	data.IncHost = types.BoolPointerValue(rule.Inchost)
	data.NoCase = types.BoolPointerValue(rule.Nocase)
	data.Negate = types.BoolPointerValue(rule.Negate)
	data.IncQuery = types.BoolPointerValue(rule.Inchost)
	data.SetOnMatch = types.Int32Value(int32(rule.Setonmatch))
	data.OnlyOnFlag = types.Int32Value(int32(rule.Onlyonflag))
	data.OnlyOnNoFlag = types.Int32Value(int32(rule.Onlyonnoflag))
	data.MustFail = types.BoolPointerValue(rule.Mustfail)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
