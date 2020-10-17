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

	"github.com/ta9mi1shi1/toggl-go/reports"
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

func TestGetWeekly(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               context.Context
		out              error
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/weekly.json",
			in:               context.Background(),
			out:              nil,
		},
		{
			name:             "401 Unauthorized",
			httpStatus:       http.StatusUnauthorized,
			testdataFilePath: "testdata/401_unauthorized.json",
			in:               context.Background(),
			out: &reports.ReportsError{
				Err: struct {
					Message string `json:"message"`
					Tip     string `json:"tip"`
					Code    int    `json:"code"`
				}{
					Message: "api token missing",
					Tip:     "You can find your API Token in your profile at https://www.toggl.com",
					Code:    http.StatusUnauthorized,
				},
			},
		},
		{
			name:             "429 Too Many Requests",
			httpStatus:       http.StatusTooManyRequests,
			testdataFilePath: "testdata/429_too_many_requests.html",
			in:               context.Background(),
			out: &reports.ReportsError{
				Err: struct {
					Message string `json:"message"`
					Tip     string `json:"tip"`
					Code    int    `json:"code"`
				}{
					Message: "Too Many Requests",
					Tip:     "Add delay between requests",
					Code:    http.StatusTooManyRequests,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/weekly.json",
			in:               nil,
			out:              reports.ErrContextNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer, testdata := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			actualWeeklyReport := new(weeklyReport)
			client := reports.NewClient(apiToken, baseURL(mockServer.URL))
			err := client.GetWeekly(
				c.in,
				&reports.WeeklyRequestParameters{
					StandardRequestParameters: &reports.StandardRequestParameters{
						UserAgent:   userAgent,
						WorkspaceId: workspaceId,
					},
				},
				actualWeeklyReport,
			)

			if err == nil {
				expectedWeeklyReport := new(weeklyReport)
				if err := json.Unmarshal(testdata, expectedWeeklyReport); err != nil {
					t.Error(err.Error())
				}
				if !reflect.DeepEqual(actualWeeklyReport, expectedWeeklyReport) {
					t.Errorf("\nwant: %+v\ngot : %+v\n", expectedWeeklyReport, actualWeeklyReport)
				}
			} else {
				if !reflect.DeepEqual(actualWeeklyReport, &weeklyReport{}) {
					t.Errorf("\nwant: %+v\ngot : %+v\n", &weeklyReport{}, actualWeeklyReport)
				}
			}

			var reportsError reports.Error
			if errors.As(err, &reportsError) {
				if !reflect.DeepEqual(reportsError, c.out) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out, reportsError)
				}
			} else {
				if !errors.Is(err, c.out) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out, err)
				}
			}
		})
	}
}

func TestGetWeeklyEncodeRequestParameters(t *testing.T) {
	expectedQueryString := url.Values{
		"user_agent":   []string{userAgent},
		"workspace_id": []string{workspaceId},
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
				WorkspaceId: workspaceId,
			},
			Grouping: "users",
		},
		new(summaryReport),
	)
}
