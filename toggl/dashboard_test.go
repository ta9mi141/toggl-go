package toggl_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/it-akumi/toggl-go/toggl"
)

func TestGetDashboard(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			workspaceId int
		}
		out struct {
			dashboard *toggl.Dashboard
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/dashboard/get_200_ok.json",
			in: struct {
				ctx         context.Context
				workspaceId int
			}{
				ctx:         context.Background(),
				workspaceId: 1234567,
			},
			out: struct {
				dashboard *toggl.Dashboard
				err       error
			}{
				dashboard: &toggl.Dashboard{
					Activity: []struct {
						UserId      int       `json:"user_id"`
						ProjectId   int       `json:"project_id"`
						Duration    int       `json:"duration"`
						Description string    `json:"description"`
						Stop        time.Time `json:"stop"`
						Tid         int       `json:"tid"`
					}{
						{
							UserId:      1234567,
							ProjectId:   12345678,
							Duration:    -1580912718,
							Description: "toggl-go",
						},
						{
							UserId:      1234567,
							ProjectId:   12345678,
							Duration:    1413,
							Description: "toggl-go",
							Stop:        time.Date(2020, time.February, 5, 0, 24, 23, 0, time.FixedZone("", 0)),
						},
						{
							UserId:      1234567,
							ProjectId:   87654321,
							Duration:    3426,
							Description: "og-lggot",
							Stop:        time.Date(2020, time.February, 4, 13, 20, 44, 0, time.FixedZone("", 0)),
						},
						{
							UserId:      1234567,
							ProjectId:   87654321,
							Duration:    178,
							Description: "og-lggot",
							Stop:        time.Date(2020, time.February, 3, 1, 48, 17, 0, time.FixedZone("", 0)),
						},
						{
							UserId:      1234567,
							ProjectId:   12345678,
							Duration:    4510,
							Description: "toggl-go",
							Stop:        time.Date(2020, time.February, 2, 7, 51, 23, 0, time.FixedZone("", 0)),
						},
					},
					MostActiveUser: []struct {
						UserId   int `json:"user_id"`
						Duration int `json:"duration"`
					}{
						{
							UserId:   1234567,
							Duration: 123456,
						},
					},
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/dashboard/get_400_bad_request.json",
			in: struct {
				ctx         context.Context
				workspaceId int
			}{
				ctx:         context.Background(),
				workspaceId: 1234567,
			},
			out: struct {
				dashboard *toggl.Dashboard
				err       error
			}{
				dashboard: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/dashboard/get_403_forbidden.json",
			in: struct {
				ctx         context.Context
				workspaceId int
			}{
				ctx:         context.Background(),
				workspaceId: 1234567,
			},
			out: struct {
				dashboard *toggl.Dashboard
				err       error
			}{
				dashboard: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/dashboard/get_404_not_found.json",
			in: struct {
				ctx         context.Context
				workspaceId int
			}{
				ctx:         context.Background(),
				workspaceId: 1234567,
			},
			out: struct {
				dashboard *toggl.Dashboard
				err       error
			}{
				dashboard: nil,
				err: &toggl.TogglError{
					Message: "404 page not found\n",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/dashboard/get_200_ok.json",
			in: struct {
				ctx         context.Context
				workspaceId int
			}{
				ctx:         nil,
				workspaceId: 1234567,
			},
			out: struct {
				dashboard *toggl.Dashboard
				err       error
			}{
				dashboard: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualDashboard, err := client.GetDashboard(c.in.ctx, c.in.workspaceId)
			if !reflect.DeepEqual(actualDashboard, c.out.dashboard) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.dashboard, actualDashboard)
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

func TestGetDashboardUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/dashboard/" + strconv.Itoa(workspaceId)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetDashboard(context.Background(), workspaceId)
}
