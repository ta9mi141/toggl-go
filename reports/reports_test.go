package reports

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	expectedAPIToken := "api_token"
	client := NewClient(expectedAPIToken)

	if client.apiToken != expectedAPIToken {
		t.Error("client.apiToken = " + client.apiToken + ", [Expected: " + expectedAPIToken + "]")
	}

	if client.url.String() != defaultBaseURL {
		t.Error("client.url.String() = " + client.url.String() + ", [Expected: " + defaultBaseURL + "]")
	}

	expectedContentType := "application/json"
	if client.header.Get("Content-type") != expectedContentType {
		t.Error(`client.header.Get("Content-type") = ` + client.header.Get("Content-type") + ", [Expected: " + expectedContentType + "]")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	expectedAPIToken := "api_token"
	expectedTimeout := "5s"
	timeout, _ := time.ParseDuration(expectedTimeout)
	client := NewClient(expectedAPIToken, HTTPClient(&http.Client{Timeout: timeout}))

	if client.apiToken != expectedAPIToken {
		t.Error("client.apiToken = " + client.apiToken + ", [Expected: " + expectedAPIToken + "]")
	}

	if client.httpClient.Timeout.String() != expectedTimeout {
		t.Error("client.httpClient.Timeout = " + client.httpClient.Timeout.String() + ", [Expected: " + expectedTimeout + "]")
	}

	if client.url.String() != defaultBaseURL {
		t.Error("client.url.String() = " + client.url.String() + ", [Expected: " + defaultBaseURL + "]")
	}

	expectedContentType := "application/json"
	if client.header.Get("Content-type") != expectedContentType {
		t.Error(`client.header.Get("Content-type") = ` + client.header.Get("Content-type") + ", [Expected: " + expectedContentType + "]")
	}
}

func setupMockServer_200_Ok(t *testing.T, testdataFilePath string) (*httptest.Server, []byte) {
	testdata, err := ioutil.ReadFile(testdataFilePath)
	if err != nil {
		t.Error(err.Error())
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(testdata))
	}))

	return mockServer, testdata
}

func setupMockServer_401_Unauthorized(t *testing.T) (*httptest.Server, []byte) {
	errorTestData, err := ioutil.ReadFile("testdata/401_unauthorized.json")
	if err != nil {
		t.Error(err.Error())
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, string(errorTestData))
	}))

	return mockServer, errorTestData
}

const (
	apiToken    string = "api_token"
	userAgent   string = "user_agent"
	workSpaceId string = "workspace_id"
)
