package toggl_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/ta9mi1shi1/toggl-go/toggl"
)

func TestCreateTimeEntry(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			timeEntry *toggl.TimeEntry
		}
		out struct {
			timeEntry *toggl.TimeEntry
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/create_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Pid:         123456789,
					Start:       time.Date(2013, time.March, 5, 7, 58, 58, 0, time.FixedZone("", 0)),
					Duration:    1200,
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					CreatedWith: "toggl-go",
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: &toggl.TimeEntry{
					Id:          1234567890,
					Wid:         2345678,
					Pid:         123456789,
					Start:       time.Date(2013, time.March, 5, 7, 58, 58, 0, time.FixedZone("", 0)),
					Stop:        time.Date(2013, time.March, 5, 8, 18, 58, 0, time.FixedZone("", 0)),
					Duration:    1200,
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					Duronly:     false,
					At:          time.Date(2018, time.September, 23, 8, 47, 51, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/time_entries/create_400_bad_request.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Pid:         123456789,
					Start:       time.Date(2013, time.March, 5, 7, 58, 58, 0, time.FixedZone("", 0)),
					Duration:    1200,
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					CreatedWith: "",
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "created_with needs to be provided an a valid string\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/time_entries/create_403_forbidden.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Pid:         123456789,
					Start:       time.Date(2013, time.March, 5, 7, 58, 58, 0, time.FixedZone("", 0)),
					Duration:    1200,
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					CreatedWith: "toggl-go",
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/create_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: nil,
				timeEntry: &toggl.TimeEntry{
					Pid:         123456789,
					Start:       time.Date(2013, time.March, 5, 7, 58, 58, 0, time.FixedZone("", 0)),
					Duration:    1200,
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					CreatedWith: "toggl-go",
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without time entry",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/create_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx:       context.Background(),
				timeEntry: nil,
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrTimeEntryNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTimeEntry, err := client.CreateTimeEntry(c.in.ctx, c.in.timeEntry)
			if !reflect.DeepEqual(actualTimeEntry, c.out.timeEntry) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.timeEntry, actualTimeEntry)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestCreateTimeEntryConvertParamsToRequestBody(t *testing.T) {
	expectedTimeEntryRequest := &toggl.TimeEntry{
		Description: "Meeting with possible clients",
		Tags:        []string{"billed"},
		Duration:    1200,
		Start:       time.Date(2013, time.March, 5, 7, 58, 58, 0, time.FixedZone("", 0)),
		Pid:         123,
		CreatedWith: "toggl-go",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualTimeEntryRequest := new(toggl.TimeEntry)
		if err := json.Unmarshal(requestBody, actualTimeEntryRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualTimeEntryRequest, expectedTimeEntryRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedTimeEntryRequest, actualTimeEntryRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.CreateTimeEntry(context.Background(), expectedTimeEntryRequest)
}

func TestUpdateTimeEntry(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			timeEntry *toggl.TimeEntry
		}
		out struct {
			timeEntry *toggl.TimeEntry
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/update_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Description: "Meeting with possible clients",
					Tags:        []string{""},
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: &toggl.TimeEntry{
					Id:          1234567890,
					Wid:         1234567,
					Pid:         123456789,
					Start:       time.Date(2020, time.October, 5, 7, 58, 58, 0, time.FixedZone("", 0)),
					Stop:        time.Date(2020, time.October, 5, 8, 18, 58, 0, time.FixedZone("", 0)),
					Duration:    1200,
					Description: "Meeting with possible clients",
					Duronly:     false,
					At:          time.Date(2020, time.October, 11, 1, 23, 45, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/time_entries/update_403_forbidden.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Description: "Meeting with possible clients",
					Tags:        []string{""},
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/time_entries/update_404_not_found.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Description: "Meeting with possible clients",
					Tags:        []string{""},
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "Time entry not found/no access to it",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/update_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: nil,
				timeEntry: &toggl.TimeEntry{
					Description: "Meeting with possible clients",
					Tags:        []string{""},
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without time entry",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/update_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx:       context.Background(),
				timeEntry: nil,
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrTimeEntryNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTimeEntry, err := client.UpdateTimeEntry(c.in.ctx, c.in.timeEntry)
			if !reflect.DeepEqual(actualTimeEntry, c.out.timeEntry) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.timeEntry, actualTimeEntry)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestUpdateTimeEntryConvertParamsToRequestBody(t *testing.T) {
	expectedTimeEntryRequest := &toggl.TimeEntry{
		Description: "Meeting with possible clients",
		Tags:        []string{""},
		Duration:    1240,
		Start:       time.Date(2013, time.March, 5, 7, 58, 58, 0, time.FixedZone("", 0)),
		Stop:        time.Date(2013, time.March, 5, 8, 58, 58, 0, time.FixedZone("", 0)),
		Pid:         123,
		Duronly:     true,
		CreatedWith: "toggl-go",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualTimeEntryRequest := new(toggl.TimeEntry)
		if err := json.Unmarshal(requestBody, actualTimeEntryRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualTimeEntryRequest, expectedTimeEntryRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedTimeEntryRequest, actualTimeEntryRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateTimeEntry(context.Background(), expectedTimeEntryRequest)
}

func TestUpdateTimeEntryUseURLIncludingTimeEntryId(t *testing.T) {
	timeEntryId := 1234567890
	expectedRequestURI := "/api/v8/time_entries/" + strconv.Itoa(timeEntryId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateTimeEntry(context.Background(), &toggl.TimeEntry{
		Id: timeEntryId,
	})
}

func TestDeleteTimeEntry(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			timeEntry *toggl.TimeEntry
		}
		out error
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/delete_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: nil,
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/time_entries/delete_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: &toggl.TogglError{
				Message: "strconv.ParseInt: parsing \"test\": invalid syntax\n",
				Code:    400,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/time_entries/delete_403_forbidden.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: &toggl.TogglError{
				Message: "",
				Code:    403,
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/delete_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: nil,
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: toggl.ErrContextNotFound,
		},
		{
			name:             "Without time entry",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/delete_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx:       context.Background(),
				timeEntry: nil,
			},
			out: toggl.ErrTimeEntryNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			err := client.DeleteTimeEntry(c.in.ctx, c.in.timeEntry)

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out, togglError)
				}
			} else {
				if !errors.Is(err, c.out) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out, err)
				}
			}
		})
	}
}

func TestDeleteTimeEntryUseURLIncludingTimeEntryId(t *testing.T) {
	timeEntryId := 1234567890
	expectedRequestURI := "/api/v8/time_entries/" + strconv.Itoa(timeEntryId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_ = client.DeleteTimeEntry(context.Background(), &toggl.TimeEntry{
		Id: timeEntryId,
	})
}

func TestGetTimeEntry(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			timeEntry *toggl.TimeEntry
		}
		out struct {
			timeEntry *toggl.TimeEntry
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/get_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: &toggl.TimeEntry{
					Id:          1234567890,
					Wid:         1234567,
					Pid:         123456789,
					Start:       time.Date(2018, time.October, 23, 10, 11, 59, 0, time.FixedZone("", 0)),
					Stop:        time.Date(2018, time.October, 23, 10, 16, 20, 0, time.FixedZone("", 0)),
					Duration:    261,
					Description: "Test",
					Duronly:     false,
					At:          time.Date(2018, time.October, 23, 10, 16, 20, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/time_entries/get_403_forbidden.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/time_entries/get_404_not_found.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "null\n",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/get_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: nil,
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without time entry",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/get_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx:       context.Background(),
				timeEntry: nil,
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrTimeEntryNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTimeEntry, err := client.GetTimeEntry(c.in.ctx, c.in.timeEntry)
			if !reflect.DeepEqual(actualTimeEntry, c.out.timeEntry) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.timeEntry, actualTimeEntry)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestGetTimeEntryUseURLIncludingTimeEntryId(t *testing.T) {
	timeEntryId := 1234567890
	expectedRequestURI := "/api/v8/time_entries/" + strconv.Itoa(timeEntryId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetTimeEntry(context.Background(), &toggl.TimeEntry{
		Id: timeEntryId,
	})
}

func TestGetRunningTimeEntry(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx context.Context
		}
		out struct {
			timeEntry *toggl.TimeEntry
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/get_current_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: &toggl.TimeEntry{
					Id:          1234567890,
					Wid:         1234567,
					Pid:         123456789,
					Start:       time.Date(2018, time.October, 22, 9, 58, 20, 0, time.FixedZone("", 0)),
					Duration:    -1603274300,
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					Duronly:     false,
					At:          time.Date(2018, time.October, 22, 9, 58, 20, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "200 OK (null)",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/get_current_200_ok_null.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/time_entries/get_current_403_forbidden.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/get_current_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: nil,
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTimeEntry, err := client.GetRunningTimeEntry(c.in.ctx)
			if !reflect.DeepEqual(actualTimeEntry, c.out.timeEntry) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.timeEntry, actualTimeEntry)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestStart(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			timeEntry *toggl.TimeEntry
		}
		out struct {
			timeEntry *toggl.TimeEntry
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/start_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					Pid:         123456789,
					CreatedWith: "toggl-go",
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: &toggl.TimeEntry{
					Id:          1234567890,
					Wid:         1234567,
					Pid:         123456789,
					Start:       time.Date(2018, time.October, 29, 1, 23, 45, 0, time.FixedZone("", 0)),
					Duration:    -1603274300,
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					Duronly:     false,
					At:          time.Date(2018, time.October, 29, 1, 23, 45, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/time_entries/start_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					Pid:         123456789,
					CreatedWith: "toggl-go",
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "User cannot access the selected project",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/time_entries/start_403_forbidden.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					Pid:         123456789,
					CreatedWith: "toggl-go",
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/start_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: nil,
				timeEntry: &toggl.TimeEntry{
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					Pid:         123456789,
					CreatedWith: "toggl-go",
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without time entry",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/start_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx:       context.Background(),
				timeEntry: nil,
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrTimeEntryNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTimeEntry, err := client.Start(c.in.ctx, c.in.timeEntry)
			if !reflect.DeepEqual(actualTimeEntry, c.out.timeEntry) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.timeEntry, actualTimeEntry)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestStartConvertParamsToRequestBody(t *testing.T) {
	expectedStartRequest := &toggl.TimeEntry{
		Description: "Meeting with possible clients",
		Tags:        []string{"billed"},
		Pid:         123,
		CreatedWith: "toggl-go",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualStartRequest := new(toggl.TimeEntry)
		if err := json.Unmarshal(requestBody, actualStartRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualStartRequest, expectedStartRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedStartRequest, actualStartRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.Start(context.Background(), expectedStartRequest)
}

func TestStop(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			timeEntry *toggl.TimeEntry
		}
		out struct {
			timeEntry *toggl.TimeEntry
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/stop_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: &toggl.TimeEntry{
					Id:          1234567890,
					Wid:         1234567,
					Pid:         123456789,
					Start:       time.Date(2018, time.October, 22, 9, 58, 20, 0, time.FixedZone("", 0)),
					Stop:        time.Date(2018, time.October, 22, 10, 8, 4, 0, time.FixedZone("", 0)),
					Duration:    584,
					Description: "Meeting with possible clients",
					Tags:        []string{"billed"},
					Duronly:     false,
					At:          time.Date(2018, time.October, 22, 10, 8, 4, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/time_entries/stop_403_forbidden.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/time_entries/stop_404_not_found.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: context.Background(),
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err: &toggl.TogglError{
					Message: "Time entry not found/no access to it",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/stop_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx: nil,
				timeEntry: &toggl.TimeEntry{
					Id: 1234567890,
				},
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without time entry",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/time_entries/stop_200_ok.json",
			in: struct {
				ctx       context.Context
				timeEntry *toggl.TimeEntry
			}{
				ctx:       context.Background(),
				timeEntry: nil,
			},
			out: struct {
				timeEntry *toggl.TimeEntry
				err       error
			}{
				timeEntry: nil,
				err:       toggl.ErrTimeEntryNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTimeEntry, err := client.Stop(c.in.ctx, c.in.timeEntry)
			if !reflect.DeepEqual(actualTimeEntry, c.out.timeEntry) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.timeEntry, actualTimeEntry)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestStopUseURLIncludingTimeEntryId(t *testing.T) {
	timeEntryId := 1234567890
	expectedRequestURI := "/api/v8/time_entries/" + strconv.Itoa(timeEntryId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.Stop(context.Background(), &toggl.TimeEntry{
		Id: timeEntryId,
	})
}
