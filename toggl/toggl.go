/*
Package toggl is a library of Toggl API v8 for Go programming language.

See API documentation for more details.
https://github.com/toggl/toggl_api_docs/blob/master/toggl_api.md
*/
package toggl

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const (
	basicAuthPassword string = "api_token" // Defined in Toggl API v8
	defaultBaseURL    string = "https://toggl.com"
)

// Error wraps error interface with status code.
// Use errors.As method or type assertions with Error's StatusCode method
// to get detailed information about the error.
type Error interface {
	error
	// StatusCode returns HTTP status code of the error
	StatusCode() int
}

// TogglError represents a response of unsuccessful request.
type TogglError struct {
	Message string
	Code    int
}

func (e TogglError) Error() string {
	return e.Message
}

func (e TogglError) StatusCode() int {
	return e.Code
}

var (
	// ErrContextNotFound is returned when the provided context is nil.
	ErrContextNotFound = errors.New("The provided ctx must be non-nil")
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

func (c *Client) buildURL(endpoint string) string {
	c.URL.Path = endpoint
	return c.URL.String()
}

func (c *Client) httpGet(ctx context.Context, url string, resp interface{}) error {
	return nil
}

func (c *Client) httpPost(ctx context.Context, url string, req, resp interface{}) error {
	return nil
}

func (c *Client) httpPut(ctx context.Context, url string, req, resp interface{}) error {
	return nil
}

func (c *Client) httpDelete(ctx context.Context, url string) error {
	return nil
}

func checkResponse(resp *http.Response, err error) (*http.Response, error) {
	return nil, nil
}

func decodeJSON(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
