package toggl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"reflect"
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

func newMockServer(t *testing.T, apiSpecificPath string, statusCode int, testdataFile string) *httptest.Server {
	testdata, err := ioutil.ReadFile(testdataFile)
	if err != nil {
		t.Fatal(err.Error())
	}

	mux := http.NewServeMux()
	pattern := path.Join("/", apiSpecificPath) // mockServer returns 404 page not found if pattern does not start with "/".
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, string(testdata))
	})

	mockServer := httptest.NewServer(mux)
	// The caller should call Close to shut down the server.
	return mockServer
}

func newMockServerToAssertRequestBody(t *testing.T, newRequestBody, expectedRequestBody interface{}) *httptest.Server {
	// The caller should call Close to shut down the server.
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawRequestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err.Error())
		}
		if err := json.Unmarshal(rawRequestBody, newRequestBody); err != nil {
			t.Helper()
			t.Errorf("%s\nnewRequestBody: %+#v\nrawRequestBody: %+#v\n", err.Error(), newRequestBody, rawRequestBody)
		}
		actualRequestBody := newRequestBody // Rename for readability
		if !reflect.DeepEqual(actualRequestBody, expectedRequestBody) {
			t.Helper()
			errorf(t, actualRequestBody, expectedRequestBody)
		}
	}))
}

// withBaseURL makes client testable by configurable URL.
func withBaseURL(baseURL string) Option {
	return baseURLOption(baseURL)
}

type baseURLOption string

func (b baseURLOption) apply(c *Client) {
	baseURL, _ := url.Parse(string(b))
	c.baseURL = baseURL
}
