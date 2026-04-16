// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

//go:build generate

package tools

import (
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)

// Format Terraform code in documentation examples.
//go:generate terraform fmt -recursive ../examples/

// Generate provider documentation from schema + examples.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .. -provider-name anytype
