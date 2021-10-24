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

// Option is an option for a Toggl API v8 client.
type Option func(*Client)

// WithHTTPClient returns a Option that specifies the HTTP client for communication.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

var (
	// ErrAPITokenNotFound is returned when the provided API token is empty.
	ErrAPITokenNotFound = errors.New("the provided API token must be non-empty")

	// ErrContextNotFound is returned when the provided context is nil.
	ErrContextNotFound = errors.New("the provided ctx must be non-nil")
)
