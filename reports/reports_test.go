package reports

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

const (
	apiToken string = "api_token"
)

func errorf(t *testing.T, got, want interface{}) {
	t.Errorf("\ngot : %+#v\nwant: %+#v\n", got, want)
}

func TestNewClient(t *testing.T) {
	client := NewClient(apiToken)

	if client.baseURL.String() != defaultBaseURL {
		errorf(t, client.baseURL.String(), defaultBaseURL)
	}
	if !reflect.DeepEqual(client.httpClient, http.DefaultClient) {
		errorf(t, client.httpClient, http.DefaultClient)
	}
	if client.apiToken != apiToken {
		errorf(t, client.apiToken, apiToken)
	}
}

func TestNewClientWithHTTPClient(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy.example.com:8080")
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	client := NewClient(apiToken, WithHTTPClient(httpClient))

	if !reflect.DeepEqual(client.httpClient, httpClient) {
		errorf(t, client.httpClient, httpClient)
	}
}
