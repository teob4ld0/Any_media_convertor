// Package client provides a thin wrapper around net/http for making API requests.
package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client wraps net/http.Client with convenience methods.
type Client struct {
	inner *http.Client
}

// New returns a Client with a 30-second timeout.
func New() *Client {
	return &Client{
		inner: &http.Client{Timeout: 30 * time.Second},
	}
}

// Get performs an HTTP GET and returns the raw response body.
// headers is an optional map of extra request headers.
func (c *Client) Get(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building GET request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.do(req)
}

// Post performs an HTTP POST with a plain-text body and returns the raw response body.
func (c *Client) Post(url string, headers map[string]string, body string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building POST request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.do(req)
}

// do executes req, reads the body, and returns an error on non-2xx status.
func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.inner.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request to %s: %w", req.URL, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, req.URL)
	}
	return data, nil
}