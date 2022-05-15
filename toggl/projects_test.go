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
)

func TestCreateProject(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		testdataFile string
		in           struct {
			project *Project
		}
		out struct {
			project *Project
			err     error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/projects/create_200_ok.json",
			in: struct {
				project *Project
			}{
				project: &Project{
					WID:  Int(2345678),
					Name: String("An awesome project"),
				},
			},
			out: struct {
				project *Project
				err     error
			}{
				project: &Project{
					ID:            Int(123456789),
					WID:           Int(2345678),
					Name:          String("An awesome project"),
					Billable:      Bool(false),
					IsPrivate:     Bool(true),
					Active:        Bool(true),
					Template:      Bool(false),
					At:            Time(time.Date(2021, time.April, 28, 1, 23, 45, 0, time.FixedZone("", 0))),
					Color:         String("2"),
					AutoEstimates: Bool(false),
					HexColor:      String("#d94182"),
				},
				err: nil,
			},
		},
		{
			name:         "400 Bad Request",
			statusCode:   http.StatusBadRequest,
			testdataFile: "testdata/projects/create_400_bad_request.json",
			in: struct {
				project *Project
			}{
				project: &Project{},
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "\"Project name must be present\"\n",
					header: http.Header{
						"Content-Length": []string{"31"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/projects/create_403_forbidden",
			in: struct {
				project *Project
			}{
				project: &Project{
					Name: String("An awesome project"),
					WID:  Int(2345678),
				},
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServer(t, projectsPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			project, err := client.CreateProject(context.Background(), tt.in.project)

			if !reflect.DeepEqual(project, tt.out.project) {
				errorf(t, project, tt.out.project)
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

func TestCreateProjectRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *Project
		out  string
	}{
		{
			name: "int and string",
			in: &Project{
				WID:  Int(2345678),
				Name: String("An awesome project"),
			},
			out: "{\"project\":{\"wid\":2345678,\"name\":\"An awesome project\"}}",
		},
		{
			name: "int and string and bool (true)",
			in: &Project{
				WID:       Int(2345678),
				Name:      String("An awesome project"),
				IsPrivate: Bool(true),
			},
			out: "{\"project\":{\"wid\":2345678,\"name\":\"An awesome project\",\"is_private\":true}}",
		},
		{
			name: "int and string and bool (false)",
			in: &Project{
				WID:       Int(2345678),
				Name:      String("An awesome project"),
				IsPrivate: Bool(false),
			},
			out: "{\"project\":{\"wid\":2345678,\"name\":\"An awesome project\",\"is_private\":false}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			_, _ = client.CreateProject(context.Background(), tt.in)
		})
	}
}

func TestGetProject(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		testdataFile string
		in           struct {
			id int
		}
		out struct {
			project *Project
			err     error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/projects/get_200_ok.json",
			in: struct {
				id int
			}{
				id: 123456789,
			},
			out: struct {
				project *Project
				err     error
			}{
				project: &Project{
					ID:            Int(123456789),
					WID:           Int(2345678),
					CID:           Int(34567890),
					Name:          String("Project1"),
					Billable:      Bool(false),
					IsPrivate:     Bool(true),
					Active:        Bool(true),
					Template:      Bool(false),
					At:            Time(time.Date(2022, time.April, 29, 1, 23, 45, 0, time.FixedZone("", 0))),
					CreatedAt:     Time(time.Date(2020, time.September, 13, 5, 43, 21, 0, time.FixedZone("", 0))),
					Color:         String("6"),
					AutoEstimates: Bool(false),
					ActualHours:   Int(3),
					HexColor:      String("#06a893"),
				},
				err: nil,
			},
		},
		{
			name:         "400 Bad Request",
			statusCode:   http.StatusBadRequest,
			testdataFile: "testdata/projects/get_400_bad_request.json",
			in: struct {
				id int
			}{
				id: 123456789,
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "\"Invalid project_id\"\n",
					header: http.Header{
						"Content-Length": []string{"21"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name:         "404 Not Found",
			statusCode:   http.StatusNotFound,
			testdataFile: "testdata/projects/get_404_not_found.json",
			in: struct {
				id int
			}{
				id: 123456789,
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
				err: &errorResponse{
					statusCode: 404,
					message:    "\"\"\n",
					header: http.Header{
						"Content-Length": []string{"3"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiSpecificPath := path.Join(projectsPath, strconv.Itoa(tt.in.id))
			mockServer := newMockServer(t, apiSpecificPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaces, err := client.GetProject(context.Background(), tt.in.id)

			if !reflect.DeepEqual(workspaces, tt.out.project) {
				errorf(t, workspaces, tt.out.project)
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

func TestUpdateProject(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		testdataFile string
		in           struct {
			id      int
			project *Project
		}
		out struct {
			project *Project
			err     error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/projects/update_200_ok.json",
			in: struct {
				id      int
				project *Project
			}{
				id: 123456789,
				project: &Project{
					WID:  Int(1234567),
					Name: String("Changed the name"),
				},
			},
			out: struct {
				project *Project
				err     error
			}{
				project: &Project{
					ID:            Int(123456789),
					WID:           Int(1234567),
					Name:          String("Changed the name"),
					Billable:      Bool(false),
					IsPrivate:     Bool(true),
					Active:        Bool(true),
					Template:      Bool(false),
					At:            Time(time.Date(2021, time.April, 13, 1, 23, 45, 0, time.FixedZone("", 0))),
					CreatedAt:     Time(time.Date(2021, time.March, 13, 1, 23, 45, 0, time.FixedZone("", 0))),
					Color:         String("2"),
					AutoEstimates: Bool(false),
					HexColor:      String("#d94182"),
				},
				err: nil,
			},
		},
		{
			name:         "400 Bad Request",
			statusCode:   http.StatusBadRequest,
			testdataFile: "testdata/projects/update_400_bad_request.json",
			in: struct {
				id      int
				project *Project
			}{
				id:      123456789,
				project: &Project{},
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "\"Project name must be present\"\n",
					header: http.Header{
						"Content-Length": []string{"31"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name:         "404 Not Found",
			statusCode:   http.StatusNotFound,
			testdataFile: "testdata/projects/update_404_not_found.json",
			in: struct {
				id      int
				project *Project
			}{
				id: 123456789,
				project: &Project{
					Name: String("Changed the name"),
					WID:  Int(1234567),
				},
			},
			out: struct {
				project *Project
				err     error
			}{
				project: nil,
				err: &errorResponse{
					statusCode: 404,
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
			mockServer := newMockServer(t, projectsPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			project, err := client.UpdateProject(context.Background(), tt.in.id, tt.in.project)

			if !reflect.DeepEqual(project, tt.out.project) {
				errorf(t, project, tt.out.project)
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

func TestUpdateProjectRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *Project
		out  string
	}{
		{
			name: "int and string",
			in: &Project{
				WID:  Int(2345678),
				Name: String("Changed the name"),
			},
			out: "{\"project\":{\"wid\":2345678,\"name\":\"Changed the name\"}}",
		},
		{
			name: "int and string and bool (true)",
			in: &Project{
				WID:       Int(2345678),
				Name:      String("Changed the name"),
				IsPrivate: Bool(true),
			},
			out: "{\"project\":{\"wid\":2345678,\"name\":\"Changed the name\",\"is_private\":true}}",
		},
		{
			name: "int and string and bool (false)",
			in: &Project{
				WID:       Int(2345678),
				Name:      String("Changed the name"),
				IsPrivate: Bool(false),
			},
			out: "{\"project\":{\"wid\":2345678,\"name\":\"Changed the name\",\"is_private\":false}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			projectID := 123456789
			_, _ = client.UpdateProject(context.Background(), projectID, tt.in)
		})
	}
}
