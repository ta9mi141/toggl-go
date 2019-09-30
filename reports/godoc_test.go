package reports_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/it-akumi/toggl-go/reports"
)

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
