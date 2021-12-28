/*
Package toggl is a library of Toggl API v8 and Toggl Reports API v2 for the Go programming language.

See API documentation for more details.
https://github.com/toggl/toggl_api_docs/blob/master/toggl_api.md
*/
package toggl

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const (
	defaultBaseURL string = "https://api.track.toggl.com/api/v8/"

	basicAuthPassword string = "api_token" // Defined in Toggl API
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
		option.apply(newClient)
	}

	return newClient
}

// Option is an option for a Toggl API v8 client.
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

var (
	// ErrAuthenticationFailure is returned when the API returns 403.
	ErrAuthenticationFailure = errors.New("authentication failed")
)

func (c *Client) httpGet(ctx context.Context, path string, respBody interface{}) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return c.do(req, respBody)
}

func (c *Client) newRequest(ctx context.Context, httpMethod, apiSpecificPath string, reqBody interface{}) (*http.Request, error) {
	url := *c.baseURL
	url.Path = path.Join(url.Path, apiSpecificPath)

	requestBody := io.Reader(nil)

	req, err := http.NewRequestWithContext(ctx, httpMethod, url.String(), requestBody)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiToken, basicAuthPassword)

	return req, nil
}

func (c *Client) do(req *http.Request, respBody interface{}) error {
	resp, err := checkResponse(c.httpClient.Do(req))
	if err != nil {
		return err
	}
	err = decodeJSON(resp, respBody)
	if err != nil {
		return err
	}
	return nil
}

func checkResponse(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}

	switch {
	case resp.StatusCode == http.StatusForbidden:
		// Since resp.Body is empty if the API returns 403 Forbidden,
		// an explicit error message is needed to indicate why the request failed.
		return nil, ErrAuthenticationFailure
	case resp.StatusCode <= 199 || 300 <= resp.StatusCode:
		message, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(message))
	default:
		return resp, nil
	}
}

func decodeJSON(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
