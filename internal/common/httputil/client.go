package httputil

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

type (
	// Client calls a JSON HTTP API.
	Client struct {
		baseURL string
		headers http.Header
		client  *http.Client
	}
	// Config configures an HTTP Client.
	Config struct {
		BaseURL    string
		Timeout    time.Duration
		Headers    http.Header
		HTTPClient *http.Client
	}
)

// New returns a Client for the given configuration.
func New(cfg Config) *Client {
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: cmp.Or(cfg.Timeout, defaultTimeout)}
	}
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		headers: cfg.Headers.Clone(),
		client:  client,
	}
}

const (
	defaultTimeout = 30 * time.Second
	maxErrorBody   = 4 << 10
)

// Request describes a single call. Zero values are fine for everything but Path.
type Request struct {
	Method  string      // Defaults to [http.MethodGet].
	Path    string      // Path is joined onto the client's base URL.
	Query   url.Values  // Appended as the URL query string.
	Body    any         // Marshalled as JSON when non-nil.
	Headers http.Header // Merged over the client's default headers.
}

// Error is returned when the API responds with a non-2xx status.
type Error struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *Error) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("httputil: unexpected status %s", e.Status)
	}
	return fmt.Sprintf("httputil: unexpected status %s: %s", e.Status, e.Body)
}

// Get performs a GET request and decodes the response into T.
func (c *Client) Get[T any](ctx context.Context, path string, query url.Values) (T, error) {
	return c.Do[T](ctx, Request{Path: path, Query: query})
}

// Post performs a POST request with a JSON body and decodes the response into T.
func (c *Client) Post[T any](ctx context.Context, path string, body any) (T, error) {
	return c.Do[T](ctx, Request{Method: http.MethodPost, Path: path, Body: body})
}

// Do performs a request and decodes the JSON response into T.
func (c *Client) Do[T any](ctx context.Context, req Request) (T, error) {
	var out T

	body, err := encodeBody(req.Body)
	if err != nil {
		return out, err
	}

	uri, err := c.url(req)
	if err != nil {
		return out, err
	}

	method := cmp.Or(req.Method, http.MethodGet)
	httpReq, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return out, fmt.Errorf("httputil: building %s request: %w", method, err)
	}
	c.applyHeaders(httpReq, req)

	res, err := c.client.Do(httpReq)
	if err != nil {
		return out, fmt.Errorf("httputil: %s %s: %w", method, httpReq.URL.Redacted(), err)
	}
	defer res.Body.Close() //nolint:errcheck

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		snippet, _ := io.ReadAll(io.LimitReader(res.Body, maxErrorBody))
		return out, &Error{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Body:       strings.TrimSpace(string(snippet)),
		}
	}

	if res.StatusCode == http.StatusNoContent {
		return out, nil
	}

	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return out, fmt.Errorf("httputil: decoding %s response: %w", method, err)
	}

	return out, nil
}

func (c *Client) url(req Request) (string, error) {
	uri := c.baseURL
	if req.Path != "" {
		joined, err := url.JoinPath(c.baseURL, req.Path)
		if err != nil {
			return "", fmt.Errorf("httputil: joining path %q: %w", req.Path, err)
		}
		uri = joined
	}
	if len(req.Query) > 0 {
		uri += "?" + req.Query.Encode()
	}
	return uri, nil
}

const applicationJSON = "application/json"

// applyHeaders merges the client defaults then the per-request overrides,
// cloning each value slice so the request never aliases shared storage.
func (c *Client) applyHeaders(httpReq *http.Request, req Request) {
	for _, header := range []http.Header{c.headers, req.Headers} {
		for key, values := range header {
			httpReq.Header[key] = slices.Clone(values)
		}
	}

	httpReq.Header.Set("Accept", applicationJSON)
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", applicationJSON)
	}
}

func encodeBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("httputil: encoding request body: %w", err)
	}
	return bytes.NewReader(b), nil
}
