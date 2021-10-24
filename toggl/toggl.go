/*
Package toggl is a library of Toggl API v8 for the Go programming language.

See API documentation for more details.
https://github.com/toggl/toggl_api_docs/blob/master/toggl_api.md
*/
package toggl

import (
	"errors"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL string = "https://api.track.toggl.com"
)

// Client is a client for interacting with Toggl API v8.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client

	apiToken string
}

// NewClient creates a new Toggl API v8 client.
func NewClient(options ...Option) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	newClient := &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
	for _, option := range options {
		option(newClient)
	}
	return newClient
}

// Option is an option for a Toggl API v8 client.
type Option func(*Client)

// WithHTTPClient returns a Option that specifies the HTTP client for communication.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithAPIToken returns a Option that specifies an API token for authentication.
func WithAPIToken(apiToken string) Option {
	return func(c *Client) {
		c.apiToken = apiToken
	}
}

var (
	// ErrContextNotFound is returned when the provided context is nil.
	ErrContextNotFound = errors.New("the provided context must be non-nil")

	// ErrAuthenticationFailure is returned when the API returns 403.
	ErrAuthenticationFailure = errors.New("authentication failed")
)
