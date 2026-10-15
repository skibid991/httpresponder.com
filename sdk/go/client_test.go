package httpresponder

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestEncoding(t *testing.T) {
	c := New()
	target := "https://example.com/a?x=1&y=two#part"
	s, e := c.URL("/test", map[string]any{"url": target, "data": "+/= 🌎", "enabled": false, "headers": map[string]string{"X-Test": "yes"}})
	if e != nil {
		t.Fatal(e)
	}
	u, _ := url.Parse(s)
	q := u.Query()
	if q.Get("url") != target || q.Get("data") != "+/= 🌎" || q.Get("enabled") != "false" {
		t.Fatal(q)
	}
	var h map[string]string
	if e = json.Unmarshal([]byte(q.Get("headers")), &h); e != nil || h["X-Test"] != "yes" {
		t.Fatal(q)
	}
}
func TestResponses(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/redirect":
			w.Header().Set("Location", "/ok")
			w.WriteHeader(302)
		case "/error":
			w.WriteHeader(418)
		default:
			if r.Method == "POST" {
				w.Header().Set("X-Auth", r.Header.Get("Authorization"))
				var v map[string]string
				_ = json.NewDecoder(r.Body).Decode(&v)
				_ = json.NewEncoder(w).Encode(v)
			} else {
				_, _ = w.Write([]byte{0, 255, 128})
			}
		}
	}))
	defer s.Close()
	c := New()
	c.BaseURL = s.URL
	ctx := context.Background()
	r, e := c.Request(ctx, "/ok", Options{})
	if e != nil || len(r.Body) != 3 || r.Body[1] != 255 {
		t.Fatal(r, e)
	}
	_, e = c.Request(ctx, "/error", Options{})
	var apiErr *APIError
	if !errors.As(e, &apiErr) || apiErr.StatusCode != 418 {
		t.Fatal(e)
	}
	r, e = c.Request(ctx, "/redirect", Options{AllowHTTPError: true})
	if e != nil || r.StatusCode != 302 {
		t.Fatal(r, e)
	}
	r, e = c.Request(ctx, "/test", Options{Method: "POST", Body: map[string]string{"text": "🌎"}, Headers: map[string]string{"Authorization": "Bearer sample"}})
	if e != nil || r.Header.Get("X-Auth") != "Bearer sample" {
		t.Fatal(r, e)
	}
	var body map[string]string
	if e = r.JSON(&body); e != nil || body["text"] != "🌎" {
		t.Fatal(body, e)
	}
}
func TestContextAndPath(t *testing.T) {
	c := New()
	if _, e := c.URL("//example.com/", nil); e == nil {
		t.Fatal("expected bad path")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e := c.Request(ctx, "/test", Options{}); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
