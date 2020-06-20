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

	"github.com/it-akumi/toggl-go/toggl"
)

func TestGetWorkspaces(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx context.Context
		}
		out struct {
			workspaces []*toggl.Workspace
			err        error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspaces_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				workspaces []*toggl.Workspace
				err        error
			}{
				workspaces: []*toggl.Workspace{
					{
						Id:                          1234567,
						Name:                        "Sample workspace",
						Premium:                     false,
						Admin:                       true,
						DefaultHourlyRate:           0,
						DefaultCurrency:             "USD",
						OnlyAdminsMayCreateProjects: false,
						OnlyAdminsSeeBillableRates:  false,
						Rounding:                    0,
						RoundingMinutes:             0,
						At:                          time.Date(2017, time.July, 3, 9, 31, 1, 0, time.FixedZone("", 0)),
						LogoURL:                     "",
					},
					{
						Id:                          9876543,
						Name:                        "toggl-go",
						Premium:                     false,
						Admin:                       true,
						DefaultHourlyRate:           0,
						DefaultCurrency:             "USD",
						OnlyAdminsMayCreateProjects: false,
						OnlyAdminsSeeBillableRates:  false,
						Rounding:                    1,
						RoundingMinutes:             0,
						At:                          time.Date(2020, time.June, 10, 6, 2, 29, 0, time.FixedZone("", 0)),
						LogoURL:                     "",
					},
				},
				err: nil,
			},
		},
		{
			name:             "401 Unauthorized",
			httpStatus:       http.StatusUnauthorized,
			testdataFilePath: "testdata/workspaces/get_workspaces_401_unauthorized.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				workspaces []*toggl.Workspace
				err        error
			}{
				workspaces: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    401,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspaces/get_workspaces_403_forbidden.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				workspaces []*toggl.Workspace
				err        error
			}{
				workspaces: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspaces_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: nil,
			},
			out: struct {
				workspaces []*toggl.Workspace
				err        error
			}{
				workspaces: nil,
				err:        toggl.ErrContextNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaces, err := client.GetWorkspaces(c.in.ctx)
			if !reflect.DeepEqual(actualWorkspaces, c.out.workspaces) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.workspaces, actualWorkspaces)
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

func TestGetWorkspace(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			workspace *toggl.Workspace
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: &toggl.Workspace{
					Id:                          1234567,
					Name:                        "toggl-go",
					Premium:                     false,
					Admin:                       true,
					DefaultHourlyRate:           0,
					DefaultCurrency:             "USD",
					OnlyAdminsMayCreateProjects: false,
					OnlyAdminsSeeBillableRates:  false,
					Rounding:                    1,
					RoundingMinutes:             0,
					At:                          time.Date(2020, time.June, 10, 6, 2, 29, 0, time.FixedZone("", 0)),
					LogoURL:                     "",
				},
				err: nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspaces/get_workspace_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: nil,
				err: toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/workspaces/get_workspace_404_not_found.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: nil,
				err: toggl.TogglError{
					Message: "null\n",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       nil,
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: nil,
				err:       toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspace, err := client.GetWorkspace(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualWorkspace, c.out.workspace) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.workspace, actualWorkspace)
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

func TestGetWorkspaceUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetWorkspace(context.Background(), &toggl.Workspace{
		Id: workspaceId,
	})
}

func TestUpdateWorkspace(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			workspace *toggl.Workspace
			err       error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/update_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					Id:   1234567,
					Name: "updated",
				},
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: &toggl.Workspace{
					Id:                          1234567,
					Name:                        "updated",
					Premium:                     false,
					Admin:                       true,
					DefaultHourlyRate:           0,
					DefaultCurrency:             "USD",
					OnlyAdminsMayCreateProjects: false,
					OnlyAdminsSeeBillableRates:  false,
					Rounding:                    1,
					RoundingMinutes:             0,
					At:                          time.Date(2020, time.June, 10, 6, 2, 29, 0, time.FixedZone("", 0)),
					LogoURL:                     "",
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/workspaces/update_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					Id:   1234567,
					Name: "updated",
				},
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: nil,
				err: toggl.TogglError{
					Message: "workspace ID must be a positive integer\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspaces/update_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					Id:   1234567,
					Name: "updated",
				},
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: nil,
				err: toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/update_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx: nil,
				workspace: &toggl.Workspace{
					Id:   1234567,
					Name: "updated",
				},
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: nil,
				err:       toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/update_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				workspace *toggl.Workspace
				err       error
			}{
				workspace: nil,
				err:       toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspace, err := client.UpdateWorkspace(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualWorkspace, c.out.workspace) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.workspace, actualWorkspace)
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

