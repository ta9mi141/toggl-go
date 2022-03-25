/*
Package toggl is a library of Toggl API v8 and Toggl Reports API v2 for the Go programming language.

See API documentation for more details.
https://github.com/toggl/toggl_api_docs/blob/master/toggl_api.md
*/
package toggl

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

const (
	defaultBaseURL string = "https://api.track.toggl.com/"

	basicAuthPassword string = "api_token" // Defined in Toggl API
)

// Client is a client for interacting with Toggl API.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client

	apiToken string
}

// NewClient creates a new Toggl API client.
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

// Option is an option for a Toggl API client.
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

func (c *Client) httpGet(ctx context.Context, apiSpecificPath string, respBody interface{}) error {
	req, err := c.newRequest(ctx, http.MethodGet, apiSpecificPath, nil)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.do(req, respBody)
}

func (c *Client) httpPut(ctx context.Context, apiSpecificPath string, reqBody, respBody interface{}) error {
	req, err := c.newRequest(ctx, http.MethodPut, apiSpecificPath, reqBody)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.do(req, respBody)
}

func (c *Client) newRequest(ctx context.Context, httpMethod, apiSpecificPath string, reqBody interface{}) (*http.Request, error) {
	url := *c.baseURL
	url.Path = path.Join(url.Path, apiSpecificPath)

	requestBody := io.Reader(nil)
	if httpMethod == http.MethodPut {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		requestBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, url.String(), requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	req.SetBasicAuth(c.apiToken, basicAuthPassword)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, respBody interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = checkResponse(resp)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = decodeJSON(resp, respBody)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 204:
		return nil
	}

	errorResponse := &errorResponse{statusCode: resp.StatusCode, header: resp.Header}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "")
	}
	errorResponse.message = string(body)

	return errorResponse
}

func decodeJSON(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
