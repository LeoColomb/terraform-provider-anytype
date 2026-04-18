// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateObject(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/spaces/space-1/objects" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body CreateObjectRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.TypeKey != "note" {
			t.Errorf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(objectResponse{Object: Object{
			ID:      "obj-1",
			SpaceID: "space-1",
			Name:    body.Name,
			Object:  "object",
		}})
	}))
	got, err := c.CreateObject(context.Background(), "space-1", CreateObjectRequest{
		TypeKey: "note",
		Name:    "Welcome",
		Body:    "# Welcome",
	})
	if err != nil {
		t.Fatalf("CreateObject: %v", err)
	}
	if got.ID != "obj-1" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestListMembers(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/spaces/space-1/members" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(membersResponse{Data: []Member{{ID: "m1", Role: "owner"}}})
	}))
	got, err := c.ListMembers(context.Background(), "space-1", ListMembersOptions{})
	if err != nil {
		t.Fatalf("ListMembers: %v", err)
	}
	if len(got) != 1 || got[0].Role != "owner" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestListTemplates(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/spaces/space-1/types/type-1/templates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(templatesResponse{Data: []Template{{ID: "t1"}}})
	}))
	got, err := c.ListTemplates(context.Background(), "space-1", "type-1", ListTemplatesOptions{})
	if err != nil {
		t.Fatalf("ListTemplates: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("unexpected: %+v", got)
	}
}
