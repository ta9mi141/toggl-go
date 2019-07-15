package reports

import (
	"encoding/json"
	"reflect"
	"testing"
)

type weeklyReport struct {
	WeekTotals []interface{} `json:"week_totals"`
	Data       []struct {
		Title struct {
			Project string `json:"project"`
			Color   string `json:"color"`
			User    string `json:"user"`
		} `json:"title"`
		Totals  []interface{} `json:"totals`
		Details []struct {
			Title struct {
				Project string `json:"project"`
				Color   string `json:"color"`
				User    string `json:"user"`
			} `json:"title"`
			Totals []interface{} `json:"totals`
		} `json:"details"`
	} `json:"data"`
}

func TestGetWeekly_200_Ok(t *testing.T) {
	mockServer, weeklyTestData := setupMockServer_200_Ok(t, "testdata/weekly.json")
	defer mockServer.Close()

	actualWeeklyReport := new(weeklyReport)
	client := NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetWeekly(&WeeklyRequestParameters{
		StandardRequestParameters: &StandardRequestParameters{
			UserAgent:   userAgent,
			WorkSpaceId: workSpaceId,
		},
	}, actualWeeklyReport)
	if err != nil {
		t.Error("GetWeekly returns error though it gets '200 OK'")
	}

	expectedWeeklyReport := new(weeklyReport)
	if err := json.Unmarshal(weeklyTestData, expectedWeeklyReport); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualWeeklyReport, expectedWeeklyReport) {
		t.Error("GetWeekly fails to decode weeklyReport")
	}
}

func TestGetWeekly_401_Unauthorized(t *testing.T) {
	mockServer, unauthorizedTestData := setupMockServer_401_Unauthorized(t)
	defer mockServer.Close()

	client := NewClient(apiToken, baseURL(mockServer.URL))
	actualReportsError := client.GetWeekly(&WeeklyRequestParameters{
		StandardRequestParameters: &StandardRequestParameters{
			UserAgent:   userAgent,
			WorkSpaceId: workSpaceId,
		},
	}, new(weeklyReport))
	if actualReportsError == nil {
		t.Error("GetWeekly doesn't return error though it gets '401 Unauthorized'")
	}

	expectedReportsError := new(ReportsError)
	if err := json.Unmarshal(unauthorizedTestData, expectedReportsError); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualReportsError, expectedReportsError) {
		t.Error("GetWeekly fails to decode ReportsError though it returns error as expected")
	}
}
