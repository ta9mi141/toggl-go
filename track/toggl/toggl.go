/*
Package toggl is a library of Toggl API v9 for the Go programming language.

See API documentation for more details.
https://developers.track.toggl.com/docs/
*/
package toggl

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
	"github.com/ta9mi141/toggl-go/track/internal"
)

// Client is a client for interacting with Toggl API v9.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client

	apiToken string
}

// NewClient creates a new Toggl API v9 client.
func NewClient(options ...Option) *Client {
	baseURL, _ := url.Parse(internal.DefaultBaseURL)
	newClient := &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option.apply(newClient)
	}

	return newClient
}

// Option is an option for a Toggl API v9 client.
type Option interface {
	apply(*Client)
}

// WithHTTPClient returns a Option that specifies the HTTP client for communication.
func WithHTTPClient(httpClient *http.Client) Option {
	return &httpClientOption{httpClient: httpClient}
}

type httpClientOption struct {
	httpClient *http.Client
}

func (h *httpClientOption) apply(c *Client) {
	c.httpClient = h.httpClient
}

// WithAPIToken returns a Option that specifies an API token for authentication.
func WithAPIToken(apiToken string) Option {
	return apiTokenOption(apiToken)
}

type apiTokenOption string

func (a apiTokenOption) apply(c *Client) {
	c.apiToken = string(a)
}

// withBaseURL makes client testable by configurable URL.
func withBaseURL(baseURL string) Option {
	return baseURLOption(baseURL)
}

type baseURLOption string

func (b baseURLOption) apply(c *Client) {
	baseURL, _ := url.Parse(string(b))
	c.baseURL = baseURL
}

func (c *Client) httpGet(ctx context.Context, apiSpecificPath string, query, respBody any) error {
	req, err := c.newRequest(ctx, http.MethodGet, apiSpecificPath, query)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.do(req, respBody)
}

func (c *Client) httpPost(ctx context.Context, apiSpecificPath string, reqBody, respBody any) error {
	req, err := c.newRequest(ctx, http.MethodPost, apiSpecificPath, reqBody)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.do(req, respBody)
}

func (c *Client) httpPut(ctx context.Context, apiSpecificPath string, reqBody, respBody any) error {
	req, err := c.newRequest(ctx, http.MethodPut, apiSpecificPath, reqBody)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.do(req, respBody)
}

func (c *Client) httpDelete(ctx context.Context, apiSpecificPath string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, apiSpecificPath, nil)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.do(req, nil)
}

func (c *Client) newRequest(ctx context.Context, httpMethod, apiSpecificPath string, input any) (*http.Request, error) {
	url := c.baseURL
	url.Path = path.Join(url.Path, apiSpecificPath)

	req, err := internal.NewRequest(ctx, httpMethod, url, input)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	req.SetBasicAuth(c.apiToken, internal.BasicAuthPassword)

	return req, nil
}

func (c *Client) do(req *http.Request, respBody any) error {
	return internal.Do(c.httpClient, req, respBody)
}