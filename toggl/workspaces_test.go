package toggl

import (
	"context"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestGetWorkspaces(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		testdataFile string
		out          struct {
			workspaces []*Workspace
			err        error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/workspaces/get_workspaces_200_ok.json",
			out: struct {
				workspaces []*Workspace
				err        error
			}{
				workspaces: []*Workspace{
					{
						ID:                          3134975,
						Name:                        "John's personal ws",
						Profile:                     0,
						Premium:                     false,
						Admin:                       true,
						DefaultHourlyRate:           0,
						DefaultCurrency:             "USD",
						OnlyAdminsMayCreateProjects: false,
						OnlyAdminsSeeBillableRates:  false,
						OnlyAdminsSeeTeamDashboard:  false,
						ProjectsBillableByDefault:   true,
						Rounding:                    1,
						RoundingMinutes:             0,
						APIToken:                    "1234567890abcdefghijklmnopqrstuv",
						At:                          time.Date(2013, time.August, 28, 16, 22, 21, 0, time.FixedZone("", 0)),
						IcalEnabled:                 true,
					},
					{
						ID:                          7777777,
						Name:                        "My Company Inc",
						Profile:                     100,
						Premium:                     true,
						Admin:                       true,
						DefaultHourlyRate:           0,
						DefaultCurrency:             "USD",
						OnlyAdminsMayCreateProjects: false,
						OnlyAdminsSeeBillableRates:  false,
						OnlyAdminsSeeTeamDashboard:  false,
						ProjectsBillableByDefault:   true,
						Rounding:                    1,
						RoundingMinutes:             0,
						APIToken:                    "67890abcdefghijklmnopqrstuv12345",
						At:                          time.Date(2013, time.August, 28, 16, 22, 21, 0, time.FixedZone("", 0)),
						IcalEnabled:                 true,
						LogoURL:                     "https://assets.toggl.com/images/workspace.jpg",
						IcalURL:                     "/ical/workspace_user/abcdefghijklmn1234567890opqrstuv",
					},
				},
				err: nil,
			},
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/get_workspaces_403_forbidden.json",
			out: struct {
				workspaces []*Workspace
				err        error
			}{
				workspaces: nil,
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
			mockServer := newMockServer(t, "workspaces", tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaces, err := client.GetWorkspaces(context.Background())

			if !reflect.DeepEqual(workspaces, tt.out.workspaces) {
				errorf(t, workspaces, tt.out.workspaces)
			}
			if !reflect.DeepEqual(err, tt.out.err) {
				errorf(t, err, tt.out.err)
			}
		})
	}
}

func TestGetWorkspace(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		testdataFile string
		in           struct {
			id int
		}
		out struct {
			workspace *Workspace
			err       error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/workspaces/get_workspace_200_ok.json",
			in: struct {
				id int
			}{
				id: 3134975,
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: &Workspace{
					ID:                          3134975,
					Name:                        "John's personal ws",
					Profile:                     135,
					Premium:                     true,
					Admin:                       true,
					DefaultHourlyRate:           150,
					DefaultCurrency:             "USD",
					OnlyAdminsMayCreateProjects: false,
					OnlyAdminsSeeBillableRates:  false,
					OnlyAdminsSeeTeamDashboard:  false,
					ProjectsBillableByDefault:   true,
					Rounding:                    1,
					RoundingMinutes:             0,
					APIToken:                    "1234567890abcdefghijklmnopqrstuv",
					At:                          time.Date(2013, time.August, 28, 16, 22, 21, 0, time.FixedZone("", 3*60*60)),
					LogoURL:                     "my_logo.png",
					IcalURL:                     "/ical/workspace_user/9876543210abcdefghijklmnopqrstuv",
					IcalEnabled:                 true,
				},
				err: nil,
			},
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/get_workspace_403_forbidden.json",
			in: struct {
				id int
			}{
				id: 3134975,
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
			name:         "404 Not Found",
			statusCode:   http.StatusNotFound,
			testdataFile: "testdata/workspaces/get_workspace_404_not_found.json",
			in: struct {
				id int
			}{
				id: 1234567,
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
				err: &errorResponse{
					statusCode: 404,
					message:    "null",
					header: http.Header{
						"Content-Length": []string{"4"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiSpecificPath := path.Join("workspaces", strconv.Itoa(tt.in.id))
			mockServer := newMockServer(t, apiSpecificPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspace, err := client.GetWorkspace(context.Background(), tt.in.id)

			if !reflect.DeepEqual(workspace, tt.out.workspace) {
				errorf(t, workspace, tt.out.workspace)
			}
			if !reflect.DeepEqual(err, tt.out.err) {
				errorf(t, err, tt.out.err)
			}
		})
	}
}
