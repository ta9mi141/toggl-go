package internal

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
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
