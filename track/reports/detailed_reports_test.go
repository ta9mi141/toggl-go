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

func TestSearchDetailedReport(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			detailedReport *DetailedReport
			err            error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/detailed_reports/search_detailed_report_200_ok.json",
			},
			out: struct {
				detailedReport *DetailedReport
				err            error
			}{
				detailedReport: &DetailedReport{
					{
						UserID:                track.Ptr(1234567),
						Username:              track.Ptr("Toggl"),
						ProjectID:             track.Ptr(123456789),
						TaskID:                nil,
						Billable:              track.Ptr(false),
						Description:           track.Ptr("Awesome Description"),
						TagIDs:                []*int{},
						BillableAmountInCents: nil,
						HourlyRateInCents:     nil,
						Currency:              track.Ptr("USD"),
						TimeEntries: []*timeEntry{
							&timeEntry{
								ID:      track.Ptr(1234567890),
								Seconds: track.Ptr(8040),
								Start:   track.Ptr(time.Date(2020, time.January, 2, 9, 59, 9, 0, time.FixedZone("", 0))),
								Stop:    track.Ptr(time.Date(2020, time.January, 2, 12, 13, 9, 0, time.FixedZone("", 0))),
								At:      track.Ptr(time.Date(2020, time.January, 2, 14, 30, 36, 0, time.FixedZone("", 0))),
							},
						},
						RowNumber: track.Ptr(1),
					},
					{
						UserID:                track.Ptr(1234567),
						Username:              track.Ptr("Toggl"),
						ProjectID:             track.Ptr(234567890),
						TaskID:                nil,
						Billable:              track.Ptr(false),
						Description:           track.Ptr("NewDescription"),
						TagIDs:                []*int{},
						BillableAmountInCents: nil,
						HourlyRateInCents:     nil,
						Currency:              track.Ptr("USD"),
						TimeEntries: []*timeEntry{
							&timeEntry{
								ID:      track.Ptr(2345678901),
								Seconds: track.Ptr(30),
								Start:   track.Ptr(time.Date(2020, time.January, 2, 13, 17, 57, 0, time.FixedZone("", 0))),
								Stop:    track.Ptr(time.Date(2020, time.January, 2, 13, 18, 27, 0, time.FixedZone("", 0))),
								At:      track.Ptr(time.Date(2020, time.January, 2, 14, 18, 38, 0, time.FixedZone("", 0))),
							},
						},
						RowNumber: track.Ptr(2),
					},
					{
						UserID:                track.Ptr(1234567),
						Username:              track.Ptr("Toggl"),
						ProjectID:             track.Ptr(234567890),
						TaskID:                nil,
						Billable:              track.Ptr(false),
						Description:           track.Ptr("NewDescription"),
						TagIDs:                []*int{},
						BillableAmountInCents: nil,
						HourlyRateInCents:     nil,
						Currency:              track.Ptr("USD"),
						TimeEntries: []*timeEntry{
							&timeEntry{
								ID:      track.Ptr(3456789012),
								Seconds: track.Ptr(8),
								Start:   track.Ptr(time.Date(2020, time.January, 2, 13, 24, 49, 0, time.FixedZone("", 0))),
								Stop:    track.Ptr(time.Date(2020, time.January, 2, 13, 24, 57, 0, time.FixedZone("", 0))),
								At:      track.Ptr(time.Date(2020, time.January, 2, 14, 25, 7, 0, time.FixedZone("", 0))),
							},
						},
						RowNumber: track.Ptr(3),
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
				testdataFile: "testdata/detailed_reports/search_detailed_report_400_bad_request.json",
			},
			out: struct {
				detailedReport *DetailedReport
				err            error
			}{
				detailedReport: nil,
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
				testdataFile: "testdata/detailed_reports/search_detailed_report_401_unauthorized",
			},
			out: struct {
				detailedReport *DetailedReport
				err            error
			}{
				detailedReport: nil,
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
				testdataFile: "testdata/detailed_reports/search_detailed_report_403_forbidden.txt",
			},
			out: struct {
				detailedReport *DetailedReport
				err            error
			}{
				detailedReport: nil,
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
			apiSpecificPath := path.Join(reportsPath, strconv.Itoa(workspaceID), "search/time_entries")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			detailedReport, err := apiClient.SearchDetailedReport(context.Background(), workspaceID, &SearchDetailedReportRequestBody{})

			if !reflect.DeepEqual(detailedReport, tt.out.detailedReport) {
				internal.Errorf(t, detailedReport, tt.out.detailedReport)
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

func TestSearchDetailedReportRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *SearchDetailedReportRequestBody
		out  string
	}{
		{
			name: "string",
			in: &SearchDetailedReportRequestBody{
				StartDate: track.Ptr("2006-01-02"),
			},
			out: "{\"start_date\":\"2006-01-02\"}",
		},
		{
			name: "string and bool",
			in: &SearchDetailedReportRequestBody{
				Billable:  track.Ptr(true),
				StartDate: track.Ptr("2006-01-02"),
			},
			out: "{\"billable\":true,\"start_date\":\"2006-01-02\"}",
		},
		{
			name: "string, bool, and array of integer",
			in: &SearchDetailedReportRequestBody{
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
			_, _ = apiClient.SearchDetailedReport(context.Background(), workspaceID, tt.in)
		})
	}
}
