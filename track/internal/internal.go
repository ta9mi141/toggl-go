package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	DefaultBaseURL    string = "https://api.track.toggl.com/"
	BasicAuthPassword string = "api_token" // Defined in Toggl Track API
)

func NewRequest(ctx context.Context, httpMethod string, url *url.URL, input any) (*http.Request, error) {
	requestBody := io.Reader(nil)
	switch httpMethod {
	case http.MethodPost, http.MethodPut:
		b, err := json.Marshal(input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal input")
		}
		requestBody = bytes.NewReader(b)
	case http.MethodGet:
		values, err := query.Values(input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode input")
		}
		url.RawQuery = values.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, url.String(), requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new request with context")
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func Do(client *http.Client, req *http.Request, respBody any) error {
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send a request")
	}

	err = checkResponse(resp)
	if err != nil {
		return errors.Wrap(err, "failed to complete a request")
	}

	switch req.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		err = decodeJSON(resp, respBody)
		if err != nil {
			return errors.Wrap(err, "failed to decode response body")
		}
	}

	return nil
}

func checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 204:
		return nil
	}

	errorResponse := &ErrorResponse{StatusCode: resp.StatusCode, Header: resp.Header}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read error response")
	}
	errorResponse.Message = string(body)

	return errorResponse
}

func decodeJSON(resp *http.Response, out any) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
