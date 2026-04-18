// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := New(Config{
		Endpoint: srv.URL,
		APIKey:   "test-key",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestCreateSpace(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/spaces" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization header = %q", got)
		}
		if got := r.Header.Get("Anytype-Version"); got != APIVersion {
			t.Errorf("Anytype-Version header = %q", got)
		}

		var body CreateSpaceRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Name != "My Space" {
			t.Errorf("name = %q", body.Name)
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(spaceResponse{
			Space: Space{
				ID:          "space-id",
				Name:        body.Name,
				Description: body.Description,
				Object:      "space",
			},
		})
	}))

	got, err := c.CreateSpace(context.Background(), CreateSpaceRequest{
		Name:        "My Space",
		Description: "wiki",
	})
	if err != nil {
		t.Fatalf("CreateSpace: %v", err)
	}
	if got.ID != "space-id" || got.Name != "My Space" {
		t.Errorf("unexpected response: %+v", got)
	}
}

func TestGetSpaceNotFound(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))

	_, err := c.GetSpace(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAPIErrorNon2xx(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))

	_, err := c.GetSpace(context.Background(), "id")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %v", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d", apiErr.StatusCode)
	}
	if !strings.Contains(apiErr.Error(), "boom") {
		t.Errorf("error body not included: %v", apiErr)
	}
}

func TestListSpacesQuery(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("limit = %q", got)
		}
		if got := r.URL.Query().Get("offset"); got != "10" {
			t.Errorf("offset = %q", got)
		}
		_ = json.NewEncoder(w).Encode(spacesResponse{Data: []Space{{ID: "a"}, {ID: "b"}}})
	}))

	got, err := c.ListSpaces(context.Background(), ListSpacesOptions{Limit: 50, Offset: 10})
	if err != nil {
		t.Fatalf("ListSpaces: %v", err)
	}
	if len(got) != 2 || got[0].ID != "a" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestNewInvalidEndpoint(t *testing.T) {
	if _, err := New(Config{Endpoint: "::not a url"}); err == nil {
		t.Fatal("expected error for malformed endpoint")
	}
	if _, err := New(Config{Endpoint: "no-scheme-or-host"}); err == nil {
		t.Fatal("expected error for endpoint without scheme")
	}
}
