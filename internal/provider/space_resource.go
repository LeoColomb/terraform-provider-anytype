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
	"github.com/LeoColomb/terraform-provider-anytype/internal/generated/resource_schemas"
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

// Schema starts from the code-generated schema (validators and descriptions
// derived from the Anytype OpenAPI) and layers the Terraform-specific
// adjustments on top: `id` becomes Computed-only, the polymorphic response
// envelope is flattened to top-level attributes, and "Anytype does not yet
// support deleting spaces" is surfaced in the resource description.
func (r *spaceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := resource_schemas.SpaceResourceSchema(ctx)
	s.MarkdownDescription = "Manages an [Anytype space](https://anytype.io). Anytype does not " +
		"currently support deleting spaces through the API; on `terraform destroy` the space " +
		"is removed from state but remains in your Anytype account."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	// The generated schema wraps the response body under a `space` nested
	// attribute. Flatten those read-only fields to top level so the resource
	// stays ergonomic (anytype_space.foo.network_id rather than
	// anytype_space.foo.space.network_id).
	flattenResponseEnvelope(s.Attributes, "space")

	// Re-introduce `icon` (dropped by the OpenAPI generator because it is a
	// `oneOf`). The Anytype API does not accept icon writes on Create/Update
	// Space, so this attribute is read-only on the resource.
	s.Attributes["icon"] = iconResourceAttributeReadOnly()

	resp.Schema = s
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
	Icon        *iconModel   `tfsdk:"icon"`
}

func (m *spaceResourceModel) fromAPI(s *client.Space) {
	m.ID = types.StringValue(s.ID)
	m.Name = types.StringValue(s.Name)
	m.Description = types.StringValue(s.Description)
	m.NetworkID = types.StringValue(s.NetworkID)
	m.GatewayURL = types.StringValue(s.GatewayURL)
	m.Object = types.StringValue(s.Object)
	m.Icon = iconFromAPI(s.Icon)
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
