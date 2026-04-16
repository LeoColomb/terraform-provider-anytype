// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

// Package client is a thin HTTP client for the Anytype API.
//
// It is intentionally minimal: the Terraform provider uses only the
// endpoints it needs for the resources it manages. The surface of this
// package mirrors the operations declared in the Anytype OpenAPI
// specification (see codegen/openapi.yaml).
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultEndpoint is the local Anytype desktop application API.
	DefaultEndpoint = "http://127.0.0.1:31009"

	// APIVersion is the Anytype-Version header value this provider targets.
	APIVersion = "2025-11-08"

	defaultTimeout = 30 * time.Second
)

// ErrNotFound indicates the requested resource does not exist.
var ErrNotFound = errors.New("anytype: resource not found")

// Client is a typed wrapper around the Anytype HTTP API.
type Client struct {
	endpoint   *url.URL
	apiKey     string
	apiVersion string
	userAgent  string
	http       *http.Client
}

// Config configures the client.
type Config struct {
	Endpoint   string
	APIKey     string
	APIVersion string
	UserAgent  string
	HTTPClient *http.Client
}

// New returns a new Anytype API client.
func New(cfg Config) (*Client, error) {
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint %q: %w", endpoint, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid endpoint %q: missing scheme or host", endpoint)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	apiVersion := cfg.APIVersion
	if apiVersion == "" {
		apiVersion = APIVersion
	}

	return &Client{
		endpoint:   u,
		apiKey:     cfg.APIKey,
		apiVersion: apiVersion,
		userAgent:  cfg.UserAgent,
		http:       httpClient,
	}, nil
}

// APIError represents a non-2xx response from the Anytype API.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("anytype API error: %s: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("anytype API error: %s", e.Status)
}

func (e *APIError) IsNotFound() bool { return e.StatusCode == http.StatusNotFound }

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	u := *c.endpoint
	u.Path = strings.TrimRight(u.Path, "/") + path
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Anytype-Version", c.apiVersion)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(respBody),
		}
		if apiErr.IsNotFound() {
			return ErrNotFound
		}
		return apiErr
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}
	return nil
}
