package toggl

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/ta9mi141/toggl-go/track/internal"
)

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client.baseURL.String() != defaultBaseURL {
		internal.Errorf(t, client.baseURL.String(), defaultBaseURL)
	}
	if !reflect.DeepEqual(client.httpClient, http.DefaultClient) {
		internal.Errorf(t, client.httpClient, http.DefaultClient)
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
		internal.Errorf(t, client.httpClient, httpClient)
	}
}

func TestNewClientWithAPIToken(t *testing.T) {
	client := NewClient(WithAPIToken(internal.APIToken))

	if client.apiToken != internal.APIToken {
		internal.Errorf(t, client.apiToken, internal.APIToken)
	}
}
