package toggl

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

func TestGetTimeEntries(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			timeEntries []*TimeEntry
			err         error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/time_entries/get_time_entries_200_ok.json",
			},
			out: struct {
				timeEntries []*TimeEntry
				err         error
			}{
				timeEntries: []*TimeEntry{
					{
						ID:              track.Ptr(1234567890),
						WorkspaceID:     track.Ptr(1234567),
						ProjectID:       track.Ptr(123456789),
						TaskID:          nil,
						Billable:        track.Ptr(false),
						Start:           track.Ptr(time.Date(2012, time.March, 4, 5, 6, 20, 0, time.Local)),
						Stop:            track.Ptr(time.Date(2012, time.March, 4, 5, 6, 23, 0, time.UTC)),
						Duration:        track.Ptr(3),
						Description:     track.Ptr("test time entry"),
						Tags:            []*string{track.Ptr("billed")},
						TagIDs:          []*int{track.Ptr(1234567)},
						Duronly:         track.Ptr(false),
						At:              track.Ptr(time.Date(2022, time.March, 4, 5, 6, 7, 0, time.Local)),
						ServerDeletedAt: nil,
						UserID:          track.Ptr(9876543),
						UID:             track.Ptr(9876543),
						WID:             track.Ptr(1234567),
						PID:             track.Ptr(123456789),
					},
					{
						ID:              track.Ptr(2345678901),
						WorkspaceID:     track.Ptr(1234567),
						ProjectID:       track.Ptr(234567890),
						TaskID:          nil,
						Billable:        track.Ptr(false),
						Start:           track.Ptr(time.Date(2022, time.January, 2, 3, 47, 41, 0, time.Local)),
						Stop:            track.Ptr(time.Date(2022, time.January, 2, 3, 48, 1, 0, time.UTC)),
						Duration:        track.Ptr(20),
						Description:     track.Ptr("test time entry"),
						Tags:            []*string{track.Ptr("billed"), track.Ptr("toggl-go")},
						TagIDs:          []*int{track.Ptr(1234567), track.Ptr(2345678)},
						Duronly:         track.Ptr(false),
						At:              track.Ptr(time.Date(2022, time.March, 4, 5, 6, 7, 0, time.Local)),
						ServerDeletedAt: nil,
						UserID:          track.Ptr(9876543),
						UID:             track.Ptr(9876543),
						WID:             track.Ptr(1234567),
						PID:             track.Ptr(234567890),
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
				testdataFile: "testdata/time_entries/get_time_entries_400_bad_request.json",
			},
			out: struct {
				timeEntries []*TimeEntry
				err         error
			}{
				timeEntries: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"Since is expected to be an unix timestamp, integer value\"\n",
					Header: http.Header{
						"Content-Length": []string{"59"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
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
				testdataFile: "testdata/time_entries/get_time_entries_403_forbidden",
			},
			out: struct {
				timeEntries []*TimeEntry
				err         error
			}{
				timeEntries: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusInternalServerError,
				testdataFile: "testdata/time_entries/get_time_entries_500_internal_server_error",
			},
			out: struct {
				timeEntries []*TimeEntry
				err         error
			}{
				timeEntries: nil,
				err: &internal.ErrorResponse{
					StatusCode: 500,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiSpecificPath := path.Join(mePath, "time_entries")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			timeEntries, err := client.GetTimeEntries(context.Background(), nil)

			if !reflect.DeepEqual(timeEntries, tt.out.timeEntries) {
				internal.Errorf(t, timeEntries, tt.out.timeEntries)
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

func TestGetTimeEntriesQuery(t *testing.T) {
	tests := []struct {
		name string
		in   *GetTimeEntriesQuery
		out  string
	}{
		{
			name: "GetTimeEntriesQuery is nil",
			in:   nil,
			out:  "",
		},
		{
			name: "before=2022-07-01",
			in:   &GetTimeEntriesQuery{Before: track.Ptr("2022-07-01")},
			out:  "before=2022-07-01",
		},
		{
			name: "since=1656687597",
			in:   &GetTimeEntriesQuery{Since: track.Ptr(1656687597)},
			out:  "since=1656687597",
		},
		{
			name: "end_date=2022-07-07&start_date=2022-07-01",
			in: &GetTimeEntriesQuery{
				StartDate: track.Ptr("2022-07-01"),
				EndDate:   track.Ptr("2022-07-07"),
			},
			out: "end_date=2022-07-07&start_date=2022-07-01",
		},
		{
			name: "GetTimeEntriesQuery is empty",
			in:   &GetTimeEntriesQuery{},
			out:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertQuery(t, tt.out)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			_, _ = client.GetTimeEntries(context.Background(), tt.in)
		})
	}
}

func TestGetCurrentTimeEntry(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			timeEntry *TimeEntry
			err       error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/time_entries/get_current_time_entry_200_ok.json",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: &TimeEntry{
					ID:              track.Ptr(1234567890),
					WorkspaceID:     track.Ptr(1234567),
					ProjectID:       track.Ptr(123456789),
					TaskID:          nil,
					Billable:        track.Ptr(false),
					Start:           track.Ptr(time.Date(2020, time.January, 23, 4, 56, 31, 0, time.Local)),
					Stop:            nil,
					Duration:        track.Ptr(-1579722991),
					Description:     track.Ptr("running time entry"),
					Tags:            []*string{track.Ptr("toggl-go")},
					TagIDs:          []*int{track.Ptr(1234567)},
					Duronly:         track.Ptr(false),
					At:              track.Ptr(time.Date(2020, time.January, 23, 4, 56, 34, 0, time.Local)),
					ServerDeletedAt: nil,
					UserID:          track.Ptr(1234567),
					UID:             track.Ptr(1234567),
					WID:             track.Ptr(1234567),
					PID:             track.Ptr(123456789),
				},
				err: nil,
			},
		},
		{
			name: "200 OK (null)",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/time_entries/get_current_time_entry_200_ok_null.json",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       nil,
			},
		},
		{
			name: "403 Forbidden",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusForbidden,
				testdataFile: "testdata/time_entries/get_current_time_entry_403_forbidden",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusInternalServerError,
				testdataFile: "testdata/time_entries/get_current_time_entry_500_internal_server_error",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 500,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiSpecificPath := path.Join(mePath, "time_entries/current")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			timeEntry, err := client.GetCurrentTimeEntry(context.Background())

			if !reflect.DeepEqual(timeEntry, tt.out.timeEntry) {
				internal.Errorf(t, timeEntry, tt.out.timeEntry)
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

func TestCreateTimeEntry(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			timeEntry *TimeEntry
			err       error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/time_entries/create_time_entry_200_ok.json",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: &TimeEntry{
					ID:              track.Ptr(1234567890),
					WorkspaceID:     track.Ptr(1234567),
					ProjectID:       track.Ptr(123456789),
					TaskID:          nil,
					Billable:        track.Ptr(false),
					Start:           track.Ptr(time.Date(2021, time.July, 6, 5, 4, 3, 0, time.UTC)),
					Stop:            track.Ptr(time.Date(2021, time.July, 6, 5, 9, 3, 0, time.UTC)),
					Duration:        track.Ptr(300),
					Description:     track.Ptr("created manually"),
					Tags:            nil,
					TagIDs:          nil,
					Duronly:         track.Ptr(false),
					At:              track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.Local)),
					ServerDeletedAt: nil,
					UserID:          track.Ptr(1234567),
					UID:             track.Ptr(1234567),
					WID:             track.Ptr(1234567),
					PID:             track.Ptr(123456789),
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
				testdataFile: "testdata/time_entries/create_time_entry_400_bad_request.json",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"JSON is not valid\"\n",
					Header: http.Header{
						"Content-Length": []string{"20"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
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
				testdataFile: "testdata/time_entries/create_time_entry_403_forbidden",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusInternalServerError,
				testdataFile: "testdata/time_entries/create_time_entry_500_internal_server_error",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 500,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceID := 1234567
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "time_entries")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			timeEntry, err := client.CreateTimeEntry(context.Background(), workspaceID, &CreateTimeEntryRequestBody{})

			if !reflect.DeepEqual(timeEntry, tt.out.timeEntry) {
				internal.Errorf(t, timeEntry, tt.out.timeEntry)
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

func TestCreateTimeEntryRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *CreateTimeEntryRequestBody
		out  string
	}{
		{
			name: "int, string, and time",
			in: &CreateTimeEntryRequestBody{
				WorkspaceID: track.Ptr(1234567),
				Start:       track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    track.Ptr(300),
				CreatedWith: track.Ptr("toggl-go"),
				Description: track.Ptr("created manually"),
				ProjectID:   track.Ptr(123456789),
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"created manually\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and bool",
			in: &CreateTimeEntryRequestBody{
				WorkspaceID: track.Ptr(1234567),
				Start:       track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    track.Ptr(300),
				CreatedWith: track.Ptr("toggl-go"),
				Description: track.Ptr("created manually"),
				ProjectID:   track.Ptr(123456789),
				Billable:    track.Ptr(false),
			},
			out: "{\"billable\":false,\"created_with\":\"toggl-go\",\"description\":\"created manually\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and slice of string",
			in: &CreateTimeEntryRequestBody{
				WorkspaceID: track.Ptr(1234567),
				Start:       track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    track.Ptr(300),
				CreatedWith: track.Ptr("toggl-go"),
				Description: track.Ptr("created manually"),
				ProjectID:   track.Ptr(123456789),
				Tags:        []*string{track.Ptr("tag1"), track.Ptr("tag2")},
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"created manually\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"tags\":[\"tag1\",\"tag2\"],\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and slice of int",
			in: &CreateTimeEntryRequestBody{
				WorkspaceID: track.Ptr(1234567),
				Start:       track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    track.Ptr(300),
				CreatedWith: track.Ptr("toggl-go"),
				Description: track.Ptr("created manually"),
				ProjectID:   track.Ptr(123456789),
				TagIDs:      []*int{track.Ptr(1234567), track.Ptr(9876543)},
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"created manually\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"tag_ids\":[1234567,9876543],\"workspace_id\":1234567}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			workspaceID := 1234567
			_, _ = client.CreateTimeEntry(context.Background(), workspaceID, tt.in)
		})
	}
}

func TestUpdateTimeEntry(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			timeEntry *TimeEntry
			err       error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/time_entries/update_time_entry_200_ok.json",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: &TimeEntry{
					ID:              track.Ptr(1234567890),
					WorkspaceID:     track.Ptr(1324567),
					ProjectID:       track.Ptr(234567890),
					TaskID:          nil,
					Billable:        track.Ptr(false),
					Start:           track.Ptr(time.Date(2022, time.July, 6, 5, 43, 31, 0, time.UTC)),
					Stop:            track.Ptr(time.Date(2022, time.July, 6, 5, 44, 37, 0, time.UTC)),
					Duration:        track.Ptr(66),
					Description:     track.Ptr("updated time entry"),
					Tags:            []*string{track.Ptr("toggl-go")},
					TagIDs:          []*int{track.Ptr(3456789)},
					Duronly:         track.Ptr(false),
					At:              track.Ptr(time.Date(2022, time.July, 7, 12, 34, 56, 0, time.Local)),
					ServerDeletedAt: nil,
					UserID:          track.Ptr(1234567),
					UID:             track.Ptr(1234567),
					WID:             track.Ptr(1324567),
					PID:             track.Ptr(234567890),
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
				testdataFile: "testdata/time_entries/update_time_entry_400_bad_request.json",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"JSON is not valid\"\n",
					Header: http.Header{
						"Content-Length": []string{"20"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
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
				testdataFile: "testdata/time_entries/update_time_entry_403_forbidden",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "404 Not Found",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusNotFound,
				testdataFile: "testdata/time_entries/update_time_entry_404_not_found.json",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 404,
					Message:    "\"Time entry not found\"\n",
					Header: http.Header{
						"Content-Length": []string{"23"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusInternalServerError,
				testdataFile: "testdata/time_entries/update_time_entry_500_internal_server_error",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &internal.ErrorResponse{
					StatusCode: 500,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceID := 1234567
			timeEntryID := 1234567890
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "time_entries", strconv.Itoa(timeEntryID))
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			timeEntry, err := client.UpdateTimeEntry(context.Background(), workspaceID, timeEntryID, &UpdateTimeEntryRequestBody{})

			if !reflect.DeepEqual(timeEntry, tt.out.timeEntry) {
				internal.Errorf(t, timeEntry, tt.out.timeEntry)
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

func TestUpdateTimeEntryRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *UpdateTimeEntryRequestBody
		out  string
	}{
		{
			name: "int, string, and time",
			in: &UpdateTimeEntryRequestBody{
				WorkspaceID: track.Ptr(1234567),
				Start:       track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    track.Ptr(300),
				CreatedWith: track.Ptr("toggl-go"),
				Description: track.Ptr("updated time entry"),
				ProjectID:   track.Ptr(123456789),
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"updated time entry\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and bool",
			in: &UpdateTimeEntryRequestBody{
				WorkspaceID: track.Ptr(1234567),
				Start:       track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    track.Ptr(300),
				CreatedWith: track.Ptr("toggl-go"),
				Description: track.Ptr("updated time entry"),
				ProjectID:   track.Ptr(123456789),
				Billable:    track.Ptr(false),
			},
			out: "{\"billable\":false,\"created_with\":\"toggl-go\",\"description\":\"updated time entry\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and slice of string",
			in: &UpdateTimeEntryRequestBody{
				WorkspaceID: track.Ptr(1234567),
				Start:       track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    track.Ptr(300),
				CreatedWith: track.Ptr("toggl-go"),
				Description: track.Ptr("updated time entry"),
				ProjectID:   track.Ptr(123456789),
				Tags:        []*string{track.Ptr("tag1"), track.Ptr("tag2")},
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"updated time entry\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"tags\":[\"tag1\",\"tag2\"],\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and slice of int",
			in: &UpdateTimeEntryRequestBody{
				WorkspaceID: track.Ptr(1234567),
				Start:       track.Ptr(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    track.Ptr(300),
				CreatedWith: track.Ptr("toggl-go"),
				Description: track.Ptr("updated time entry"),
				ProjectID:   track.Ptr(123456789),
				TagIDs:      []*int{track.Ptr(1234567), track.Ptr(9876543)},
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"updated time entry\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"tag_ids\":[1234567,9876543],\"workspace_id\":1234567}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			workspaceID := 1234567
			timeEntryID := 1234567890
			_, _ = client.UpdateTimeEntry(context.Background(), workspaceID, timeEntryID, tt.in)
		})
	}
}

func TestDeleteTimeEntry(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			err error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/time_entries/delete_time_entry_200_ok.json",
			},
			out: struct {
				err error
			}{
				err: nil,
			},
		},
		{
			name: "403 Forbidden",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusForbidden,
				testdataFile: "testdata/time_entries/delete_time_entry_403_forbidden",
			},
			out: struct {
				err error
			}{
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "404 Not Found",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusNotFound,
				testdataFile: "testdata/time_entries/delete_time_entry_404_not_found.json",
			},
			out: struct {
				err error
			}{
				err: &internal.ErrorResponse{
					StatusCode: 404,
					Message:    "\"Time entry not found\"\n",
					Header: http.Header{
						"Content-Length": []string{"23"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusInternalServerError,
				testdataFile: "testdata/time_entries/delete_time_entry_500_internal_server_error",
			},
			out: struct {
				err error
			}{
				err: &internal.ErrorResponse{
					StatusCode: 500,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceID := 1234567
			timeEntryID := 1234567890
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "time_entries", strconv.Itoa(timeEntryID))
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			err := client.DeleteTimeEntry(context.Background(), workspaceID, timeEntryID)

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
