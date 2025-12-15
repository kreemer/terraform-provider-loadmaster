// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kreemer/loadmaster-go-client/api"
)

var _ action.Action = &VirtualServiceRestartAction{}
var _ action.ActionWithConfigure = &VirtualServiceRestartAction{}

func NewVirtualServiceRestartAction() action.Action {
	return &VirtualServiceRestartAction{}
}

type VirtualServiceRestartAction struct {
	client *api.Client
}

type VirtualServiceRestartActionModel struct {
	VirtualServiceId types.String `tfsdk:"virtual_service_id"`
}

func (a *VirtualServiceRestartAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_service_restart"
}

func (e *VirtualServiceRestartAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This action will restart a `VirtualService` which is necessary for some configuration changes to take effect.",
		Attributes: map[string]schema.Attribute{
			"virtual_service_id": schema.StringAttribute{
				MarkdownDescription: "The id of the virtual service. This is also called `VIndex` in the LoadMaster API.",
				Required:            true,
			},
		},
	}
}

func (r *VirtualServiceRestartAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

func (e *VirtualServiceRestartAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data VirtualServiceRestartActionModel

	// Read action config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operation := ClientBackoff(func() (*api.VirtualServiceResponse, error) {
		return e.client.ShowVirtualService(data.VirtualServiceId.ValueString())
	})
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error restarting virtual service",
			fmt.Sprintf("Could not restart virtual service with id %s: %s", data.VirtualServiceId.ValueString(), err.Error()),
		)
		return
	}

	if !*response.Enable {
		resp.Diagnostics.AddWarning(
			"Virtual Service is disabled",
			fmt.Sprintf("The virtual service with id %s is currently disabled. Restart action skipped.", data.VirtualServiceId.ValueString()),
		)
		return
	}

	operation = ClientBackoff(func() (*api.VirtualServiceResponse, error) {
		_, err := e.client.ModifyVirtualService(data.VirtualServiceId.ValueString(), api.VirtualServiceParameters{
			VirtualServiceParametersBasicProperties: &api.VirtualServiceParametersBasicProperties{
				Enable: bool2ptr(false),
			},
		})

		if err != nil {
			return nil, err
		}

		_, err = e.client.ModifyVirtualService(data.VirtualServiceId.ValueString(), api.VirtualServiceParameters{
			VirtualServiceParametersBasicProperties: &api.VirtualServiceParametersBasicProperties{
				Enable: bool2ptr(true),
			},
		})

		return nil, err
	})
	_, err = backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error restarting virtual service",
			fmt.Sprintf("Could not restart virtual service with id %s: %s", data.VirtualServiceId.ValueString(), err.Error()),
		)
		return
	}
}
