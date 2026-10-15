// Package httpresponder provides a dependency-free client for httpresponder.com.
package httpresponder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (r *Response) Text() string     { return string(r.Body) }
func (r *Response) JSON(v any) error { return json.Unmarshal(r.Body, v) }

type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	body := e.Body
	if len(body) > 500 {
		body = body[:500]
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, body)
}

type Options struct {
	Query          map[string]any
	Method         string
	Headers        map[string]string
	Body           any
	RawBody        []byte
	AllowHTTPError bool
}

func New() *Client {
	return &Client{BaseURL: "https://httpresponder.com", HTTPClient: &http.Client{Timeout: 45 * time.Second}}
}
func (c *Client) URL(path string, query map[string]any) (string, error) {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return "", fmt.Errorf("use an API path beginning with one slash")
	}
	base := c.BaseURL
	if base == "" {
		base = "https://httpresponder.com"
	}
	u, err := url.Parse(strings.TrimRight(base, "/") + path)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range query {
		if v == nil {
			continue
		}
		switch x := v.(type) {
		case string:
			q.Set(k, x)
		default:
			b, e := json.Marshal(v)
			if e != nil {
				return "", e
			}
			q.Set(k, string(b))
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// Request buffers the response body and leaves redirects visible to callers.
func (c *Client) Request(ctx context.Context, path string, o Options) (*Response, error) {
	target, err := c.URL(path, o.Query)
	if err != nil {
		return nil, err
	}
	method := o.Method
	if method == "" {
		method = "GET"
	}
	body := o.RawBody
	if o.Body != nil {
		body, err = json.Marshal(o.Body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range o.Headers {
		req.Header.Set(k, v)
	}
	if o.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 45 * time.Second}
	}
	copyClient := *client
	copyClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	res, err := copyClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	out := &Response{StatusCode: res.StatusCode, Header: res.Header, Body: data}
	if !o.AllowHTTPError && (res.StatusCode < 200 || res.StatusCode >= 300) {
		return out, &APIError{StatusCode: res.StatusCode, Body: data}
	}
	return out, nil
}
func copyMap(in map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
func (c *Client) Endpoint(ctx context.Context, path string, o Options) (*Response, error) {
	o.AllowHTTPError = true
	return c.Request(ctx, path, o)
}
func (c *Client) Decode(ctx context.Context, data, contentType string) (*Response, error) {
	q := map[string]any{"data": data}
	if contentType != "" {
		q["content_type"] = contentType
	}
	return c.Request(ctx, "/base64", Options{Query: q})
}
func (c *Client) Encode(ctx context.Context, data string) (*Response, error) {
	return c.Request(ctx, "/base64", Options{Query: map[string]any{"data": data, "mode": "encode"}})
}
func (c *Client) Respond(ctx context.Context, params map[string]any) (*Response, error) {
	return c.Endpoint(ctx, "/respond", Options{Query: params})
}
func (c *Client) Status(ctx context.Context, code int) (*Response, error) {
	return c.Endpoint(ctx, fmt.Sprintf("/status/%d", code), Options{})
}
func (c *Client) Echo(ctx context.Context, o Options) (*Response, error) {
	return c.Endpoint(ctx, "/anything", o)
}
