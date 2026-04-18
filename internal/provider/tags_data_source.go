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
	_ datasource.DataSource              = &tagsDataSource{}
	_ datasource.DataSourceWithConfigure = &tagsDataSource{}
)

// NewTagsDataSource is the constructor registered with the provider.
func NewTagsDataSource() datasource.DataSource {
	return &tagsDataSource{}
}

type tagsDataSource struct {
	client *client.Client
}

type tagsDataSourceModel struct {
	SpaceID    types.String      `tfsdk:"space_id"`
	PropertyID types.String      `tfsdk:"property_id"`
	Tags       []tagSummaryModel `tfsdk:"tags"`
}

type tagSummaryModel struct {
	ID     types.String `tfsdk:"id"`
	Key    types.String `tfsdk:"key"`
	Name   types.String `tfsdk:"name"`
	Color  types.String `tfsdk:"color"`
	Object types.String `tfsdk:"object"`
}

func (d *tagsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

func (d *tagsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// The generated schema for this list endpoint collides with other
	// PaginatedResponse list schemas (see codegen/generator_config.yml); the
	// list data source is hand-written.
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all tags defined on an Anytype property.",
		Attributes: map[string]schema.Attribute{
			"space_id":    schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the space."},
			"property_id": schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the property."},
			"tags": schema.ListNestedAttribute{
				MarkdownDescription: "The list of tags defined on the property.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":     schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the tag."},
						"key":    schema.StringAttribute{Computed: true, MarkdownDescription: "The key of the tag."},
						"name":   schema.StringAttribute{Computed: true, MarkdownDescription: "The name of the tag."},
						"color":  schema.StringAttribute{Computed: true, MarkdownDescription: "The color of the tag."},
						"object": schema.StringAttribute{Computed: true, MarkdownDescription: "The data model of the object."},
					},
				},
			},
		},
	}
}

func (d *tagsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data tagsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := d.client.ListTags(ctx, data.SpaceID.ValueString(), data.PropertyID.ValueString(), client.ListTagsOptions{Limit: 1000})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Anytype tags", err.Error())
		return
	}
	data.Tags = make([]tagSummaryModel, 0, len(list))
	for _, t := range list {
		data.Tags = append(data.Tags, tagSummaryModel{
			ID:     types.StringValue(t.ID),
			Key:    types.StringValue(t.Key),
			Name:   types.StringValue(t.Name),
			Color:  types.StringValue(t.Color),
			Object: types.StringValue(t.Object),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
