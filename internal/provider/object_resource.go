// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/LeoColomb/terraform-provider-anytype/internal/client"
	"github.com/LeoColomb/terraform-provider-anytype/internal/generated/resource_schemas"
)

var (
	_ resource.Resource                = &objectResource{}
	_ resource.ResourceWithConfigure   = &objectResource{}
	_ resource.ResourceWithImportState = &objectResource{}
)

// NewObjectResource is the constructor registered with the provider.
func NewObjectResource() resource.Resource {
	return &objectResource{}
}

type objectResource struct {
	client *client.Client
}

func (r *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

func (r *objectResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := resource_schemas.ObjectResourceSchema(ctx)
	s.MarkdownDescription = "Manages an [Anytype object](https://anytype.io) inside a space. " +
		"Objects are concrete instances of an `anytype_type`. This resource manages the " +
		"name, markdown body, and type-key mapping; the polymorphic `properties` array " +
		"is not exposed yet."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the object.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	s.Attributes["space_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space the object belongs to.",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	// Exactly one of `type_id` / `type_key` must be set. `type_id` is the
	// preferred way to reference a managed `anytype_type` — the provider
	// resolves its `key` automatically, so consumers never have to reach into
	// `anytype_type.foo.key`. `type_key` remains available for built-in
	// types that are not managed as Terraform resources.
	s.Attributes["type_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the `anytype_type` of the object. " +
			"Mutually exclusive with `type_key`. Changing this forces a new resource.",
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.ExactlyOneOf(
				path.MatchRoot("type_id"),
				path.MatchRoot("type_key"),
			),
		},
	}
	s.Attributes["type_key"] = schema.StringAttribute{
		MarkdownDescription: "The key of the type of the object. " +
			"Mutually exclusive with `type_id`. Changing this forces a new resource.",
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	s.Attributes["template_id"] = schema.StringAttribute{
		MarkdownDescription: "Optional template ID to seed the object body on create. Immutable after create.",
		Optional:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	s.Attributes["name"] = schema.StringAttribute{
		MarkdownDescription: "The name of the object.",
		Optional:            true,
		Computed:            true,
	}
	s.Attributes["body"] = schema.StringAttribute{
		MarkdownDescription: "The markdown body of the object.",
		Optional:            true,
		Computed:            true,
	}
	// `format` is the query-string read format; not user-facing on a resource.
	delete(s.Attributes, "format")

	flattenResponseEnvelope(s.Attributes, "object")

	// Re-introduce `icon` (dropped by the OpenAPI generator due to oneOf).
	s.Attributes["icon"] = iconResourceAttribute()

	// After flattening, `space_id` / `name` from the nested envelope would
	// otherwise shadow our Required overrides — flattenResponseEnvelope
	// skips conflicts, so top-level wins. The read-only `markdown` from the
	// envelope is kept; our top-level `body` carries the create-time payload.

	resp.Schema = s
}

func (r *objectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type objectResourceModel struct {
	ID         types.String `tfsdk:"id"`
	SpaceID    types.String `tfsdk:"space_id"`
	TypeID     types.String `tfsdk:"type_id"`
	TypeKey    types.String `tfsdk:"type_key"`
	TemplateID types.String `tfsdk:"template_id"`
	Name       types.String `tfsdk:"name"`
	Body       types.String `tfsdk:"body"`
	Markdown   types.String `tfsdk:"markdown"`
	Snippet    types.String `tfsdk:"snippet"`
	Layout     types.String `tfsdk:"layout"`
	Object     types.String `tfsdk:"object"`
	Archived   types.Bool   `tfsdk:"archived"`
	Icon       *iconModel   `tfsdk:"icon"`
}

func (m *objectResourceModel) fromAPI(o *client.Object) {
	m.ID = types.StringValue(o.ID)
	m.Name = types.StringValue(o.Name)
	m.Layout = types.StringValue(o.Layout)
	m.Markdown = types.StringValue(o.Markdown)
	m.Snippet = types.StringValue(o.Snippet)
	m.Object = types.StringValue(o.Object)
	m.Archived = types.BoolValue(o.Archived)
	m.Icon = iconFromAPI(o.Icon)
}

func (r *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan objectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// When the user references a type via `type_id`, look up the `type_key`
	// required by the Anytype API so they don't have to wire both attributes
	// themselves.
	if plan.TypeKey.IsNull() || plan.TypeKey.IsUnknown() || plan.TypeKey.ValueString() == "" {
		if plan.TypeID.IsNull() || plan.TypeID.IsUnknown() || plan.TypeID.ValueString() == "" {
			resp.Diagnostics.AddError("Missing type reference",
				"Exactly one of `type_id` or `type_key` must be set on anytype_object.")
			return
		}
		t, err := r.client.GetType(ctx, plan.SpaceID.ValueString(), plan.TypeID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to resolve type_id", err.Error())
			return
		}
		plan.TypeKey = types.StringValue(t.Key)
	}

	created, err := r.client.CreateObject(ctx, plan.SpaceID.ValueString(), client.CreateObjectRequest{
		TypeKey:    plan.TypeKey.ValueString(),
		Name:       plan.Name.ValueString(),
		Body:       plan.Body.ValueString(),
		TemplateID: plan.TemplateID.ValueString(),
		Icon:       iconToAPI(plan.Icon),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Anytype object", err.Error())
		return
	}
	plan.fromAPI(created)
	// type_id / type_key are not echoed back by the Create response; keep
	// whatever the plan supplied so both are populated in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state objectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	got, err := r.client.GetObject(ctx, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype object", err.Error())
		return
	}
	state.fromAPI(got)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update := client.UpdateObjectRequest{}
	if !plan.Name.Equal(state.Name) {
		n := plan.Name.ValueString()
		update.Name = &n
	}
	if !plan.Body.Equal(state.Body) {
		b := plan.Body.ValueString()
		update.Markdown = &b
	}
	if !iconsEqual(plan.Icon, state.Icon) {
		update.Icon = iconToAPI(plan.Icon)
	}

	updated, err := r.client.UpdateObject(ctx, state.SpaceID.ValueString(), state.ID.ValueString(), update)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update Anytype object", err.Error())
		return
	}
	plan.fromAPI(updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state objectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteObject(ctx, state.SpaceID.ValueString(), state.ID.ValueString()); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Anytype object", err.Error())
		return
	}
}

// ImportState accepts "<space_id>/<object_id>".
func (r *objectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	spaceID, objectID, ok := splitCompositeID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected import ID in the form <space_id>/<object_id>",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), spaceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), objectID)...)
}
