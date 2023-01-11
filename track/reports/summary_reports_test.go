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

func TestSearchSummaryReport(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			summaryReport *SummaryReport
			err           error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/summary_reports/search_summary_report_200_ok.json",
			},
			out: struct {
				summaryReport *SummaryReport
				err           error
			}{
				summaryReport: &SummaryReport{
					Groups: []*group{
						&group{
							ID: track.Ptr(123456789),
							SubGroups: []*subGroup{
								&subGroup{
									ID:      nil,
									Title:   track.Ptr("Description 1"),
									Seconds: track.Ptr(123),
								},
							},
						},
						&group{
							ID: track.Ptr(234567891),
							SubGroups: []*subGroup{
								&subGroup{
									ID:      nil,
									Title:   track.Ptr("Description 2"),
									Seconds: track.Ptr(456),
								},
							},
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
				testdataFile: "testdata/summary_reports/search_summary_report_400_bad_request.json",
			},
			out: struct {
				summaryReport *SummaryReport
				err           error
			}{
				summaryReport: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"Maximum allowed date range is 365 days\"\n",
					Header: http.Header{
						"Content-Length": []string{"41"},
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
				testdataFile: "testdata/summary_reports/search_summary_report_401_unauthorized",
			},
			out: struct {
				summaryReport *SummaryReport
				err           error
			}{
				summaryReport: nil,
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
				testdataFile: "testdata/summary_reports/search_summary_report_403_forbidden.txt",
			},
			out: struct {
				summaryReport *SummaryReport
				err           error
			}{
				summaryReport: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "Incorrect username and/or password\n",
					Header: http.Header{
						"Content-Length": []string{"35"},
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
			apiSpecificPath := path.Join(reportsPath, strconv.Itoa(workspaceID), "summary/time_entries")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			summaryReport, err := apiClient.SearchSummaryReport(context.Background(), workspaceID, &SearchSummaryReportRequestBody{})

			if !reflect.DeepEqual(summaryReport, tt.out.summaryReport) {
				internal.Errorf(t, summaryReport, tt.out.summaryReport)
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

func TestSearchSummaryReportRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *SearchSummaryReportRequestBody
		out  string
	}{
		{
			name: "string",
			in: &SearchSummaryReportRequestBody{
				StartDate: track.Ptr("2006-01-02"),
			},
			out: "{\"start_date\":\"2006-01-02\"}",
		},
		{
			name: "string and bool",
			in: &SearchSummaryReportRequestBody{
				Billable:  track.Ptr(true),
				StartDate: track.Ptr("2006-01-02"),
			},
			out: "{\"billable\":true,\"start_date\":\"2006-01-02\"}",
		},
		{
			name: "string, bool, and array of integer",
			in: &SearchSummaryReportRequestBody{
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
			_, _ = apiClient.SearchSummaryReport(context.Background(), workspaceID, tt.in)
		})
	}
}

func TestLoadProjectSummary(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			projectSummary *ProjectSummary
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
				testdataFile: "testdata/summary_reports/load_project_summary_200_ok.json",
			},
			out: struct {
				projectSummary *ProjectSummary
				err            error
			}{
				projectSummary: &ProjectSummary{
					Seconds:    track.Ptr(123),
					Resolution: nil,
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
				testdataFile: "testdata/summary_reports/load_project_summary_400_bad_request.json",
			},
			out: struct {
				projectSummary *ProjectSummary
				err            error
			}{
				projectSummary: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"end_date should be within 2006-01-01 to 2030-01-01\"\n",
					Header: http.Header{
						"Content-Length": []string{"53"},
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
				testdataFile: "testdata/summary_reports/load_project_summary_401_unauthorized",
			},
			out: struct {
				projectSummary *ProjectSummary
				err            error
			}{
				projectSummary: nil,
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
				testdataFile: "testdata/summary_reports/load_project_summary_403_forbidden.txt",
			},
			out: struct {
				projectSummary *ProjectSummary
				err            error
			}{
				projectSummary: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "Incorrect username and/or password\n",
					Header: http.Header{
						"Content-Length": []string{"35"},
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
			projectID := 12345678
			apiSpecificPath := path.Join(reportsPath, strconv.Itoa(workspaceID), "projects", strconv.Itoa(projectID), "summary")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			projectSummary, err := apiClient.LoadProjectSummary(context.Background(), workspaceID, projectID, &LoadProjectSummaryRequestBody{})

			if !reflect.DeepEqual(projectSummary, tt.out.projectSummary) {
				internal.Errorf(t, projectSummary, tt.out.projectSummary)
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

func TestLoadProjectSummaryRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *LoadProjectSummaryRequestBody
		out  string
	}{
		{
			name: "string",
			in: &LoadProjectSummaryRequestBody{
				EndDate:   track.Ptr("2007-01-02"),
				StartDate: track.Ptr("2006-01-02"),
			},
			out: "{\"end_date\":\"2007-01-02\",\"start_date\":\"2006-01-02\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			workspaceID := 1234567
			projectID := 12345678
			_, _ = apiClient.LoadProjectSummary(context.Background(), workspaceID, projectID, tt.in)
		})
	}
}
