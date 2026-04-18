// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateTag(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/spaces/space-1/properties/prop-1/tags" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body CreateTagRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Color != "red" || body.Name != "High" {
			t.Errorf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(tagResponse{Tag: Tag{
			ID:     "tag-1",
			Name:   body.Name,
			Color:  body.Color,
			Object: "tag",
		}})
	}))
	got, err := c.CreateTag(context.Background(), "space-1", "prop-1", CreateTagRequest{Name: "High", Color: "red"})
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}
	if got.ID != "tag-1" {
		t.Errorf("unexpected: %+v", got)
	}
}
