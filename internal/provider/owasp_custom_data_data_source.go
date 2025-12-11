// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var (
	_ datasource.DataSource              = &OwaspCustomDataDataSource{}
	_ datasource.DataSourceWithConfigure = &OwaspCustomDataDataSource{}
)

func NewOwaspCustomDataDataSource() datasource.DataSource {
	return &OwaspCustomDataDataSource{}
}

type OwaspCustomDataDataSource struct {
	client *api.Client
}

type OwaspCustomDataDataSourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Data     types.String `tfsdk:"data"`
}

func (d *OwaspCustomDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_owasp_custom_data"
}

func (d *OwaspCustomDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about a owasp custom data.",
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data.",
				Required:            true,
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "The content of the custom data.",
				Computed:            true,
			},
		},
	}
}

func (d *OwaspCustomDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OwaspCustomDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OwaspCustomDataDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operation := ClientBackoff(func() (*api.LoadMasterDataResponse, error) {
		return d.client.ShowOwaspCustomData(data.Filename.ValueString())
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read owasp custom data, got error: %s", err))
		return
	}

	content := ""
	regex := regexp.MustCompile("^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$")
	if !regex.MatchString(response.Data) {
		content = response.Data
	} else {
		// Decoding shenanigans
		content_bytes, err := base64.StdEncoding.DecodeString(response.Data)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to decode owasp custom data, got error: %s", err))
			return
		}
		content = strings.TrimSuffix(strings.TrimPrefix(string(content_bytes), OwaspCustomDataResource{}.getMarker()), "\r\n")
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.Filename = types.StringValue(data.Filename.ValueString())
	data.Data = types.StringValue(content)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
