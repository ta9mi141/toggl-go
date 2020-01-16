package reports_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/it-akumi/toggl-go/reports"
)

const (
	apiToken    string = "api_token"
	userAgent   string = "user_agent"
	workSpaceId string = "workspace_id"
)

func TestNewClient(t *testing.T) {
	expectedAPIToken := apiToken
	client := reports.NewClient(expectedAPIToken)

	if client.APIToken != expectedAPIToken {
		t.Error("client.APIToken = " + client.APIToken + ", [Expected: " + expectedAPIToken + "]")
	}
}

func TestNewClientWithHTTPClient(t *testing.T) {
	expectedTimeout := "5s"
	timeout, _ := time.ParseDuration(expectedTimeout)
	client := reports.NewClient(apiToken, reports.HTTPClient(&http.Client{Timeout: timeout}))

	if client.HTTPClient.Timeout.String() != expectedTimeout {
		t.Error("client.HTTPClient.Timeout = " + client.HTTPClient.Timeout.String() + ", [Expected: " + expectedTimeout + "]")
	}
}

func ExampleNewClient() {
	client := reports.NewClient("YOUR_API_TOKEN")
	fmt.Println(client.APIToken)
	// Output: YOUR_API_TOKEN
}

func ExampleNewClient_option() {
	client := reports.NewClient("YOUR_API_TOKEN", reports.HTTPClient(
		&http.Client{Timeout: 5 * time.Second},
	))
	fmt.Println(client.HTTPClient.Timeout)
	// Output: 5s
}

// baseURL makes client testable by configurable URL.
func baseURL(rawurl string) reports.Option {
	return func(c *reports.Client) {
		url, _ := url.Parse(rawurl)
		c.URL = url
	}
}

// setupMockServer returns mockServer and testdata.
// Since testdata is same as the response of mockServer,
// users don't need to send a request to mockServer to get the response.
func setupMockServer(t *testing.T, httpStatus int, testdataFilePath string) (*httptest.Server, []byte) {
	testdata, err := ioutil.ReadFile(testdataFilePath)
	if err != nil {
		t.Error(err.Error())
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
		fmt.Fprintf(w, string(testdata))
	}))

	return mockServer, testdata
}
