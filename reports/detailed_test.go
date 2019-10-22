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

type detailedReport struct {
	TotalCount int `json:"total_count"`
	PerPage    int `json:"per_page"`
	Data       []struct {
		User        string `json:"user"`
		Project     string `json:"project"`
		Description string `json:"description"`
	} `json:"data"`
}

func TestGetDetailed_200_Ok(t *testing.T) {
	mockServer, detailedTestData := setupMockServer_200_Ok(t, "testdata/detailed.json")
	defer mockServer.Close()

	actualDetailedReport := new(detailedReport)
	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetDetailed(
		context.Background(),
		&reports.DetailedRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		actualDetailedReport,
	)
	if err != nil {
		t.Error("GetDetailed returns error though it gets '200 OK'")
	}

	expectedDetailedReport := new(detailedReport)
	if err := json.Unmarshal(detailedTestData, expectedDetailedReport); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualDetailedReport, expectedDetailedReport) {
		t.Error("GetDetailed fails to decode detailedReport")
	}
}

func TestGetDetailed_401_Unauthorized(t *testing.T) {
	mockServer, unauthorizedTestData := setupMockServer_401_Unauthorized(t)
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	actualError := client.GetDetailed(
		context.Background(),
		&reports.DetailedRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(detailedReport),
	)
	if actualError == nil {
		t.Error("GetDetailed doesn't return error though it gets '401 Unauthorized'")
	}

	var actualReportsError reports.Error
	if errors.As(actualError, &actualReportsError) {
		expectedReportsError := new(reports.ReportsError)
		if err := json.Unmarshal(unauthorizedTestData, expectedReportsError); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualReportsError, expectedReportsError) {
			t.Error("GetDetailed fails to decode ReportsError though it returns reports.Error as expected")
		}
	} else {
		t.Error(actualError.Error())
	}
}

func TestGetDetailed_429_Too_Many_Requests(t *testing.T) {
	mockServer, _ := setupMockServer_429_Too_Many_Requests(t)
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	actualError := client.GetDetailed(
		context.Background(),
		&reports.DetailedRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(detailedReport),
	)
	if actualError == nil {
		t.Error("GetDetailed doesn't return error though it gets '429 Too Many Requests'")
	}

	var reportsError reports.Error
	if errors.As(actualError, &reportsError) {
		if reportsError.StatusCode() != http.StatusTooManyRequests {
			t.Error("GetDetailed fails to return '429 Too Many Requests' though it returns reports.Error as expected")
		}
	} else {
		t.Error(actualError.Error())
	}
}

func TestGetDetailed_WithoutContext(t *testing.T) {
	mockServer, _ := setupMockServer_200_Ok(t, "testdata/detailed.json")
	defer mockServer.Close()

	client := reports.NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetDetailed(
		nil,
		&reports.DetailedRequestParameters{
			StandardRequestParameters: &reports.StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(detailedReport),
	)
	if err == nil {
		t.Error("GetDetailed doesn't return error though it gets nil context")
	}
}
