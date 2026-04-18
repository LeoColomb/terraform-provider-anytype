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
	_ resource.Resource                = &tagResource{}
	_ resource.ResourceWithConfigure   = &tagResource{}
	_ resource.ResourceWithImportState = &tagResource{}
)

// NewTagResource is the constructor registered with the provider.
func NewTagResource() resource.Resource {
	return &tagResource{}
}

type tagResource struct {
	client *client.Client
}

func (r *tagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *tagResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := resource_schemas.TagResourceSchema(ctx)
	s.MarkdownDescription = "Manages a [tag](https://anytype.io) on a `select` or `multi_select` " +
		"Anytype property. Tags carry the set of allowed values for the property."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the tag.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	s.Attributes["space_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space the tag belongs to.",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	s.Attributes["property_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the property the tag belongs to.",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	s.Attributes["key"] = schema.StringAttribute{
		MarkdownDescription: "Optional custom key for the tag. If omitted, Anytype generates one.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	flattenResponseEnvelope(s.Attributes, "tag")

	resp.Schema = s
}

func (r *tagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type tagResourceModel struct {
	ID         types.String `tfsdk:"id"`
	SpaceID    types.String `tfsdk:"space_id"`
	PropertyID types.String `tfsdk:"property_id"`
	Key        types.String `tfsdk:"key"`
	Name       types.String `tfsdk:"name"`
	Color      types.String `tfsdk:"color"`
	Object     types.String `tfsdk:"object"`
}

func (m *tagResourceModel) fromAPI(t *client.Tag) {
	m.ID = types.StringValue(t.ID)
	m.Key = types.StringValue(t.Key)
	m.Name = types.StringValue(t.Name)
	m.Color = types.StringValue(t.Color)
	m.Object = types.StringValue(t.Object)
}

func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateTag(ctx, plan.SpaceID.ValueString(), plan.PropertyID.ValueString(), client.CreateTagRequest{
		Name:  plan.Name.ValueString(),
		Color: plan.Color.ValueString(),
		Key:   plan.Key.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Anytype tag", err.Error())
		return
	}
	plan.fromAPI(created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	got, err := r.client.GetTag(ctx, state.SpaceID.ValueString(), state.PropertyID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype tag", err.Error())
		return
	}
	state.fromAPI(got)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state tagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update := client.UpdateTagRequest{}
	if !plan.Name.Equal(state.Name) {
		n := plan.Name.ValueString()
		update.Name = &n
	}
	if !plan.Color.Equal(state.Color) {
		c := plan.Color.ValueString()
		update.Color = &c
	}
	if !plan.Key.Equal(state.Key) && !plan.Key.IsUnknown() && !plan.Key.IsNull() {
		k := plan.Key.ValueString()
		update.Key = &k
	}

	updated, err := r.client.UpdateTag(ctx, state.SpaceID.ValueString(), state.PropertyID.ValueString(), state.ID.ValueString(), update)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update Anytype tag", err.Error())
		return
	}
	plan.fromAPI(updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteTag(ctx, state.SpaceID.ValueString(), state.PropertyID.ValueString(), state.ID.ValueString()); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Anytype tag", err.Error())
		return
	}
}

// ImportState accepts "<space_id>/<property_id>/<tag_id>".
func (r *tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected import ID in the form <space_id>/<property_id>/<tag_id>",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("property_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
}
