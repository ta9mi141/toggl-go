package track

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"
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
					ID:                 Int(1234567),
					APIToken:           String("abcdefghijklmnopqrstuvwxyz123456"),
					Email:              String("example@toggl.com"),
					Fullname:           String("Example Toggl"),
					Timezone:           String("Asia/Tokyo"),
					DefaultWorkspaceID: Int(1234567),
					BeginningOfWeek:    Int(1),
					ImageURL:           String("https://assets.track.toggl.com/images/profile.png"),
					CreatedAt:          Time(time.Date(2012, time.March, 4, 1, 23, 45, 210809000, time.UTC)),
					UpdatedAt:          Time(time.Date(2012, time.May, 6, 2, 34, 56, 346231000, time.UTC)),
					OpenIDEnabled:      Bool(false),
					At:                 Time(time.Date(2012, time.June, 7, 8, 9, 10, 810517000, time.UTC)),
					IntercomHash:       String("1234567890abcdefghijklmnopqustuvwxyz1234567890avcdefghijklmnopqr"),
					HasPassword:        Bool(true),
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
				testdataFile: "testdata/me/get_me_500_internal_server_error",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
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
			mockServer := newMockServer(t, mePath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			me, err := client.GetMe(context.Background())

			if !reflect.DeepEqual(me, tt.out.me) {
				errorf(t, me, tt.out.me)
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

func TestPutMe(t *testing.T) {
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
				testdataFile: "testdata/me/put_me_200_ok.json",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: &Me{
					ID:                 Int(1234567),
					APIToken:           String("abcdefghijklmnopqrstuvwxyz123456"),
					Email:              String("example@toggl.com"),
					Fullname:           String("Example Toggl"),
					Timezone:           String("Asia/Tokyo"),
					DefaultWorkspaceID: Int(1234567),
					BeginningOfWeek:    Int(0),
					ImageURL:           String("https://assets.track.toggl.com/images/profile.png"),
					CreatedAt:          Time(time.Date(2012, time.March, 4, 1, 23, 45, 210809000, time.UTC)),
					UpdatedAt:          Time(time.Date(2012, time.May, 6, 2, 34, 56, 346231000, time.UTC)),
					OpenIDEnabled:      Bool(false),
					At:                 Time(time.Date(2012, time.June, 7, 8, 9, 10, 810517000, time.UTC)),
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
				testdataFile: "testdata/me/put_me_400_bad_request.json",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
				err: &errorResponse{
					statusCode: 400,
					message:    "\"Invalid beginning_of_week\"\n",
					header: http.Header{
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
				testdataFile: "testdata/me/put_me_403_forbidden",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
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
				testdataFile: "testdata/me/put_me_500_internal_server_error",
			},
			out: struct {
				me  *Me
				err error
			}{
				me: nil,
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
			mockServer := newMockServer(t, mePath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			me, err := client.PutMe(context.Background(), &PutMeRequestBody{})

			if !reflect.DeepEqual(me, tt.out.me) {
				errorf(t, me, tt.out.me)
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

func TestPutMeRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *PutMeRequestBody
		out  string
	}{
		{
			name: "int",
			in: &PutMeRequestBody{
				BeginningOfWeek: Int(0),
			},
			out: "{\"beginning_of_week\":0}",
		},
		{
			name: "string",
			in: &PutMeRequestBody{
				Fullname: String("Awesome Name"),
			},
			out: "{\"fullname\":\"Awesome Name\"}",
		},
		{
			name: "int and string",
			in: &PutMeRequestBody{
				CurrentPassword:    String("vulnerable password"),
				DefaultWorkspaceID: Int(1234567),
				Password:           String("secure password"),
			},
			out: "{\"current_password\":\"vulnerable password\",\"default_workspace_id\":1234567,\"password\":\"secure password\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := newMockServerToAssertRequestBody(t, tt.out)
			defer mockServer.Close()
			client := NewClient(WithAPIToken(apiToken), withBaseURL(mockServer.URL))
			_, _ = client.PutMe(context.Background(), tt.in)
		})
	}
}
