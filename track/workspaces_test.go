package track

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

func TestGetWorkspace(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			workspace *Workspace
			err       error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/workspaces/get_workspace_200_ok.json",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: &Workspace{
					ID:                          Int(1234567),
					OrganizationID:              Int(2345678),
					Name:                        String("Workspace1"),
					Profile:                     Int(0),
					Premium:                     Bool(false),
					BusinessWs:                  Bool(false),
					Admin:                       Bool(true),
					SuspendedAt:                 nil,
					ServerDeletedAt:             nil,
					DefaultHourlyRate:           nil,
					RateLastUpdated:             nil,
					DefaultCurrency:             String("USD"),
					OnlyAdminsMayCreateProjects: Bool(false),
					OnlyAdminsMayCreateTags:     Bool(false),
					OnlyAdminsSeeBillableRates:  Bool(false),
					OnlyAdminsSeeTeamDashboard:  Bool(false),
					ProjectsBillableByDefault:   Bool(true),
					ReportsCollapse:             Bool(true),
					Rounding:                    Int(1),
					RoundingMinutes:             Int(0),
					APIToken:                    String("1234567890abcdefghijklmnopqrstuv"),
					At:                          Time(time.Date(2020, time.January, 23, 4, 5, 06, 0, time.FixedZone("", 0))),
					LogoURL:                     String("https://assets.toggl.com/images/workspace.jpg"),
					IcalURL:                     String("/ical/workspace_user/2345678901abcdefghijklmnopqrstuv"),
					IcalEnabled:                 Bool(true),
					CsvUpload:                   nil,
					Subscription:                nil,
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
				testdataFile: "testdata/workspaces/get_workspace_400_bad_request.json",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "\"Missing or invalid workspace_id\"\n",
					header: http.Header{
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
				testdataFile: "testdata/workspaces/get_workspace_403_forbidden",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
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
		{
			name: "500 Internal Server Error",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusInternalServerError,
				testdataFile: "testdata/workspaces/get_workspace_500_internal_server_error",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
				err: &errorResponse{
					statusCode: 500,
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
			workspaceID := 1234567
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID))
			mockServer := newMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspace, err := client.GetWorkspace(context.Background(), workspaceID)

			if !reflect.DeepEqual(workspace, tt.out.workspace) {
				errorf(t, workspace, tt.out.workspace)
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
