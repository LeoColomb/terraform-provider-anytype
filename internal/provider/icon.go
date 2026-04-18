// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/LeoColomb/terraform-provider-anytype/internal/client"
)

// iconModel is the Terraform state representation of an Anytype icon. Anytype
// expresses icons as a `oneOf` (EmojiIcon | FileIcon | NamedIcon) discriminated
// by `format`. The Terraform Plugin Framework code generator does not yet
// support polymorphic schemas, so we expose a single nested object with the
// union of fields and rely on the discriminator to pick the right variant on
// the wire (see internal/client/icon.go).
type iconModel struct {
	Format types.String `tfsdk:"format"`
	Emoji  types.String `tfsdk:"emoji"`
	File   types.String `tfsdk:"file"`
	Name   types.String `tfsdk:"name"`
	Color  types.String `tfsdk:"color"`
}

const iconMarkdown = "Polymorphic Anytype icon. The `format` field selects which " +
	"variant fields are used: `emoji` requires `emoji`; `file` requires `file` " +
	"(a CID); `icon` requires `name` (one of the IconName enum) and an optional `color`."

// iconResourceAttribute returns the schema definition for a writable icon
// nested attribute on a resource.
func iconResourceAttribute() resourceschema.SingleNestedAttribute {
	return resourceschema.SingleNestedAttribute{
		MarkdownDescription: iconMarkdown,
		Optional:            true,
		Computed:            true,
		Attributes: map[string]resourceschema.Attribute{
			"format": resourceschema.StringAttribute{
				MarkdownDescription: "Icon discriminator: one of `emoji`, `file`, `icon`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(client.IconFormatEmoji, client.IconFormatFile, client.IconFormatIcon),
				},
			},
			"emoji": resourceschema.StringAttribute{
				MarkdownDescription: "Emoji character. Only used when `format = \"emoji\"`.",
				Optional:            true,
				Computed:            true,
			},
			"file": resourceschema.StringAttribute{
				MarkdownDescription: "Content-addressed file ID (CID). Only used when `format = \"file\"`.",
				Optional:            true,
				Computed:            true,
			},
			"name": resourceschema.StringAttribute{
				MarkdownDescription: "Named icon, one of the IconName enum. Only used when `format = \"icon\"`.",
				Optional:            true,
				Computed:            true,
			},
			"color": resourceschema.StringAttribute{
				MarkdownDescription: "Color of the named icon. Only used when `format = \"icon\"`.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// iconResourceAttributeReadOnly returns a Computed-only icon nested attribute
// for resources where the API does not accept writes (e.g. anytype_space).
func iconResourceAttributeReadOnly() resourceschema.SingleNestedAttribute {
	return resourceschema.SingleNestedAttribute{
		MarkdownDescription: iconMarkdown + " Read-only on this resource — the Anytype API does not accept icon writes here.",
		Computed:            true,
		Attributes: map[string]resourceschema.Attribute{
			"format": resourceschema.StringAttribute{Computed: true, MarkdownDescription: "Icon discriminator."},
			"emoji":  resourceschema.StringAttribute{Computed: true, MarkdownDescription: "Emoji character."},
			"file":   resourceschema.StringAttribute{Computed: true, MarkdownDescription: "Content-addressed file ID."},
			"name":   resourceschema.StringAttribute{Computed: true, MarkdownDescription: "Named icon."},
			"color":  resourceschema.StringAttribute{Computed: true, MarkdownDescription: "Color of the named icon."},
		},
	}
}

// iconDataSourceAttribute returns the Computed-only icon nested attribute used
// by data sources.
func iconDataSourceAttribute() datasourceschema.SingleNestedAttribute {
	return datasourceschema.SingleNestedAttribute{
		MarkdownDescription: iconMarkdown,
		Computed:            true,
		Attributes: map[string]datasourceschema.Attribute{
			"format": datasourceschema.StringAttribute{Computed: true, MarkdownDescription: "Icon discriminator."},
			"emoji":  datasourceschema.StringAttribute{Computed: true, MarkdownDescription: "Emoji character."},
			"file":   datasourceschema.StringAttribute{Computed: true, MarkdownDescription: "Content-addressed file ID."},
			"name":   datasourceschema.StringAttribute{Computed: true, MarkdownDescription: "Named icon."},
			"color":  datasourceschema.StringAttribute{Computed: true, MarkdownDescription: "Color of the named icon."},
		},
	}
}

// iconFromAPI converts a *client.Icon (possibly nil) into a *iconModel suitable
// for storage on a resource/data source model.
func iconFromAPI(i *client.Icon) *iconModel {
	if i == nil || (*i == client.Icon{}) {
		return nil
	}
	return &iconModel{
		Format: types.StringValue(i.Format),
		Emoji:  stringOrNull(i.Emoji),
		File:   stringOrNull(i.File),
		Name:   stringOrNull(i.Name),
		Color:  stringOrNull(i.Color),
	}
}

// iconToAPI converts a *iconModel from Terraform plan/state into the wire
// payload expected by the API. Returns nil if the icon is unset.
func iconToAPI(m *iconModel) *client.Icon {
	if m == nil {
		return nil
	}
	if m.Format.IsNull() || m.Format.IsUnknown() {
		return nil
	}
	return &client.Icon{
		Format: m.Format.ValueString(),
		Emoji:  m.Emoji.ValueString(),
		File:   m.File.ValueString(),
		Name:   m.Name.ValueString(),
		Color:  m.Color.ValueString(),
	}
}

// iconsEqual reports whether two iconModel pointers carry the same payload.
// nil and a model whose Format is null/unknown are treated as equal so the
// "no icon" cases match.
func iconsEqual(a, b *iconModel) bool {
	an := a == nil || a.Format.IsNull() || a.Format.IsUnknown()
	bn := b == nil || b.Format.IsNull() || b.Format.IsUnknown()
	if an && bn {
		return true
	}
	if an != bn {
		return false
	}
	return a.Format.Equal(b.Format) &&
		a.Emoji.Equal(b.Emoji) &&
		a.File.Equal(b.File) &&
		a.Name.Equal(b.Name) &&
		a.Color.Equal(b.Color)
}

func stringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
