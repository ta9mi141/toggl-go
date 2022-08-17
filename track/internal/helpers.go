package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const (
	APIToken string = "api_token"
)

// It's highly likely that the got and want are compared just before calling errorf,
// so for the caller, the natural order of the arguments is got, want.
func Errorf(t *testing.T, got, want any) {
	t.Helper()
	// The order of the arguments in t.Errorf is swapped
	// because it's easier to read the error message when want is before got.
	t.Errorf("\nwant: %+#v\ngot : %+#v\n", want, got)
}

func NewMockServer(t *testing.T, apiSpecificPath string, statusCode int, testdataFile string) *httptest.Server {
	testdata, err := os.ReadFile(testdataFile)
	if err != nil {
		t.Fatal(err.Error())
	}

	mux := http.NewServeMux()
	pattern := path.Join("/", apiSpecificPath) // mockServer returns 404 page not found if pattern does not start with "/".
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Type to emulate behavior of Toggl API.
		// Since changing the header map after a call to WriteHeader has no effect,
		// Content-Type must be set before a call to WriteHeader.
		switch filepath.Ext(testdataFile) {
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		case ".txt":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}

		w.WriteHeader(statusCode)
		fmt.Fprint(w, string(testdata))
	})

	mockServer := httptest.NewServer(mux)
	// The caller should call Close to shut down the server.
	return mockServer
}

func NewMockServerToAssertRequestBody(t *testing.T, expectedRequestBody string) *httptest.Server {
	// The caller should call Close to shut down the server.
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawRequestBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err.Error())
		}
		actualRequestBody := string(rawRequestBody)
		if actualRequestBody != expectedRequestBody {
			Errorf(t, actualRequestBody, expectedRequestBody)
		}
	}))
}

func NewMockServerToAssertQuery(t *testing.T, expectedQuery string) *httptest.Server {
	// The caller should call Close to shut down the server.
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualQuery := r.URL.Query().Encode()
		if actualQuery != expectedQuery {
			Errorf(t, actualQuery, expectedQuery)
		}
	}))

}
