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

func TestCreateProject(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx     context.Context
			project *toggl.Project
		}
		out struct {
			project *toggl.Project
			err     error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/create_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx: context.Background(),
				project: &toggl.Project{
					Name: "toggl-go",
					Wid:  1234567,
				},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: &toggl.Project{
					Id:        123456789,
					Name:      "toggl-go",
					Wid:       1234567,
					Active:    true,
					IsPrivate: true,
					Template:  false,
					At:        time.Date(2020, time.May, 17, 6, 58, 8, 0, time.FixedZone("", 0)),
					Color:     "6",
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/projects/create_400_bad_request.txt",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx: context.Background(),
				project: &toggl.Project{
					Name: "toggl-go",
					Wid:  1234567,
				},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err: &toggl.TogglError{
					Message: "unexpected end of JSON input\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/projects/create_403_forbidden.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx: context.Background(),
				project: &toggl.Project{
					Name: "toggl-go",
					Wid:  1234567,
				},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/create_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx: nil,
				project: &toggl.Project{
					Name: "toggl-go",
					Wid:  1234567,
				},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err:     toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without project",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/create_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: nil,
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err:     toggl.ErrProjectNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualProject, err := client.CreateProject(c.in.ctx, c.in.project)
			if !reflect.DeepEqual(actualProject, c.out.project) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.project, actualProject)
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

func TestCreateProjectConvertParamsToRequestBody(t *testing.T) {
	expectedProjectRequest := &toggl.Project{
		Name:      "toggl-go",
		Wid:       1234567,
		IsPrivate: true,
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualProjectRequest := new(toggl.Project)
		if err := json.Unmarshal(requestBody, actualProjectRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualProjectRequest, expectedProjectRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedProjectRequest, actualProjectRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.CreateProject(context.Background(), expectedProjectRequest)
}

func TestUpdateProject(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx     context.Context
			project *toggl.Project
		}
		out struct {
			project *toggl.Project
			err     error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/update_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx: context.Background(),
				project: &toggl.Project{
					Wid:  1234567,
					Name: "updated",
				},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: &toggl.Project{
					Id:        123456789,
					Name:      "updated",
					Wid:       1234567,
					Active:    true,
					IsPrivate: true,
					Template:  false,
					At:        time.Date(2020, time.May, 17, 7, 1, 10, 0, time.FixedZone("", 0)),
					Color:     "5",
					CreatedAt: time.Date(2020, time.May, 17, 7, 1, 10, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/projects/update_400_bad_request.txt",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx: context.Background(),
				project: &toggl.Project{
					Name: "updated",
					Wid:  1234567,
				},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err: &toggl.TogglError{
					Message: "unexpected end of JSON input\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/projects/update_403_forbidden.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx: context.Background(),
				project: &toggl.Project{
					Name: "updated",
					Wid:  1234567,
				},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/update_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx: nil,
				project: &toggl.Project{
					Name: "updated",
					Wid:  1234567,
				},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err:     toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without project",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/update_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: nil,
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err:     toggl.ErrProjectNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualProject, err := client.UpdateProject(c.in.ctx, c.in.project)
			if !reflect.DeepEqual(actualProject, c.out.project) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.project, actualProject)
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

func TestUpdateProjectConvertParamsToRequestBody(t *testing.T) {
	expectedProjectRequest := &toggl.Project{
		Name: "updated",
		Wid:  1234567,
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualProjectRequest := new(toggl.Project)
		if err := json.Unmarshal(requestBody, actualProjectRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualProjectRequest, expectedProjectRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedProjectRequest, actualProjectRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateProject(context.Background(), expectedProjectRequest)
}

func TestUpdateProjectUseURLIncludingProjectId(t *testing.T) {
	projectId := 123456789
	expectedRequestURI := "/api/v8/projects/" + strconv.Itoa(projectId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateProject(context.Background(), &toggl.Project{
		Id:        projectId,
		Name:      "toggl-go",
		IsPrivate: false,
	})
}

func TestDeleteProjects(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx      context.Context
			projects []*toggl.Project
		}
		out error
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/delete_200_ok.json",
			in: struct {
				ctx      context.Context
				projects []*toggl.Project
			}{
				ctx: context.Background(),
				projects: []*toggl.Project{
					{Id: 123456789},
				},
			},
			out: nil,
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/projects/delete_400_bad_request.txt",
			in: struct {
				ctx      context.Context
				projects []*toggl.Project
			}{
				ctx: context.Background(),
				projects: []*toggl.Project{
					{Id: 0},
				},
			},
			out: &toggl.TogglError{
				Message: "project_id must be a positive integer\n",
				Code:    400,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/projects/delete_403_forbidden.json",
			in: struct {
				ctx      context.Context
				projects []*toggl.Project
			}{
				ctx: context.Background(),
				projects: []*toggl.Project{
					{Id: 123456789},
				},
			},
			out: &toggl.TogglError{
				Message: "",
				Code:    403,
			},
		},
		{
			name:             "200 OK Mass Actions",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/mass_delete_200_ok.json",
			in: struct {
				ctx      context.Context
				projects []*toggl.Project
			}{
				ctx: context.Background(),
				projects: []*toggl.Project{
					{Id: 123456789},
					{Id: 234567891},
					{Id: 345678912},
				},
			},
			out: nil,
		},
		{
			name:             "400 Bad Request Mass Actions",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/projects/mass_delete_400_bad_request.txt",
			in: struct {
				ctx      context.Context
				projects []*toggl.Project
			}{
				ctx: context.Background(),
				projects: []*toggl.Project{
					{Id: 123456789},
					{Id: 234567891},
					{Id: 0},
				},
			},
			out: &toggl.TogglError{
				Message: "project_id must be a positive integer\n",
				Code:    400,
			},
		},
		{
			name:             "403 Forbidden Mass Actions",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/projects/mass_delete_403_forbidden.json",
			in: struct {
				ctx      context.Context
				projects []*toggl.Project
			}{
				ctx: context.Background(),
				projects: []*toggl.Project{
					{Id: 123456789},
					{Id: 234567891},
					{Id: 345678912},
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
			testdataFilePath: "testdata/projects/delete_200_ok.json",
			in: struct {
				ctx      context.Context
				projects []*toggl.Project
			}{
				ctx: nil,
				projects: []*toggl.Project{
					{Id: 123456789},
				},
			},
			out: toggl.ErrContextNotFound,
		},
		{
			name:             "Without projects",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/delete_200_ok.json",
			in: struct {
				ctx      context.Context
				projects []*toggl.Project
			}{
				ctx:      context.Background(),
				projects: nil,
			},
			out: toggl.ErrProjectNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			err := client.DeleteProjects(c.in.ctx, c.in.projects)

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

func TestDeleteProjectsUseURLIncludingProjectId(t *testing.T) {
	projectId := 123456789
	expectedRequestURI := "/api/v8/projects/" + strconv.Itoa(projectId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_ = client.DeleteProjects(context.Background(), []*toggl.Project{
		{
			Id: projectId,
		},
	})
}

func TestDeleteProjectsUseURLIncludingProjectIds(t *testing.T) {
	projectIds := []int{123456789, 234567891, 345678912}

	var projects []*toggl.Project
	expectedRequestURI := "/api/v8/projects/"
	for i, id := range projectIds {
		projects = append(projects, &toggl.Project{Id: id})
		if i > 0 {
			expectedRequestURI += ","
		}
		expectedRequestURI += strconv.Itoa(id)
	}
	expectedRequestURI += "?"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_ = client.DeleteProjects(context.Background(), projects)
}

func TestGetProject(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx     context.Context
			project *toggl.Project
		}
		out struct {
			project *toggl.Project
			err     error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/get_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: &toggl.Project{Id: 123456789},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: &toggl.Project{
					Id:        123456789,
					Name:      "toggl-go",
					Wid:       1234567,
					Active:    true,
					IsPrivate: true,
					Template:  false,
					At:        time.Date(2020, time.May, 17, 7, 1, 10, 0, time.FixedZone("", 0)),
					Color:     "5",
					CreatedAt: time.Date(2020, time.May, 17, 7, 1, 10, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/projects/get_403_forbidden.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: &toggl.Project{Id: 123456789},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/projects/get_404_not_found.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: &toggl.Project{Id: 123456789},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err: &toggl.TogglError{
					Message: "null\n",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/get_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     nil,
				project: &toggl.Project{Id: 123456789},
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err:     toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without project",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/get_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: nil,
			},
			out: struct {
				project *toggl.Project
				err     error
			}{
				project: nil,
				err:     toggl.ErrProjectNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualProject, err := client.GetProject(c.in.ctx, c.in.project)
			if !reflect.DeepEqual(actualProject, c.out.project) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.project, actualProject)
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

func TestGetProjectUseURLIncludingProjectId(t *testing.T) {
	projectId := 123456789
	expectedRequestURI := "/api/v8/projects/" + strconv.Itoa(projectId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetProject(context.Background(), &toggl.Project{
		Id: projectId,
	})
}

func TestGetProjectUsers(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx     context.Context
			project *toggl.Project
		}
		out struct {
			projectUsers []*toggl.ProjectUser
			err          error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/get_project_users_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: &toggl.Project{Id: 123456789},
			},
			out: struct {
				projectUsers []*toggl.ProjectUser
				err          error
			}{
				projectUsers: []*toggl.ProjectUser{
					{
						Pid:     123456789,
						Uid:     1234567,
						Wid:     1234567,
						Manager: true,
						At:      time.Date(2020, time.May, 17, 7, 1, 10, 0, time.FixedZone("", 0)),
					},
				},
				err: nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/projects/get_project_users_403_forbidden.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: &toggl.Project{Id: 123456789},
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
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/projects/get_project_users_404_not_found.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: &toggl.Project{Id: 123456789},
			},
			out: struct {
				projectUsers []*toggl.ProjectUser
				err          error
			}{
				projectUsers: nil,
				err: &toggl.TogglError{
					Message: "null\n",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/get_project_users_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     nil,
				project: &toggl.Project{Id: 123456789},
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
			name:             "Without project",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/projects/get_project_users_200_ok.json",
			in: struct {
				ctx     context.Context
				project *toggl.Project
			}{
				ctx:     context.Background(),
				project: nil,
			},
			out: struct {
				projectUsers []*toggl.ProjectUser
				err          error
			}{
				projectUsers: nil,
				err:          toggl.ErrProjectNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualProjectUsers, err := client.GetProjectUsers(c.in.ctx, c.in.project)
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

func TestGetProjectUsersUseURLIncludingProjectId(t *testing.T) {
	projectId := 123456789
	expectedRequestURI := "/api/v8/projects/" + strconv.Itoa(projectId) + "/project_users" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetProjectUsers(context.Background(), &toggl.Project{
		Id: projectId,
	})
}
