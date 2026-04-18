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
	_ datasource.DataSource              = &objectDataSource{}
	_ datasource.DataSourceWithConfigure = &objectDataSource{}
)

// NewObjectDataSource is the constructor registered with the provider.
func NewObjectDataSource() datasource.DataSource {
	return &objectDataSource{}
}

type objectDataSource struct {
	client *client.Client
}

type objectDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	SpaceID  types.String `tfsdk:"space_id"`
	Name     types.String `tfsdk:"name"`
	Markdown types.String `tfsdk:"markdown"`
	Snippet  types.String `tfsdk:"snippet"`
	Layout   types.String `tfsdk:"layout"`
	Object   types.String `tfsdk:"object"`
	Archived types.Bool   `tfsdk:"archived"`
}

func (d *objectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

func (d *objectDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := datasource_schemas.ObjectDataSourceSchema(ctx)
	s.MarkdownDescription = "Look up a single Anytype object by ID in a given space."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the object.",
		Required:            true,
	}
	s.Attributes["space_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space.",
		Required:            true,
	}
	// `format` is the ?format=md read option — not useful on a data source.
	delete(s.Attributes, "format")

	flattenResponseEnvelopeDS(s.Attributes, "object")

	resp.Schema = s
}

func (d *objectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *objectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data objectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	o, err := d.client.GetObject(ctx, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Anytype object not found",
				fmt.Sprintf("no object with id %q in space %q", data.ID.ValueString(), data.SpaceID.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype object", err.Error())
		return
	}

	data.Name = types.StringValue(o.Name)
	data.Markdown = types.StringValue(o.Markdown)
	data.Snippet = types.StringValue(o.Snippet)
	data.Layout = types.StringValue(o.Layout)
	data.Object = types.StringValue(o.Object)
	data.Archived = types.BoolValue(o.Archived)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
