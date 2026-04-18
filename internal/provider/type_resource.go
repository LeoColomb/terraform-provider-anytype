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
// RequiresReplace, the response envelope is flattened to top level, and the
// generated CustomType is stripped from `properties` so the resource model
// can use an ordinary slice.
func (r *typeResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := resource_schemas.TypeResourceSchema(ctx)
	s.MarkdownDescription = "Manages an [Anytype type](https://anytype.io) inside a space. " +
		"Types define the shape of objects and can link to `anytype_property` resources."

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

	// Strip the generated CustomType from `properties` so the resource model
	// can use a plain []propertyLinkModel. Validators are preserved.
	if gen, ok := s.Attributes["properties"].(schema.ListNestedAttribute); ok {
		inner := make(map[string]schema.Attribute, len(gen.NestedObject.Attributes))
		for name, child := range gen.NestedObject.Attributes {
			if sa, ok := child.(schema.StringAttribute); ok {
				sa.CustomType = nil
				inner[name] = sa
			} else {
				inner[name] = child
			}
		}
		s.Attributes["properties"] = schema.ListNestedAttribute{
			MarkdownDescription: gen.MarkdownDescription,
			Optional:            true,
			Computed:            false,
			NestedObject: schema.NestedAttributeObject{
				Attributes: inner,
			},
		}
	}

	flattenResponseEnvelope(s.Attributes, "type")

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
	Properties []propertyLinkModel `tfsdk:"properties"`
}

type propertyLinkModel struct {
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

	created, err := r.client.CreateType(ctx, plan.SpaceID.ValueString(), client.CreateTypeRequest{
		Key:        plan.Key.ValueString(),
		Name:       plan.Name.ValueString(),
		PluralName: plan.PluralName.ValueString(),
		Layout:     plan.Layout.ValueString(),
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
	if !propertyLinksEqual(plan.Properties, state.Properties) {
		links := plan.propertyLinks()
		update.Properties = &links
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

func propertyLinksEqual(a, b []propertyLinkModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Key.Equal(b[i].Key) ||
			!a[i].Name.Equal(b[i].Name) ||
			!a[i].Format.Equal(b[i].Format) {
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
