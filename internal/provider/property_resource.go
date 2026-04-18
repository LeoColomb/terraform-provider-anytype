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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/LeoColomb/terraform-provider-anytype/internal/client"
	"github.com/LeoColomb/terraform-provider-anytype/internal/generated/resource_schemas"
)

var (
	_ resource.Resource                = &propertyResource{}
	_ resource.ResourceWithConfigure   = &propertyResource{}
	_ resource.ResourceWithImportState = &propertyResource{}
)

// NewPropertyResource is the constructor registered with the provider.
func NewPropertyResource() resource.Resource {
	return &propertyResource{}
}

type propertyResource struct {
	client *client.Client
}

func (r *propertyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property"
}

// Schema starts from the code-generated schema (validators — notably the
// OneOf enums for `format` and tag `color` — and descriptions come from the
// Anytype OpenAPI) and layers the Terraform-specific adjustments on top:
// `id` is Computed-only, `space_id` and `format` are Required with
// RequiresReplace, the response envelope is flattened, and the generated
// CustomType is stripped from `tags` so the resource model can use an
// ordinary slice.
func (r *propertyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := resource_schemas.PropertyResourceSchema(ctx)
	s.MarkdownDescription = "Manages an [Anytype property](https://anytype.io) inside a space. " +
		"Properties describe typed fields that can be attached to `anytype_type` definitions."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the property.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	s.Attributes["space_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space the property belongs to.",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	s.Attributes["key"] = schema.StringAttribute{
		MarkdownDescription: "The snake_case key of the property. If omitted, Anytype generates one.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	// Keep the generated OneOf validator on `format` but force replacement
	// when it changes (the API does not allow changing a property format).
	if gen, ok := s.Attributes["format"].(schema.StringAttribute); ok {
		gen.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
		s.Attributes["format"] = gen
	}

	// Strip the generated CustomType from `tags` and make it strictly
	// Optional + RequiresReplace (tags are create-time only in this provider).
	if gen, ok := s.Attributes["tags"].(schema.ListNestedAttribute); ok {
		inner := make(map[string]schema.Attribute, len(gen.NestedObject.Attributes))
		for name, child := range gen.NestedObject.Attributes {
			if sa, ok := child.(schema.StringAttribute); ok {
				sa.CustomType = nil
				inner[name] = sa
			} else {
				inner[name] = child
			}
		}
		s.Attributes["tags"] = schema.ListNestedAttribute{
			MarkdownDescription: "Initial tags to seed for `select` / `multi_select` properties. " +
				"Tags are immutable once created; changes force the property to be replaced.",
			Optional: true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.RequiresReplace(),
			},
			NestedObject: schema.NestedAttributeObject{Attributes: inner},
		}
	}

	flattenResponseEnvelope(s.Attributes, "property")

	resp.Schema = s
}

func (r *propertyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// propertyResourceModel is the Terraform state representation of a property.
type propertyResourceModel struct {
	ID      types.String   `tfsdk:"id"`
	SpaceID types.String   `tfsdk:"space_id"`
	Key     types.String   `tfsdk:"key"`
	Name    types.String   `tfsdk:"name"`
	Format  types.String   `tfsdk:"format"`
	Object  types.String   `tfsdk:"object"`
	Tags    []tagSeedModel `tfsdk:"tags"`
}

type tagSeedModel struct {
	Name  types.String `tfsdk:"name"`
	Color types.String `tfsdk:"color"`
	Key   types.String `tfsdk:"key"`
}

func (m *propertyResourceModel) fromAPI(p *client.Property) {
	m.ID = types.StringValue(p.ID)
	m.Key = types.StringValue(p.Key)
	m.Name = types.StringValue(p.Name)
	m.Format = types.StringValue(p.Format)
	m.Object = types.StringValue(p.Object)
}

func (r *propertyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan propertyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags := make([]client.CreateTagRequest, 0, len(plan.Tags))
	for _, t := range plan.Tags {
		tags = append(tags, client.CreateTagRequest{
			Name:  t.Name.ValueString(),
			Color: t.Color.ValueString(),
			Key:   t.Key.ValueString(),
		})
	}

	created, err := r.client.CreateProperty(ctx, plan.SpaceID.ValueString(), client.CreatePropertyRequest{
		Key:    plan.Key.ValueString(),
		Name:   plan.Name.ValueString(),
		Format: plan.Format.ValueString(),
		Tags:   tags,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Anytype property", err.Error())
		return
	}

	plan.fromAPI(created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *propertyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state propertyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := r.client.GetProperty(ctx, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype property", err.Error())
		return
	}

	state.fromAPI(got)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *propertyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state propertyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update := client.UpdatePropertyRequest{}
	if !plan.Name.Equal(state.Name) {
		n := plan.Name.ValueString()
		update.Name = &n
	}
	if !plan.Key.Equal(state.Key) && !plan.Key.IsUnknown() && !plan.Key.IsNull() {
		k := plan.Key.ValueString()
		update.Key = &k
	}

	updated, err := r.client.UpdateProperty(ctx, state.SpaceID.ValueString(), state.ID.ValueString(), update)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update Anytype property", err.Error())
		return
	}

	plan.fromAPI(updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *propertyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state propertyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteProperty(ctx, state.SpaceID.ValueString(), state.ID.ValueString()); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Anytype property", err.Error())
		return
	}
}

// ImportState accepts "<space_id>/<property_id>".
func (r *propertyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	spaceID, propertyID, ok := splitCompositeID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected import ID in the form <space_id>/<property_id>",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), spaceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), propertyID)...)
}
