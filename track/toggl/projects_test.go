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

func TestGetWorkspaceProjects(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			projects []*Project
			err      error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/projects/get_workspace_projects_200_ok.json",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: []*Project{
					{
						ID:                  track.Ptr(123456789),
						WorkspaceID:         track.Ptr(1234567),
						ClientID:            nil,
						Name:                track.Ptr("Project1"),
						IsPrivate:           track.Ptr(false),
						Active:              track.Ptr(true),
						At:                  track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 0, time.FixedZone("", 0))),
						CreatedAt:           track.Ptr(time.Date(2021, time.January, 2, 3, 4, 5, 0, time.FixedZone("", 0))),
						ServerDeletedAt:     nil,
						Color:               track.Ptr("#abcdef"),
						Billable:            nil,
						Template:            nil,
						AutoEstimates:       nil,
						EstimatedHours:      nil,
						Rate:                nil,
						RateLastUpdated:     nil,
						Currency:            nil,
						Recurring:           track.Ptr(false),
						RecurringParameters: nil,
						CurrentPeriod:       nil,
						FixedFee:            nil,
						ActualHours:         track.Ptr(0),
						WID:                 track.Ptr(1234567),
						CID:                 nil,
					},
					{
						ID:                  track.Ptr(234567890),
						WorkspaceID:         track.Ptr(1234567),
						ClientID:            nil,
						Name:                track.Ptr("Project2"),
						IsPrivate:           track.Ptr(true),
						Active:              track.Ptr(true),
						At:                  track.Ptr(time.Date(2021, time.February, 3, 4, 5, 6, 0, time.FixedZone("", 0))),
						CreatedAt:           track.Ptr(time.Date(2021, time.February, 3, 4, 5, 6, 0, time.FixedZone("", 0))),
						ServerDeletedAt:     nil,
						Color:               track.Ptr("#123456"),
						Billable:            nil,
						Template:            nil,
						AutoEstimates:       nil,
						EstimatedHours:      nil,
						Rate:                nil,
						RateLastUpdated:     nil,
						Currency:            nil,
						Recurring:           track.Ptr(false),
						RecurringParameters: nil,
						CurrentPeriod:       nil,
						FixedFee:            nil,
						ActualHours:         track.Ptr(0),
						WID:                 track.Ptr(1234567),
						CID:                 nil,
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
				testdataFile: "testdata/projects/get_workspace_projects_400_bad_request.json",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"Missing or invalid workspace_id\"\n",
					Header: http.Header{
						"Content-Length": []string{"34"},
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
				testdataFile: "testdata/projects/get_workspace_projects_403_forbidden",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
				testdataFile: "testdata/projects/get_workspace_projects_500_internal_server_error",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "projects")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			projects, err := client.GetWorkspaceProjects(context.Background(), workspaceID, &GetWorkspaceProjectsQuery{})

			if !reflect.DeepEqual(projects, tt.out.projects) {
				internal.Errorf(t, projects, tt.out.projects)
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

func TestGetWorkspaceProjectsQuery(t *testing.T) {
	tests := []struct {
		name string
		in   *GetWorkspaceProjectsQuery
		out  string
	}{
		{
			name: "GetWorkspaceProjectsQuery is nil",
			in:   nil,
			out:  "",
		},
		{
			name: "active=true",
			in:   &GetWorkspaceProjectsQuery{Active: track.Ptr(true)},
			out:  "active=true",
		},
		{
			name: "active=true&name=MyProject",
			in:   &GetWorkspaceProjectsQuery{Active: track.Ptr(true), Name: track.Ptr("MyProject")},
			out:  "active=true&name=MyProject",
		},
		{
			name: "active=true&name=MyProject&page=2",
			in:   &GetWorkspaceProjectsQuery{Active: track.Ptr(true), Name: track.Ptr("MyProject"), Page: track.Ptr(2)},
			out:  "active=true&name=MyProject&page=2",
		},
		{
			name: "GetWorkspaceProjectsQuery is empty",
			in:   &GetWorkspaceProjectsQuery{},
			out:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertQuery(t, tt.out)
			defer mockServer.Close()

			workspaceID := 1234567
			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			_, _ = client.GetWorkspaceProjects(context.Background(), workspaceID, tt.in)
		})
	}
}

func TestGetWorkspaceProject(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			project *Project
			err     error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/projects/get_workspace_project_200_ok.json",
			},
			out: struct {
				project *Project
				err     error
			}{
				project: &Project{
					ID:                  track.Ptr(123456789),
					WorkspaceID:         track.Ptr(1234567),
					ClientID:            nil,
					Name:                track.Ptr("MyProject"),
					IsPrivate:           track.Ptr(false),
					Active:              track.Ptr(true),
					At:                  track.Ptr(time.Date(2021, time.February, 3, 4, 5, 6, 0, time.FixedZone("", 0))),
					CreatedAt:           track.Ptr(time.Date(2021, time.February, 3, 4, 5, 6, 0, time.FixedZone("", 0))),
					ServerDeletedAt:     nil,
					Color:               track.Ptr("#456abc"),
					Billable:            nil,
					Template:            nil,
					AutoEstimates:       nil,
					EstimatedHours:      nil,
					Rate:                nil,
					RateLastUpdated:     nil,
					Currency:            nil,
					Recurring:           track.Ptr(false),
					RecurringParameters: nil,
					CurrentPeriod:       nil,
					FixedFee:            nil,
					ActualHours:         track.Ptr(0),
					WID:                 track.Ptr(1234567),
					CID:                 nil,
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
				testdataFile: "testdata/projects/get_workspace_project_400_bad_request.json",
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"We're expecting an integer as part of the url for project_id\"\n",
					Header: http.Header{
						"Content-Length": []string{"63"},
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
				testdataFile: "testdata/projects/get_workspace_project_403_forbidden",
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
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
				testdataFile: "testdata/projects/get_workspace_project_500_internal_server_error",
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
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
			projectID := 123456789
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "projects", strconv.Itoa(projectID))
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			project, err := client.GetWorkspaceProject(context.Background(), workspaceID, projectID, &GetWorkspaceProjectQuery{})

			if !reflect.DeepEqual(project, tt.out.project) {
				internal.Errorf(t, project, tt.out.project)
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

func TestGetWorkspaceProjectQuery(t *testing.T) {
	tests := []struct {
		name string
		in   *GetWorkspaceProjectQuery
		out  string
	}{
		{
			name: "GetWorkspaceProjectQuery is nil",
			in:   nil,
			out:  "",
		},
		{
			name: "with_first_time_entry=true",
			in:   &GetWorkspaceProjectQuery{WithFirstTimeEntry: track.Ptr(true)},
			out:  "with_first_time_entry=true",
		},
		{
			name: "GetWorkspaceProjectQuery is empty",
			in:   &GetWorkspaceProjectQuery{},
			out:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertQuery(t, tt.out)
			defer mockServer.Close()

			workspaceID := 1234567
			projectID := 123456789
			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			_, _ = client.GetWorkspaceProject(context.Background(), workspaceID, projectID, tt.in)
		})
	}
}
