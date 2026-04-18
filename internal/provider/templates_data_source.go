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
	_ datasource.DataSource              = &templatesDataSource{}
	_ datasource.DataSourceWithConfigure = &templatesDataSource{}
)

// NewTemplatesDataSource is the constructor registered with the provider.
func NewTemplatesDataSource() datasource.DataSource {
	return &templatesDataSource{}
}

type templatesDataSource struct {
	client *client.Client
}

type templatesDataSourceModel struct {
	SpaceID   types.String           `tfsdk:"space_id"`
	TypeID    types.String           `tfsdk:"type_id"`
	Templates []templateSummaryModel `tfsdk:"templates"`
}

type templateSummaryModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Layout   types.String `tfsdk:"layout"`
	Snippet  types.String `tfsdk:"snippet"`
	Object   types.String `tfsdk:"object"`
	Archived types.Bool   `tfsdk:"archived"`
}

func (d *templatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_templates"
}

func (d *templatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all templates defined on a given Anytype type.",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the space."},
			"type_id":  schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the type."},
			"templates": schema.ListNestedAttribute{
				MarkdownDescription: "The list of templates attached to the type.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the template."},
						"name":     schema.StringAttribute{Computed: true, MarkdownDescription: "The name of the template."},
						"layout":   schema.StringAttribute{Computed: true, MarkdownDescription: "The layout of the template."},
						"snippet":  schema.StringAttribute{Computed: true, MarkdownDescription: "The snippet of the template."},
						"object":   schema.StringAttribute{Computed: true, MarkdownDescription: "The data model of the object."},
						"archived": schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the template is archived."},
					},
				},
			},
		},
	}
}

func (d *templatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *templatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data templatesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := d.client.ListTemplates(ctx, data.SpaceID.ValueString(), data.TypeID.ValueString(), client.ListTemplatesOptions{Limit: 1000})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Anytype templates", err.Error())
		return
	}
	data.Templates = make([]templateSummaryModel, 0, len(list))
	for _, t := range list {
		data.Templates = append(data.Templates, templateSummaryModel{
			ID:       types.StringValue(t.ID),
			Name:     types.StringValue(t.Name),
			Layout:   types.StringValue(t.Layout),
			Snippet:  types.StringValue(t.Snippet),
			Object:   types.StringValue(t.Object),
			Archived: types.BoolValue(t.Archived),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