func TestUpdateWorkspaceConvertParamsToRequestBody(t *testing.T) {
	expectedWorkspaceRequest := &toggl.Workspace{
		Id:   1234567,
		Name: "updated",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualWorkspaceRequest := new(toggl.Project)
		if err := json.Unmarshal(requestBody, actualWorkspaceRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualWorkspaceRequest, expectedWorkspaceRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedWorkspaceRequest, actualWorkspaceRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateWorkspace(context.Background(), expectedWorkspaceRequest)
}

func TestUpdateWorkspaceUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateWorkspace(context.Background(), &toggl.Workspace{
		Id: workspaceId,
	})
}

func TestGetWorkspaceUsers(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			users []*toggl.User
			err   error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_users_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				users []*toggl.User
				err   error
			}{
				users: []*toggl.User{
					{
						Id:                     1234567,
						DefaultWid:             9876543,
						Email:                  "john@swift.com",
						Fullname:               "John Swift",
						JQueryTimeofdayFormat:  "H:i",
						JQueryDateFormat:       "m/d/Y",
						TimeofdayFormat:        "H:mm",
						DateFormat:             "MM/DD/YYYY",
						StoreStartAndStopTime:  true,
						BeginningOfWeek:        1,
						Language:               "en_US",
						ImageUrl:               "https://assets.toggl.com/avatars/abcdefghijklmnopqrstuvwxyz012345.png",
						SidebarPiechart:        true,
						At:                     time.Date(2013, time.March, 7, 14, 21, 38, 0, time.FixedZone("", 0)),
						SendProductEmails:      true,
						SendWeeklyReport:       true,
						SendTimerNotifications: true,
						OpenidEnabled:          false,
						Timezone:               "Asia/Tokyo",
					},
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/workspaces/get_workspace_users_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				users []*toggl.User
				err   error
			}{
				users: nil,
				err: toggl.TogglError{
					Message: "Missing or invalid workspace_id\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspaces/get_workspace_users_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				users []*toggl.User
				err   error
			}{
				users: nil,
				err: toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_users_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       nil,
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				users []*toggl.User
				err   error
			}{
				users: nil,
				err:   toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_users_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				users []*toggl.User
				err   error
			}{
				users: nil,
				err:   toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaceUsers, err := client.GetWorkspaceUsers(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualWorkspaceUsers, c.out.users) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.users, actualWorkspaceUsers)
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

func TestGetWorkspaceUsersUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "/users" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetWorkspaceUsers(context.Background(), &toggl.Workspace{
		Id: workspaceId,
	})
}

func TestGetWorkspaceClients(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			clients []*toggl.TogglClient
			err     error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_clients_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				clients []*toggl.TogglClient
				err     error
			}{
				clients: []*toggl.TogglClient{
					{
						Id:   12345678,
						Name: "toggl-go",
						Wid:  9876543,
						At:   time.Date(2020, time.June, 10, 6, 54, 51, 0, time.FixedZone("", 0)),
					},
					{
						Id:   23456789,
						Name: "sample-client",
						Wid:  9876543,
						At:   time.Date(2020, time.June, 10, 6, 54, 47, 0, time.FixedZone("", 0)),
					},
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/workspaces/get_workspace_clients_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				clients []*toggl.TogglClient
				err     error
			}{
				clients: nil,
				err: toggl.TogglError{
					Message: "Missing or invalid workspace_id\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspaces/get_workspace_clients_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				clients []*toggl.TogglClient
				err     error
			}{
				clients: nil,
				err: toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_clients_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       nil,
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				clients []*toggl.TogglClient
				err     error
			}{
				clients: nil,
				err:     toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_clients_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				clients []*toggl.TogglClient
				err     error
			}{
				clients: nil,
				err:     toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaceClients, err := client.GetWorkspaceClients(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualWorkspaceClients, c.out.clients) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.clients, actualWorkspaceClients)
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

func TestGetWorkspaceClientsUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "/clients" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetWorkspaceClients(context.Background(), &toggl.Workspace{
		Id: workspaceId,
	})
}

func TestGetWorkspaceGroups(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			groups []*toggl.Group
			err    error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_groups_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				groups []*toggl.Group
				err    error
			}{
				groups: []*toggl.Group{
					{
						Id:   123456,
						Name: "toggl-go",
						Wid:  1234567,
						At:   time.Date(2020, time.June, 10, 6, 59, 43, 0, time.FixedZone("", 0)),
					},
					{
						Id:   234567,
						Name: "sample-group",
						Wid:  1234567,
						At:   time.Date(2020, time.June, 10, 6, 59, 38, 0, time.FixedZone("", 0)),
					},
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/workspaces/get_workspace_groups_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				groups []*toggl.Group
				err    error
			}{
				groups: nil,
				err: toggl.TogglError{
					Message: "Missing or invalid workspace_id\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspaces/get_workspace_groups_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				groups []*toggl.Group
				err    error
			}{
				groups: nil,
				err: toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_groups_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       nil,
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				groups []*toggl.Group
				err    error
			}{
				groups: nil,
				err:    toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_groups_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				groups []*toggl.Group
				err    error
			}{
				groups: nil,
				err:    toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaceGroups, err := client.GetWorkspaceGroups(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualWorkspaceGroups, c.out.groups) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.groups, actualWorkspaceGroups)
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

func TestGetWorkspaceGroupsUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "/groups" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetWorkspaceGroups(context.Background(), &toggl.Workspace{
		Id: workspaceId,
	})
}

func TestGetWorkspaceProjects(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			projects []*toggl.Project
			err      error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_projects_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				projects []*toggl.Project
				err      error
			}{
				projects: []*toggl.Project{
					{
						Id:        123456789,
						Name:      "sample-project",
						Wid:       9876543,
						Active:    true,
						IsPrivate: true,
						Template:  false,
						At:        time.Date(2020, time.June, 10, 6, 51, 48, 0, time.FixedZone("", 0)),
						Color:     "2",
						CreatedAt: time.Date(2020, time.June, 10, 6, 51, 48, 0, time.FixedZone("", 0)),
					},
					{
						Id:        234567890,
						Name:      "toggl-go",
						Wid:       9876543,
						Active:    true,
						IsPrivate: true,
						Template:  false,
						At:        time.Date(2020, time.June, 10, 6, 51, 54, 0, time.FixedZone("", 0)),
						Color:     "10",
						CreatedAt: time.Date(2020, time.June, 10, 6, 51, 54, 0, time.FixedZone("", 0)),
					},
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/workspaces/get_workspace_projects_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				projects []*toggl.Project
				err      error
			}{
				projects: nil,
				err: toggl.TogglError{
					Message: "Missing or invalid workspace_id\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspaces/get_workspace_projects_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				projects []*toggl.Project
				err      error
			}{
				projects: nil,
				err: toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_projects_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       nil,
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				projects []*toggl.Project
				err      error
			}{
				projects: nil,
				err:      toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_projects_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				projects []*toggl.Project
				err      error
			}{
				projects: nil,
				err:      toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaceProjects, err := client.GetWorkspaceProjects(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualWorkspaceProjects, c.out.projects) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.projects, actualWorkspaceProjects)
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

func TestGetWorkspaceProjectsUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "/projects" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetWorkspaceProjects(context.Background(), &toggl.Workspace{
		Id: workspaceId,
	})
}

func TestGetWorkspaceProjectsUseURLIncludingQueryStrings(t *testing.T) {
	workspaceId := 1234567
	onlyTemplates := "true"
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "/projects?only_templates=" + onlyTemplates
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetWorkspaceProjects(context.Background(), &toggl.Workspace{Id: workspaceId}, toggl.OnlyTemplates(onlyTemplates))
}

func TestGetWorkspaceTags(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			tags []*toggl.Tag
			err  error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_tags_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				tags []*toggl.Tag
				err  error
			}{
				tags: []*toggl.Tag{
					{
						Id:   1234567,
						Name: "sample-tag",
						Wid:  9876543,
					},
					{
						Id:   1234568,
						Name: "toggl-go",
						Wid:  9876543,
					},
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/workspaces/get_workspace_tags_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				tags []*toggl.Tag
				err  error
			}{
				tags: nil,
				err: toggl.TogglError{
					Message: "Missing or invalid workspace_id\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspaces/get_workspace_tags_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				tags []*toggl.Tag
				err  error
			}{
				tags: nil,
				err: toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_tags_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       nil,
				workspace: &toggl.Workspace{Id: 1234567},
			},
			out: struct {
				tags []*toggl.Tag
				err  error
			}{
				tags: nil,
				err:  toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspaces/get_workspace_tags_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				tags []*toggl.Tag
				err  error
			}{
				tags: nil,
				err:  toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaceTags, err := client.GetWorkspaceTags(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualWorkspaceTags, c.out.tags) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.tags, actualWorkspaceTags)
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

func TestGetWorkspaceTagsUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "/tags" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetWorkspaceTags(context.Background(), &toggl.Workspace{
		Id: workspaceId,
	})
}
