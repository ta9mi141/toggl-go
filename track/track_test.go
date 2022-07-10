package track

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client.baseURL.String() != defaultBaseURL {
		errorf(t, client.baseURL.String(), defaultBaseURL)
	}
	if !reflect.DeepEqual(client.httpClient, http.DefaultClient) {
		errorf(t, client.httpClient, http.DefaultClient)
	}
}

func TestNewClientWithHTTPClient(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy.example.com:8080")
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	client := NewClient(WithHTTPClient(httpClient))

	if !reflect.DeepEqual(client.httpClient, httpClient) {
		errorf(t, client.httpClient, httpClient)
	}
}

func TestNewClientWithAPIToken(t *testing.T) {
	client := NewClient(WithAPIToken(apiToken))

	if client.apiToken != apiToken {
		errorf(t, client.apiToken, apiToken)
	}
}
