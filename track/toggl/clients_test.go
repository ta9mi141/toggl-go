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

func TestGetClients(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			clients []*Client
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
				testdataFile: "testdata/clients/get_clients_200_ok.json",
			},
			out: struct {
				clients []*Client
				err     error
			}{
				clients: []*Client{
					{
						ID:   track.Ptr(12345678),
						WID:  track.Ptr(1234567),
						Name: track.Ptr("test client"),
						At:   track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 0, time.FixedZone("", 0))),
					},
					{
						ID:   track.Ptr(23456789),
						WID:  track.Ptr(2345678),
						Name: track.Ptr("new client"),
						At:   track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 0, time.FixedZone("", 0))),
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
				testdataFile: "testdata/clients/get_clients_400_bad_request.json",
			},
			out: struct {
				clients []*Client
				err     error
			}{
				clients: nil,
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
				testdataFile: "testdata/clients/get_clients_403_forbidden",
			},
			out: struct {
				clients []*Client
				err     error
			}{
				clients: nil,
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
				testdataFile: "testdata/clients/get_clients_500_internal_server_error",
			},
			out: struct {
				clients []*Client
				err     error
			}{
				clients: nil,
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
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "clients")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			clients, err := apiClient.GetClients(context.Background(), workspaceID)

			if !reflect.DeepEqual(clients, tt.out.clients) {
				internal.Errorf(t, clients, tt.out.clients)
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

func TestGetClient(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			client *Client
			err    error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/clients/get_client_200_ok.json",
			},
			out: struct {
				client *Client
				err    error
			}{
				client: &Client{
					ID:   track.Ptr(12345678),
					WID:  track.Ptr(2345678),
					Name: track.Ptr("test client"),
					At:   track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 0, time.FixedZone("", 0))),
				},
				err: nil,
			},
		},
		{
			name: "403 Forbidden",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusForbidden,
				testdataFile: "testdata/clients/get_client_403_forbidden",
			},
			out: struct {
				client *Client
				err    error
			}{
				client: nil,
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
			name: "404 Not Found",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusNotFound,
				testdataFile: "testdata/clients/get_client_404_not_found.json",
			},
			out: struct {
				client *Client
				err    error
			}{
				client: nil,
				err: &internal.ErrorResponse{
					StatusCode: 404,
					Message:    "\"No client with ID 0 was found\"\n",
					Header: http.Header{
						"Content-Length": []string{"32"},
						"Content-Type":   []string{"application/json; charset=utf-8"},
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
				testdataFile: "testdata/clients/get_client_500_internal_server_error",
			},
			out: struct {
				client *Client
				err    error
			}{
				client: nil,
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
			clientID := 12345678
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "clients", strconv.Itoa(clientID))
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			client, err := apiClient.GetClient(context.Background(), workspaceID, clientID)

			if !reflect.DeepEqual(client, tt.out.client) {
				internal.Errorf(t, client, tt.out.client)
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
