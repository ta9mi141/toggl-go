package toggl

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/ta9mi141/toggl-go/track/internal"
)

func TestMain(m *testing.M) {
	time.Local = time.FixedZone("", 0)
	m.Run()
}

func TestNewAPIClient(t *testing.T) {
	apiClient := NewAPIClient()

	if apiClient.baseURL.String() != internal.DefaultBaseURL {
		internal.Errorf(t, apiClient.baseURL.String(), internal.DefaultBaseURL)
	}
	if !reflect.DeepEqual(apiClient.httpClient, http.DefaultClient) {
		internal.Errorf(t, apiClient.httpClient, http.DefaultClient)
	}
}

func TestNewAPIClientWithHTTPClient(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy.example.com:8080")
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	apiClient := NewAPIClient(WithHTTPClient(httpClient))

	if !reflect.DeepEqual(apiClient.httpClient, httpClient) {
		internal.Errorf(t, apiClient.httpClient, httpClient)
	}
}

func TestNewAPIClientWithAPIToken(t *testing.T) {
	apiClient := NewAPIClient(WithAPIToken(internal.APIToken))

	if apiClient.apiToken != internal.APIToken {
		internal.Errorf(t, apiClient.apiToken, internal.APIToken)
	}
}
