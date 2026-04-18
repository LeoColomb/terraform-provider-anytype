// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/path"

// frameworkPath is a tiny helper that returns an attribute path for the
// top-level attribute name. Keeping it as a dedicated helper keeps call
// sites readable even when building more complex paths later.
func frameworkPath(name string) path.Path {
	return path.Root(name)
}
