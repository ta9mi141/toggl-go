package track

type temporaryError interface {
	IsTemporaryError() (bool, string)
}

// IsTemporary checks if the error is temporary and also returns the value of Retry-After header.
func IsTemporary(err error) (bool, string) {
	if e, ok := err.(temporaryError); ok {
		return e.IsTemporaryError()
	}
	return false, ""
}

type timeoutError interface {
	IsTimeoutError() bool
}

// IsTimeout checks if the error was caused by a timeout.
func IsTimeout(err error) bool {
	e, ok := err.(timeoutError)
	return ok && e.IsTimeoutError()
}
