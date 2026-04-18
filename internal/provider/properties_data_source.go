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
	_ datasource.DataSource              = &propertiesDataSource{}
	_ datasource.DataSourceWithConfigure = &propertiesDataSource{}
)

// NewPropertiesDataSource is the constructor registered with the provider.
func NewPropertiesDataSource() datasource.DataSource {
	return &propertiesDataSource{}
}

type propertiesDataSource struct {
	client *client.Client
}

type propertiesDataSourceModel struct {
	SpaceID    types.String           `tfsdk:"space_id"`
	Properties []propertySummaryModel `tfsdk:"properties"`
}

type propertySummaryModel struct {
	ID     types.String `tfsdk:"id"`
	Key    types.String `tfsdk:"key"`
	Name   types.String `tfsdk:"name"`
	Format types.String `tfsdk:"format"`
	Object types.String `tfsdk:"object"`
}

func (d *propertiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_properties"
}

func (d *propertiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all Anytype properties defined in a given space.",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{MarkdownDescription: "The ID of the space.", Required: true},
			"properties": schema.ListNestedAttribute{
				MarkdownDescription: "The list of properties in the space.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":     schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the property."},
						"key":    schema.StringAttribute{Computed: true, MarkdownDescription: "The snake_case key."},
						"name":   schema.StringAttribute{Computed: true, MarkdownDescription: "The name."},
						"format": schema.StringAttribute{Computed: true, MarkdownDescription: "The property format."},
						"object": schema.StringAttribute{Computed: true, MarkdownDescription: "The data model of the object."},
					},
				},
			},
		},
	}
}

func (d *propertiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *propertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data propertiesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := d.client.ListProperties(ctx, data.SpaceID.ValueString(), client.ListPropertiesOptions{Limit: 1000})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Anytype properties", err.Error())
		return
	}

	data.Properties = make([]propertySummaryModel, 0, len(list))
	for _, p := range list {
		data.Properties = append(data.Properties, propertySummaryModel{
			ID:     types.StringValue(p.ID),
			Key:    types.StringValue(p.Key),
			Name:   types.StringValue(p.Name),
			Format: types.StringValue(p.Format),
			Object: types.StringValue(p.Object),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
