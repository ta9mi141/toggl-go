package reports

import (
	"testing"
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
