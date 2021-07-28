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

	"github.com/ta9mi1shi1/toggl-go/toggl"
)

func TestInviteUsersToWorkspace(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
			users     []*toggl.User
		}
		out struct {
			workspaceUsers []*toggl.WorkspaceUser
			err            error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/invite_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
				users     []*toggl.User
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
				users: []*toggl.User{
					{
						Email: "test.user@toggl.com",
					},
					{
						Email: "jane.swift@toggl.com",
					},
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: []*toggl.WorkspaceUser{
					{
						ID:        1234567,
						UID:       2345678,
						Admin:     false,
						Active:    false,
						InviteURL: "https://toggl.com/",
					},
				},
				err: errors.New("User jane.swift@toggl.com is already in the workspace"),
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/workspace_users/invite_400_bad_request.txt",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
				users     []*toggl.User
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
				users: []*toggl.User{
					{
						Email: "",
					},
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err: &toggl.TogglError{
					Message: "Emails should not be blank\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspace_users/invite_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
				users     []*toggl.User
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
				users: []*toggl.User{
					{
						Email: "test.user@toggl.com",
					},
					{
						Email: "jane.swift@toggl.com",
					},
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/invite_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
				users     []*toggl.User
			}{
				ctx: nil,
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
				users: []*toggl.User{
					{
						Email: "test.user@toggl.com",
					},
					{
						Email: "jane.swift@toggl.com",
					},
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err:            toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/invite_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
				users     []*toggl.User
			}{
				ctx:       context.Background(),
				workspace: nil,
				users: []*toggl.User{
					{
						Email: "test.user@toggl.com",
					},
					{
						Email: "jane.swift@toggl.com",
					},
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err:            toggl.ErrWorkspaceNotFound,
			},
		},
		{
			name:             "Without users",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/invite_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
				users     []*toggl.User
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
				users: nil,
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err:            toggl.ErrUserNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaceUsers, err := client.InviteUsersToWorkspace(c.in.ctx, c.in.workspace, c.in.users)
			if !reflect.DeepEqual(actualWorkspaceUsers, c.out.workspaceUsers) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.workspaceUsers, actualWorkspaceUsers)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				// Since errors.New returns a distinct error value even if the text is identical,
				// equality of errors generated dynamically by errors.New cannot be compared by errors.Is.
				if !errors.Is(err, c.out.err) && (err.Error() != c.out.err.Error()) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestInviteUsersToWorkspaceConvertParamsToRequestBody(t *testing.T) {
	workspace := &toggl.Workspace{
		ID: 3456789,
	}
	expectedRequest := struct {
		Emails []string `json:"emails"`
	}{
		Emails: []string{"test.user@toggl.com", "jane.swift@toggl.com"},
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		var actualRequest struct {
			Emails []string `json:"emails"`
		}
		if err := json.Unmarshal(requestBody, &actualRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualRequest, expectedRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequest, actualRequest)
		}
	}))

	users := []*toggl.User{}
	for _, email := range expectedRequest.Emails {
		users = append(users, &toggl.User{Email: email})
	}
	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.InviteUsersToWorkspace(context.Background(), workspace, users)
}

func TestInviteUsersToWorkspaceUseURLIncludingWorkspaceID(t *testing.T) {
	workspaceID := 1234567
	users := []*toggl.User{
		{
			Email: "test.user@toggl.com",
		},
	}
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceID) + "/invite" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.InviteUsersToWorkspace(context.Background(), &toggl.Workspace{ID: workspaceID}, users)
}

func TestUpdateWorkspaceUser(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx           context.Context
			workspaceUser *toggl.WorkspaceUser
		}
		out struct {
			workspaceUser *toggl.WorkspaceUser
			err           error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/update_200_ok.json",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx: context.Background(),
				workspaceUser: &toggl.WorkspaceUser{
					ID:    1234567,
					Admin: false,
				},
			},
			out: struct {
				workspaceUser *toggl.WorkspaceUser
				err           error
			}{
				workspaceUser: &toggl.WorkspaceUser{
					ID:     1234567,
					UID:    2345678,
					Admin:  false,
					Active: true,
				},
				err: nil,
			},
		},
		{
			name:             "401 Unauthorized",
			httpStatus:       http.StatusUnauthorized,
			testdataFilePath: "testdata/workspace_users/update_401_unauthorized.json",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx: context.Background(),
				workspaceUser: &toggl.WorkspaceUser{
					ID:    1234567,
					Admin: true,
				},
			},
			out: struct {
				workspaceUser *toggl.WorkspaceUser
				err           error
			}{
				workspaceUser: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    401,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/workspace_users/update_404_not_found.txt",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx: context.Background(),
				workspaceUser: &toggl.WorkspaceUser{
					Admin: true,
				},
			},
			out: struct {
				workspaceUser *toggl.WorkspaceUser
				err           error
			}{
				workspaceUser: nil,
				err: &toggl.TogglError{
					Message: "Workspace not found/accessible\n",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/update_200_ok.json",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx: nil,
				workspaceUser: &toggl.WorkspaceUser{
					Admin: true,
				},
			},
			out: struct {
				workspaceUser *toggl.WorkspaceUser
				err           error
			}{
				workspaceUser: nil,
				err:           toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace user",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/update_200_ok.json",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx:           context.Background(),
				workspaceUser: nil,
			},
			out: struct {
				workspaceUser *toggl.WorkspaceUser
				err           error
			}{
				workspaceUser: nil,
				err:           toggl.ErrWorkspaceUserNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaceUser, err := client.UpdateWorkspaceUser(c.in.ctx, c.in.workspaceUser)
			if !reflect.DeepEqual(actualWorkspaceUser, c.out.workspaceUser) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.workspaceUser, actualWorkspaceUser)
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

func TestUpdateWorkspaceUserConvertParamsToRequestBody(t *testing.T) {
	expectedWorkspaceUserRequest := &toggl.WorkspaceUser{
		Admin: false,
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualWorkspaceUserRequest := new(toggl.WorkspaceUser)
		if err := json.Unmarshal(requestBody, actualWorkspaceUserRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualWorkspaceUserRequest, expectedWorkspaceUserRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedWorkspaceUserRequest, actualWorkspaceUserRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateWorkspaceUser(context.Background(), expectedWorkspaceUserRequest)
}

func TestDeleteWorkspaceUser(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx           context.Context
			workspaceUser *toggl.WorkspaceUser
		}
		out error
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/delete_200_ok.json",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx: context.Background(),
				workspaceUser: &toggl.WorkspaceUser{
					ID: 1234567,
				},
			},
			out: nil,
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/workspace_users/delete_400_bad_request.txt",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx: context.Background(),
				workspaceUser: &toggl.WorkspaceUser{
					ID: 1234567,
				},
			},
			out: &toggl.TogglError{
				Message: "Cannot access workspace users\n",
				Code:    400,
			},
		},
		{
			name:             "401 Unauthorized",
			httpStatus:       http.StatusUnauthorized,
			testdataFilePath: "testdata/workspace_users/delete_401_unauthorized.json",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx: context.Background(),
				workspaceUser: &toggl.WorkspaceUser{
					ID: 1234567,
				},
			},
			out: &toggl.TogglError{
				Message: "",
				Code:    401,
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/delete_200_ok.json",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx: nil,
				workspaceUser: &toggl.WorkspaceUser{
					ID: 1234567,
				},
			},
			out: toggl.ErrContextNotFound,
		},
		{
			name:             "Without workspace user",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/delete_200_ok.json",
			in: struct {
				ctx           context.Context
				workspaceUser *toggl.WorkspaceUser
			}{
				ctx:           context.Background(),
				workspaceUser: nil,
			},
			out: toggl.ErrWorkspaceUserNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			err := client.DeleteWorkspaceUser(c.in.ctx, c.in.workspaceUser)

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

func TestDeleteWorkspaceUserUseURLIncludingWorkspaceUserID(t *testing.T) {
	workspaceUserID := 1234567
	expectedRequestURI := "/api/v8/workspace_users/" + strconv.Itoa(workspaceUserID) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_ = client.DeleteWorkspaceUser(context.Background(), &toggl.WorkspaceUser{
		ID: workspaceUserID,
	})
}

func TestGetWorkspaceUsersAsWorkspaceUser(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx       context.Context
			workspace *toggl.Workspace
		}
		out struct {
			workspaceUsers []*toggl.WorkspaceUser
			err            error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/get_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: []*toggl.WorkspaceUser{
					{
						ID:     1234567,
						UID:    2345678,
						Admin:  false,
						Active: true,
					},
					{
						ID:     9876543,
						UID:    8765432,
						Admin:  true,
						Active: true,
					},
				},
				err: nil,
			},
		},
		{
			name:             "401 Unauthorized",
			httpStatus:       http.StatusUnauthorized,
			testdataFilePath: "testdata/workspace_users/get_401_unauthorized.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    401,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/workspace_users/get_403_forbidden.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx: context.Background(),
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/get_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx: nil,
				workspace: &toggl.Workspace{
					ID: 3456789,
				},
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err:            toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without workspace",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/workspace_users/get_200_ok.json",
			in: struct {
				ctx       context.Context
				workspace *toggl.Workspace
			}{
				ctx:       context.Background(),
				workspace: nil,
			},
			out: struct {
				workspaceUsers []*toggl.WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
				err:            toggl.ErrWorkspaceNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualWorkspaceUsers, err := client.GetWorkspaceUsersAsWorkspaceUser(c.in.ctx, c.in.workspace)
			if !reflect.DeepEqual(actualWorkspaceUsers, c.out.workspaceUsers) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.workspaceUsers, actualWorkspaceUsers)
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

func TestGetWorkspaceUsersAsWorkspaceUserUseURLIncludingWorkspaceID(t *testing.T) {
	workspaceID := 1234567
	expectedRequestURI := "/api/v8/workspaces/" + strconv.Itoa(workspaceID) + "/workspace_users" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetWorkspaceUsersAsWorkspaceUser(context.Background(), &toggl.Workspace{
		ID: workspaceID,
	})
}
