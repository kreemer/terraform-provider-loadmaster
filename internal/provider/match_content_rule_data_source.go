// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var (
	_ datasource.DataSource              = &MatchContentRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &MatchContentRuleDataSource{}
)

func NewMatchContentRuleDataSource() datasource.DataSource {
	return &MatchContentRuleDataSource{}
}

type MatchContentRuleDataSource struct {
	client *api.Client
}

type MatchContentRuleDataSourceModel struct {
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

func (d *MatchContentRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_match_content_rule"
}

func (d *MatchContentRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about a match content rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the rule.",
				Required:            true,
			},
			"pattern": schema.StringAttribute{
				MarkdownDescription: "The pattern to be matched.",
				Computed:            true,
			},
			"match_type": schema.StringAttribute{
				MarkdownDescription: "The type of matching to be performed by the rule. Should be either `regex`, `prefix` or `postfix`.",
				Computed:            true,
			},
			"inc_host": schema.BoolAttribute{
				MarkdownDescription: "Prepend the hostname to request URI before performing the match.",
				Computed:            true,
			},
			"no_case": schema.BoolAttribute{
				MarkdownDescription: "Ignore case when comparing the strings.",
				Computed:            true,
			},
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Invert the sense of the match.",
				Computed:            true,
			},
			"inc_query": schema.BoolAttribute{
				MarkdownDescription: "Append the query string to the URI before performing a match.",
				Computed:            true,
			},
			"header": schema.StringAttribute{
				MarkdownDescription: "The header field name that should be matched. If no header field is set, the default is to match in the URL. Set this to `body` to match on the body of a request.",
				Computed:            true,
			},
			"set_on_match": schema.Int32Attribute{
				MarkdownDescription: "If the rule is successfully matched, set the specified flag.",
				Computed:            true,
			},
			"only_on_flag": schema.Int32Attribute{
				MarkdownDescription: "Only try to execute this rule if the specified flag is set.",
				Computed:            true,
			},
			"only_on_no_flag": schema.Int32Attribute{
				MarkdownDescription: "Only try to execute this rule if the specified flag is not set.",
				Computed:            true,
			},
			"must_fail": schema.BoolAttribute{
				MarkdownDescription: "If this rule is matched, then always fail to connect.",
				Computed:            true,
			},
		},
	}
}

func (d *MatchContentRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *MatchContentRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MatchContentRuleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.ShowRule(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read match content rule, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.MatchContentRules[len(response.MatchContentRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(rule.Header)
	data.Pattern = types.StringValue(rule.Pattern)
	data.MatchType = types.StringValue(rule.MatchType)
	data.IncHost = types.BoolPointerValue(rule.IncHost)
	data.NoCase = types.BoolPointerValue(rule.CaseIndependent)
	data.Negate = types.BoolPointerValue(rule.Negate)
	data.IncQuery = types.BoolPointerValue(rule.IncHost)
	data.SetOnMatch = types.Int32PointerValue(rule.SetOnMatch)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)
	data.MustFail = types.BoolPointerValue(rule.MustFail)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
