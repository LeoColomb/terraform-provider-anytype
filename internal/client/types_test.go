// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateType(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/spaces/space-1/types" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body CreateTypeRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Name != "Account" || body.Layout != "basic" {
			t.Errorf("unexpected body: %+v", body)
		}
		if len(body.Properties) != 1 || body.Properties[0].Format != "select" {
			t.Errorf("expected linked property, got %+v", body.Properties)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(typeResponse{Type: Type{
			ID:         "type-1",
			Key:        "account",
			Name:       body.Name,
			PluralName: body.PluralName,
			Layout:     body.Layout,
			Object:     "type",
		}})
	}))

	got, err := c.CreateType(context.Background(), "space-1", CreateTypeRequest{
		Name:       "Account",
		PluralName: "Accounts",
		Layout:     "basic",
		Properties: []PropertyLink{{Key: "status", Name: "Status", Format: "select"}},
	})
	if err != nil {
		t.Fatalf("CreateType: %v", err)
	}
	if got.ID != "type-1" || got.Key != "account" {
		t.Errorf("unexpected response: %+v", got)
	}
}

func TestListTypesQuery(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/spaces/space-1/types" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("limit = %q", got)
		}
		_ = json.NewEncoder(w).Encode(typesResponse{Data: []Type{{ID: "a"}, {ID: "b"}}})
	}))

	got, err := c.ListTypes(context.Background(), "space-1", ListTypesOptions{Limit: 5})
	if err != nil {
		t.Fatalf("ListTypes: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("unexpected: %+v", got)
	}
}
