/*
Package webhooks is a library of Toggl Webhooks API for the Go programming language.

See API documentation for more details.
https://developers.track.toggl.com/docs/webhooks_start
*/
package webhooks

import (
	"net/http"
	"net/url"

	"github.com/ta9mi141/toggl-go/track/internal"
)

// APIClient is a client for interacting with Toggl Webhooks API.
type APIClient struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiToken   string
}

// NewAPIClient creates a new Toggl Webhooks API client.
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

// Option is an option for a Toggl Webhooks API client.
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
