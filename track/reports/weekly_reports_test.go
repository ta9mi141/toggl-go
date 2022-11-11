package reports

import (
	"context"
	"errors"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/ta9mi141/toggl-go/track"
	"github.com/ta9mi141/toggl-go/track/internal"
)

func TestSearchWeeklyReport(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			weeklyReport *WeeklyReport
			err          error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/weekly_reports/search_weekly_report_200_ok.json",
			},
			out: struct {
				weeklyReport *WeeklyReport
				err          error
			}{
				weeklyReport: &WeeklyReport{
					{
						UserID:    track.Ptr(1234567),
						ProjectID: track.Ptr(123456789),
						Seconds: []*int{
							track.Ptr(0),
							track.Ptr(1234),
							track.Ptr(0),
							track.Ptr(56),
							track.Ptr(0),
							track.Ptr(0),
							track.Ptr(0),
						},
					},
					{
						UserID:    track.Ptr(1234567),
						ProjectID: track.Ptr(234567890),
						Seconds: []*int{
							track.Ptr(0),
							track.Ptr(0),
							track.Ptr(0),
							track.Ptr(7890),
							track.Ptr(0),
							track.Ptr(0),
							track.Ptr(0),
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "400 Bad Request",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusBadRequest,
				testdataFile: "testdata/weekly_reports/search_weekly_report_400_bad_request.json",
			},
			out: struct {
				weeklyReport *WeeklyReport
				err          error
			}{
				weeklyReport: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"At least one parameter must be set\"\n",
					Header: http.Header{
						"Content-Length": []string{"37"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "401 Unauthorized",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusUnauthorized,
				testdataFile: "testdata/weekly_reports/search_weekly_report_401_unauthorized",
			},
			out: struct {
				weeklyReport *WeeklyReport
				err          error
			}{
				weeklyReport: nil,
				err: &internal.ErrorResponse{
					StatusCode: 401,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "403 Forbidden",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusForbidden,
				testdataFile: "testdata/weekly_reports/search_weekly_report_403_forbidden.txt",
			},
			out: struct {
				weeklyReport *WeeklyReport
				err          error
			}{
				weeklyReport: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Content-Type":   []string{"text/plain; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceID := 1234567
			apiSpecificPath := path.Join(reportsPath, strconv.Itoa(workspaceID), "weekly/time_entries")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			weeklyReport, err := apiClient.SearchWeeklyReport(context.Background(), workspaceID, &SearchWeeklyReportRequestBody{})

			if !reflect.DeepEqual(weeklyReport, tt.out.weeklyReport) {
				internal.Errorf(t, weeklyReport, tt.out.weeklyReport)
			}

			errorResp := new(internal.ErrorResponse)
			if errors.As(err, &errorResp) {
				if !reflect.DeepEqual(errorResp, tt.out.err) {
					internal.Errorf(t, errorResp, tt.out.err)
				}
			} else {
				if !reflect.DeepEqual(err, tt.out.err) {
					internal.Errorf(t, err, tt.out.err)
				}
			}
		})
	}
}

func TestSearchWeeklyReportRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *SearchWeeklyReportRequestBody
		out  string
	}{
		{
			name: "string",
			in: &SearchWeeklyReportRequestBody{
				StartDate: track.Ptr("2006-01-02"),
			},
			out: "{\"start_date\":\"2006-01-02\"}",
		},
		{
			name: "string and bool",
			in: &SearchWeeklyReportRequestBody{
				Billable:  track.Ptr(true),
				StartDate: track.Ptr("2006-01-02"),
			},
			out: "{\"billable\":true,\"start_date\":\"2006-01-02\"}",
		},
		{
			name: "string, bool, and array of integer",
			in: &SearchWeeklyReportRequestBody{
				Billable:   track.Ptr(true),
				ProjectIDs: []*int{track.Ptr(123456789), track.Ptr(234567890)},
				StartDate:  track.Ptr("2006-01-02"),
			},
			out: "{\"billable\":true,\"project_ids\":[123456789,234567890],\"start_date\":\"2006-01-02\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			workspaceID := 1234567
			_, _ = apiClient.SearchWeeklyReport(context.Background(), workspaceID, tt.in)
		})
	}
}
