// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/LeoColomb/terraform-provider-anytype/internal/client"
	"github.com/LeoColomb/terraform-provider-anytype/internal/generated/datasource_schemas"
)

var (
	_ datasource.DataSource              = &typeDataSource{}
	_ datasource.DataSourceWithConfigure = &typeDataSource{}
)

// NewTypeDataSource is the constructor registered with the provider.
func NewTypeDataSource() datasource.DataSource {
	return &typeDataSource{}
}

type typeDataSource struct {
	client *client.Client
}

type typeDataSourceModel struct {
	ID         types.String              `tfsdk:"id"`
	SpaceID    types.String              `tfsdk:"space_id"`
	Key        types.String              `tfsdk:"key"`
	Name       types.String              `tfsdk:"name"`
	PluralName types.String              `tfsdk:"plural_name"`
	Layout     types.String              `tfsdk:"layout"`
	Object     types.String              `tfsdk:"object"`
	Archived   types.Bool                `tfsdk:"archived"`
	Icon       *iconModel                `tfsdk:"icon"`
	Properties []typePropertyDSReference `tfsdk:"properties"`
}

type typePropertyDSReference struct {
	ID     types.String `tfsdk:"id"`
	Key    types.String `tfsdk:"key"`
	Name   types.String `tfsdk:"name"`
	Format types.String `tfsdk:"format"`
}

func (d *typeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_type"
}

// Schema starts from the code-generated data source schema and adds the
// `properties` list that the OpenAPI generator drops because of a Go
// type-name collision (see codegen/generator_config.yml).
func (d *typeDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := datasource_schemas.TypeDataSourceSchema(ctx)
	s.MarkdownDescription = "Look up a single Anytype type by ID in a given space."

	flattenResponseEnvelopeDS(s.Attributes, "type")

	// Re-introduce `icon` (dropped by the OpenAPI generator due to oneOf).
	s.Attributes["icon"] = iconDataSourceAttribute()

	s.Attributes["properties"] = schema.ListNestedAttribute{
		MarkdownDescription: "Properties currently linked to this type.",
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id":     schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the property."},
				"key":    schema.StringAttribute{Computed: true, MarkdownDescription: "The snake_case key of the property."},
				"name":   schema.StringAttribute{Computed: true, MarkdownDescription: "The name of the property."},
				"format": schema.StringAttribute{Computed: true, MarkdownDescription: "The property format."},
			},
		},
	}

	resp.Schema = s
}

func (d *typeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = c
}

func (d *typeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data typeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	t, err := d.client.GetType(ctx, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Anytype type not found",
				fmt.Sprintf("no type with id %q in space %q", data.ID.ValueString(), data.SpaceID.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype type", err.Error())
		return
	}

	data.Key = types.StringValue(t.Key)
	data.Name = types.StringValue(t.Name)
	data.PluralName = types.StringValue(t.PluralName)
	data.Layout = types.StringValue(t.Layout)
	data.Object = types.StringValue(t.Object)
	data.Archived = types.BoolValue(t.Archived)
	data.Icon = iconFromAPI(t.Icon)

	data.Properties = make([]typePropertyDSReference, 0, len(t.Properties))
	for _, p := range t.Properties {
		data.Properties = append(data.Properties, typePropertyDSReference{
			ID:     types.StringValue(p.ID),
			Key:    types.StringValue(p.Key),
			Name:   types.StringValue(p.Name),
			Format: types.StringValue(p.Format),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
