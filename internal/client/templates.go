// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Template is an ObjectWithBody served by the templates endpoints. We expose
// the same scalar fields as Object for consistency.
type Template = Object

type templateResponse struct {
	Template Template `json:"template"`
}

type templatesResponse struct {
	Data []Template `json:"data"`
}

// GetTemplate calls GET /v1/spaces/{space_id}/types/{type_id}/templates/{template_id}.
func (c *Client) GetTemplate(ctx context.Context, spaceID, typeID, templateID string) (*Template, error) {
	var out templateResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) +
		"/types/" + url.PathEscape(typeID) +
		"/templates/" + url.PathEscape(templateID)
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Template, nil
}

// ListTemplatesOptions configures pagination for ListTemplates.
type ListTemplatesOptions struct {
	Offset int
	Limit  int
}

// ListTemplates calls GET /v1/spaces/{space_id}/types/{type_id}/templates.
func (c *Client) ListTemplates(ctx context.Context, spaceID, typeID string, opts ListTemplatesOptions) ([]Template, error) {
	q := url.Values{}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	var out templatesResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) +
		"/types/" + url.PathEscape(typeID) + "/templates"
	if err := c.do(ctx, http.MethodGet, path, q, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
