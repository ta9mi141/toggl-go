package reports_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/it-akumi/toggl-go/reports"
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
	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetWeekly(
		context.Background(),
		&reports.WeeklyRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		actualWeeklyReport,
	)
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

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	actualError := client.GetWeekly(
		context.Background(),
		&reports.WeeklyRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(weeklyReport),
	)
	if actualError == nil {
		t.Error("GetWeekly doesn't return error though it gets '401 Unauthorized'")
	}

	var actualReportsError reports.Error
	if errors.As(actualError, &actualReportsError) {
		expectedReportsError := new(reports.ReportsError)
		if err := json.Unmarshal(unauthorizedTestData, expectedReportsError); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualReportsError, expectedReportsError) {
			t.Error("GetWeekly fails to decode ReportsError though it returns reports.Error as expected")
		}
	} else {
		t.Error(actualError.Error())
	}
}

func TestGetWeekly_429_Too_Many_Requests(t *testing.T) {
	mockServer, _ := setupMockServer_429_Too_Many_Requests(t)
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	actualError := client.GetWeekly(
		context.Background(),
		&reports.WeeklyRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(weeklyReport),
	)
	if actualError == nil {
		t.Error("GetWeekly doesn't return error though it gets '429 Too Many Requests'")
	}

	var reportsError reports.Error
	if errors.As(actualError, &reportsError) {
		if reportsError.StatusCode() != http.StatusTooManyRequests {
			t.Error("GetWeekly fails to return '429 Too Many Requests' though it returns reports.Error as expected")
		}
	} else {
		t.Error(actualError.Error())
	}
}

func TestGetWeekly_WithoutContext(t *testing.T) {
	mockServer, _ := setupMockServer_200_Ok(t, "testdata/weekly.json")
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetWeekly(
		nil,
		&reports.WeeklyRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(weeklyReport),
	)
	if err == nil {
		t.Error("GetWeekly doesn't return error though it gets nil context")
	}
}
