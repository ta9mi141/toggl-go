package track

import (
	"context"
	"errors"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"testing"
	"time"
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
						ID:              Int(1234567890),
						WorkspaceID:     Int(1234567),
						ProjectID:       Int(123456789),
						TaskID:          nil,
						Billable:        Bool(false),
						Start:           Time(time.Date(2012, time.March, 4, 5, 6, 20, 0, time.FixedZone("", 0))),
						Stop:            Time(time.Date(2012, time.March, 4, 5, 6, 23, 0, time.UTC)),
						Duration:        Int(3),
						Description:     String("test time entry"),
						Tags:            []*string{String("billed")},
						TagIDs:          []*int{Int(1234567)},
						Duronly:         Bool(false),
						At:              Time(time.Date(2022, time.March, 4, 5, 6, 7, 0, time.FixedZone("", 0))),
						ServerDeletedAt: nil,
						UserID:          Int(9876543),
						UID:             Int(9876543),
						WID:             Int(1234567),
						PID:             Int(123456789),
					},
					{
						ID:              Int(2345678901),
						WorkspaceID:     Int(1234567),
						ProjectID:       Int(234567890),
						TaskID:          nil,
						Billable:        Bool(false),
						Start:           Time(time.Date(2022, time.January, 2, 3, 47, 41, 0, time.FixedZone("", 0))),
						Stop:            Time(time.Date(2022, time.January, 2, 3, 48, 1, 0, time.UTC)),
						Duration:        Int(20),
						Description:     String("test time entry"),
						Tags:            []*string{String("billed"), String("toggl-go")},
						TagIDs:          []*int{Int(1234567), Int(2345678)},
						Duronly:         Bool(false),
						At:              Time(time.Date(2022, time.March, 4, 5, 6, 7, 0, time.FixedZone("", 0))),
						ServerDeletedAt: nil,
						UserID:          Int(9876543),
						UID:             Int(9876543),
						WID:             Int(1234567),
						PID:             Int(234567890),
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
				err: &errorResponse{
					statusCode: 400,
					message:    "\"Since is expected to be an unix timestamp, integer value\"\n",
					header: http.Header{
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
				err: &errorResponse{
					statusCode: 403,
					message:    "",
					header: http.Header{
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
				err: &errorResponse{
					statusCode: 500,
					message:    "",
					header: http.Header{
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
			mockServer := newMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			timeEntries, err := client.GetTimeEntries(context.Background(), nil)

			if !reflect.DeepEqual(timeEntries, tt.out.timeEntries) {
				errorf(t, timeEntries, tt.out.timeEntries)
			}

			errorResp := new(errorResponse)
			if errors.As(err, &errorResp) {
				if !reflect.DeepEqual(errorResp, tt.out.err) {
					errorf(t, errorResp, tt.out.err)
				}
			} else {
				if !reflect.DeepEqual(err, tt.out.err) {
					errorf(t, err, tt.out.err)
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
			in:   &GetTimeEntriesQuery{Before: String("2022-07-01")},
			out:  "before=2022-07-01",
		},
		{
			name: "since=1656687597",
			in:   &GetTimeEntriesQuery{Since: Int(1656687597)},
			out:  "since=1656687597",
		},
		{
			name: "end_date=2022-07-07&start_date=2022-07-01",
			in: &GetTimeEntriesQuery{
				StartDate: String("2022-07-01"),
				EndDate:   String("2022-07-07"),
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
			mockServer := newMockServerToAssertQuery(t, tt.out)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
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
					ID:              Int(1234567890),
					WorkspaceID:     Int(1234567),
					ProjectID:       Int(123456789),
					TaskID:          nil,
					Billable:        Bool(false),
					Start:           Time(time.Date(2020, time.January, 23, 4, 56, 31, 0, time.FixedZone("", 0))),
					Stop:            nil,
					Duration:        Int(-1579722991),
					Description:     String("running time entry"),
					Tags:            []*string{String("toggl-go")},
					TagIDs:          []*int{Int(1234567)},
					Duronly:         Bool(false),
					At:              Time(time.Date(2020, time.January, 23, 4, 56, 34, 0, time.FixedZone("", 0))),
					ServerDeletedAt: nil,
					UserID:          Int(1234567),
					UID:             Int(1234567),
					WID:             Int(1234567),
					PID:             Int(123456789),
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
				err: &errorResponse{
					statusCode: 403,
					message:    "",
					header: http.Header{
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
				err: &errorResponse{
					statusCode: 500,
					message:    "",
					header: http.Header{
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
			mockServer := newMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			timeEntry, err := client.GetCurrentTimeEntry(context.Background())

			if !reflect.DeepEqual(timeEntry, tt.out.timeEntry) {
				errorf(t, timeEntry, tt.out.timeEntry)
			}

			errorResp := new(errorResponse)
			if errors.As(err, &errorResp) {
				if !reflect.DeepEqual(errorResp, tt.out.err) {
					errorf(t, errorResp, tt.out.err)
				}
			} else {
				if !reflect.DeepEqual(err, tt.out.err) {
					errorf(t, err, tt.out.err)
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
					ID:              Int(1234567890),
					WorkspaceID:     Int(1234567),
					ProjectID:       Int(123456789),
					TaskID:          nil,
					Billable:        Bool(false),
					Start:           Time(time.Date(2021, time.July, 6, 5, 4, 3, 0, time.UTC)),
					Stop:            Time(time.Date(2021, time.July, 6, 5, 9, 3, 0, time.UTC)),
					Duration:        Int(300),
					Description:     String("created manually"),
					Tags:            nil,
					TagIDs:          nil,
					Duronly:         Bool(false),
					At:              Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.FixedZone("", 0))),
					ServerDeletedAt: nil,
					UserID:          Int(1234567),
					UID:             Int(1234567),
					WID:             Int(1234567),
					PID:             Int(123456789),
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
				err: &errorResponse{
					statusCode: 400,
					message:    "\"JSON is not valid\"\n",
					header: http.Header{
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
				err: &errorResponse{
					statusCode: 403,
					message:    "",
					header: http.Header{
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
				err: &errorResponse{
					statusCode: 500,
					message:    "",
					header: http.Header{
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
			mockServer := newMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			timeEntry, err := client.CreateTimeEntry(context.Background(), workspaceID, &CreateTimeEntryRequestBody{})

			if !reflect.DeepEqual(timeEntry, tt.out.timeEntry) {
				errorf(t, timeEntry, tt.out.timeEntry)
			}

			errorResp := new(errorResponse)
			if errors.As(err, &errorResp) {
				if !reflect.DeepEqual(errorResp, tt.out.err) {
					errorf(t, errorResp, tt.out.err)
				}
			} else {
				if !reflect.DeepEqual(err, tt.out.err) {
					errorf(t, err, tt.out.err)
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
				WorkspaceID: Int(1234567),
				Start:       Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    Int(300),
				CreatedWith: String("toggl-go"),
				Description: String("created manually"),
				ProjectID:   Int(123456789),
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"created manually\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and bool",
			in: &CreateTimeEntryRequestBody{
				WorkspaceID: Int(1234567),
				Start:       Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    Int(300),
				CreatedWith: String("toggl-go"),
				Description: String("created manually"),
				ProjectID:   Int(123456789),
				Billable:    Bool(false),
			},
			out: "{\"billable\":false,\"created_with\":\"toggl-go\",\"description\":\"created manually\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and slice of string",
			in: &CreateTimeEntryRequestBody{
				WorkspaceID: Int(1234567),
				Start:       Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    Int(300),
				CreatedWith: String("toggl-go"),
				Description: String("created manually"),
				ProjectID:   Int(123456789),
				Tags:        []*string{String("tag1"), String("tag2")},
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"created manually\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"tags\":[\"tag1\",\"tag2\"],\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and slice of int",
			in: &CreateTimeEntryRequestBody{
				WorkspaceID: Int(1234567),
				Start:       Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    Int(300),
				CreatedWith: String("toggl-go"),
				Description: String("created manually"),
				ProjectID:   Int(123456789),
				TagIDs:      []*int{Int(1234567), Int(9876543)},
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"created manually\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"tag_ids\":[1234567,9876543],\"workspace_id\":1234567}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
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
					ID:              Int(1234567890),
					WorkspaceID:     Int(1324567),
					ProjectID:       Int(234567890),
					TaskID:          nil,
					Billable:        Bool(false),
					Start:           Time(time.Date(2022, time.July, 6, 5, 43, 31, 0, time.UTC)),
					Stop:            Time(time.Date(2022, time.July, 6, 5, 44, 37, 0, time.UTC)),
					Duration:        Int(66),
					Description:     String("updated time entry"),
					Tags:            []*string{String("toggl-go")},
					TagIDs:          []*int{Int(3456789)},
					Duronly:         Bool(false),
					At:              Time(time.Date(2022, time.July, 7, 12, 34, 56, 0, time.FixedZone("", 0))),
					ServerDeletedAt: nil,
					UserID:          Int(1234567),
					UID:             Int(1234567),
					WID:             Int(1324567),
					PID:             Int(234567890),
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
				err: &errorResponse{
					statusCode: 400,
					message:    "\"JSON is not valid\"\n",
					header: http.Header{
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
				err: &errorResponse{
					statusCode: 403,
					message:    "",
					header: http.Header{
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
				testdataFile: "testdata/time_entries/update_time_entry_500_internal_server_error",
			},
			out: struct {
				timeEntry *TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &errorResponse{
					statusCode: 500,
					message:    "",
					header: http.Header{
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
			mockServer := newMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			timeEntry, err := client.UpdateTimeEntry(context.Background(), workspaceID, timeEntryID, &UpdateTimeEntryRequestBody{})

			if !reflect.DeepEqual(timeEntry, tt.out.timeEntry) {
				errorf(t, timeEntry, tt.out.timeEntry)
			}

			errorResp := new(errorResponse)
			if errors.As(err, &errorResp) {
				if !reflect.DeepEqual(errorResp, tt.out.err) {
					errorf(t, errorResp, tt.out.err)
				}
			} else {
				if !reflect.DeepEqual(err, tt.out.err) {
					errorf(t, err, tt.out.err)
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
				WorkspaceID: Int(1234567),
				Start:       Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    Int(300),
				CreatedWith: String("toggl-go"),
				Description: String("updated time entry"),
				ProjectID:   Int(123456789),
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"updated time entry\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and bool",
			in: &UpdateTimeEntryRequestBody{
				WorkspaceID: Int(1234567),
				Start:       Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    Int(300),
				CreatedWith: String("toggl-go"),
				Description: String("updated time entry"),
				ProjectID:   Int(123456789),
				Billable:    Bool(false),
			},
			out: "{\"billable\":false,\"created_with\":\"toggl-go\",\"description\":\"updated time entry\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and slice of string",
			in: &UpdateTimeEntryRequestBody{
				WorkspaceID: Int(1234567),
				Start:       Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    Int(300),
				CreatedWith: String("toggl-go"),
				Description: String("updated time entry"),
				ProjectID:   Int(123456789),
				Tags:        []*string{String("tag1"), String("tag2")},
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"updated time entry\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"tags\":[\"tag1\",\"tag2\"],\"workspace_id\":1234567}",
		},
		{
			name: "int, string, time, and slice of int",
			in: &UpdateTimeEntryRequestBody{
				WorkspaceID: Int(1234567),
				Start:       Time(time.Date(2022, time.July, 6, 5, 4, 3, 0, time.UTC)),
				Duration:    Int(300),
				CreatedWith: String("toggl-go"),
				Description: String("updated time entry"),
				ProjectID:   Int(123456789),
				TagIDs:      []*int{Int(1234567), Int(9876543)},
			},
			out: "{\"created_with\":\"toggl-go\",\"description\":\"updated time entry\",\"duration\":300,\"project_id\":123456789,\"start\":\"2022-07-06T05:04:03Z\",\"tag_ids\":[1234567,9876543],\"workspace_id\":1234567}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaceID := 1234567
			timeEntryID := 1234567890
			_, _ = client.UpdateTimeEntry(context.Background(), workspaceID, timeEntryID, tt.in)
		})
	}
}
