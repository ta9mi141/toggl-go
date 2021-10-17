package toggl

import (
	"testing"
)

const (
	apiToken string = "api_token"
)

// It's highly likely that the got and want are compared just before calling errorf,
// so for the caller, the natural order of the arguments is got, want.
func errorf(t *testing.T, got, want interface{}) {
	t.Helper()
	// The order of the arguments in t.Errorf is swapped
	// because it's easier to read the error message when want is before got.
	t.Errorf("\nwant: %+#v\ngot : %+#v\n", want, got)
}
