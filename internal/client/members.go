// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Member mirrors the Member schema from the Anytype OpenAPI (icon omitted).
type Member struct {
	ID         string `json:"id,omitempty"`
	Identity   string `json:"identity,omitempty"`
	GlobalName string `json:"global_name,omitempty"`
	Name       string `json:"name,omitempty"`
	Role       string `json:"role,omitempty"`
	Status     string `json:"status,omitempty"`
	Object     string `json:"object,omitempty"`
}

type memberResponse struct {
	Member Member `json:"member"`
}

type membersResponse struct {
	Data []Member `json:"data"`
}

// GetMember calls GET /v1/spaces/{space_id}/members/{member_id}.
func (c *Client) GetMember(ctx context.Context, spaceID, memberID string) (*Member, error) {
	var out memberResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/members/" + url.PathEscape(memberID)
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Member, nil
}

// ListMembersOptions configures pagination for ListMembers.
type ListMembersOptions struct {
	Offset int
	Limit  int
}

// ListMembers calls GET /v1/spaces/{space_id}/members.
func (c *Client) ListMembers(ctx context.Context, spaceID string, opts ListMembersOptions) ([]Member, error) {
	q := url.Values{}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	var out membersResponse
	path := "/v1/spaces/" + url.PathEscape(spaceID) + "/members"
	if err := c.do(ctx, http.MethodGet, path, q, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
