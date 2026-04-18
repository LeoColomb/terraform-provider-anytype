// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	_ resource.Resource                = &typeResource{}
	_ resource.ResourceWithConfigure   = &typeResource{}
	_ resource.ResourceWithImportState = &typeResource{}
)

// NewTypeResource is the constructor registered with the provider.
func NewTypeResource() resource.Resource {
	return &typeResource{}
}

type typeResource struct {
	client *client.Client
}

func (r *typeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_type"
}

// Schema starts from the code-generated schema (validators and descriptions
// derived from the Anytype OpenAPI) and layers the Terraform-specific
// adjustments on top: `id` is Computed-only, `space_id` is Required with
// RequiresReplace, the response envelope is flattened to top level, and
// `properties` is replaced with a nested block that references existing
// `anytype_property` resources by `id` — the backend-required
// `key` / `name` / `format` triplet is resolved by the provider so users
// never have to re-declare those attributes on the consuming type.
func (r *typeResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := resource_schemas.TypeResourceSchema(ctx)
	s.MarkdownDescription = "Manages an [Anytype type](https://anytype.io) inside a space. " +
		"Types define the shape of objects and link to `anytype_property` resources by `id`."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the type.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	s.Attributes["space_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space the type belongs to.",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	s.Attributes["key"] = schema.StringAttribute{
		MarkdownDescription: "The snake_case key of the type. If omitted, Anytype generates one.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	// Replace the generated `properties` schema with a nested block keyed on
	// property `id`. The `key` / `name` / `format` attributes required by the
	// Anytype API are Computed here: the provider resolves them by calling
	// GetProperty on the referenced IDs, so consumers can simply write
	// `{ id = anytype_property.foo.id }` without repeating those fields.
	s.Attributes["properties"] = schema.ListNestedAttribute{
		MarkdownDescription: "Properties linked to this type. Each entry references an " +
			"existing `anytype_property` by `id`; the provider resolves the backend-required " +
			"`key` / `name` / `format` triplet automatically.",
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "The ID of the `anytype_property` to link.",
					Required:            true,
				},
				"key": schema.StringAttribute{
					MarkdownDescription: "The key of the linked property (resolved from `id`).",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"name": schema.StringAttribute{
					MarkdownDescription: "The name of the linked property (resolved from `id`).",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"format": schema.StringAttribute{
					MarkdownDescription: "The format of the linked property (resolved from `id`).",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
	}

	flattenResponseEnvelope(s.Attributes, "type")

	// Re-introduce `icon` (dropped by the OpenAPI generator due to oneOf).
	s.Attributes["icon"] = iconResourceAttribute()

	resp.Schema = s
}

func (r *typeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// typeResourceModel is the Terraform state representation of a type.
type typeResourceModel struct {
	ID         types.String        `tfsdk:"id"`
	SpaceID    types.String        `tfsdk:"space_id"`
	Key        types.String        `tfsdk:"key"`
	Name       types.String        `tfsdk:"name"`
	PluralName types.String        `tfsdk:"plural_name"`
	Layout     types.String        `tfsdk:"layout"`
	Object     types.String        `tfsdk:"object"`
	Archived   types.Bool          `tfsdk:"archived"`
	Icon       *iconModel          `tfsdk:"icon"`
	Properties []propertyLinkModel `tfsdk:"properties"`
}

// propertyLinkModel is the nested-object state for `anytype_type.properties`.
// `ID` is the user-provided reference to an `anytype_property`; the remaining
// fields are resolved by the provider and exposed as Computed so the user
// does not have to repeat them in configuration.
type propertyLinkModel struct {
	ID     types.String `tfsdk:"id"`
	Key    types.String `tfsdk:"key"`
	Name   types.String `tfsdk:"name"`
	Format types.String `tfsdk:"format"`
}

func (m *typeResourceModel) fromAPI(t *client.Type) {
	m.ID = types.StringValue(t.ID)
	m.Key = types.StringValue(t.Key)
	m.Name = types.StringValue(t.Name)
	m.PluralName = types.StringValue(t.PluralName)
	m.Layout = types.StringValue(t.Layout)
	m.Object = types.StringValue(t.Object)
	m.Archived = types.BoolValue(t.Archived)
	m.Icon = iconFromAPI(t.Icon)

	// Refresh the Computed triplet on each linked property from the API
	// response so drift is surfaced when a linked property is renamed /
	// changes key. The user-authored `id` is preserved when it matches the
	// API entry by position.
	if len(t.Properties) == 0 {
		m.Properties = nil
		return
	}
	out := make([]propertyLinkModel, len(t.Properties))
	for i, p := range t.Properties {
		out[i] = propertyLinkModel{
			ID:     types.StringValue(p.ID),
			Key:    types.StringValue(p.Key),
			Name:   types.StringValue(p.Name),
			Format: types.StringValue(p.Format),
		}
	}
	m.Properties = out
}

// resolveProperties looks up each linked property by ID to obtain the
// `key` / `name` / `format` triplet the API requires and fills in the
// Computed attributes in-place.
func (m *typeResourceModel) resolveProperties(ctx context.Context, c *client.Client) error {
	for i := range m.Properties {
		id := m.Properties[i].ID.ValueString()
		if id == "" {
			return fmt.Errorf("properties[%d].id must be a non-empty property ID", i)
		}
		p, err := c.GetProperty(ctx, m.SpaceID.ValueString(), id)
		if err != nil {
			return fmt.Errorf("properties[%d] (id=%s): %w", i, id, err)
		}
		m.Properties[i].Key = types.StringValue(p.Key)
		m.Properties[i].Name = types.StringValue(p.Name)
		m.Properties[i].Format = types.StringValue(p.Format)
	}
	return nil
}

func (m *typeResourceModel) propertyLinks() []client.PropertyLink {
	if len(m.Properties) == 0 {
		return nil
	}
	out := make([]client.PropertyLink, 0, len(m.Properties))
	for _, p := range m.Properties {
		out = append(out, client.PropertyLink{
			Key:    p.Key.ValueString(),
			Name:   p.Name.ValueString(),
			Format: p.Format.ValueString(),
		})
	}
	return out
}

func (r *typeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan typeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := plan.resolveProperties(ctx, r.client); err != nil {
		resp.Diagnostics.AddError("Unable to resolve linked properties", err.Error())
		return
	}

	created, err := r.client.CreateType(ctx, plan.SpaceID.ValueString(), client.CreateTypeRequest{
		Key:        plan.Key.ValueString(),
		Name:       plan.Name.ValueString(),
		PluralName: plan.PluralName.ValueString(),
		Layout:     plan.Layout.ValueString(),
		Icon:       iconToAPI(plan.Icon),
		Properties: plan.propertyLinks(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Anytype type", err.Error())
		return
	}

	plan.fromAPI(created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *typeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state typeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := r.client.GetType(ctx, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype type", err.Error())
		return
	}

	state.fromAPI(got)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *typeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state typeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertiesChanged := !propertyLinksEqual(plan.Properties, state.Properties)
	if propertiesChanged {
		if err := plan.resolveProperties(ctx, r.client); err != nil {
			resp.Diagnostics.AddError("Unable to resolve linked properties", err.Error())
			return
		}
	}

	update := client.UpdateTypeRequest{}
	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		update.Name = &name
	}
	if !plan.PluralName.Equal(state.PluralName) {
		pn := plan.PluralName.ValueString()
		update.PluralName = &pn
	}
	if !plan.Layout.Equal(state.Layout) {
		l := plan.Layout.ValueString()
		update.Layout = &l
	}
	if !plan.Key.Equal(state.Key) && !plan.Key.IsUnknown() && !plan.Key.IsNull() {
		k := plan.Key.ValueString()
		update.Key = &k
	}
	if propertiesChanged {
		links := plan.propertyLinks()
		update.Properties = &links
	}
	if !iconsEqual(plan.Icon, state.Icon) {
		update.Icon = iconToAPI(plan.Icon)
	}

	updated, err := r.client.UpdateType(ctx, state.SpaceID.ValueString(), state.ID.ValueString(), update)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update Anytype type", err.Error())
		return
	}

	plan.fromAPI(updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *typeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state typeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteType(ctx, state.SpaceID.ValueString(), state.ID.ValueString()); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Anytype type", err.Error())
		return
	}
}

// ImportState accepts "<space_id>/<type_id>" since types are scoped by space.
func (r *typeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	spaceID, typeID, ok := splitCompositeID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected import ID in the form <space_id>/<type_id>",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), spaceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), typeID)...)
}

// propertyLinksEqual compares plan and state entries by the user-authored
// `id`. The Computed triplet is derived from that id, so equality on id
// implies equality on the rest of the link.
func propertyLinksEqual(a, b []propertyLinkModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].ID.Equal(b[i].ID) {
			return false
		}
	}
	return true
}

// splitCompositeID splits a "<space_id>/<child_id>" import identifier.
func splitCompositeID(id string) (string, string, bool) {
	i := strings.LastIndex(id, "/")
	if i <= 0 || i == len(id)-1 {
		return "", "", false
	}
	return id[:i], id[i+1:], true
}
