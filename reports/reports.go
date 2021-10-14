/*
Package reports is a library of Toggl Reports API v2 for the Go programming language.

See API documentation for more details.
https://github.com/toggl/toggl_api_docs/blob/master/reports.md
*/
package reports

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL string = "https://api.track.toggl.com"
)

// Client is a client for interacting with Toggl Reports API v2.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client

	apiToken string
}

// NewClient creates a new Toggl Reports API v2 client.
func NewClient(apiToken string, options ...Option) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	newClient := &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		apiToken:   apiToken,
	}
	for _, option := range options {
		option(newClient)
	}
	return newClient
}

// Option is an option for a Toggl Reports API v2 client.
type Option func(*Client)

// WithHTTPClient returns a Option that specifies the HTTP client for communication.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
