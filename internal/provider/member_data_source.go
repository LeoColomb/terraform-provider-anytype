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
	_ datasource.DataSource              = &memberDataSource{}
	_ datasource.DataSourceWithConfigure = &memberDataSource{}
)

// NewMemberDataSource is the constructor registered with the provider.
func NewMemberDataSource() datasource.DataSource {
	return &memberDataSource{}
}

type memberDataSource struct {
	client *client.Client
}

type memberDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	SpaceID    types.String `tfsdk:"space_id"`
	Identity   types.String `tfsdk:"identity"`
	GlobalName types.String `tfsdk:"global_name"`
	Name       types.String `tfsdk:"name"`
	Role       types.String `tfsdk:"role"`
	Status     types.String `tfsdk:"status"`
	Object     types.String `tfsdk:"object"`
}

func (d *memberDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_member"
}

func (d *memberDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := datasource_schemas.MemberDataSourceSchema(ctx)
	s.MarkdownDescription = "Look up a single Anytype space member by ID or identity."

	s.Attributes["id"] = schema.StringAttribute{
		MarkdownDescription: "The ID (or identity) of the member.",
		Required:            true,
	}
	s.Attributes["space_id"] = schema.StringAttribute{
		MarkdownDescription: "The ID of the space.",
		Required:            true,
	}
	flattenResponseEnvelopeDS(s.Attributes, "member")

	resp.Schema = s
}

func (d *memberDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *memberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data memberDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	m, err := d.client.GetMember(ctx, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Anytype member not found",
				fmt.Sprintf("no member with id %q in space %q", data.ID.ValueString(), data.SpaceID.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError("Unable to read Anytype member", err.Error())
		return
	}

	data.Identity = types.StringValue(m.Identity)
	data.GlobalName = types.StringValue(m.GlobalName)
	data.Name = types.StringValue(m.Name)
	data.Role = types.StringValue(m.Role)
	data.Status = types.StringValue(m.Status)
	data.Object = types.StringValue(m.Object)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
