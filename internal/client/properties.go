// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Property mirrors the Property schema from the Anytype OpenAPI response body.
type Property struct {
	ID     string `json:"id,omitempty"`
	Key    string `json:"key,omitempty"`
	Name   string `json:"name,omitempty"`
	Format string `json:"format,omitempty"`
	Object string `json:"object,omitempty"`
}

// CreateTagRequest maps to CreateTagRequest in the OpenAPI. It is also used
// when seeding tags on a select/multi_select property at creation time.
type CreateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Key   string `json:"key,omitempty"`
}

// CreatePropertyRequest maps to CreatePropertyRequest in the OpenAPI.
type CreatePropertyRequest struct {
	Name   string             `json:"name"`
	Format string             `json:"format"`
	Key    string             `json:"key,omitempty"`
	Tags   []CreateTagRequest `json:"tags,omitempty"`
}

// UpdatePropertyRequest maps to UpdatePropertyRequest in the OpenAPI.
type UpdatePropertyRequest struct {
	Name *string `json:"name,omitempty"`
	Key  *string `json:"key,omitempty"`
}

type propertyResponse struct {
	Property Property `json:"property"`
}

type propertiesResponse struct {
	Data []Property `json:"data"`
}

// CreateProperty calls POST /v1/spaces/{space_id}/properties.
func (c *Client) CreateProperty(ctx context.Context, spaceID string, req CreatePropertyRequest) (*Property, error) {
	var out propertyResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/properties"
	if err := c.do(ctx, http.MethodPost, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Property, nil
}

// GetProperty calls GET /v1/spaces/{space_id}/properties/{property_id}.
func (c *Client) GetProperty(ctx context.Context, spaceID, propertyID string) (*Property, error) {
	var out propertyResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/properties/" + url.PathEscape(propertyID)
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Property, nil
}

// UpdateProperty calls PATCH /v1/spaces/{space_id}/properties/{property_id}.
func (c *Client) UpdateProperty(ctx context.Context, spaceID, propertyID string, req UpdatePropertyRequest) (*Property, error) {
	var out propertyResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/properties/" + url.PathEscape(propertyID)
	if err := c.do(ctx, http.MethodPatch, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Property, nil
}

// DeleteProperty calls DELETE /v1/spaces/{space_id}/properties/{property_id}.
// The API marks the property as archived.
func (c *Client) DeleteProperty(ctx context.Context, spaceID, propertyID string) error {
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/properties/" + url.PathEscape(propertyID)
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListPropertiesOptions configures pagination for ListProperties.
type ListPropertiesOptions struct {
	Offset int
	Limit  int
}

// ListProperties calls GET /v1/spaces/{space_id}/properties.
func (c *Client) ListProperties(ctx context.Context, spaceID string, opts ListPropertiesOptions) ([]Property, error) {
	q := url.Values{}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	var out propertiesResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/properties"
	if err := c.do(ctx, http.MethodGet, path, q, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
