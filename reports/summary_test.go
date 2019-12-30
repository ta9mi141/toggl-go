package reports_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
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

func TestGetSummaryEncodeRequestParameters(t *testing.T) {
	expectedQueryString := url.Values{
		"user_agent":   []string{userAgent},
		"workspace_id": []string{workSpaceId},
		"grouping":     []string{"projects"},
	}

	mockServer := setupMockServer_TestingQueryString(t, expectedQueryString)
	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	_ = client.GetSummary(
		context.Background(),
		&reports.SummaryRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
			Grouping: "projects",
		},
		new(summaryReport),
	)
}

func TestGetSummaryHandle_200_Ok(t *testing.T) {
	mockServer, summaryTestData := setupMockServer(t, http.StatusOK, "testdata/summary.json")
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

func TestGetSummaryHandle_401_Unauthorized(t *testing.T) {
	mockServer, unauthorizedTestData := setupMockServer(t, http.StatusUnauthorized, "testdata/401_unauthorized.json")
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	actualError := client.GetSummary(
		context.Background(),
		&reports.SummaryRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(summaryReport),
	)
	if actualError == nil {
		t.Error("GetSummary doesn't return error though it gets '401 Unauthorized'")
	}

	var actualReportsError reports.Error
	if errors.As(actualError, &actualReportsError) {
		expectedReportsError := new(reports.ReportsError)
		if err := json.Unmarshal(unauthorizedTestData, expectedReportsError); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualReportsError, expectedReportsError) {
			t.Error("GetSummary fails to decode ReportsError though it returns reports.Error as expected")
		}
	} else {
		t.Error(actualError.Error())
	}
}

func TestGetSummaryHandle_429_TooManyRequests(t *testing.T) {
	mockServer, _ := setupMockServer(t, http.StatusTooManyRequests, "testdata/429_too_many_requests.html")
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	actualError := client.GetSummary(
		context.Background(),
		&reports.SummaryRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(summaryReport),
	)
	if actualError == nil {
		t.Error("GetSummary doesn't return error though it gets '429 Too Many Requests'")
	}

	var reportsError reports.Error
	if errors.As(actualError, &reportsError) {
		if reportsError.StatusCode() != http.StatusTooManyRequests {
			t.Error("GetSummary fails to return '429 Too Many Requests' though it returns reports.Error as expected")
		}
	} else {
		t.Error(actualError.Error())
	}
}

func TestGetSummaryWithoutContextReturnError(t *testing.T) {
	mockServer, _ := setupMockServer(t, http.StatusOK, "testdata/summary.json")
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
