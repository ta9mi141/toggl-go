package reports_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ta9mi1shi1/toggl-go/reports"
)

// Users must define a data structure for each report
type summary struct {
	Data []struct {
		Title struct {
			Project string
		}
		Items []struct {
			Title struct {
				TimeEntry string `json:"time_entry"`
			}
			Time int
		}
	}
}

func Example() {
	client := reports.NewClient("YOUR_API_TOKEN")
	summaryReport := new(summary)
	err := client.GetSummary(
		context.Background(),
		&reports.SummaryRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   "YOUR_USER_AGENT",
				WorkspaceID: "YOUR_WORKSPACE_ID",
				Since:       time.Now().AddDate(0, 0, -7),
				Until:       time.Now(),
			},
		},
		summaryReport,
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, datum := range summaryReport.Data {
		fmt.Println(datum.Title.Project)
		for _, item := range datum.Items {
			fmt.Printf("Time entry: %v, Time: %d hours", item.Title.TimeEntry, item.Time/1000/3600)
		}
	}
}
