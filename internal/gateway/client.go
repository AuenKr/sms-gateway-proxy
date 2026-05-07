package gateway

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/fx"

	"sms-gateway/internal/config"
)

const DefaultHTTPTimeout = 30 * time.Second

type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

type NewHTTPClientResult struct {
	fx.Out

	HTTPClient *http.Client
}

func NewHTTPClient() NewHTTPClientResult {
	return NewHTTPClientResult{HTTPClient: &http.Client{Timeout: DefaultHTTPTimeout}}
}

type NewClientParams struct {
	fx.In

	Config     config.Config
	HTTPClient *http.Client
}

type NewClientResult struct {
	fx.Out

	Client *Client
}

func NewClient(in NewClientParams) NewClientResult {
	return NewClientResult{
		Client: &Client{
			baseURL:    strings.TrimRight(in.Config.GatewayBaseURL, "/"),
			username:   in.Config.Username,
			password:   in.Config.Password,
			httpClient: in.HTTPClient,
		},
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
