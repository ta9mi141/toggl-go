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
			mockServer := newMockServer(t, workspacesPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaces, err := client.GetWorkspaces(context.Background())

			if !reflect.DeepEqual(workspaces, tt.out.workspaces) {
				errorf(t, workspaces, tt.out.workspaces)
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
					message:    "null\n",
					header: http.Header{
						"Content-Length": []string{"5"},
						"Content-Type":   []string{"text/plain; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(tt.in.id))
			mockServer := newMockServer(t, apiSpecificPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspace, err := client.GetWorkspace(context.Background(), tt.in.id)

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
		name         string
		statusCode   int
		testdataFile string
		in           struct {
			id int
		}
		out struct {
			users []*User
			err   error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/workspaces/get_workspace_users_200_ok.json",
			in: struct {
				id int
			}{
				id: 777,
			},
			out: struct {
				users []*User
				err   error
			}{
				users: []*User{
					{
						ID:                    123123,
						DefaultWID:            777,
						Email:                 "john@swift.com",
						Fullname:              "John Swift",
						JqueryTimeofdayFormat: "h:i A",
						JqueryDateFormat:      "m/d/Y",
						TimeofdayFormat:       "h:mm A",
						DateFormat:            "MM/DD/YYYY",
						StoreStartAndStopTime: true,
						BeginningOfWeek:       0,
						Language:              "en_US",
						ImageURL:              "https://www.toggl.com/system/avatars/123123/small/open-uri20121116-2767-b1qr8l.png",
						SidebarPiechart:       false,
						At:                    time.Date(2013, time.March, 6, 8, 57, 12, 0, time.FixedZone("", 0)),
						CreatedAt:             time.Date(2013, time.March, 6, 8, 57, 12, 0, time.FixedZone("", 0)),
						Retention:             9,
						RecordTimeline:        true,
						RenderTimeline:        true,
						TimelineEnabled:       true,
						TimelineExperiment:    true,
						ShouldUpgrade:         true,
						Timezone:              "Etc/UTC",
						OpenIDEnabled:         false,
						SendProductEmails:     true,
						SendWeeklyReport:      true,
						SendTimeNotifications: true,
						Invitation:            struct{}{},
						DurationFormat:        "improved",
					},
					{
						ID:                    321321,
						Email:                 "Happy@worker.com",
						Fullname:              "Happy Worker",
						JqueryTimeofdayFormat: "h:i A",
						JqueryDateFormat:      "m/d/Y",
						TimeofdayFormat:       "h:mm A",
						DateFormat:            "MM/DD/YYYY",
						StoreStartAndStopTime: true,
						BeginningOfWeek:       1,
						Language:              "en_US",
						ImageURL:              "https://www.toggl.com/images/profile.png",
						SidebarPiechart:       false,
						At:                    time.Date(2013, time.March, 6, 8, 46, 7, 0, time.FixedZone("", 0)),
						CreatedAt:             time.Date(2013, time.March, 6, 7, 52, 3, 0, time.FixedZone("", 0)),
						Retention:             0,
						RecordTimeline:        true,
						RenderTimeline:        true,
						TimelineEnabled:       true,
						TimelineExperiment:    true,
						ShouldUpgrade:         true,
						Timezone:              "Etc/UTC",
						OpenIDEnabled:         false,
						SendProductEmails:     true,
						SendWeeklyReport:      true,
						SendTimeNotifications: true,
						Invitation:            struct{}{},
						DurationFormat:        "improved",
					},
				},
				err: nil,
			},
		},
		{
			name:         "400 Bad Request",
			statusCode:   http.StatusBadRequest,
			testdataFile: "testdata/workspaces/get_workspace_users_400_bad_request.txt",
			in: struct {
				id int
			}{
				id: 777,
			},
			out: struct {
				users []*User
				err   error
			}{
				users: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "Missing or invalid workspace_id\n",
					header: http.Header{
						"Content-Length": []string{"32"},
						"Content-Type":   []string{"text/plain; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/get_workspace_users_403_forbidden.json",
			in: struct {
				id int
			}{
				id: 777,
			},
			out: struct {
				users []*User
				err   error
			}{
				users: nil,
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
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(tt.in.id), "users")
			mockServer := newMockServer(t, apiSpecificPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			users, err := client.GetWorkspaceUsers(context.Background(), tt.in.id)

			if !reflect.DeepEqual(users, tt.out.users) {
				errorf(t, users, tt.out.users)
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
		name         string
		statusCode   int
		testdataFile string
		in           struct {
			id        int
			workspace *Workspace
		}
		out struct {
			workspace *Workspace
			err       error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/workspaces/update_200_ok.json",
			in: struct {
				id        int
				workspace *Workspace
			}{
				id: 3134975,
				workspace: &Workspace{
					Name: "John's ws",
				},
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: &Workspace{
					ID:                          3134975,
					Name:                        "John's ws",
					Profile:                     0,
					Premium:                     true,
					Admin:                       true,
					DefaultHourlyRate:           50,
					DefaultCurrency:             "USD",
					OnlyAdminsMayCreateProjects: false,
					OnlyAdminsSeeBillableRates:  true,
					OnlyAdminsSeeTeamDashboard:  false,
					ProjectsBillableByDefault:   true,
					Rounding:                    1,
					RoundingMinutes:             60,
					APIToken:                    "1234567890abcdefghijklmnopqrstuv",
					At:                          time.Date(2013, time.August, 28, 16, 22, 21, 0, time.FixedZone("", 3*60*60)),
					LogoURL:                     "my_logo.png",
					IcalEnabled:                 true,
				},
				err: nil,
			},
		},
		{
			name:         "400 Bad Request",
			statusCode:   http.StatusBadRequest,
			testdataFile: "testdata/workspaces/update_400_bad_request.txt",
			in: struct {
				id        int
				workspace *Workspace
			}{
				id: 3134975,
				workspace: &Workspace{
					Name: "John's ws",
				},
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "workspace missing from json structure\n",
					header: http.Header{
						"Content-Length": []string{"38"},
						"Content-Type":   []string{"text/plain; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/update_403_forbidden.json",
			in: struct {
				id        int
				workspace *Workspace
			}{
				id: 3134975,
				workspace: &Workspace{
					Name: "John's ws",
				},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(tt.in.id))
			mockServer := newMockServer(t, apiSpecificPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspace, err := client.UpdateWorkspace(context.Background(), tt.in.id, tt.in.workspace)

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
		in   *Workspace
		out  *workspaceRequest
	}{
		{
			name: "string",
			in: &Workspace{
				Name: "John's ws",
			},
			out: &workspaceRequest{
				Workspace: Workspace{
					Name: "John's ws",
				},
			},
		},
		{
			name: "string and float64",
			in: &Workspace{
				Name:              "John's ws",
				DefaultHourlyRate: 50,
			},
			out: &workspaceRequest{
				Workspace: Workspace{
					Name:              "John's ws",
					DefaultHourlyRate: 50,
				},
			},
		},
		{
			name: "string, float64, and bool",
			in: &Workspace{
				Name:                        "John's ws",
				DefaultHourlyRate:           50,
				OnlyAdminsMayCreateProjects: false,
			},
			out: &workspaceRequest{
				Workspace: Workspace{
					Name:                        "John's ws",
					DefaultHourlyRate:           50,
					OnlyAdminsMayCreateProjects: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestBody(t, new(workspaceRequest), tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaceID := 3134975
			_, _ = client.UpdateWorkspace(context.Background(), workspaceID, tt.in)
		})
	}
}
