package internal

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int
	Message    string
	Header     http.Header
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%d %s", e.StatusCode, e.Message)
}

func (e *ErrorResponse) IsTemporaryError() (bool, string) {
	isTemporary := e.StatusCode == http.StatusTooManyRequests || e.StatusCode == http.StatusServiceUnavailable
	retryAfter := e.Header.Get("Retry-After")
	return isTemporary, retryAfter
}

func (e *ErrorResponse) IsTimeoutError() bool {
	return e.StatusCode == http.StatusRequestTimeout || e.StatusCode == http.StatusGatewayTimeout
}
