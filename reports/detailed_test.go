package reports

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
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
	client := NewClient(&Config{ApiToken: apiToken}, baseURL(mockServer.URL))
	err := client.GetDetailed(
		context.Background(),
		&DetailedRequestParameters{
			StandardRequestParameters: &StandardRequestParameters{
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

	client := NewClient(&Config{ApiToken: apiToken}, baseURL(mockServer.URL))
	actualReportsError := client.GetDetailed(
		context.Background(),
		&DetailedRequestParameters{
			StandardRequestParameters: &StandardRequestParameters{
				UserAgent:   userAgent,
				WorkSpaceId: workSpaceId,
			},
		},
		new(detailedReport),
	)
	if actualReportsError == nil {
		t.Error("GetDetailed doesn't return error though it gets '401 Unauthorized'")
	}

	expectedReportsError := new(ReportsError)
	if err := json.Unmarshal(unauthorizedTestData, expectedReportsError); err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(actualReportsError, expectedReportsError) {
		t.Error("GetDetailed fails to decode ReportsError though it returns error as expected")
	}
}

func TestGetDetailed_WithoutContext(t *testing.T) {
	mockServer, _ := setupMockServer_200_Ok(t, "testdata/detailed.json")
	defer mockServer.Close()

	client := NewClient(&Config{ApiToken: apiToken}, baseURL(mockServer.URL))
	err := client.GetDetailed(
		nil,
		&DetailedRequestParameters{
			StandardRequestParameters: &StandardRequestParameters{
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
