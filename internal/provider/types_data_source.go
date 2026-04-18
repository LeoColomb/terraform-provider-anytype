// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/LeoColomb/terraform-provider-anytype/internal/client"
)

var (
	_ datasource.DataSource              = &typesDataSource{}
	_ datasource.DataSourceWithConfigure = &typesDataSource{}
)

// NewTypesDataSource is the constructor registered with the provider.
func NewTypesDataSource() datasource.DataSource {
	return &typesDataSource{}
}

type typesDataSource struct {
	client *client.Client
}

type typesDataSourceModel struct {
	SpaceID types.String       `tfsdk:"space_id"`
	Types   []typeSummaryModel `tfsdk:"types"`
}

type typeSummaryModel struct {
	ID         types.String `tfsdk:"id"`
	Key        types.String `tfsdk:"key"`
	Name       types.String `tfsdk:"name"`
	PluralName types.String `tfsdk:"plural_name"`
	Layout     types.String `tfsdk:"layout"`
	Object     types.String `tfsdk:"object"`
	Archived   types.Bool   `tfsdk:"archived"`
}

func (d *typesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_types"
}

func (d *typesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all Anytype types defined in a given space.",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{MarkdownDescription: "The ID of the space.", Required: true},
			"types": schema.ListNestedAttribute{
				MarkdownDescription: "The list of types in the space.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the type."},
						"key":         schema.StringAttribute{Computed: true, MarkdownDescription: "The key of the type."},
						"name":        schema.StringAttribute{Computed: true, MarkdownDescription: "The singular name."},
						"plural_name": schema.StringAttribute{Computed: true, MarkdownDescription: "The plural name."},
						"layout":      schema.StringAttribute{Computed: true, MarkdownDescription: "The layout."},
						"object":      schema.StringAttribute{Computed: true, MarkdownDescription: "The data model of the object."},
						"archived":    schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the type is archived."},
					},
				},
			},
		},
	}
}

func (d *typesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *typesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data typesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := d.client.ListTypes(ctx, data.SpaceID.ValueString(), client.ListTypesOptions{Limit: 1000})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Anytype types", err.Error())
		return
	}

	data.Types = make([]typeSummaryModel, 0, len(list))
	for _, t := range list {
		data.Types = append(data.Types, typeSummaryModel{
			ID:         types.StringValue(t.ID),
			Key:        types.StringValue(t.Key),
			Name:       types.StringValue(t.Name),
			PluralName: types.StringValue(t.PluralName),
			Layout:     types.StringValue(t.Layout),
			Object:     types.StringValue(t.Object),
			Archived:   types.BoolValue(t.Archived),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
