// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var (
	_ datasource.DataSource              = &DeleteHeaderRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &DeleteHeaderRuleDataSource{}
)

func NewDeleteHeaderRuleDataSource() datasource.DataSource {
	return &DeleteHeaderRuleDataSource{}
}

type DeleteHeaderRuleDataSource struct {
	client *api.Client
}

type DeleteHeaderRuleDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Header       types.String `tfsdk:"header"`
	OnlyOnFlag   types.Int32  `tfsdk:"only_on_flag"`
	OnlyOnNoFlag types.Int32  `tfsdk:"only_on_no_flag"`
}

func (d *DeleteHeaderRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delete_header_rule"
}

func (d *DeleteHeaderRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about a delete header rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the rule.",
				Required:            true,
			},
			"header": schema.StringAttribute{
				MarkdownDescription: "Name of the header field to be removed.",
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
		},
	}
}

func (d *DeleteHeaderRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DeleteHeaderRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeleteHeaderRuleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operation := ClientBackoff(func() (*api.RuleResponse, error) {
		return d.client.ShowRule(data.Id.ValueString())
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read delete header rule, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	rule := response.DeleteHeaderRules[len(response.DeleteHeaderRules)-1]
	data.Id = types.StringValue(rule.Name)
	data.Header = types.StringPointerValue(&rule.Pattern)
	data.OnlyOnFlag = types.Int32PointerValue(rule.OnlyOnFlag)
	data.OnlyOnNoFlag = types.Int32PointerValue(rule.OnlyOnNoFlag)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
