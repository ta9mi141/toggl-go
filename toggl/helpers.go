package toggl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func newMockServer(t *testing.T, path string, statusCode int, testdataFile string) *httptest.Server {
	testdata, err := ioutil.ReadFile(testdataFile)
	if err != nil {
		t.Fatal(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, string(testdata))
	})

	mockServer := httptest.NewServer(mux)
	// The caller should call Close to shut down the server.
	return mockServer
}

// withBaseURL makes client testable by configurable URL.
func withBaseURL(baseURL string) Option {
	return baseURLOption(baseURL)
}

type baseURLOption string

func (b baseURLOption) apply(c *Client) {
	c.setBaseURL(string(b))
}
