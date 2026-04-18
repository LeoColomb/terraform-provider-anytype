// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

//go:build generate

package tools

import (
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"

	// Provider code generation (OpenAPI -> IR -> Plugin Framework schemas).
	_ "github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework"
	_ "github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi"
)

// Format Terraform code in documentation examples.
//go:generate terraform fmt -recursive ../examples/

// Generate provider documentation from schema + examples.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .. -provider-name anytype
