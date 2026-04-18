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
	_ datasource.DataSource              = &templateDataSource{}
	_ datasource.DataSourceWithConfigure = &templateDataSource{}
)

// NewTemplateDataSource is the constructor registered with the provider.
func NewTemplateDataSource() datasource.DataSource {
	return &templateDataSource{}
}

type templateDataSource struct {
	client *client.Client
}

type templateDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	SpaceID  types.String `tfsdk:"space_id"`
	TypeID   types.String `tfsdk:"type_id"`
	Name     types.String `tfsdk:"name"`
	Markdown types.String `tfsdk:"markdown"`
	Snippet  types.String `tfsdk:"snippet"`
	Layout   types.String `tfsdk:"layout"`
	Object   types.String `tfsdk:"object"`
	Archived types.Bool   `tfsdk:"archived"`
}

func (d *templateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (d *templateDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := datasource_schemas.TemplateDataSourceSchema(ctx)
	s.MarkdownDescription = "Look up a single Anytype template (a named ObjectWithBody attached to a type)."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the template.",
		Required:            true,
	}
	s.Attributes["space_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space.",
		Required:            true,
	}
	s.Attributes["type_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the type the template belongs to.",
		Required:            true,
	}
	flattenResponseEnvelopeDS(s.Attributes, "template")

	resp.Schema = s
}

func (d *templateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *templateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data templateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	t, err := d.client.GetTemplate(ctx, data.SpaceID.ValueString(), data.TypeID.ValueString(), data.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Anytype template not found",
				fmt.Sprintf("no template with id %q in space %q, type %q",
					data.ID.ValueString(), data.SpaceID.ValueString(), data.TypeID.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype template", err.Error())
		return
	}

	data.Name = types.StringValue(t.Name)
	data.Markdown = types.StringValue(t.Markdown)
	data.Snippet = types.StringValue(t.Snippet)
	data.Layout = types.StringValue(t.Layout)
	data.Object = types.StringValue(t.Object)
	data.Archived = types.BoolValue(t.Archived)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
