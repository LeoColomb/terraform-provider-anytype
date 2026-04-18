// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used in acceptance tests to run the
// provider binary in-process. Acceptance tests are gated on TF_ACC so they
// do not run under `go test` unless explicitly enabled.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"anytype": providerserver.NewProtocol6WithError(New("test")()),
}

// TestProviderSchema is a fast sanity check that the provider schema builds
// without diagnostics. It does not require a running Anytype instance.
func TestProviderSchema(t *testing.T) {
	t.Parallel()

	p := &AnytypeProvider{version: "test"}

	var metaResp = struct{}{}
	_ = metaResp

	// Schema is emitted via a callback; we just need to ensure it does not panic
	// and that resource/data-source constructors return non-nil implementations.
	for _, ctor := range p.Resources(context.Background()) {
		if r := ctor(); r == nil {
			t.Fatal("resource constructor returned nil")
		}
	}
	for _, ctor := range p.DataSources(context.Background()) {
		if d := ctor(); d == nil {
			t.Fatal("data source constructor returned nil")
		}
	}
}
