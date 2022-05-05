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
