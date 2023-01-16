package reports

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

func TestListProjects(t *testing.T) {
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
				testdataFile: "testdata/utils/list_projects_200_ok.json",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: []*Project{
					{
						ID:       track.Ptr(12345678),
						Name:     track.Ptr("Project1"),
						ClientID: nil,
						Color:    track.Ptr("#000000"),
						Active:   track.Ptr(true),
						Currency: nil,
						Billable: nil,
					},
					{
						ID:       track.Ptr(23456789),
						Name:     track.Ptr("Project2"),
						ClientID: nil,
						Color:    track.Ptr("#ffffff"),
						Active:   track.Ptr(true),
						Currency: nil,
						Billable: nil,
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
				testdataFile: "testdata/utils/list_projects_400_bad_request.json",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
			name: "401 Unauthorized",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusUnauthorized,
				testdataFile: "testdata/utils/list_projects_401_unauthorized",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
				err: &internal.ErrorResponse{
					StatusCode: 401,
					Message:    "",
					Header: http.Header{
						"Content-Length": []string{"0"},
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
				testdataFile: "testdata/utils/list_projects_403_forbidden.txt",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "Incorrect username and/or password\n",
					Header: http.Header{
						"Content-Length": []string{"35"},
						"Content-Type":   []string{"text/plain; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceID := 1234567
			apiSpecificPath := path.Join(reportsPath, strconv.Itoa(workspaceID), "filters/projects")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			projects, err := apiClient.ListProjects(context.Background(), workspaceID, &ListProjectsRequestBody{})

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

func TestListProjectsRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *ListProjectsRequestBody
		out  string
	}{
		{
			name: "string",
			in: &ListProjectsRequestBody{
				Name: track.Ptr("Project1"),
			},
			out: "{\"name\":\"Project1\"}",
		},
		{
			name: "bool and string",
			in: &ListProjectsRequestBody{
				IsActive: track.Ptr(true),
				Name:     track.Ptr("Project1"),
			},
			out: "{\"is_active\":true,\"name\":\"Project1\"}",
		},
		{
			name: "bool, string, and array of integer",
			in: &ListProjectsRequestBody{
				IDs:      []*int{track.Ptr(12345678), track.Ptr(23456789)},
				IsActive: track.Ptr(true),
				Name:     track.Ptr("Project1"),
			},
			out: "{\"ids\":[12345678,23456789],\"is_active\":true,\"name\":\"Project1\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			workspaceID := 1234567
			_, _ = apiClient.ListProjects(context.Background(), workspaceID, tt.in)
		})
	}
}
