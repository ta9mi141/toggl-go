package internal

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	DefaultBaseURL    string = "https://api.track.toggl.com/"
	BasicAuthPassword string = "api_token" // Defined in Toggl Track API
)

func CheckResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 204:
		return nil
	}

	errorResponse := &ErrorResponse{StatusCode: resp.StatusCode, Header: resp.Header}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "")
	}
	errorResponse.Message = string(body)

	return errorResponse
}

func DecodeJSON(resp *http.Response, out any) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
