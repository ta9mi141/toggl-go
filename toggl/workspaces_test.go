package toggl

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetWorkspaces(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		testdataFile string
		in           struct {
			ctx context.Context
		}
		out struct {
			workspaces []*Workspace
			err        error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/workspaces/get_workspaces_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
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
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				workspaces []*Workspace
				err        error
			}{
				workspaces: nil,
				err:        ErrAuthenticationFailure,
			},
		},
		{
			name:         "Without context",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/workspaces/get_workspaces_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: nil,
			},
			out: struct {
				workspaces []*Workspace
				err        error
			}{
				workspaces: nil,
				err:        ErrContextNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServer(t, workspacesEndpoint, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaces, err := client.GetWorkspaces(tt.in.ctx)

			if !reflect.DeepEqual(workspaces, tt.out.workspaces) {
				errorf(t, workspaces, tt.out.workspaces)
			}
			if !errors.Is(err, tt.out.err) {
				errorf(t, err, tt.out.err)
			}
		})
	}
}
