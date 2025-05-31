package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseURL = "https://api.topstepx.com"

	DefaultTimeout = 30 * time.Second
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	userAgent  string
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		userAgent: "topstepx-go-client/1.0.0",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) GetToken() string {
	return c.token
}

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Query   url.Values
	Body    interface{}
}

type Response struct {
	*http.Response
	Body []byte
}

func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = req.Path
	if req.Query != nil {
		u.RawQuery = req.Query.Encode()
	}

	var body io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(jsonBody)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", c.userAgent)

	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	resp := &Response{
		Response: httpResp,
		Body:     respBody,
	}

	if httpResp.StatusCode >= 400 {
		return resp, fmt.Errorf("HTTP %d: %s", httpResp.StatusCode, string(respBody))
	}

	return resp, nil
}

func (c *Client) DoJSON(ctx context.Context, req *Request, v interface{}) error {
	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}

	if v != nil && len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, v); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
