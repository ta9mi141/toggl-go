/*
Package toggl is a library of Toggl API v8 for Go programming language.

See API documentation for more details.
https://github.com/toggl/toggl_api_docs/blob/master/toggl_api.md
*/
package toggl

import (
	"net/http"
	"net/url"
)

const (
	basicAuthPassword string = "api_token" // Defined in Toggl API v8
	defaultBaseURL    string = "https://toggl.com"
)

// Client implements the basic request and response handling used by all types of APIs.
type Client struct {
	HTTPClient *http.Client
	APIToken   string
	Email      string
	Password   string
	URL        *url.URL
	header     http.Header
}

// Option represents optional parameters of NewClient.
type Option func(c *Client)

// HTTPClient sets an HTTP client to use when sending requests.
// By default, http.DefaultClient is used.
func HTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// APIToken sets an API token to authenticate yourself when sending requests.
func APIToken(apiToken string) Option {
	return func(c *Client) {
		c.APIToken = apiToken
	}
}

// Email sets an email to authenticate yourself when sending requests.
// If users use email for authentication, they must set password too.
func Email(email string) Option {
	return func(c *Client) {
		c.Email = email
	}
}

// Password sets a password to authenticate yourself when sending requests.
// If users use password for authentication, they must set email too.
func Password(password string) Option {
	return func(c *Client) {
		c.Password = password
	}
}

// NewClient returns a pointer to a new initialized client.
// Users must set APIToken or the pair of Email and Password for authentication.
func NewClient(options ...Option) *Client {
	url, _ := url.Parse(defaultBaseURL)
	newClient := &Client{
		HTTPClient: http.DefaultClient,
		URL:        url,
		header:     make(http.Header),
	}
	newClient.header.Set("Content-type", "application/json")
	for _, option := range options {
		option(newClient)
	}
	return newClient
}
