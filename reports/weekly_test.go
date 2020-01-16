package reports_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		Totals  []interface{} `json:"totals"`
		Details []struct {
			Title struct {
				Project string `json:"project"`
				Color   string `json:"color"`
				User    string `json:"user"`
			} `json:"title"`
			Totals []interface{} `json:"totals"`
		} `json:"details"`
	} `json:"data"`
}

func TestGetWeeklyEncodeRequestParameters(t *testing.T) {
	expectedQueryString := url.Values{
		"user_agent":   []string{userAgent},
		"workspace_id": []string{workSpaceId},
		"grouping":     []string{"users"},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualQueryString := r.URL.Query()
		if !reflect.DeepEqual(actualQueryString, expectedQueryString) {
			t.Error("Actual query string (" + actualQueryString.Encode() + ") is not as expected.")
		}
	}))

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	_ = client.GetWeekly(
		context.Background(),
		&reports.WeeklyRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
			Grouping: "users",
		},
		new(summaryReport),
	)
}

func TestGetWeeklyHandle_200_Ok(t *testing.T) {
	mockServer, weeklyTestData := setupMockServer(t, http.StatusOK, "testdata/weekly.json")
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
		t.Error("GetWeekly returns an error though it gets '200 OK'")
	}

	expectedWeeklyReport := new(weeklyReport)
	if err := json.Unmarshal(weeklyTestData, expectedWeeklyReport); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualWeeklyReport, expectedWeeklyReport) {
		t.Error("GetWeekly fails to decode weeklyReport")
	}
}

func TestGetWeeklyHandle_401_Unauthorized(t *testing.T) {
	mockServer, unauthorizedTestData := setupMockServer(t, http.StatusUnauthorized, "testdata/401_unauthorized.json")
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
		t.Error("GetWeekly doesn't return an error though it gets '401 Unauthorized'")
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

func TestGetWeeklyHandle_429_TooManyRequests(t *testing.T) {
	mockServer, _ := setupMockServer(t, http.StatusTooManyRequests, "testdata/429_too_many_requests.html")
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
		t.Error("GetWeekly doesn't return an error though it gets '429 Too Many Requests'")
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

func TestGetWeeklyWithoutContextReturnError(t *testing.T) {
	mockServer, _ := setupMockServer(t, http.StatusOK, "testdata/weekly.json")
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
		t.Error("GetWeekly doesn't return an error though it gets nil context")
	}
}
