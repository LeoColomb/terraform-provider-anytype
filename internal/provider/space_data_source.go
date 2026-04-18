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
	Icon        *iconModel   `tfsdk:"icon"`
}

func (d *spaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (d *spaceDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := datasource_schemas.SpaceDataSourceSchema(ctx)
	s.MarkdownDescription = "Look up a single Anytype space by ID."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space.",
		Required:            true,
	}
	flattenResponseEnvelopeDS(s.Attributes, "space")

	// Re-introduce `icon` (dropped by the OpenAPI generator due to oneOf).
	s.Attributes["icon"] = iconDataSourceAttribute()

	resp.Schema = s
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
	data.Icon = iconFromAPI(space.Icon)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
