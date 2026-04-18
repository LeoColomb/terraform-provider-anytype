// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreatePropertyWithTags(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/spaces/space-1/properties" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body CreatePropertyRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Format != "select" || len(body.Tags) != 2 {
			t.Errorf("unexpected body: %+v", body)
		}
		if body.Tags[0].Color != "yellow" {
			t.Errorf("tag color = %q", body.Tags[0].Color)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(propertyResponse{Property: Property{
			ID:     "prop-1",
			Key:    "priority",
			Name:   body.Name,
			Format: body.Format,
			Object: "property",
		}})
	}))

	got, err := c.CreateProperty(context.Background(), "space-1", CreatePropertyRequest{
		Name:   "Priority",
		Format: "select",
		Tags: []CreateTagRequest{
			{Name: "Low", Color: "yellow"},
			{Name: "High", Color: "red"},
		},
	})
	if err != nil {
		t.Fatalf("CreateProperty: %v", err)
	}
	if got.ID != "prop-1" || got.Format != "select" {
		t.Errorf("unexpected response: %+v", got)
	}
}

func TestListProperties(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/spaces/space-1/properties" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(propertiesResponse{Data: []Property{{ID: "a"}, {ID: "b"}, {ID: "c"}}})
	}))

	got, err := c.ListProperties(context.Background(), "space-1", ListPropertiesOptions{})
	if err != nil {
		t.Fatalf("ListProperties: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("unexpected: %+v", got)
	}
}
