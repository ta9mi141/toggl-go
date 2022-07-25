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
					ID:                      Int(1234567),
					Name:                    String("test organization"),
					PricingPlanID:           Int(0),
					CreatedAt:               Time(time.Date(2020, time.January, 23, 4, 56, 07, 678184000, time.UTC)),
					At:                      Time(time.Date(2020, time.January, 23, 4, 56, 07, 678184000, time.UTC)),
					ServerDeletedAt:         nil,
					IsMultiWorkspaceEnabled: Bool(false),
					SuspendedAt:             nil,
					UserCount:               Int(1),
					TrialInfo: &TrialInfo{
						Trial:             Bool(false),
						TrialAvailable:    Bool(true),
						TrialEndDate:      nil,
						NextPaymentDate:   nil,
						LastPricingPlanID: nil,
					},
					IsChargify:    Bool(false),
					IsUnified:     Bool(false),
					MaxWorkspaces: Int(20),
					Admin:         Bool(true),
					Owner:         Bool(true),
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
