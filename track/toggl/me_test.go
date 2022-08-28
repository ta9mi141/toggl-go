package toggl

import (
	"context"
	"errors"
	"net/http"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/ta9mi141/toggl-go/track"
	"github.com/ta9mi141/toggl-go/track/internal"
)

func TestGetMe(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			me  *Me
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
				testdataFile: "testdata/me/get_me_200_ok.json",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: &Me{
					ID:                 track.Ptr(1234567),
					APIToken:           track.Ptr("abcdefghijklmnopqrstuvwxyz123456"),
					Email:              track.Ptr("example@toggl.com"),
					Fullname:           track.Ptr("Example Toggl"),
					Timezone:           track.Ptr("Asia/Tokyo"),
					DefaultWorkspaceID: track.Ptr(1234567),
					BeginningOfWeek:    track.Ptr(1),
					ImageURL:           track.Ptr("https://assets.track.toggl.com/images/profile.png"),
					CreatedAt:          track.Ptr(time.Date(2012, time.March, 4, 1, 23, 45, 210809000, time.UTC)),
					UpdatedAt:          track.Ptr(time.Date(2012, time.May, 6, 2, 34, 56, 346231000, time.UTC)),
					OpenIDEnabled:      track.Ptr(false),
					At:                 track.Ptr(time.Date(2012, time.June, 7, 8, 9, 10, 810517000, time.UTC)),
					IntercomHash:       track.Ptr("1234567890abcdefghijklmnopqustuvwxyz1234567890avcdefghijklmnopqr"),
					HasPassword:        track.Ptr(true),
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
				testdataFile: "testdata/me/get_me_403_forbidden",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
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
				testdataFile: "testdata/me/get_me_500_internal_server_error",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
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
			mockServer := internal.NewMockServer(t, mePath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			me, err := apiClient.GetMe(context.Background())

			if !reflect.DeepEqual(me, tt.out.me) {
				internal.Errorf(t, me, tt.out.me)
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

func TestUpdateMe(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			me  *Me
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
				testdataFile: "testdata/me/update_me_200_ok.json",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: &Me{
					ID:                 track.Ptr(1234567),
					APIToken:           track.Ptr("abcdefghijklmnopqrstuvwxyz123456"),
					Email:              track.Ptr("example@toggl.com"),
					Fullname:           track.Ptr("Example Toggl"),
					Timezone:           track.Ptr("Asia/Tokyo"),
					DefaultWorkspaceID: track.Ptr(1234567),
					BeginningOfWeek:    track.Ptr(0),
					ImageURL:           track.Ptr("https://assets.track.toggl.com/images/profile.png"),
					CreatedAt:          track.Ptr(time.Date(2012, time.March, 4, 1, 23, 45, 210809000, time.UTC)),
					UpdatedAt:          track.Ptr(time.Date(2012, time.May, 6, 2, 34, 56, 346231000, time.UTC)),
					OpenIDEnabled:      track.Ptr(false),
					At:                 track.Ptr(time.Date(2012, time.June, 7, 8, 9, 10, 810517000, time.UTC)),
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
				testdataFile: "testdata/me/update_me_400_bad_request.json",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"Invalid beginning_of_week\"\n",
					Header: http.Header{
						"Content-Length": []string{"28"},
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
				testdataFile: "testdata/me/update_me_403_forbidden",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
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
				testdataFile: "testdata/me/update_me_500_internal_server_error",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
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
			mockServer := internal.NewMockServer(t, mePath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			me, err := apiClient.UpdateMe(context.Background(), &UpdateMeRequestBody{})

			if !reflect.DeepEqual(me, tt.out.me) {
				internal.Errorf(t, me, tt.out.me)
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

func TestUpdateMeRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *UpdateMeRequestBody
		out  string
	}{
		{
			name: "int",
			in: &UpdateMeRequestBody{
				BeginningOfWeek: track.Ptr(0),
			},
			out: "{\"beginning_of_week\":0}",
		},
		{
			name: "string",
			in: &UpdateMeRequestBody{
				Fullname: track.Ptr("Awesome Name"),
			},
			out: "{\"fullname\":\"Awesome Name\"}",
		},
		{
			name: "int and string",
			in: &UpdateMeRequestBody{
				CurrentPassword:    track.Ptr("vulnerable password"),
				DefaultWorkspaceID: track.Ptr(1234567),
				Password:           track.Ptr("secure password"),
			},
			out: "{\"current_password\":\"vulnerable password\",\"default_workspace_id\":1234567,\"password\":\"secure password\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			_, _ = apiClient.UpdateMe(context.Background(), tt.in)
		})
	}
}

func TestGetMyOrganizations(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			organizations []*Organization
			err           error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/me/get_my_organizations_200_ok.json",
			},
			out: struct {
				organizations []*Organization
				err           error
			}{
				organizations: []*Organization{
					{
						ID:                      track.Ptr(1234567),
						Name:                    track.Ptr("Organization 1"),
						PricingPlanID:           track.Ptr(0),
						CreatedAt:               track.Ptr(time.Date(2018, time.January, 23, 4, 56, 15, 288620000, time.UTC)),
						At:                      track.Ptr(time.Date(2019, time.January, 23, 4, 56, 15, 288620000, time.UTC)),
						ServerDeletedAt:         nil,
						IsMultiWorkspaceEnabled: track.Ptr(false),
						SuspendedAt:             nil,
						UserCount:               track.Ptr(1),
						TrialInfo: &TrialInfo{
							Trial:             track.Ptr(false),
							TrialAvailable:    track.Ptr(true),
							TrialEndDate:      nil,
							NextPaymentDate:   nil,
							LastPricingPlanID: nil,
						},
						IsChargify:    track.Ptr(false),
						MaxWorkspaces: track.Ptr(20),
						Admin:         track.Ptr(true),
						Owner:         track.Ptr(true),
					},
					{
						ID:                      track.Ptr(2345678),
						Name:                    track.Ptr("Organization 2"),
						PricingPlanID:           track.Ptr(0),
						CreatedAt:               track.Ptr(time.Date(2020, time.February, 3, 4, 5, 6, 678184000, time.UTC)),
						At:                      track.Ptr(time.Date(2021, time.February, 3, 4, 5, 6, 678184000, time.UTC)),
						ServerDeletedAt:         nil,
						IsMultiWorkspaceEnabled: track.Ptr(false),
						SuspendedAt:             nil,
						UserCount:               track.Ptr(1),
						TrialInfo: &TrialInfo{
							Trial:             track.Ptr(false),
							TrialAvailable:    track.Ptr(true),
							TrialEndDate:      nil,
							NextPaymentDate:   nil,
							LastPricingPlanID: nil,
						},
						IsChargify:    track.Ptr(false),
						MaxWorkspaces: track.Ptr(20),
						Admin:         track.Ptr(true),
						Owner:         track.Ptr(true),
					},
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
				testdataFile: "testdata/me/get_my_organizations_403_forbidden",
			},
			out: struct {
				organizations []*Organization
				err           error
			}{
				organizations: nil,
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
				testdataFile: "testdata/me/get_my_organizations_500_internal_server_error",
			},
			out: struct {
				organizations []*Organization
				err           error
			}{
				organizations: nil,
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
			apiSpecificPath := path.Join(mePath, "organizations")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			organizations, err := apiClient.GetMyOrganizations(context.Background())

			if !reflect.DeepEqual(organizations, tt.out.organizations) {
				internal.Errorf(t, organizations, tt.out.organizations)
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

func TestGetMyProjects(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			projects []*Project
			err      error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/me/get_my_projects_200_ok.json",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: []*Project{
					{
						ID:                  track.Ptr(123456789),
						WorkspaceID:         track.Ptr(1234567),
						ClientID:            nil,
						Name:                track.Ptr("Project1"),
						IsPrivate:           track.Ptr(true),
						Active:              track.Ptr(true),
						At:                  track.Ptr(time.Date(2013, time.March, 4, 5, 6, 7, 0, time.Local)),
						CreatedAt:           track.Ptr(time.Date(2012, time.March, 4, 5, 6, 7, 0, time.Local)),
						ServerDeletedAt:     nil,
						Color:               track.Ptr("#e36a00"),
						Billable:            nil,
						Template:            nil,
						AutoEstimates:       nil,
						EstimatedHours:      nil,
						Rate:                nil,
						RateLastUpdated:     nil,
						Currency:            nil,
						Recurring:           track.Ptr(false),
						RecurringParameters: nil,
						CurrentPeriod:       nil,
						FixedFee:            nil,
						ActualHours:         track.Ptr(0),
						WID:                 track.Ptr(1234567),
						CID:                 nil,
					},
					{
						ID:                  track.Ptr(987654321),
						WorkspaceID:         track.Ptr(9876543),
						ClientID:            nil,
						Name:                track.Ptr("Project2"),
						IsPrivate:           track.Ptr(true),
						Active:              track.Ptr(true),
						At:                  track.Ptr(time.Date(2021, time.January, 23, 4, 56, 7, 0, time.Local)),
						CreatedAt:           track.Ptr(time.Date(2020, time.January, 23, 4, 56, 7, 0, time.Local)),
						ServerDeletedAt:     nil,
						Color:               track.Ptr("#c9806b"),
						Billable:            nil,
						Template:            nil,
						AutoEstimates:       nil,
						EstimatedHours:      nil,
						Rate:                nil,
						RateLastUpdated:     nil,
						Currency:            nil,
						Recurring:           track.Ptr(false),
						RecurringParameters: nil,
						CurrentPeriod:       nil,
						FixedFee:            nil,
						ActualHours:         track.Ptr(0),
						WID:                 track.Ptr(9876543),
						CID:                 nil,
					},
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
				testdataFile: "testdata/me/get_my_projects_403_forbidden",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
				testdataFile: "testdata/me/get_my_projects_404_not_found.json",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
				err: &internal.ErrorResponse{
					StatusCode: 404,
					Message:    "\"Invalid include_archived\"\n",
					Header: http.Header{
						"Content-Length": []string{"27"},
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
				testdataFile: "testdata/me/get_my_projects_500_internal_server_error",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
			apiSpecificPath := path.Join(mePath, "projects")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			projects, err := apiClient.GetMyProjects(context.Background(), nil)

			if !reflect.DeepEqual(projects, tt.out.projects) {
				internal.Errorf(t, projects, tt.out.projects)
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

func TestGetMyProjectsQuery(t *testing.T) {
	tests := []struct {
		name string
		in   *GetMyProjectsQuery
		out  string
	}{
		{
			name: "GetMyProjectsQuery is nil",
			in:   nil,
			out:  "",
		},
		{
			name: "include_archived=true",
			in:   &GetMyProjectsQuery{IncludeArchived: track.Ptr("true")},
			out:  "include_archived=true",
		},
		{
			name: "include_archived=false",
			in:   &GetMyProjectsQuery{IncludeArchived: track.Ptr("false")},
			out:  "include_archived=false",
		},
		{
			name: "GetMyProjectsQuery is empty",
			in:   &GetMyProjectsQuery{},
			out:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertQuery(t, tt.out)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			_, _ = apiClient.GetMyProjects(context.Background(), tt.in)
		})
	}
}

func TestGetMyProjectsPaginated(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			projects []*Project
			err      error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/me/get_my_projects_paginated_200_ok.json",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: []*Project{
					{
						ID:                  track.Ptr(123456789),
						WorkspaceID:         track.Ptr(1234567),
						ClientID:            nil,
						Name:                track.Ptr("Project1"),
						IsPrivate:           track.Ptr(true),
						Active:              track.Ptr(true),
						At:                  track.Ptr(time.Date(2013, time.March, 4, 5, 6, 7, 0, time.Local)),
						CreatedAt:           track.Ptr(time.Date(2012, time.March, 4, 5, 6, 7, 0, time.Local)),
						ServerDeletedAt:     nil,
						Color:               track.Ptr("#e36a00"),
						Billable:            nil,
						Template:            nil,
						AutoEstimates:       nil,
						EstimatedHours:      nil,
						Rate:                nil,
						RateLastUpdated:     nil,
						Currency:            nil,
						Recurring:           track.Ptr(false),
						RecurringParameters: nil,
						CurrentPeriod:       nil,
						FixedFee:            nil,
						ActualHours:         track.Ptr(0),
						WID:                 track.Ptr(1234567),
						CID:                 nil,
					},
					{
						ID:                  track.Ptr(987654321),
						WorkspaceID:         track.Ptr(9876543),
						ClientID:            nil,
						Name:                track.Ptr("Project2"),
						IsPrivate:           track.Ptr(true),
						Active:              track.Ptr(true),
						At:                  track.Ptr(time.Date(2021, time.January, 23, 4, 56, 7, 0, time.Local)),
						CreatedAt:           track.Ptr(time.Date(2020, time.January, 23, 4, 56, 7, 0, time.Local)),
						ServerDeletedAt:     nil,
						Color:               track.Ptr("#c9806b"),
						Billable:            nil,
						Template:            nil,
						AutoEstimates:       nil,
						EstimatedHours:      nil,
						Rate:                nil,
						RateLastUpdated:     nil,
						Currency:            nil,
						Recurring:           track.Ptr(false),
						RecurringParameters: nil,
						CurrentPeriod:       nil,
						FixedFee:            nil,
						ActualHours:         track.Ptr(0),
						WID:                 track.Ptr(9876543),
						CID:                 nil,
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
				testdataFile: "testdata/me/get_my_projects_paginated_400_bad_request.json",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
				err: &internal.ErrorResponse{
					StatusCode: 400,
					Message:    "\"Invalid start_project_id\"\n",
					Header: http.Header{
						"Content-Length": []string{"27"},
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
				testdataFile: "testdata/me/get_my_projects_paginated_403_forbidden",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
				testdataFile: "testdata/me/get_my_projects_paginated_500_internal_server_error",
			},
			out: struct {
				projects []*Project
				err      error
			}{
				projects: nil,
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
			apiSpecificPath := path.Join(mePath, "projects/paginated")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			projects, err := apiClient.GetMyProjectsPaginated(context.Background(), nil)

			if !reflect.DeepEqual(projects, tt.out.projects) {
				internal.Errorf(t, projects, tt.out.projects)
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

func TestGetMyProjectsPaginatedQuery(t *testing.T) {
	tests := []struct {
		name string
		in   *GetMyProjectsPaginatedQuery
		out  string
	}{
		{
			name: "GetMyProjectsPaginatedQuery is nil",
			in:   nil,
			out:  "",
		},
		{
			name: "start_project_id=12345",
			in:   &GetMyProjectsPaginatedQuery{StartProjectID: track.Ptr(12345)},
			out:  "start_project_id=12345",
		},
		{
			name: "start_project_id=0",
			in:   &GetMyProjectsPaginatedQuery{StartProjectID: track.Ptr(0)},
			out:  "start_project_id=0",
		},
		{
			name: "GetMyProjectsPaginatedQuery is empty",
			in:   &GetMyProjectsPaginatedQuery{},
			out:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := internal.NewMockServerToAssertQuery(t, tt.out)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			_, _ = apiClient.GetMyProjectsPaginated(context.Background(), tt.in)
		})
	}
}

func TestGetMyTags(t *testing.T) {
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
				testdataFile: "testdata/me/get_my_tags_200_ok.json",
			},
			out: struct {
				tags []*Tag
				err  error
			}{
				tags: []*Tag{
					{
						ID:          track.Ptr(1234567),
						WorkspaceID: track.Ptr(1234567),
						Name:        track.Ptr("tag"),
						At:          track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 678900000, time.UTC)),
					},
					{
						ID:          track.Ptr(23456789),
						WorkspaceID: track.Ptr(2345678),
						Name:        track.Ptr("newtag"),
						At:          track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 678901000, time.UTC)),
					},
					{
						ID:          track.Ptr(3456789),
						WorkspaceID: track.Ptr(1234567),
						Name:        track.Ptr("toggl-go"),
						At:          track.Ptr(time.Date(2020, time.January, 2, 3, 4, 5, 678901000, time.UTC)),
					},
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
				testdataFile: "testdata/me/get_my_tags_403_forbidden",
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
				testdataFile: "testdata/me/get_my_tags_500_internal_server_error",
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
			mockServer := internal.NewMockServer(t, mePath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(WithAPIToken(internal.APIToken), withBaseURL(mockServer.URL))
			tags, err := apiClient.GetMyTags(context.Background())

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
