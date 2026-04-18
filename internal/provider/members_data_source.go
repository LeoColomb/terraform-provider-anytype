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
	_ datasource.DataSource              = &membersDataSource{}
	_ datasource.DataSourceWithConfigure = &membersDataSource{}
)

// NewMembersDataSource is the constructor registered with the provider.
func NewMembersDataSource() datasource.DataSource {
	return &membersDataSource{}
}

type membersDataSource struct {
	client *client.Client
}

type membersDataSourceModel struct {
	SpaceID types.String         `tfsdk:"space_id"`
	Members []memberSummaryModel `tfsdk:"members"`
}

type memberSummaryModel struct {
	ID         types.String `tfsdk:"id"`
	Identity   types.String `tfsdk:"identity"`
	GlobalName types.String `tfsdk:"global_name"`
	Name       types.String `tfsdk:"name"`
	Role       types.String `tfsdk:"role"`
	Status     types.String `tfsdk:"status"`
	Object     types.String `tfsdk:"object"`
}

func (d *membersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_members"
}

func (d *membersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all members of an Anytype space.",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the space."},
			"members": schema.ListNestedAttribute{
				MarkdownDescription: "The list of members in the space.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the member."},
						"identity":    schema.StringAttribute{Computed: true, MarkdownDescription: "The network identity of the member."},
						"global_name": schema.StringAttribute{Computed: true, MarkdownDescription: "The global name of the member."},
						"name":        schema.StringAttribute{Computed: true, MarkdownDescription: "The name of the member."},
						"role":        schema.StringAttribute{Computed: true, MarkdownDescription: "The role of the member (`viewer`, `editor`, `owner`, `no_permission`)."},
						"status":      schema.StringAttribute{Computed: true, MarkdownDescription: "The membership status."},
						"object":      schema.StringAttribute{Computed: true, MarkdownDescription: "The data model of the object."},
					},
				},
			},
		},
	}
}

func (d *membersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *membersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data membersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := d.client.ListMembers(ctx, data.SpaceID.ValueString(), client.ListMembersOptions{Limit: 1000})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Anytype members", err.Error())
		return
	}
	data.Members = make([]memberSummaryModel, 0, len(list))
	for _, m := range list {
		data.Members = append(data.Members, memberSummaryModel{
			ID:         types.StringValue(m.ID),
			Identity:   types.StringValue(m.Identity),
			GlobalName: types.StringValue(m.GlobalName),
			Name:       types.StringValue(m.Name),
			Role:       types.StringValue(m.Role),
			Status:     types.StringValue(m.Status),
			Object:     types.StringValue(m.Object),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
