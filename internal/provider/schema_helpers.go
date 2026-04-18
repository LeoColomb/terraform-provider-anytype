// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// flattenResponseEnvelope hoists every Computed child of the SingleNestedAttribute
// named `envelope` (typically "space", "type", "property", ...) into the top-level
// attributes map and removes the wrapper. If an attribute of the same name already
// exists at the top level, the existing attribute is preserved (top-level wins).
//
// The Anytype OpenAPI wraps every mutation response under `{ space: {...} }` /
// `{ type: {...} }` / etc., which the framework code generator faithfully
// translates into a SingleNestedAttribute. For Terraform UX we surface those
// fields as a flat resource — `anytype_space.foo.network_id` reads better than
// `anytype_space.foo.space.network_id`.
func flattenResponseEnvelope(attrs map[string]resourceschema.Attribute, envelope string) {
	wrapper, ok := attrs[envelope].(resourceschema.SingleNestedAttribute)
	if !ok {
		return
	}
	for name, child := range wrapper.Attributes {
		if _, exists := attrs[name]; exists {
			continue
		}
		attrs[name] = child
	}
	delete(attrs, envelope)
}

// flattenResponseEnvelopeDS mirrors flattenResponseEnvelope for data source
// schemas, which live in a separate package.
func flattenResponseEnvelopeDS(attrs map[string]datasourceschema.Attribute, envelope string) {
	wrapper, ok := attrs[envelope].(datasourceschema.SingleNestedAttribute)
	if !ok {
		return
	}
	for name, child := range wrapper.Attributes {
		if _, exists := attrs[name]; exists {
			continue
		}
		attrs[name] = child
	}
	delete(attrs, envelope)
}
