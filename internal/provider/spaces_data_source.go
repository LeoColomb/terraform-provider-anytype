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
	_ datasource.DataSource              = &spacesDataSource{}
	_ datasource.DataSourceWithConfigure = &spacesDataSource{}
)

// NewSpacesDataSource is the constructor registered with the provider.
func NewSpacesDataSource() datasource.DataSource {
	return &spacesDataSource{}
}

type spacesDataSource struct {
	client *client.Client
}

type spacesDataSourceModel struct {
	Spaces []spaceSummaryModel `tfsdk:"spaces"`
}

type spaceSummaryModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	NetworkID   types.String `tfsdk:"network_id"`
	GatewayURL  types.String `tfsdk:"gateway_url"`
	Object      types.String `tfsdk:"object"`
}

func (d *spacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_spaces"
}

func (d *spacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all Anytype spaces the authenticated user can access.",
		Attributes: map[string]schema.Attribute{
			"spaces": schema.ListNestedAttribute{
				MarkdownDescription: "The list of accessible spaces.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the space."},
						"name":        schema.StringAttribute{Computed: true, MarkdownDescription: "The name of the space."},
						"description": schema.StringAttribute{Computed: true, MarkdownDescription: "The description of the space."},
						"network_id":  schema.StringAttribute{Computed: true, MarkdownDescription: "The Anytype network the space belongs to."},
						"gateway_url": schema.StringAttribute{Computed: true, MarkdownDescription: "Gateway URL used to serve files and media."},
						"object":      schema.StringAttribute{Computed: true, MarkdownDescription: "The data model of the object."},
					},
				},
			},
		},
	}
}

func (d *spacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *spacesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	spaces, err := d.client.ListSpaces(ctx, client.ListSpacesOptions{Limit: 1000})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Anytype spaces", err.Error())
		return
	}

	state := spacesDataSourceModel{
		Spaces: make([]spaceSummaryModel, 0, len(spaces)),
	}
	for _, s := range spaces {
		state.Spaces = append(state.Spaces, spaceSummaryModel{
			ID:          types.StringValue(s.ID),
			Name:        types.StringValue(s.Name),
			Description: types.StringValue(s.Description),
			NetworkID:   types.StringValue(s.NetworkID),
			GatewayURL:  types.StringValue(s.GatewayURL),
			Object:      types.StringValue(s.Object),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
