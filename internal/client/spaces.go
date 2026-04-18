// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Space mirrors the Space schema from the Anytype OpenAPI. The Anytype API
// does not accept an `icon` on CreateSpace/UpdateSpace, so Icon is read-only.
type Space struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Object      string `json:"object,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
	GatewayURL  string `json:"gateway_url,omitempty"`
	Icon        *Icon  `json:"icon,omitempty"`
}

// CreateSpaceRequest maps to CreateSpaceRequest in the OpenAPI.
type CreateSpaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateSpaceRequest maps to UpdateSpaceRequest in the OpenAPI.
type UpdateSpaceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type spaceResponse struct {
	Space Space `json:"space"`
}

type spacesResponse struct {
	Data []Space `json:"data"`
}

// CreateSpace calls POST /v1/spaces.
func (c *Client) CreateSpace(ctx context.Context, req CreateSpaceRequest) (*Space, error) {
	var out spaceResponse
	if err := c.do(ctx, http.MethodPost, "/v1/spaces", nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Space, nil
}

// GetSpace calls GET /v1/spaces/{space_id}.
func (c *Client) GetSpace(ctx context.Context, id string) (*Space, error) {
	var out spaceResponse
	if err := c.do(ctx, http.MethodGet, "/v1/spaces/"+url.PathEscape(id), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Space, nil
}

// UpdateSpace calls PATCH /v1/spaces/{space_id}.
func (c *Client) UpdateSpace(ctx context.Context, id string, req UpdateSpaceRequest) (*Space, error) {
	var out spaceResponse
	if err := c.do(ctx, http.MethodPatch, "/v1/spaces/"+url.PathEscape(id), nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Space, nil
}

// ListSpacesOptions configures pagination for ListSpaces.
type ListSpacesOptions struct {
	Offset int
	Limit  int
}

// ListSpaces calls GET /v1/spaces.
func (c *Client) ListSpaces(ctx context.Context, opts ListSpacesOptions) ([]Space, error) {
	q := url.Values{}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	var out spacesResponse
	if err := c.do(ctx, http.MethodGet, "/v1/spaces", q, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
