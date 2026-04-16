// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/LeoColomb/terraform-provider-anytype/internal/client"
)

var (
	_ resource.Resource                = &spaceResource{}
	_ resource.ResourceWithConfigure   = &spaceResource{}
	_ resource.ResourceWithImportState = &spaceResource{}
)

// NewSpaceResource is the constructor registered with the provider.
func NewSpaceResource() resource.Resource {
	return &spaceResource{}
}

type spaceResource struct {
	client *client.Client
}

func (r *spaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

// Schema is a hand-tuned variant of the generated SpaceResourceSchema that:
//   - marks `id` as Computed-only (server-assigned)
//   - exposes the read-only network_id / gateway_url / object attributes that
//     the OAS generator could not infer because the Space response schema
//     contains a polymorphic `icon` that is skipped.
func (r *spaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an [Anytype space](https://anytype.io). Anytype does not " +
			"currently support deleting spaces through the API; on `terraform destroy` the space " +
			"is removed from state but remains in your Anytype account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the space.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the space.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the space.",
				Optional:            true,
				Computed:            true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The Anytype network the space belongs to.",
				Computed:            true,
			},
			"gateway_url": schema.StringAttribute{
				MarkdownDescription: "The gateway URL used to serve files and media for this space.",
				Computed:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The data model of the object (`space` or `chat`).",
				Computed:            true,
			},
		},
	}
}

func (r *spaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected provider data",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

// spaceResourceModel is the Terraform state representation of a space.
type spaceResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	NetworkID   types.String `tfsdk:"network_id"`
	GatewayURL  types.String `tfsdk:"gateway_url"`
	Object      types.String `tfsdk:"object"`
}

func (m *spaceResourceModel) fromAPI(s *client.Space) {
	m.ID = types.StringValue(s.ID)
	m.Name = types.StringValue(s.Name)
	m.Description = types.StringValue(s.Description)
	m.NetworkID = types.StringValue(s.NetworkID)
	m.GatewayURL = types.StringValue(s.GatewayURL)
	m.Object = types.StringValue(s.Object)
}

func (r *spaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan spaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, err := r.client.CreateSpace(ctx, client.CreateSpaceRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Anytype space", err.Error())
		return
	}

	plan.fromAPI(space)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *spaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state spaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, err := r.client.GetSpace(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype space", err.Error())
		return
	}

	state.fromAPI(space)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *spaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state spaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update := client.UpdateSpaceRequest{}
	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		update.Name = &name
	}
	if !plan.Description.Equal(state.Description) {
		desc := plan.Description.ValueString()
		update.Description = &desc
	}

	space, err := r.client.UpdateSpace(ctx, state.ID.ValueString(), update)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update Anytype space", err.Error())
		return
	}

	plan.fromAPI(space)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *spaceResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Anytype has no public endpoint for deleting spaces. Removing from state
	// leaves the space intact in Anytype itself, which matches user expectation
	// (destroy does not purge the underlying workspace).
	resp.Diagnostics.AddWarning(
		"Space left in Anytype",
		"The Anytype API does not support deleting spaces. The resource has been removed "+
			"from Terraform state, but the space still exists in your Anytype account.",
	)
}

func (r *spaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
