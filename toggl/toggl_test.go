package toggl_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ta9mi1shi1/toggl-go/toggl"
)

const (
	apiToken string = "api_token"
	email    string = "email"
	password string = "password"
)

func TestNewClient_WithAPIToken(t *testing.T) {
	expectedAPIToken := apiToken
	client := toggl.NewClient(toggl.APIToken(expectedAPIToken))

	if client.APIToken != expectedAPIToken {
		t.Error("client.APIToken = " + client.APIToken + ", [Expected: " + expectedAPIToken + "]")
	}
}

func TestNewClient_WithEmailAndPassword(t *testing.T) {
	expectedEmail := email
	expectedPassword := password
	client := toggl.NewClient(
		toggl.Email(expectedEmail),
		toggl.Password(expectedPassword),
	)

	if client.Email != expectedEmail {
		t.Error("client.Email = " + client.Email + ", [Expected: " + expectedEmail + "]")
	}
	if client.Password != expectedPassword {
		t.Error("client.Password = " + client.Password + ", [Expected: " + expectedPassword + "]")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	expectedTimeout := "5s"
	timeout, _ := time.ParseDuration(expectedTimeout)
	client := toggl.NewClient(toggl.HTTPClient(&http.Client{Timeout: timeout}))

	if client.HTTPClient.Timeout.String() != expectedTimeout {
		t.Error("client.HTTPClient.Timeout = " + client.HTTPClient.Timeout.String() + ", [Expected: " + expectedTimeout + "]")
	}
}

func ExampleNewClient() {
	client := toggl.NewClient(toggl.APIToken("YOUR_API_TOKEN"))
	fmt.Println(client.APIToken)
	// Output: YOUR_API_TOKEN
}

func ExampleNewClient_email() {
	client := toggl.NewClient(
		toggl.Email("YOUR_EMAIL"),
		toggl.Password("YOUR_PASSWORD"),
	)
	fmt.Println(client.Email, client.Password)
	// Output: YOUR_EMAIL YOUR_PASSWORD
}

// setupMockServer returns mockServer.
func setupMockServer(t *testing.T, httpStatus int, testdataFilePath string) *httptest.Server {
	testdata, err := ioutil.ReadFile(testdataFilePath)
	if err != nil {
		t.Error(err.Error())
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
		fmt.Fprint(w, string(testdata))
	}))

	return mockServer
}

// baseURL makes client testable by configurable URL.
func baseURL(rawurl string) toggl.Option {
	return func(c *toggl.Client) {
		url, _ := url.Parse(rawurl)
		c.URL = url
	}
}
