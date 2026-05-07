package gateway

import (
	"context"
	"fmt"
	"net/http"
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

func (c *Client) Forward(ctx context.Context, source *http.Request) (*http.Response, error) {
	requestURL := c.baseURL + source.URL.Path
	if source.URL.RawQuery != "" {
		requestURL += "?" + source.URL.RawQuery
	}

	req, err := http.NewRequestWithContext(ctx, source.Method, requestURL, source.Body)
	if err != nil {
		return nil, fmt.Errorf("build gateway request: %w", err)
	}

	req.Header = source.Header.Clone()
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call gateway: %w", err)
	}

	return resp, nil
}
