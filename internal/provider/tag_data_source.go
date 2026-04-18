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
	_ datasource.DataSource              = &tagDataSource{}
	_ datasource.DataSourceWithConfigure = &tagDataSource{}
)

// NewTagDataSource is the constructor registered with the provider.
func NewTagDataSource() datasource.DataSource {
	return &tagDataSource{}
}

type tagDataSource struct {
	client *client.Client
}

type tagDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	SpaceID    types.String `tfsdk:"space_id"`
	PropertyID types.String `tfsdk:"property_id"`
	Key        types.String `tfsdk:"key"`
	Name       types.String `tfsdk:"name"`
	Color      types.String `tfsdk:"color"`
	Object     types.String `tfsdk:"object"`
}

func (d *tagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (d *tagDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := datasource_schemas.TagDataSourceSchema(ctx)
	s.MarkdownDescription = "Look up a single Anytype tag by ID in a given property."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the tag.",
		Required:            true,
	}
	s.Attributes["space_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space.",
		Required:            true,
	}
	s.Attributes["property_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the property the tag belongs to.",
		Required:            true,
	}
	flattenResponseEnvelopeDS(s.Attributes, "tag")

	resp.Schema = s
}

func (d *tagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data tagDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	t, err := d.client.GetTag(ctx, data.SpaceID.ValueString(), data.PropertyID.ValueString(), data.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Anytype tag not found",
				fmt.Sprintf("no tag with id %q in property %q", data.ID.ValueString(), data.PropertyID.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype tag", err.Error())
		return
	}

	data.Key = types.StringValue(t.Key)
	data.Name = types.StringValue(t.Name)
	data.Color = types.StringValue(t.Color)
	data.Object = types.StringValue(t.Object)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
