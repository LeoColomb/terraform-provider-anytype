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
	_ datasource.DataSource              = &spaceDataSource{}
	_ datasource.DataSourceWithConfigure = &spaceDataSource{}
)

// NewSpaceDataSource is the constructor registered with the provider.
func NewSpaceDataSource() datasource.DataSource {
	return &spaceDataSource{}
}

type spaceDataSource struct {
	client *client.Client
}

type spaceDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	NetworkID   types.String `tfsdk:"network_id"`
	GatewayURL  types.String `tfsdk:"gateway_url"`
	Object      types.String `tfsdk:"object"`
}

func (d *spaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (d *spaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up a single Anytype space by ID.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{MarkdownDescription: "The ID of the space.", Required: true},
			"name":        schema.StringAttribute{MarkdownDescription: "The name of the space.", Computed: true},
			"description": schema.StringAttribute{MarkdownDescription: "The description of the space.", Computed: true},
			"network_id":  schema.StringAttribute{MarkdownDescription: "The Anytype network the space belongs to.", Computed: true},
			"gateway_url": schema.StringAttribute{MarkdownDescription: "Gateway URL used to serve files and media for this space.", Computed: true},
			"object":      schema.StringAttribute{MarkdownDescription: "The data model of the object.", Computed: true},
		},
	}
}

func (d *spaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *spaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data spaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, err := d.client.GetSpace(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError("Anytype space not found", fmt.Sprintf("no space with id %q", data.ID.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype space", err.Error())
		return
	}

	data.ID = types.StringValue(space.ID)
	data.Name = types.StringValue(space.Name)
	data.Description = types.StringValue(space.Description)
	data.NetworkID = types.StringValue(space.NetworkID)
	data.GatewayURL = types.StringValue(space.GatewayURL)
	data.Object = types.StringValue(space.Object)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
