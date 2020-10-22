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

func TestCreateProjectUser(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			projectUser *toggl.ProjectUser
		}
		out struct {
			projectUser *toggl.ProjectUser
			err         error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/create_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Pid: 123456789,
					Uid: 2345678,
				},
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: &toggl.ProjectUser{
					Id:      12345678,
					Pid:     123456789,
					Uid:     2345678,
					Wid:     3456789,
					Manager: false,
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/project_users/create_400_bad_request.txt",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Pid: 123456789,
					Uid: 2345678,
				},
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: nil,
				err: &toggl.TogglError{
					Message: "You have no permissions to add users to the project\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/project_users/create_403_forbidden.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Pid: 123456789,
					Uid: 2345678,
				},
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/create_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: nil,
				projectUser: &toggl.ProjectUser{
					Pid: 123456789,
					Uid: 2345678,
				},
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: nil,
				err:         toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without project user",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/create_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx:         context.Background(),
				projectUser: nil,
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: nil,
				err:         toggl.ErrProjectUserNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualProjectUser, err := client.CreateProjectUser(c.in.ctx, c.in.projectUser)
			if !reflect.DeepEqual(actualProjectUser, c.out.projectUser) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.projectUser, actualProjectUser)
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

func TestCreateProjectUserConvertParamsToRequestBody(t *testing.T) {
	expectedProjectUserRequest := &toggl.ProjectUser{
		Pid:     777,
		Uid:     123,
		Manager: true,
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualProjectUserRequest := new(toggl.ProjectUser)
		if err := json.Unmarshal(requestBody, actualProjectUserRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualProjectUserRequest, expectedProjectUserRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedProjectUserRequest, actualProjectUserRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.CreateProjectUser(context.Background(), expectedProjectUserRequest)
}

func TestUpdateProjectUser(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			projectUser *toggl.ProjectUser
		}
		out struct {
			projectUser *toggl.ProjectUser
			err         error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/update_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Pid:     123456789,
					Uid:     2345678,
					Wid:     3456789,
					Manager: true,
					Fields:  "fullname",
				},
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: &toggl.ProjectUser{
					Id:       98765432,
					Pid:      123456789,
					Uid:      2345678,
					Wid:      3456789,
					Manager:  true,
					Fullname: "John Swift",
					At:       time.Date(2019, time.September, 15, 1, 24, 49, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/project_users/update_400_bad_request.txt",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Pid:     123456789,
					Uid:     2345678,
					Wid:     3456789,
					Manager: true,
					Fields:  "fullname",
				},
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: nil,
				err: &toggl.TogglError{
					Message: "invalid character 'h' looking for beginning of value\n",
					Code:    400,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/project_users/update_404_not_found.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Pid:     123456789,
					Uid:     2345678,
					Wid:     3456789,
					Manager: true,
					Fields:  "fullname",
				},
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: nil,
				err: &toggl.TogglError{
					Message: "null\n",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/update_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: nil,
				projectUser: &toggl.ProjectUser{
					Pid:     123456789,
					Uid:     2345678,
					Wid:     3456789,
					Manager: true,
					Fields:  "fullname",
				},
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: nil,
				err:         toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without project user",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/update_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx:         context.Background(),
				projectUser: nil,
			},
			out: struct {
				projectUser *toggl.ProjectUser
				err         error
			}{
				projectUser: nil,
				err:         toggl.ErrProjectUserNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualProjectUsers, err := client.UpdateProjectUser(c.in.ctx, c.in.projectUser)
			if !reflect.DeepEqual(actualProjectUsers, c.out.projectUser) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.projectUser, actualProjectUsers)
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

func TestUpdateProjectUserConvertParamsToRequestBody(t *testing.T) {
	expectedProjectUserRequest := &toggl.ProjectUser{
		Id:      1234567,
		Manager: true,
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualProjectUserRequest := new(toggl.ProjectUser)
		if err := json.Unmarshal(requestBody, actualProjectUserRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualProjectUserRequest, expectedProjectUserRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedProjectUserRequest, actualProjectUserRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateProjectUser(context.Background(), expectedProjectUserRequest)
}

func TestUpdateProjectUserUseURLIncludingProjectUserId(t *testing.T) {
	projectUserId := 12345678
	expectedRequestURI := "/api/v8/project_users/" + strconv.Itoa(projectUserId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateProjectUser(context.Background(), &toggl.ProjectUser{
		Id:      projectUserId,
		Manager: true,
	})
}

func TestDeleteProjectUser(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			projectUser *toggl.ProjectUser
		}
		out error
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/delete_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Id: 12345678,
				},
			},
			out: nil,
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/project_users/delete_403_forbidden.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Id: 12345678,
				},
			},
			out: &toggl.TogglError{
				Message: "",
				Code:    403,
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/project_users/delete_404_not_found.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: context.Background(),
				projectUser: &toggl.ProjectUser{
					Id: 12345678,
				},
			},
			out: &toggl.TogglError{
				Message: "null\n",
				Code:    404,
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/delete_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx: nil,
				projectUser: &toggl.ProjectUser{
					Id: 12345678,
				},
			},
			out: toggl.ErrContextNotFound,
		},
		{
			name:             "Without project user",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/delete_200_ok.json",
			in: struct {
				ctx         context.Context
				projectUser *toggl.ProjectUser
			}{
				ctx:         context.Background(),
				projectUser: nil,
			},
			out: toggl.ErrProjectUserNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			err := client.DeleteProjectUser(c.in.ctx, c.in.projectUser)

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

func TestDeleteProjectUserUseURLIncludingProjectUserId(t *testing.T) {
	projectUserId := 12345678
	expectedRequestURI := "/api/v8/project_users/" + strconv.Itoa(projectUserId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_ = client.DeleteProjectUser(context.Background(), &toggl.ProjectUser{
		Id: projectUserId,
	})
}

func TestGetProjectUsersInWorkspace(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			projectUsers []*toggl.ProjectUser
			err          error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/get_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 4567890},
			},
			out: struct {
				projectUsers []*toggl.ProjectUser
				err          error
			}{
				projectUsers: []*toggl.ProjectUser{
					{
						Id:      12345678,
						Pid:     234567890,
						Uid:     3456789,
						Wid:     4567890,
						Manager: false,
						At:      time.Date(2018, time.March, 14, 1, 23, 45, 0, time.FixedZone("", 0)),
					},
					{
						Id:      23456789,
						Pid:     234567890,
						Uid:     4567890,
						Wid:     4567890,
						Manager: true,
						At:      time.Date(2020, time.February, 17, 9, 49, 59, 0, time.FixedZone("", 0)),
					},
					{
						Id:      34567890,
						Pid:     345678901,
						Uid:     3456789,
						Wid:     4567890,
						Manager: true,
						At:      time.Date(2017, time.January, 11, 5, 46, 47, 0, time.FixedZone("", 0)),
					},
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/project_users/get_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 4567890},
			},
			out: struct {
				projectUsers []*toggl.ProjectUser
				err          error
			}{
				projectUsers: nil,
				err: &toggl.TogglError{
					Message: "Missing or invalid workspace_id\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/project_users/get_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: &toggl.Workspace{Id: 4567890},
			},
			out: struct {
				projectUsers []*toggl.ProjectUser
				err          error
			}{
				projectUsers: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/get_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       nil,
				workspace: &toggl.Workspace{Id: 4567890},
			},
			out: struct {
				projectUsers []*toggl.ProjectUser
				err          error
			}{
				projectUsers: nil,
				err:          toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/project_users/get_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				projectUsers []*toggl.ProjectUser
				err          error
			}{
				projectUsers: nil,
				err:          toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualProjectUsers, err := client.GetProjectUsersInWorkspace(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualProjectUsers, c.out.projectUsers) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.projectUsers, actualProjectUsers)
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

func TestGetProjectUsersInWorkspaceUseURLIncludingWorkspaceId(t *testing.T) {
	workspaceId := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceId) + "/project_users" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetProjectUsersInWorkspace(context.Background(), &toggl.Workspace{
		Id: workspaceId,
	})
}
