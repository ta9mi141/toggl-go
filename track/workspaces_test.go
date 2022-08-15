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
					ID:                          Ptr(1234567),
					OrganizationID:              Ptr(2345678),
					Name:                        Ptr("Workspace1"),
					Profile:                     Ptr(0),
					Premium:                     Ptr(false),
					BusinessWs:                  Ptr(false),
					Admin:                       Ptr(true),
					SuspendedAt:                 nil,
					ServerDeletedAt:             nil,
					DefaultHourlyRate:           nil,
					RateLastUpdated:             nil,
					DefaultCurrency:             Ptr("USD"),
					OnlyAdminsMayCreateProjects: Ptr(false),
					OnlyAdminsMayCreateTags:     Ptr(false),
					OnlyAdminsSeeBillableRates:  Ptr(false),
					OnlyAdminsSeeTeamDashboard:  Ptr(false),
					ProjectsBillableByDefault:   Ptr(true),
					ReportsCollapse:             Ptr(true),
					Rounding:                    Ptr(1),
					RoundingMinutes:             Ptr(0),
					APIToken:                    Ptr("1234567890abcdefghijklmnopqrstuv"),
					At:                          Ptr(time.Date(2020, time.January, 23, 4, 5, 06, 0, time.FixedZone("", 0))),
					LogoURL:                     Ptr("https://assets.toggl.com/images/workspace.jpg"),
					IcalURL:                     Ptr("/ical/workspace_user/2345678901abcdefghijklmnopqrstuv"),
					IcalEnabled:                 Ptr(true),
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
						ID:                Ptr(1234567),
						UserID:            Ptr(2345678),
						WorkspaceID:       Ptr(3456789),
						Admin:             Ptr(true),
						OrganizationAdmin: Ptr(true),
						WorkspaceAdmin:    Ptr(true),
						Active:            Ptr(true),
						Email:             Ptr("example@toggl.com"),
						Timezone:          Ptr("Asia/Tokyo"),
						Inactive:          Ptr(false),
						At:                Ptr(time.Date(2020, time.January, 23, 4, 56, 7, 0, time.FixedZone("", 0))),
						Name:              Ptr("Toggl Track"),
						Rate:              nil,
						RateLastUpdated:   nil,
						LabourCost:        nil,
						InviteURL:         nil,
						InvitationCode:    nil,
						AvatarFileName:    nil,
						GroupIDs:          nil,
						IsDirect:          Ptr(true),
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
				testdataFile: "testdata/workspaces/get_workspace_users_403_forbidden",
			},
			out: struct {
				workspaceUsers []*WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
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
				testdataFile: "testdata/workspaces/get_workspace_users_500_internal_server_error",
			},
			out: struct {
				workspaceUsers []*WorkspaceUser
				err            error
			}{
				workspaceUsers: nil,
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
			organizationID := 1234567
			workspaceID := 2345678
			apiSpecificPath := path.Join(organizationsPath, strconv.Itoa(organizationID), "workspaces", strconv.Itoa(workspaceID))
			mockServer := newMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaceUsers, err := client.GetWorkspaceUsers(context.Background(), organizationID, workspaceID)

			if !reflect.DeepEqual(workspaceUsers, tt.out.workspaceUsers) {
				errorf(t, workspaceUsers, tt.out.workspaceUsers)
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
					ID:                          Ptr(1234567),
					OrganizationID:              Ptr(2345678),
					Name:                        Ptr("Updated Workspace"),
					Profile:                     Ptr(0),
					Premium:                     Ptr(false),
					BusinessWs:                  Ptr(false),
					Admin:                       Ptr(true),
					SuspendedAt:                 nil,
					ServerDeletedAt:             nil,
					DefaultHourlyRate:           nil,
					RateLastUpdated:             nil,
					DefaultCurrency:             Ptr("USD"),
					OnlyAdminsMayCreateProjects: Ptr(false),
					OnlyAdminsMayCreateTags:     Ptr(false),
					OnlyAdminsSeeBillableRates:  Ptr(false),
					OnlyAdminsSeeTeamDashboard:  Ptr(false),
					ProjectsBillableByDefault:   Ptr(true),
					ReportsCollapse:             Ptr(true),
					Rounding:                    Ptr(1),
					RoundingMinutes:             Ptr(0),
					APIToken:                    Ptr("1234567890abcdefghijklmnopqrstuv"),
					At:                          Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 0, time.FixedZone("", 0))),
					LogoURL:                     Ptr("https://assets.toggl.com/images/workspace.jpg"),
					IcalURL:                     Ptr("/ical/workspace_user/abcdefghijklmnopqrstuvwxyz012345"),
					IcalEnabled:                 Ptr(true),
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
				err: &errorResponse{
					statusCode: 400,
					message:    "\"JSON is not valid\"\n",
					header: http.Header{
						"Content-Length": []string{"20"},
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
				testdataFile: "testdata/workspaces/update_workspace_403_forbidden.txt",
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
				err: &errorResponse{
					statusCode: 403,
					message:    "Incorrect username and/or password",
					header: http.Header{
						"Content-Length": []string{"34"},
						"Content-Type":   []string{"text/plain; charset=utf-8"},
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
				testdataFile: "testdata/workspaces/update_workspace_500_internal_server_error",
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
			workspace, err := client.UpdateWorkspace(context.Background(), workspaceID, &UpdateWorkspaceRequestBody{})

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

func TestUpdateWorkspaceRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *UpdateWorkspaceRequestBody
		out  string
	}{
		{
			name: "int, string, and bool",
			in: &UpdateWorkspaceRequestBody{
				InitialPricingPlan:          Ptr(1),
				Name:                        Ptr("Updated Workspace"),
				OnlyAdminsMayCreateProjects: Ptr(true),
			},
			out: "{\"initial_pricing_plan\":1,\"name\":\"Updated Workspace\",\"only_admins_may_create_projects\":true}",
		},
		{
			name: "int, string, bool, and slice of int",
			in: &UpdateWorkspaceRequestBody{
				Admins:                      []*int{Ptr(1234567), Ptr(2345678)},
				InitialPricingPlan:          Ptr(1),
				Name:                        Ptr("Updated Workspace"),
				OnlyAdminsMayCreateProjects: Ptr(true),
			},
			out: "{\"admins\":[1234567,2345678],\"initial_pricing_plan\":1,\"name\":\"Updated Workspace\",\"only_admins_may_create_projects\":true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaceID := 1234567
			_, _ = client.UpdateWorkspace(context.Background(), workspaceID, tt.in)
		})
	}
}
