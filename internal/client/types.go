// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// PropertyLink mirrors PropertyLink in the Anytype OpenAPI. It is the request
// shape used when linking properties to a type on create/update.
type PropertyLink struct {
	Format string `json:"format"`
	Key    string `json:"key"`
	Name   string `json:"name"`
}

// Type mirrors the Type schema from the Anytype OpenAPI response body.
type Type struct {
	ID         string     `json:"id,omitempty"`
	Key        string     `json:"key,omitempty"`
	Name       string     `json:"name,omitempty"`
	PluralName string     `json:"plural_name,omitempty"`
	Layout     string     `json:"layout,omitempty"`
	Object     string     `json:"object,omitempty"`
	Archived   bool       `json:"archived,omitempty"`
	Icon       *Icon      `json:"icon,omitempty"`
	Properties []Property `json:"properties,omitempty"`
}

// CreateTypeRequest maps to CreateTypeRequest in the OpenAPI.
type CreateTypeRequest struct {
	Key        string         `json:"key,omitempty"`
	Name       string         `json:"name"`
	PluralName string         `json:"plural_name"`
	Layout     string         `json:"layout"`
	Icon       *Icon          `json:"icon,omitempty"`
	Properties []PropertyLink `json:"properties,omitempty"`
}

// UpdateTypeRequest maps to UpdateTypeRequest in the OpenAPI. All fields are
// pointers so the caller can distinguish "don't change" from "set to zero".
type UpdateTypeRequest struct {
	Key        *string         `json:"key,omitempty"`
	Name       *string         `json:"name,omitempty"`
	PluralName *string         `json:"plural_name,omitempty"`
	Layout     *string         `json:"layout,omitempty"`
	Icon       *Icon           `json:"icon,omitempty"`
	Properties *[]PropertyLink `json:"properties,omitempty"`
}

type typeResponse struct {
	Type Type `json:"type"`
}

type typesResponse struct {
	Data []Type `json:"data"`
}

// CreateType calls POST /v1/spaces/{space_id}/types.
func (c *Client) CreateType(ctx context.Context, spaceID string, req CreateTypeRequest) (*Type, error) {
	var out typeResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/types"
	if err := c.do(ctx, http.MethodPost, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Type, nil
}

// GetType calls GET /v1/spaces/{space_id}/types/{type_id}.
func (c *Client) GetType(ctx context.Context, spaceID, typeID string) (*Type, error) {
	var out typeResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/types/" + url.PathEscape(typeID)
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Type, nil
}

// UpdateType calls PATCH /v1/spaces/{space_id}/types/{type_id}.
func (c *Client) UpdateType(ctx context.Context, spaceID, typeID string, req UpdateTypeRequest) (*Type, error) {
	var out typeResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/types/" + url.PathEscape(typeID)
	if err := c.do(ctx, http.MethodPatch, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Type, nil
}

// DeleteType calls DELETE /v1/spaces/{space_id}/types/{type_id}. The API marks
// the type as archived rather than purging it.
func (c *Client) DeleteType(ctx context.Context, spaceID, typeID string) error {
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/types/" + url.PathEscape(typeID)
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListTypesOptions configures pagination for ListTypes.
type ListTypesOptions struct {
	Offset int
	Limit  int
}

// ListTypes calls GET /v1/spaces/{space_id}/types.
func (c *Client) ListTypes(ctx context.Context, spaceID string, opts ListTypesOptions) ([]Type, error) {
	q := url.Values{}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	var out typesResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/types"
	if err := c.do(ctx, http.MethodGet, path, q, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
