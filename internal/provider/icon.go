// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/LeoColomb/terraform-provider-anytype/internal/client"
)

// iconModel is the decoded representation used when extracting a types.Object
// into Go fields via ObjectAs. The on-model representation is types.Object so
// that Unknown plan values (Computed icons) round-trip through req.Plan.Get
// without the "target type cannot handle unknown values" error.
//
// Anytype expresses icons as a `oneOf` (EmojiIcon | FileIcon | NamedIcon)
// discriminated by `format`. The Terraform Plugin Framework code generator
// does not yet support polymorphic schemas, so we expose a single nested
// object with the union of fields and rely on the discriminator to pick the
// right variant on the wire (see internal/client/icon.go).
type iconModel struct {
	Format types.String `tfsdk:"format"`
	Emoji  types.String `tfsdk:"emoji"`
	File   types.String `tfsdk:"file"`
	Name   types.String `tfsdk:"name"`
	Color  types.String `tfsdk:"color"`
}

// iconAttrTypes is the attribute-type map for the icon nested object. It must
// stay in sync with the schema attributes declared below.
func iconAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"format": types.StringType,
		"emoji":  types.StringType,
		"file":   types.StringType,
		"name":   types.StringType,
		"color":  types.StringType,
	}
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

// iconFromAPI converts a *client.Icon (possibly nil) into a types.Object
// suitable for storage on a resource/data source model. A nil or zero-value
// icon is encoded as ObjectNull so it round-trips cleanly through state.
func iconFromAPI(i *client.Icon) types.Object {
	if i == nil || (*i == client.Icon{}) {
		return types.ObjectNull(iconAttrTypes())
	}
	obj, _ := types.ObjectValue(iconAttrTypes(), map[string]attr.Value{
		"format": types.StringValue(i.Format),
		"emoji":  stringOrNull(i.Emoji),
		"file":   stringOrNull(i.File),
		"name":   stringOrNull(i.Name),
		"color":  stringOrNull(i.Color),
	})
	return obj
}

// iconToAPI converts a types.Object from Terraform plan/state into the wire
// payload expected by the API. Returns nil when the icon is null, unknown, or
// missing its discriminator.
func iconToAPI(ctx context.Context, o types.Object) (*client.Icon, diag.Diagnostics) {
	if o.IsNull() || o.IsUnknown() {
		return nil, nil
	}
	var m iconModel
	diags := o.As(ctx, &m, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	if diags.HasError() {
		return nil, diags
	}
	if m.Format.IsNull() || m.Format.IsUnknown() {
		return nil, diags
	}
	return &client.Icon{
		Format: m.Format.ValueString(),
		Emoji:  m.Emoji.ValueString(),
		File:   m.File.ValueString(),
		Name:   m.Name.ValueString(),
		Color:  m.Color.ValueString(),
	}, diags
}

// iconsEqual reports whether two icon objects carry the same payload. Null
// and Unknown are treated as equivalent "no icon" so refresh-time fills don't
// trigger spurious update calls.
func iconsEqual(a, b types.Object) bool {
	an := a.IsNull() || a.IsUnknown()
	bn := b.IsNull() || b.IsUnknown()
	if an && bn {
		return true
	}
	if an != bn {
		return false
	}
	return a.Equal(b)
}

func stringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
