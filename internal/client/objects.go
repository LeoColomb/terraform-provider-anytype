// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Object mirrors the relevant scalar fields of the ObjectWithBody schema.
// The polymorphic `icon` and `properties` arrays are omitted — they are
// left out of the provider surface for the same reasons as the `Type`
// equivalent (see internal/client/types.go and codegen/generator_config.yml).
type Object struct {
	ID       string `json:"id,omitempty"`
	SpaceID  string `json:"space_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Layout   string `json:"layout,omitempty"`
	Markdown string `json:"markdown,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
	Object   string `json:"object,omitempty"`
	Archived bool   `json:"archived,omitempty"`
}

// CreateObjectRequest maps to CreateObjectRequest in the OpenAPI. Only the
// fields that translate cleanly to a Terraform config are exposed.
type CreateObjectRequest struct {
	TypeKey    string `json:"type_key"`
	Name       string `json:"name,omitempty"`
	Body       string `json:"body,omitempty"`
	TemplateID string `json:"template_id,omitempty"`
}

// UpdateObjectRequest maps to UpdateObjectRequest in the OpenAPI.
type UpdateObjectRequest struct {
	Name     *string `json:"name,omitempty"`
	Markdown *string `json:"markdown,omitempty"`
	TypeKey  *string `json:"type_key,omitempty"`
}

type objectResponse struct {
	Object Object `json:"object"`
}

type objectsResponse struct {
	Data []Object `json:"data"`
}

// CreateObject calls POST /v1/spaces/{space_id}/objects.
func (c *Client) CreateObject(ctx context.Context, spaceID string, req CreateObjectRequest) (*Object, error) {
	var out objectResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/objects"
	if err := c.do(ctx, http.MethodPost, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Object, nil
}

// GetObject calls GET /v1/spaces/{space_id}/objects/{object_id}.
func (c *Client) GetObject(ctx context.Context, spaceID, objectID string) (*Object, error) {
	var out objectResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/objects/" + url.PathEscape(objectID)
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Object, nil
}

// UpdateObject calls PATCH /v1/spaces/{space_id}/objects/{object_id}.
func (c *Client) UpdateObject(ctx context.Context, spaceID, objectID string, req UpdateObjectRequest) (*Object, error) {
	var out objectResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/objects/" + url.PathEscape(objectID)
	if err := c.do(ctx, http.MethodPatch, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Object, nil
}

// DeleteObject calls DELETE /v1/spaces/{space_id}/objects/{object_id}. The API
// marks the object as archived.
func (c *Client) DeleteObject(ctx context.Context, spaceID, objectID string) error {
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/objects/" + url.PathEscape(objectID)
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListObjectsOptions configures pagination for ListObjects.
type ListObjectsOptions struct {
	Offset int
	Limit  int
}

// ListObjects calls GET /v1/spaces/{space_id}/objects.
func (c *Client) ListObjects(ctx context.Context, spaceID string, opts ListObjectsOptions) ([]Object, error) {
	q := url.Values{}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	var out objectsResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/objects"
	if err := c.do(ctx, http.MethodGet, path, q, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
