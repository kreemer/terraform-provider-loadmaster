// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kreemer/loadmaster-go-client/api"
)

var _ resource.Resource = &VirtualServiceOwaspRuleResource{}
var _ resource.ResourceWithImportState = &VirtualServiceOwaspRuleResource{}

func NewVirtualServiceOwaspRuleResource() resource.Resource {
	return &VirtualServiceOwaspRuleResource{}
}

type VirtualServiceOwaspRuleResource struct {
	client *api.Client
}

type VirtualServiceOwaspRuleModel struct {
	VirtualServiceId types.String `tfsdk:"virtual_service_id"`
	Rule             types.String `tfsdk:"rule"`
	RunFirst         types.Bool   `tfsdk:"run_first"`
}

func (r *VirtualServiceOwaspRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_service_owasp_rule"
}

func (r *VirtualServiceOwaspRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages connection between a `VirtualService` and an `OwapsCustomRule`.",

		Attributes: map[string]schema.Attribute{
			"virtual_service_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the virtual service.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule": schema.StringAttribute{
				MarkdownDescription: "The name of the OWASP rule to attach to the virtual service.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"run_first": schema.BoolAttribute{
				MarkdownDescription: "Whether the rule should be run first.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *VirtualServiceOwaspRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VirtualServiceOwaspRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtualServiceOwaspRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating a resource")

	operation := func() (*api.LoadMasterResponse, error) {
		return r.client.AddVirtualServiceOwaspCustomRule(data.VirtualServiceId.ValueString(), data.Rule.ValueString(), data.RunFirst.ValueBool())
	}
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create owasp custom rule, got error: %s", err))
		return
	}
	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.RunFirst = types.BoolValue(data.RunFirst.ValueBool())

	tflog.Trace(ctx, "created a resource owasp custom rule")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualServiceOwaspRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtualServiceOwaspRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operation := func() (*api.OwaspRuleResponse, error) {
		return r.client.ShowVirtualServiceOwaspRule(data.VirtualServiceId.ValueString(), data.Rule.ValueString())
	}
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		if serr, ok := err.(*api.LoadMasterError); ok && serr.Message == "Rule not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read owasp custom rule, got error: %s", err))
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.RunFirst = types.BoolValue(response.Rule.RunFirst == "yes")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualServiceOwaspRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VirtualServiceOwaspRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.AddError("Client Error", "Unable to update owasp custom rule, got error")

}

func (r *VirtualServiceOwaspRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtualServiceOwaspRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operation := func() (*api.LoadMasterResponse, error) {
		return r.client.DeleteVirtualServiceOwaspCustomRule(data.VirtualServiceId.ValueString(), data.Rule.ValueString())
	}
	_, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete owasp custom rule, got error: %s", err))
		return
	}
}
func (r *VirtualServiceOwaspRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data VirtualServiceOwaspRuleModel

	id_list := strings.Split(req.ID, "/")

	if len(id_list) != 2 {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to parse ID: %s", req.ID))
		return
	}

	operation := func() (*api.OwaspRuleResponse, error) {
		return r.client.ShowVirtualServiceOwaspRule(id_list[0], id_list[1])
	}
	response, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read owasp custom rule for import, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.SetField(ctx, "response", response)
	tflog.Trace(ctx, "Received valid response from API")

	data.VirtualServiceId = types.StringValue(id_list[0])
	data.Rule = types.StringValue(response.Rule.Name)
	data.RunFirst = types.BoolValue(response.Rule.RunFirst == "yes")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
