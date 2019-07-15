package reports

import (
	"encoding/json"
	"reflect"
	"testing"
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
	client := NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetSummary(&SummaryRequestParameters{
		StandardRequestParameters: &StandardRequestParameters{
			UserAgent:   userAgent,
			WorkSpaceId: workSpaceId,
		},
	}, actualSummaryReport)
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

	client := NewClient(apiToken, baseURL(mockServer.URL))
	actualReportsError := client.GetSummary(&SummaryRequestParameters{
		StandardRequestParameters: &StandardRequestParameters{
			UserAgent:   userAgent,
			WorkSpaceId: workSpaceId,
		},
	}, new(summaryReport))
	if actualReportsError == nil {
		t.Error("GetSummary doesn't return error though it gets '401 Unauthorized'")
	}

	expectedReportsError := new(ReportsError)
	if err := json.Unmarshal(unauthorizedTestData, expectedReportsError); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualReportsError, expectedReportsError) {
		t.Error("GetSummary fails to decode ReportsError though it returns error as expected")
	}
}
