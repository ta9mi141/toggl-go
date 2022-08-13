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

func TestGetOrganization(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			organization *Organization
			err          error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/organizations/get_organization_200_ok.json",
			},
			out: struct {
				organization *Organization
				err          error
			}{
				organization: &Organization{
					ID:                      Ptr(1234567),
					Name:                    Ptr("test organization"),
					PricingPlanID:           Ptr(0),
					CreatedAt:               Ptr(time.Date(2020, time.January, 23, 4, 56, 07, 678184000, time.UTC)),
					At:                      Ptr(time.Date(2020, time.January, 23, 4, 56, 07, 678184000, time.UTC)),
					ServerDeletedAt:         nil,
					IsMultiWorkspaceEnabled: Ptr(false),
					SuspendedAt:             nil,
					UserCount:               Ptr(1),
					TrialInfo: &TrialInfo{
						Trial:             Ptr(false),
						TrialAvailable:    Ptr(true),
						TrialEndDate:      nil,
						NextPaymentDate:   nil,
						LastPricingPlanID: nil,
					},
					IsChargify:    Ptr(false),
					IsUnified:     Ptr(false),
					MaxWorkspaces: Ptr(20),
					Admin:         Ptr(true),
					Owner:         Ptr(true),
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
				testdataFile: "testdata/organizations/get_organization_400_bad_request.json",
			},
			out: struct {
				organization *Organization
				err          error
			}{
				organization: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "\"Missing or invalid organization_id\"\n",
					header: http.Header{
						"Content-Length": []string{"37"},
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
				testdataFile: "testdata/organizations/get_organization_403_forbidden",
			},
			out: struct {
				organization *Organization
				err          error
			}{
				organization: nil,
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
				testdataFile: "testdata/organizations/get_organization_500_internal_server_error",
			},
			out: struct {
				organization *Organization
				err          error
			}{
				organization: nil,
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
			apiSpecificPath := path.Join(organizationsPath, strconv.Itoa(organizationID))
			mockServer := newMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			organization, err := client.GetOrganization(context.Background(), organizationID)

			if !reflect.DeepEqual(organization, tt.out.organization) {
				errorf(t, organization, tt.out.organization)
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

func TestGetOrganizationUsers(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			organizationUsers []*OrganizationUser
			err               error
		}
	}{
		{
			name: "200 OK",
			in: struct {
				statusCode   int
				testdataFile string
			}{
				statusCode:   http.StatusOK,
				testdataFile: "testdata/organizations/get_organization_users_200_ok.json",
			},
			out: struct {
				organizationUsers []*OrganizationUser
				err               error
			}{
				organizationUsers: []*OrganizationUser{
					{
						ID:             Ptr(1234567),
						Name:           Ptr("Toggl Track"),
						Email:          Ptr("toggl@example.com"),
						UserID:         Ptr(2345678),
						AvatarURL:      Ptr(""),
						Admin:          Ptr(true),
						Owner:          Ptr(true),
						Joined:         Ptr(true),
						InvitationCode: nil,
						Inactive:       Ptr(false),
						CanEditEmail:   Ptr(false),
						Workspaces: []*workspace{
							{
								WorkspaceID: Ptr(3456789),
								Admin:       Ptr(true),
								Name:        Ptr("Workspace1"),
							},
						},
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
				testdataFile: "testdata/organizations/get_organization_users_400_bad_request.json",
			},
			out: struct {
				organizationUsers []*OrganizationUser
				err               error
			}{
				organizationUsers: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "\"Missing or invalid organization_id\"\n",
					header: http.Header{
						"Content-Length": []string{"37"},
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
				testdataFile: "testdata/organizations/get_organization_users_403_forbidden",
			},
			out: struct {
				organizationUsers []*OrganizationUser
				err               error
			}{
				organizationUsers: nil,
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
				testdataFile: "testdata/organizations/get_organization_users_500_internal_server_error",
			},
			out: struct {
				organizationUsers []*OrganizationUser
				err               error
			}{
				organizationUsers: nil,
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
			apiSpecificPath := path.Join(organizationsPath, strconv.Itoa(organizationID), "users")
			mockServer := newMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			organizationUsers, err := client.GetOrganizationUsers(context.Background(), organizationID, &GetOrganizationUsersQuery{})

			if !reflect.DeepEqual(organizationUsers, tt.out.organizationUsers) {
				errorf(t, organizationUsers, tt.out.organizationUsers)
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

func TestGetOrganizationUsersQuery(t *testing.T) {
	tests := []struct {
		name string
		in   *GetOrganizationUsersQuery
		out  string
	}{
		{
			name: "GetOrganizationUsersQuery is nil",
			in:   nil,
			out:  "",
		},
		{
			name: "filter=toggl",
			in:   &GetOrganizationUsersQuery{Filter: Ptr("toggl")},
			out:  "filter=toggl",
		},
		{
			name: "filter=toggl&only_admins=true",
			in: &GetOrganizationUsersQuery{
				Filter:     Ptr("toggl"),
				OnlyAdmins: Ptr("true"),
			},
			out: "filter=toggl&only_admins=true",
		},
		{
			name: "filter=toggl&only_admins=true&page=2",
			in: &GetOrganizationUsersQuery{
				Filter:     Ptr("toggl"),
				OnlyAdmins: Ptr("true"),
				Page:       Ptr(2),
			},
			out: "filter=toggl&only_admins=true&page=2",
		},
		{
			name: "GetOrganizationUsersQuery is empty",
			in:   &GetOrganizationUsersQuery{},
			out:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertQuery(t, tt.out)
			defer mockServer.Close()

			organizationID := 1234567
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			_, _ = client.GetOrganizationUsers(context.Background(), organizationID, tt.in)
		})
	}
}
