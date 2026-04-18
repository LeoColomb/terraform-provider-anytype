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
	_ datasource.DataSource              = &objectsDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsDataSource{}
)

// NewObjectsDataSource is the constructor registered with the provider.
func NewObjectsDataSource() datasource.DataSource {
	return &objectsDataSource{}
}

type objectsDataSource struct {
	client *client.Client
}

type objectsDataSourceModel struct {
	SpaceID types.String         `tfsdk:"space_id"`
	Objects []objectSummaryModel `tfsdk:"objects"`
}

type objectSummaryModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Layout   types.String `tfsdk:"layout"`
	Snippet  types.String `tfsdk:"snippet"`
	Object   types.String `tfsdk:"object"`
	Archived types.Bool   `tfsdk:"archived"`
}

func (d *objectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects"
}

func (d *objectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// PaginatedResponse list schemas collide between resources in the
	// generator; this data source is hand-written (see
	// codegen/generator_config.yml).
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all objects in an Anytype space.",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the space."},
			"objects": schema.ListNestedAttribute{
				MarkdownDescription: "The list of objects in the space.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the object."},
						"name":     schema.StringAttribute{Computed: true, MarkdownDescription: "The name of the object."},
						"layout":   schema.StringAttribute{Computed: true, MarkdownDescription: "The layout of the object."},
						"snippet":  schema.StringAttribute{Computed: true, MarkdownDescription: "The snippet of the object."},
						"object":   schema.StringAttribute{Computed: true, MarkdownDescription: "The data model of the object."},
						"archived": schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the object is archived."},
					},
				},
			},
		},
	}
}

func (d *objectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *objectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data objectsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := d.client.ListObjects(ctx, data.SpaceID.ValueString(), client.ListObjectsOptions{Limit: 1000})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Anytype objects", err.Error())
		return
	}
	data.Objects = make([]objectSummaryModel, 0, len(list))
	for _, o := range list {
		data.Objects = append(data.Objects, objectSummaryModel{
			ID:       types.StringValue(o.ID),
			Name:     types.StringValue(o.Name),
			Layout:   types.StringValue(o.Layout),
			Snippet:  types.StringValue(o.Snippet),
			Object:   types.StringValue(o.Object),
			Archived: types.BoolValue(o.Archived),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
