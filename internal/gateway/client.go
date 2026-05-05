package gateway

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

	"sms-gateway/internal/config"
)

const DefaultHTTPTimeout = 30 * time.Second

type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func NewHTTPClient() *http.Client {
	return &http.Client{Timeout: DefaultHTTPTimeout}
}

func NewClient(cfg config.Config, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimRight(cfg.GatewayBaseURL, "/"),
		username:   cfg.Username,
		password:   cfg.Password,
		httpClient: httpClient,
	}
}

func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body []byte) (Response, error) {
	requestURL := c.baseURL + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, bytes.NewReader(body))
	if err != nil {
		return Response{}, fmt.Errorf("build gateway request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("call gateway: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return Response{}, fmt.Errorf("read gateway response: %w", err)
	}

	return Response{StatusCode: resp.StatusCode, Header: resp.Header.Clone(), Body: respBody}, nil
}

func (c *Client) DoJSON(ctx context.Context, method, path string, query url.Values, payload any) (Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("marshal gateway payload: %w", err)
	}
	return c.Do(ctx, method, path, query, body)
}
