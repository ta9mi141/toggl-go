package reports_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/it-akumi/toggl-go/reports"
)

type summaryReport struct {
	Data []struct {
		Id    int `json:"id"`
		Title struct {
			Project string `json:"project"`
			Color   string `json:"color"`
			User    string `json:"user"`
		} `json:"title"`
		Time  int `json:"time"`
		Items []struct {
			Title struct {
				Project   string `json:"project"`
				User      string `json:"user"`
				TimeEntry string `json:"time_entry"`
			} `json:"title"`
			Time int `json:"time"`
		} `json:"items"`
	} `json:"data"`
}

func TestGetSummary_200_Ok(t *testing.T) {
	mockServer, summaryTestData := setupMockServer_200_Ok(t, "testdata/summary.json")
	defer mockServer.Close()

	actualSummaryReport := new(summaryReport)
	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetSummary(
		context.Background(),
		&reports.SummaryRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		actualSummaryReport,
	)
	if err != nil {
		t.Error("GetSummary returns error though it gets '200 OK'")
	}

	expectedSummaryReport := new(summaryReport)
	if err := json.Unmarshal(summaryTestData, expectedSummaryReport); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualSummaryReport, expectedSummaryReport) {
		t.Error("GetSummary fails to decode summaryReport")
	}
}

func TestGetSummary_401_Unauthorized(t *testing.T) {
	mockServer, unauthorizedTestData := setupMockServer_401_Unauthorized(t)
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	actualReportsError := client.GetSummary(
		context.Background(),
		&reports.SummaryRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(summaryReport),
	)
	if actualReportsError == nil {
		t.Error("GetSummary doesn't return error though it gets '401 Unauthorized'")
	}

	expectedReportsError := new(reports.ReportsError)
	if err := json.Unmarshal(unauthorizedTestData, expectedReportsError); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualReportsError, expectedReportsError) {
		t.Error("GetSummary fails to decode ReportsError though it returns error as expected")
	}
}

func TestGetSummary_429_Too_Many_Requests(t *testing.T) {
	mockServer, tooManyRequestsTestData := setupMockServer_429_Too_Many_Requests(t)
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	actualReportsError := client.GetSummary(
		context.Background(),
		&reports.SummaryRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(summaryReport),
	)
	if actualReportsError == nil {
		t.Error("GetSummary doesn't return error though it gets '429 Too Many Requests'")
	}

	expectedReportsError := new(reports.ReportsError)
	if err := json.Unmarshal(tooManyRequestsTestData, expectedReportsError); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualReportsError, expectedReportsError) {
		t.Error("GetSummary fails to decode ReportsError though it returns error as expected")
	}
}

func TestGetSummary_WithoutContext(t *testing.T) {
	mockServer, _ := setupMockServer_200_Ok(t, "testdata/summary.json")
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetSummary(
		nil,
		&reports.SummaryRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(summaryReport),
	)
	if err == nil {
		t.Error("GetSummary doesn't return error though it gets nil context")
	}
}
