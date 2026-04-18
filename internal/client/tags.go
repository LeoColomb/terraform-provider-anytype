// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Tag mirrors the Tag schema from the Anytype OpenAPI response body.
type Tag struct {
	ID     string `json:"id,omitempty"`
	Key    string `json:"key,omitempty"`
	Name   string `json:"name,omitempty"`
	Color  string `json:"color,omitempty"`
	Object string `json:"object,omitempty"`
}

// UpdateTagRequest maps to UpdateTagRequest in the OpenAPI.
type UpdateTagRequest struct {
	Name  *string `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
	Key   *string `json:"key,omitempty"`
}

type tagResponse struct {
	Tag Tag `json:"tag"`
}

type tagsResponse struct {
	Data []Tag `json:"data"`
}

// CreateTag calls POST /v1/spaces/{space_id}/properties/{property_id}/tags.
func (c *Client) CreateTag(ctx context.Context, spaceID, propertyID string, req CreateTagRequest) (*Tag, error) {
	var out tagResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) +
		"/properties/" + url.PathEscape(propertyID) + "/tags"
	if err := c.do(ctx, http.MethodPost, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Tag, nil
}

// GetTag calls GET /v1/spaces/{space_id}/properties/{property_id}/tags/{tag_id}.
func (c *Client) GetTag(ctx context.Context, spaceID, propertyID, tagID string) (*Tag, error) {
	var out tagResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) +
		"/properties/" + url.PathEscape(propertyID) +
		"/tags/" + url.PathEscape(tagID)
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Tag, nil
}

// UpdateTag calls PATCH /v1/spaces/{space_id}/properties/{property_id}/tags/{tag_id}.
func (c *Client) UpdateTag(ctx context.Context, spaceID, propertyID, tagID string, req UpdateTagRequest) (*Tag, error) {
	var out tagResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) +
		"/properties/" + url.PathEscape(propertyID) +
		"/tags/" + url.PathEscape(tagID)
	if err := c.do(ctx, http.MethodPatch, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out.Tag, nil
}

// DeleteTag calls DELETE /v1/spaces/{space_id}/properties/{property_id}/tags/{tag_id}.
func (c *Client) DeleteTag(ctx context.Context, spaceID, propertyID, tagID string) error {
	path := "/v1/spaces/" + url.PathEscape(spaceID) +
		"/properties/" + url.PathEscape(propertyID) +
		"/tags/" + url.PathEscape(tagID)
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListTagsOptions configures pagination for ListTags.
type ListTagsOptions struct {
	Offset int
	Limit  int
}

// ListTags calls GET /v1/spaces/{space_id}/properties/{property_id}/tags.
func (c *Client) ListTags(ctx context.Context, spaceID, propertyID string, opts ListTagsOptions) ([]Tag, error) {
	q := url.Values{}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	var out tagsResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) +
		"/properties/" + url.PathEscape(propertyID) + "/tags"
	if err := c.do(ctx, http.MethodGet, path, q, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
