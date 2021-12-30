package toggl

import (
	"fmt"
	"net/http"
)

type temporary interface {
	isTemporary() (bool, string)
}

// IsTemporary checks if the error is temporary and also returns the value of Retry-After header.
func IsTemporary(err error) (bool, string) {
	if e, ok := err.(temporary); ok {
		return e.isTemporary()
	}
	return false, ""
}

type timeout interface {
	isTimeout() bool
}

// IsTimeout checks if the error was caused by a timeout.
func IsTimeout(err error) bool {
	e, ok := err.(timeout)
	return ok && e.isTimeout()
}

type errorResponse struct {
	statusCode int
	message    string
	header     http.Header
}

func (e *errorResponse) Error() string {
	return fmt.Sprintf("%d %s", e.statusCode, e.message)
}

func (e *errorResponse) isTemporary() (bool, string) {
	isTemporary := e.statusCode == http.StatusTooManyRequests || e.statusCode == http.StatusServiceUnavailable
	retryAfter := e.header.Get("Retry-After")
	return isTemporary, retryAfter
}

func (e *errorResponse) isTimeout() bool {
	return e.statusCode == http.StatusRequestTimeout || e.statusCode == http.StatusGatewayTimeout
}
