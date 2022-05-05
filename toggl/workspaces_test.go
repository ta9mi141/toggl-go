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
						ID:                          Int(3134975),
						Name:                        String("John's personal ws"),
						Profile:                     Int(0),
						Premium:                     Bool(false),
						Admin:                       Bool(true),
						DefaultHourlyRate:           Float64(0),
						DefaultCurrency:             String("USD"),
						OnlyAdminsMayCreateProjects: Bool(false),
						OnlyAdminsSeeBillableRates:  Bool(false),
						OnlyAdminsSeeTeamDashboard:  Bool(false),
						ProjectsBillableByDefault:   Bool(true),
						Rounding:                    Int(1),
						RoundingMinutes:             Int(0),
						APIToken:                    String("1234567890abcdefghijklmnopqrstuv"),
						At:                          Time(time.Date(2013, time.August, 28, 16, 22, 21, 0, time.FixedZone("", 0))),
						IcalEnabled:                 Bool(true),
					},
					{
						ID:                          Int(7777777),
						Name:                        String("My Company Inc"),
						Profile:                     Int(100),
						Premium:                     Bool(true),
						Admin:                       Bool(true),
						DefaultHourlyRate:           Float64(0),
						DefaultCurrency:             String("USD"),
						OnlyAdminsMayCreateProjects: Bool(false),
						OnlyAdminsSeeBillableRates:  Bool(false),
						OnlyAdminsSeeTeamDashboard:  Bool(false),
						ProjectsBillableByDefault:   Bool(true),
						Rounding:                    Int(1),
						RoundingMinutes:             Int(0),
						APIToken:                    String("67890abcdefghijklmnopqrstuv12345"),
						At:                          Time(time.Date(2013, time.August, 28, 16, 22, 21, 0, time.FixedZone("", 0))),
						IcalEnabled:                 Bool(true),
						LogoURL:                     String("https://assets.toggl.com/images/workspace.jpg"),
						IcalURL:                     String("/ical/workspace_user/abcdefghijklmn1234567890opqrstuv"),
					},
				},
				err: nil,
			},
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/get_workspaces_403_forbidden",
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
					ID:                          Int(3134975),
					Name:                        String("John's personal ws"),
					Profile:                     Int(135),
					Premium:                     Bool(true),
					Admin:                       Bool(true),
					DefaultHourlyRate:           Float64(150),
					DefaultCurrency:             String("USD"),
					OnlyAdminsMayCreateProjects: Bool(false),
					OnlyAdminsSeeBillableRates:  Bool(false),
					OnlyAdminsSeeTeamDashboard:  Bool(false),
					ProjectsBillableByDefault:   Bool(true),
					Rounding:                    Int(1),
					RoundingMinutes:             Int(0),
					APIToken:                    String("1234567890abcdefghijklmnopqrstuv"),
					At:                          Time(time.Date(2013, time.August, 28, 16, 22, 21, 0, time.FixedZone("", 3*60*60))),
					LogoURL:                     String("my_logo.png"),
					IcalURL:                     String("/ical/workspace_user/9876543210abcdefghijklmnopqrstuv"),
					IcalEnabled:                 Bool(true),
				},
				err: nil,
			},
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/get_workspace_403_forbidden",
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
						ID:                    Int(123123),
						DefaultWID:            Int(777),
						Email:                 String("john@swift.com"),
						Fullname:              String("John Swift"),
						JqueryTimeofdayFormat: String("h:i A"),
						JqueryDateFormat:      String("m/d/Y"),
						TimeofdayFormat:       String("h:mm A"),
						DateFormat:            String("MM/DD/YYYY"),
						StoreStartAndStopTime: Bool(true),
						BeginningOfWeek:       Int(0),
						Language:              String("en_US"),
						ImageURL:              String("https://www.toggl.com/system/avatars/123123/small/open-uri20121116-2767-b1qr8l.png"),
						SidebarPiechart:       Bool(false),
						At:                    Time(time.Date(2013, time.March, 6, 8, 57, 12, 0, time.FixedZone("", 0))),
						CreatedAt:             Time(time.Date(2013, time.March, 6, 8, 57, 12, 0, time.FixedZone("", 0))),
						Retention:             Int(9),
						RecordTimeline:        Bool(true),
						RenderTimeline:        Bool(true),
						TimelineEnabled:       Bool(true),
						TimelineExperiment:    Bool(true),
						ShouldUpgrade:         Bool(true),
						Timezone:              String("Etc/UTC"),
						OpenIDEnabled:         Bool(false),
						SendProductEmails:     Bool(true),
						SendWeeklyReport:      Bool(true),
						SendTimeNotifications: Bool(true),
						Invitation:            &Invitation{},
						DurationFormat:        String("improved"),
					},
					{
						ID:                    Int(321321),
						Email:                 String("Happy@worker.com"),
						Fullname:              String("Happy Worker"),
						JqueryTimeofdayFormat: String("h:i A"),
						JqueryDateFormat:      String("m/d/Y"),
						TimeofdayFormat:       String("h:mm A"),
						DateFormat:            String("MM/DD/YYYY"),
						StoreStartAndStopTime: Bool(true),
						BeginningOfWeek:       Int(1),
						Language:              String("en_US"),
						ImageURL:              String("https://www.toggl.com/images/profile.png"),
						SidebarPiechart:       Bool(false),
						At:                    Time(time.Date(2013, time.March, 6, 8, 46, 7, 0, time.FixedZone("", 0))),
						CreatedAt:             Time(time.Date(2013, time.March, 6, 7, 52, 3, 0, time.FixedZone("", 0))),
						Retention:             Int(0),
						RecordTimeline:        Bool(true),
						RenderTimeline:        Bool(true),
						TimelineEnabled:       Bool(true),
						TimelineExperiment:    Bool(true),
						ShouldUpgrade:         Bool(true),
						Timezone:              String("Etc/UTC"),
						OpenIDEnabled:         Bool(false),
						SendProductEmails:     Bool(true),
						SendWeeklyReport:      Bool(true),
						SendTimeNotifications: Bool(true),
						Invitation:            &Invitation{},
						DurationFormat:        String("improved"),
					},
				},
				err: nil,
			},
		},
		{
			name:         "400 Bad Request",
			statusCode:   http.StatusBadRequest,
			testdataFile: "testdata/workspaces/get_workspace_users_400_bad_request.json",
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
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/get_workspace_users_403_forbidden",
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

func TestGetWorkspaceProjects(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		testdataFile string
		in           struct {
			id int
		}
		out struct {
			projects []*Project
			err      error
		}
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			testdataFile: "testdata/workspaces/get_workspace_projects_200_ok.json",
			in: struct {
				id int
			}{
				id: 777,
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: []*Project{
					{
						ID:            Int(123456789),
						WID:           Int(4567890),
						Name:          String("Project1"),
						Billable:      Bool(false),
						IsPrivate:     Bool(true),
						Active:        Bool(true),
						Template:      Bool(false),
						At:            Time(time.Date(2013, time.August, 13, 5, 37, 8, 0, time.FixedZone("", 0))),
						CreatedAt:     Time(time.Date(2013, time.August, 10, 4, 56, 7, 0, time.FixedZone("", 0))),
						Color:         String("13"),
						AutoEstimates: Bool(false),
						HexColor:      String("#d92b2b"),
					},
					{
						ID:            Int(234567890),
						WID:           Int(4567890),
						Name:          String("Project2"),
						Billable:      Bool(false),
						IsPrivate:     Bool(false),
						Active:        Bool(true),
						Template:      Bool(false),
						At:            Time(time.Date(2015, time.November, 26, 7, 57, 8, 0, time.FixedZone("", 0))),
						CreatedAt:     Time(time.Date(2015, time.October, 24, 4, 46, 43, 0, time.FixedZone("", 0))),
						Color:         String("6"),
						AutoEstimates: Bool(false),
						HexColor:      String("#06a893"),
					},
				},
				err: nil,
			},
		},
		{
			name:         "400 Bad Request",
			statusCode:   http.StatusBadRequest,
			testdataFile: "testdata/workspaces/get_workspace_projects_400_bad_request.json",
			in: struct {
				id int
			}{
				id: 777,
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/get_workspace_projects_403_forbidden",
			in: struct {
				id int
			}{
				id: 777,
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
			apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(tt.in.id), "projects")
			mockServer := newMockServer(t, apiSpecificPath, tt.statusCode, tt.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			projects, err := client.GetWorkspaceProjects(context.Background(), tt.in.id)

			if !reflect.DeepEqual(projects, tt.out.projects) {
				errorf(t, projects, tt.out.projects)
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

func TestGetworkspaceProjectsRequestParameters(t *testing.T) {
	tests := []struct {
		name string
		in   []requestParameter
		out  string
	}{
		{
			name: "active=false",
			in:   []requestParameter{Active("false")},
			out:  "active=false",
		},
		{
			name: "active=true&actual_hours=true",
			in:   []requestParameter{Active("true"), ActualHours(true)},
			out:  "active=true&actual_hours=true",
		},
		{
			name: "active=both&actual_hours=true&only_templates=false",
			in:   []requestParameter{Active("both"), ActualHours(true), OnlyTemplates(false)},
			out:  "active=both&actual_hours=true&only_templates=false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestParameters(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaceID := 3134975
			_, _ = client.GetWorkspaceProjects(context.Background(), workspaceID, tt.in...)
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
					Name: String("John's ws"),
				},
			},
			out: struct {
				workspace *Workspace
				err       error
			}{
				workspace: &Workspace{
					ID:                          Int(3134975),
					Name:                        String("John's ws"),
					Profile:                     Int(0),
					Premium:                     Bool(true),
					Admin:                       Bool(true),
					DefaultHourlyRate:           Float64(50),
					DefaultCurrency:             String("USD"),
					OnlyAdminsMayCreateProjects: Bool(false),
					OnlyAdminsSeeBillableRates:  Bool(true),
					OnlyAdminsSeeTeamDashboard:  Bool(false),
					ProjectsBillableByDefault:   Bool(true),
					Rounding:                    Int(1),
					RoundingMinutes:             Int(60),
					APIToken:                    String("1234567890abcdefghijklmnopqrstuv"),
					At:                          Time(time.Date(2013, time.August, 28, 16, 22, 21, 0, time.FixedZone("", 3*60*60))),
					LogoURL:                     String("my_logo.png"),
					IcalEnabled:                 Bool(true),
				},
				err: nil,
			},
		},
		{
			name:         "400 Bad Request",
			statusCode:   http.StatusBadRequest,
			testdataFile: "testdata/workspaces/update_400_bad_request.json",
			in: struct {
				id        int
				workspace *Workspace
			}{
				id: 3134975,
				workspace: &Workspace{
					Name: String("John's ws"),
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
						"Content-Type":   []string{"application/json; charset=utf-8"},
						"Date":           []string{time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)},
					},
				},
			},
		},
		{
			name:         "403 Forbidden",
			statusCode:   http.StatusForbidden,
			testdataFile: "testdata/workspaces/update_403_forbidden",
			in: struct {
				id        int
				workspace *Workspace
			}{
				id: 3134975,
				workspace: &Workspace{
					Name: String("John's ws"),
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
		out  string
	}{
		{
			name: "string",
			in: &Workspace{
				Name: String("John's ws"),
			},
			out: "{\"workspace\":{\"name\":\"John's ws\"}}",
		},
		{
			name: "string and float64",
			in: &Workspace{
				Name:              String("John's ws"),
				DefaultHourlyRate: Float64(50),
			},
			out: "{\"workspace\":{\"name\":\"John's ws\",\"default_hourly_rate\":50}}",
		},
		{
			name: "string, float64, and bool",
			in: &Workspace{
				Name:                        String("John's ws"),
				DefaultHourlyRate:           Float64(50),
				OnlyAdminsMayCreateProjects: Bool(false),
			},
			out: "{\"workspace\":{\"name\":\"John's ws\",\"default_hourly_rate\":50,\"only_admins_may_create_projects\":false}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			workspaceID := 3134975
			_, _ = client.UpdateWorkspace(context.Background(), workspaceID, tt.in)
		})
	}
}
