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

func TestGetTags(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			tags []*Tag
			err  error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/tags/get_tags_200_ok.json",
			},
			out: struct {
				tags []*Tag
				err  error
			}{
				tags: []*Tag{
					{
						ID:          track.Ptr(12345678),
						WorkspaceID: track.Ptr(1234567),
						Name:        track.Ptr("toggl-go"),
						At:          track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 678901000, time.UTC)),
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
				testdataFile: "testdata/tags/get_tags_400_bad_request.json",
			},
			out: struct {
				tags []*Tag
				err  error
			}{
				tags: nil,
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
				testdataFile: "testdata/tags/get_tags_403_forbidden",
			},
			out: struct {
				tags []*Tag
				err  error
			}{
				tags: nil,
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
				testdataFile: "testdata/tags/get_tags_500_internal_server_error",
			},
			out: struct {
				tags []*Tag
				err  error
			}{
				tags: nil,
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
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "tags")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			tags, err := client.GetTags(context.Background(), workspaceID)

			if !reflect.DeepEqual(tags, tt.out.tags) {
				internal.Errorf(t, tags, tt.out.tags)
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

func TestCreateTag(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			tag *Tag
			err error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/tags/create_tag_200_ok.json",
			},
			out: struct {
				tag *Tag
				err error
			}{
				tag: &Tag{
					ID:          track.Ptr(12345678),
					WorkspaceID: track.Ptr(1234567),
					Name:        track.Ptr("toggl-go"),
					At:          track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 678900000, time.UTC)),
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
				testdataFile: "testdata/tags/create_tag_400_bad_request.json",
			},
			out: struct {
				tag *Tag
				err error
			}{
				tag: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"tag name can't be blank\"\n",
					Header: http.Header{
						"Content-Length": []string{"26"},
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
				testdataFile: "testdata/tags/create_tag_403_forbidden",
			},
			out: struct {
				tag *Tag
				err error
			}{
				tag: nil,
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
				testdataFile: "testdata/tags/create_tag_500_internal_server_error",
			},
			out: struct {
				tag *Tag
				err error
			}{
				tag: nil,
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
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "tags")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			tag, err := client.CreateTag(context.Background(), workspaceID, &CreateTagRequestBody{})

			if !reflect.DeepEqual(tag, tt.out.tag) {
				internal.Errorf(t, tag, tt.out.tag)
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

func TestCreateTagRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *CreateTagRequestBody
		out  string
	}{
		{
			name: "string",
			in: &CreateTagRequestBody{
				Name: track.Ptr("toggl-go"),
			},
			out: "{\"name\":\"toggl-go\"}",
		},
		{
			name: "string and int",
			in: &CreateTagRequestBody{
				Name:        track.Ptr("toggl-go"),
				WorkspaceID: track.Ptr(1234567),
			},
			out: "{\"name\":\"toggl-go\",\"workspace_id\":1234567}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			workspaceID := 1234567
			_, _ = client.CreateTag(context.Background(), workspaceID, tt.in)
		})
	}
}
