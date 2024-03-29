/*
Package reports is a library of Toggl Reports API v3 for the Go programming language.

See API documentation for more details.
https://developers.track.toggl.com/docs/reports_start
*/
package reports

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
	"github.com/ta9mi141/toggl-go/track/internal"
)

const (
	reportsPath string = "reports/api/v3/workspace"
)

// APIClient is a client for interacting with Toggl Reports API v3.
type APIClient struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiToken   string
}

// NewAPIClient creates a new Toggl Reports API v3 client.
func NewAPIClient(apiToken string, options ...Option) *APIClient {
	baseURL, _ := url.Parse(internal.DefaultBaseURL)
	newAPIClient := &APIClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		apiToken:   apiToken,
	}

	for _, option := range options {
		option.apply(newAPIClient)
	}

	return newAPIClient
}

// Option is an option for a Toggl Reports API v3 client.
type Option interface {
	apply(*APIClient)
}

// WithHTTPClient returns a Option that specifies the HTTP client for communication.
func WithHTTPClient(httpClient *http.Client) Option {
	return &httpClientOption{httpClient: httpClient}
}

type httpClientOption struct {
	httpClient *http.Client
}

func (h *httpClientOption) apply(c *APIClient) {
	c.httpClient = h.httpClient
}

// withBaseURL makes client testable by configurable URL.
func withBaseURL(baseURL string) Option {
	return baseURLOption(baseURL)
}

type baseURLOption string

func (b baseURLOption) apply(c *APIClient) {
	baseURL, _ := url.Parse(string(b))
	c.baseURL = baseURL
}

func (c *APIClient) httpPost(ctx context.Context, apiSpecificPath string, reqBody, respBody any) error {
	req, err := c.newRequest(ctx, http.MethodPost, apiSpecificPath, reqBody)
	if err != nil {
		return errors.Wrap(err, "failed to create a new POST request")
	}
	return c.do(req, respBody)
}

func (c *APIClient) newRequest(ctx context.Context, httpMethod, apiSpecificPath string, input any) (*http.Request, error) {
	url := c.baseURL
	url.Path = path.Join(url.Path, apiSpecificPath)

	req, err := internal.NewRequest(ctx, httpMethod, url, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new request")
	}

	req.SetBasicAuth(c.apiToken, internal.BasicAuthPassword)

	return req, nil
}

func (c *APIClient) do(req *http.Request, respBody any) error {
	return internal.Do(c.httpClient, req, respBody)
}
