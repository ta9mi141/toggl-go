package reports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	expectedAPIToken := "api_token"
	client := NewClient(expectedAPIToken)

	if client.apiToken != expectedAPIToken {
		t.Error("client.apiToken = " + client.apiToken + ", [Expected: " + expectedAPIToken + "]")
	}

	if client.url.String() != defaultBaseURL {
		t.Error("client.url.String() = " + client.url.String() + ", [Expected: " + defaultBaseURL + "]")
	}

	expectedContentType := "application/json"
	if client.header.Get("Content-type") != expectedContentType {
		t.Error(`client.header.Get("Content-type") = ` + client.header.Get("Content-type") + ", [Expected: " + expectedContentType + "]")
	}
}

func setupMockServer_200_Ok(t *testing.T, testdataFilePath string) (*httptest.Server, []byte) {
	testdata, err := ioutil.ReadFile(testdataFilePath)
	if err != nil {
		t.Error(err.Error())
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(testdata))
	}))

	return mockServer, testdata
}

func setupMockServer_401_Unauthorized(t *testing.T) (*httptest.Server, []byte) {
	errorTestData, err := ioutil.ReadFile("testdata/401_unauthorized.json")
	if err != nil {
		t.Error(err.Error())
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, string(errorTestData))
	}))

	return mockServer, errorTestData
}

const (
	apiToken    string = "api_token"
	userAgent   string = "user_agent"
	workSpaceId string = "workspace_id"
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
	client := NewClient(apiToken, baseURL(mockServer.URL))
	err := client.GetDetailed(&DetailedRequestParameters{
		StandardRequestParameters: &StandardRequestParameters{
			UserAgent:   userAgent,
			WorkSpaceId: workSpaceId,
		},
	}, actualDetailedReport)
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

	client := NewClient(apiToken, baseURL(mockServer.URL))
	actualReportsError := client.GetDetailed(&DetailedRequestParameters{
		StandardRequestParameters: &StandardRequestParameters{
			UserAgent:   userAgent,
			WorkSpaceId: workSpaceId,
		},
	}, new(detailedReport))
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
