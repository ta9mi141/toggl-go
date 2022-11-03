/*
Package reports is a library of Toggl Reports API v3 for the Go programming language.

See API documentation for more details.
https://developers.track.toggl.com/docs/reports_start
*/
package reports

import (
	"net/http"
	"net/url"

	"github.com/ta9mi141/toggl-go/track/internal"
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
