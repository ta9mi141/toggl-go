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
					ID:                          track.Ptr(1234567),
					OrganizationID:              track.Ptr(2345678),
					Name:                        track.Ptr("Workspace1"),
					Profile:                     track.Ptr(0),
					Premium:                     track.Ptr(false),
					BusinessWs:                  track.Ptr(false),
					Admin:                       track.Ptr(true),
					SuspendedAt:                 nil,
					ServerDeletedAt:             nil,
					DefaultHourlyRate:           nil,
					RateLastUpdated:             nil,
					DefaultCurrency:             track.Ptr("USD"),
					OnlyAdminsMayCreateProjects: track.Ptr(false),
					OnlyAdminsMayCreateTags:     track.Ptr(false),
					OnlyAdminsSeeBillableRates:  track.Ptr(false),
					OnlyAdminsSeeTeamDashboard:  track.Ptr(false),
					ProjectsBillableByDefault:   track.Ptr(true),
					ReportsCollapse:             track.Ptr(true),
					Rounding:                    track.Ptr(1),
					RoundingMinutes:             track.Ptr(0),
					APIToken:                    track.Ptr("1234567890abcdefghijklmnopqrstuv"),
					At:                          track.Ptr(time.Date(2020, time.January, 23, 4, 5, 06, 0, time.Local)),
					LogoURL:                     track.Ptr("https://assets.toggl.com/images/workspace.jpg"),
					IcalURL:                     track.Ptr("/ical/workspace_user/2345678901abcdefghijklmnopqrstuv"),
					IcalEnabled:                 track.Ptr(true),
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
			name: "401 Unauthorized",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusUnauthorized,
				testdataFile: "testdata/workspaces/get_workspace_401_unauthorized",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
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
				testdataFile: "testdata/workspaces/get_workspace_403_forbidden",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceID := 1234567
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID))
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			workspace, err := apiClient.GetWorkspace(context.Background(), workspaceID)

			if !reflect.DeepEqual(workspace, tt.out.workspace) {
				internal.Errorf(t, workspace, tt.out.workspace)
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

func TestGetWorkspaceUsers(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			workspaceUsers []*WorkspaceUser
			err            error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/workspaces/get_workspace_users_200_ok.json",
			},
			out: struct {
				workspaceUsers []*WorkspaceUser
				err            error
			}{
				workspaceUsers: []*WorkspaceUser{
					{
						ID:                track.Ptr(1234567),
						UserID:            track.Ptr(2345678),
						WorkspaceID:       track.Ptr(3456789),
						Admin:             track.Ptr(true),
						OrganizationAdmin: track.Ptr(true),
						WorkspaceAdmin:    track.Ptr(true),
						Active:            track.Ptr(true),
						Email:             track.Ptr("example@toggl.com"),
						Timezone:          track.Ptr("Asia/Tokyo"),
						Inactive:          track.Ptr(false),
						At:                track.Ptr(time.Date(2020, time.January, 23, 4, 56, 7, 0, time.Local)),
						Name:              track.Ptr("Toggl Track"),
						Rate:              nil,
						RateLastUpdated:   nil,
						LabourCost:        nil,
						InviteURL:         nil,
						InvitationCode:    nil,
						AvatarFileName:    nil,
						GroupIDs:          nil,
						IsDirect:          track.Ptr(true),
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
				testdataFile: "testdata/workspaces/get_workspace_users_400_bad_request.json",
			},
			out: struct {
				workspaceUsers []*WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
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
			name: "401 Unauthorized",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusUnauthorized,
				testdataFile: "testdata/workspaces/get_workspace_users_401_unauthorized",
			},
			out: struct {
				workspaceUsers []*WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
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
				testdataFile: "testdata/workspaces/get_workspace_users_403_forbidden",
			},
			out: struct {
				workspaceUsers []*WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			organizationID := 1234567
			workspaceID := 2345678
			apiSpecificPath := path.Join(organizationsPath, strconv.Itoa(organizationID), "workspaces", strconv.Itoa(workspaceID))
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			workspaceUsers, err := apiClient.GetWorkspaceUsers(context.Background(), organizationID, workspaceID)

			if !reflect.DeepEqual(workspaceUsers, tt.out.workspaceUsers) {
				internal.Errorf(t, workspaceUsers, tt.out.workspaceUsers)
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

func TestUpdateWorkspace(t *testing.T) {
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
				testdataFile: "testdata/workspaces/update_workspace_200_ok.json",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: &Workspace{
					ID:                          track.Ptr(1234567),
					OrganizationID:              track.Ptr(2345678),
					Name:                        track.Ptr("Updated Workspace"),
					Profile:                     track.Ptr(0),
					Premium:                     track.Ptr(false),
					BusinessWs:                  track.Ptr(false),
					Admin:                       track.Ptr(true),
					SuspendedAt:                 nil,
					ServerDeletedAt:             nil,
					DefaultHourlyRate:           nil,
					RateLastUpdated:             nil,
					DefaultCurrency:             track.Ptr("USD"),
					OnlyAdminsMayCreateProjects: track.Ptr(false),
					OnlyAdminsMayCreateTags:     track.Ptr(false),
					OnlyAdminsSeeBillableRates:  track.Ptr(false),
					OnlyAdminsSeeTeamDashboard:  track.Ptr(false),
					ProjectsBillableByDefault:   track.Ptr(true),
					ReportsCollapse:             track.Ptr(true),
					Rounding:                    track.Ptr(1),
					RoundingMinutes:             track.Ptr(0),
					APIToken:                    track.Ptr("1234567890abcdefghijklmnopqrstuv"),
					At:                          track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 0, time.Local)),
					LogoURL:                     track.Ptr("https://assets.toggl.com/images/workspace.jpg"),
					IcalURL:                     track.Ptr("/ical/workspace_user/abcdefghijklmnopqrstuvwxyz012345"),
					IcalEnabled:                 track.Ptr(true),
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
				testdataFile: "testdata/workspaces/update_workspace_400_bad_request.json",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
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
				testdataFile: "testdata/workspaces/update_workspace_401_unauthorized",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
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
				testdataFile: "testdata/workspaces/update_workspace_403_forbidden.txt",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
				err: &internal.ErrorResponse{
					StatusCode: 403,
					Message:    "Incorrect username and/or password",
					Header: http.Header{
						"Content-Length": []string{"34"},
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
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID))
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			workspace, err := apiClient.UpdateWorkspace(context.Background(), workspaceID, &UpdateWorkspaceRequestBody{})

			if !reflect.DeepEqual(workspace, tt.out.workspace) {
				internal.Errorf(t, workspace, tt.out.workspace)
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

func TestUpdateWorkspaceRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *UpdateWorkspaceRequestBody
		out  string
	}{
		{
			name: "int, string, and bool",
			in: &UpdateWorkspaceRequestBody{
				InitialPricingPlan:          track.Ptr(1),
				Name:                        track.Ptr("Updated Workspace"),
				OnlyAdminsMayCreateProjects: track.Ptr(true),
			},
			out: "{\"initial_pricing_plan\":1,\"name\":\"Updated Workspace\",\"only_admins_may_create_projects\":true}",
		},
		{
			name: "int, string, bool, and slice of int",
			in: &UpdateWorkspaceRequestBody{
				Admins:                      []*int{track.Ptr(1234567), track.Ptr(2345678)},
				InitialPricingPlan:          track.Ptr(1),
				Name:                        track.Ptr("Updated Workspace"),
				OnlyAdminsMayCreateProjects: track.Ptr(true),
			},
			out: "{\"admins\":[1234567,2345678],\"initial_pricing_plan\":1,\"name\":\"Updated Workspace\",\"only_admins_may_create_projects\":true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			workspaceID := 1234567
			_, _ = apiClient.UpdateWorkspace(context.Background(), workspaceID, tt.in)
		})
	}
}
