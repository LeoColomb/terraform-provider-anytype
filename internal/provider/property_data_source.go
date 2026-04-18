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
)

var (
	_ datasource.DataSource              = &propertyDataSource{}
	_ datasource.DataSourceWithConfigure = &propertyDataSource{}
)

// NewPropertyDataSource is the constructor registered with the provider.
func NewPropertyDataSource() datasource.DataSource {
	return &propertyDataSource{}
}

type propertyDataSource struct {
	client *client.Client
}

type propertyDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	SpaceID types.String `tfsdk:"space_id"`
	Key     types.String `tfsdk:"key"`
	Name    types.String `tfsdk:"name"`
	Format  types.String `tfsdk:"format"`
	Object  types.String `tfsdk:"object"`
}

func (d *propertyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property"
}

func (d *propertyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up a single Anytype property by ID in a given space.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{MarkdownDescription: "The ID of the property.", Required: true},
			"space_id": schema.StringAttribute{MarkdownDescription: "The ID of the space.", Required: true},
			"key":      schema.StringAttribute{MarkdownDescription: "The snake_case key of the property.", Computed: true},
			"name":     schema.StringAttribute{MarkdownDescription: "The name of the property.", Computed: true},
			"format":   schema.StringAttribute{MarkdownDescription: "The property format.", Computed: true},
			"object":   schema.StringAttribute{MarkdownDescription: "The data model of the object.", Computed: true},
		},
	}
}

func (d *propertyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *propertyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data propertyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p, err := d.client.GetProperty(ctx, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Anytype property not found",
				fmt.Sprintf("no property with id %q in space %q", data.ID.ValueString(), data.SpaceID.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype property", err.Error())
		return
	}

	data.Key = types.StringValue(p.Key)
	data.Name = types.StringValue(p.Name)
	data.Format = types.StringValue(p.Format)
	data.Object = types.StringValue(p.Object)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
