package track

import (
	"context"
	"errors"
	"net/http"
	"path"
	"reflect"
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

func TestGetTimeEntriesQueries(t *testing.T) {
	tests := []struct {
		name string
		in   *GetTimeEntriesQueries
		out  string
	}{
		{
			name: "GetTimeEntriesQueries is nil",
			in:   nil,
			out:  "",
		},
		{
			name: "before=2022-07-01",
			in:   &GetTimeEntriesQueries{Before: String("2022-07-01")},
			out:  "before=2022-07-01",
		},
		{
			name: "since=1656687597",
			in:   &GetTimeEntriesQueries{Since: Int(1656687597)},
			out:  "since=1656687597",
		},
		{
			name: "start_date=2022-07-01&end_date=2022-07-07",
			in: &GetTimeEntriesQueries{
				StartDate: String("2022-07-01"),
				EndDate:   String("2022-07-07"),
			},
			out: "start_date=2022-07-01&end_date=2022-07-07",
		},
		{
			name: "GetTimeEntriesQueries is empty",
			in:   &GetTimeEntriesQueries{},
			out:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertQueries(t, tt.out)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			_, _ = client.GetTimeEntries(context.Background(), tt.in)
		})
	}
}
