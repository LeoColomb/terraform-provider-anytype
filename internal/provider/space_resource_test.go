// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccSpaceResource is a very small acceptance test. It requires a live
// Anytype instance reachable at ANYTYPE_ENDPOINT (or the local default) and a
// valid ANYTYPE_API_KEY.
func TestAccSpaceResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set; skipping acceptance test")
	}
	if os.Getenv("ANYTYPE_API_KEY") == "" {
		t.Skip("ANYTYPE_API_KEY not set; skipping acceptance test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "anytype_space" "test" {
  name        = "terraform-provider-anytype acctest"
  description = "created by acceptance test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("anytype_space.test", "name", "terraform-provider-anytype acctest"),
					resource.TestCheckResourceAttrSet("anytype_space.test", "id"),
				),
			},
			{
				Config: `
resource "anytype_space" "test" {
  name        = "terraform-provider-anytype acctest (renamed)"
  description = "updated by acceptance test"
}
`,
				Check: resource.TestCheckResourceAttr("anytype_space.test", "name", "terraform-provider-anytype acctest (renamed)"),
			},
		},
	})
}
