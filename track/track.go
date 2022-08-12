/*
Package track is a library of Toggl Track API for the Go programming language.

See API documentation for more details.
https://developers.track.toggl.com/docs/
*/
package track

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	defaultBaseURL    string = "https://api.track.toggl.com/"
	basicAuthPassword string = "api_token" // Defined in Toggl Track API
)

// Client is a client for interacting with Toggl Track API.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client

	apiToken string
}

// NewClient creates a new Toggl Track API client.
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

// Option is an option for a Toggl Track API client.
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

func (c *Client) httpGet(ctx context.Context, apiSpecificPath string, query, respBody interface{}) error {
	req, err := c.newRequest(ctx, http.MethodGet, apiSpecificPath, query)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return c.do(req, respBody)
}

func (c *Client) httpPost(ctx context.Context, apiSpecificPath string, reqBody, respBody interface{}) error {
	req, err := c.newRequest(ctx, http.MethodPost, apiSpecificPath, reqBody)
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

func (c *Client) httpDelete(ctx context.Context, apiSpecificPath string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, apiSpecificPath, nil)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return c.do(req, nil)
}

func (c *Client) newRequest(ctx context.Context, httpMethod, apiSpecificPath string, input interface{}) (*http.Request, error) {
	url := *c.baseURL
	url.Path = path.Join(url.Path, apiSpecificPath)

	requestBody := io.Reader(nil)
	switch httpMethod {
	case http.MethodPost, http.MethodPut:
		b, err := json.Marshal(input)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		requestBody = bytes.NewReader(b)
	case http.MethodGet:
		values, err := query.Values(input)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		url.RawQuery = values.Encode()
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

	switch req.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		err = decodeJSON(resp, respBody)
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	return nil
}

func checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 204:
		return nil
	}

	errorResponse := &errorResponse{statusCode: resp.StatusCode, header: resp.Header}
	body, err := io.ReadAll(resp.Body)
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
