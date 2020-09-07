package toggl_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/it-akumi/toggl-go/toggl"
)

func TestGetUser(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx context.Context
		}
		out struct {
			user *toggl.User
			err  error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/get_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: &toggl.User{
					Id:                    1234567,
					APIToken:              "1234567890abcdefghijklmnopqrstuv",
					DefaultWid:            1234567,
					Email:                 "test.user@toggl.com",
					Fullname:              "Test User",
					JQueryTimeofdayFormat: "H:i",
					JQueryDateFormat:      "m/d/Y",
					TimeofdayFormat:       "H:mm",
					DateFormat:            "MM/DD/YYYY",
					StoreStartAndStopTime: true,
					BeginningOfWeek:       1,
					Language:              "en_US",
					ImageUrl:              "https://assets.toggl.com/avatars/12345678901234567890abcdefghijkl.png",
					SidebarPiechart:       true,
					At:                    time.Date(2013, time.October, 5, 3, 21, 34, 0, time.FixedZone("", 0)),
					NewBlogPost: struct {
						Title string `json:"title"`
						URL   string `json:"url"`
					}{},
					SendProductEmails:      true,
					SendWeeklyReport:       true,
					SendTimerNotifications: true,
					OpenidEnabled:          false,
					Timezone:               "Asia/Tokyo",
					TimeEntries: []*toggl.TimeEntry{{
						Id:          1234567890,
						Description: "sample time entry",
						Wid:         1234567,
						Pid:         12345678,
						Start:       time.Date(2019, time.October, 12, 3, 21, 34, 0, time.FixedZone("", 0)),
						Duration:    -1234567890,
						Duronly:     false,
						At:          time.Date(2019, time.December, 12, 3, 21, 34, 0, time.FixedZone("", 0)),
					}},
					Projects: []*toggl.Project{{
						Id:        234567890,
						Name:      "Sample Project",
						Wid:       9876543,
						Active:    true,
						IsPrivate: true,
						At:        time.Date(2019, time.October, 12, 3, 21, 34, 0, time.FixedZone("", 0)),
						Color:     "5",
						CreatedAt: time.Date(2019, time.February, 14, 3, 21, 34, 0, time.FixedZone("", 0)),
					}},
					Tags: []*toggl.Tag{{
						Id:   8901234,
						Name: "sample-tag",
						Wid:  4321098,
					}},
					Workspaces: []*toggl.Workspace{{
						Id:                          1234567,
						Name:                        "sample workspace",
						Premium:                     false,
						Admin:                       true,
						DefaultHourlyRate:           0,
						DefaultCurrency:             "USD",
						OnlyAdminsMayCreateProjects: false,
						OnlyAdminsSeeBillableRates:  false,
						Rounding:                    1,
						RoundingMinutes:             0,
						At:                          time.Date(2019, time.December, 12, 3, 21, 34, 0, time.FixedZone("", 0)),
					}},
					Clients: []*toggl.TogglClient{{
						Id:   12345678,
						Name: "sample-client",
						Wid:  1234567,
						At:   time.Date(2019, time.December, 12, 3, 21, 34, 0, time.FixedZone("", 0)),
					}},
				},
				err: nil,
			},
		},
		{
			name:             "401 Unauthorized",
			httpStatus:       http.StatusUnauthorized,
			testdataFilePath: "testdata/users/get_401_unauthorized.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    401,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/users/get_403_forbidden.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/get_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: nil,
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err:  toggl.ErrContextNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualUser, err := client.GetUser(c.in.ctx)
			if !reflect.DeepEqual(actualUser, c.out.user) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.user, actualUser)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestGetUserUseURLIncludingQueryStrings(t *testing.T) {
	withRelatedData := "true"
	expectedRequestURI := "/api/v8/me?with_related_data=" + withRelatedData
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetUser(context.Background(), toggl.WithRelatedData(withRelatedData))
}

func TestUpdateUser(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx  context.Context
			user *toggl.User
		}
		out struct {
			user *toggl.User
			err  error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/update_200_ok.json",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx: context.Background(),
				user: &toggl.User{
					Fullname:          "Test User",
					BeginningOfWeek:   0,
					SendProductEmails: true,
				},
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: &toggl.User{
					Id:                    1234567,
					APIToken:              "1234567890abcdefghijklmnopqrstuv",
					DefaultWid:            1234567,
					Email:                 "test.user@toggl.com",
					Fullname:              "Test User",
					JQueryTimeofdayFormat: "H:i",
					JQueryDateFormat:      "m/d/Y",
					TimeofdayFormat:       "H:mm",
					DateFormat:            "MM/DD/YYYY",
					StoreStartAndStopTime: true,
					BeginningOfWeek:       0,
					Language:              "en_US",
					ImageUrl:              "https://assets.toggl.com/avatars/1234567890abcdefghijklmnopqrstuv.png",
					SidebarPiechart:       true,
					At:                    time.Date(2013, time.October, 5, 3, 21, 34, 0, time.FixedZone("", 0)),
					NewBlogPost: struct {
						Title string `json:"title"`
						URL   string `json:"url"`
					}{},
					SendProductEmails:      true,
					SendWeeklyReport:       true,
					SendTimerNotifications: true,
					OpenidEnabled:          false,
					Timezone:               "Asia/Tokyo",
				},
				err: nil,
			},
		},
		{
			name:             "401 Unauthorized",
			httpStatus:       http.StatusUnauthorized,
			testdataFilePath: "testdata/users/update_401_unauthorized.json",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx: context.Background(),
				user: &toggl.User{
					Fullname:          "Test User",
					BeginningOfWeek:   0,
					SendProductEmails: true,
				},
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    401,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/users/update_403_forbidden.json",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx: context.Background(),
				user: &toggl.User{
					Fullname:          "Test User",
					BeginningOfWeek:   0,
					SendProductEmails: true,
				},
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/update_200_ok.json",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx: nil,
				user: &toggl.User{
					Fullname:          "Test User",
					BeginningOfWeek:   0,
					SendProductEmails: true,
				},
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err:  toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without user",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/update_200_ok.json",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx:  context.Background(),
				user: nil,
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err:  toggl.ErrUserNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualUser, err := client.UpdateUser(c.in.ctx, c.in.user)
			if !reflect.DeepEqual(actualUser, c.out.user) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.user, actualUser)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestUpdateUserConvertParamsToRequestBody(t *testing.T) {
	expectedUserRequest := &toggl.User{
		SendWeeklyReport: false,
		CurrentPassword:  "old_password",
		Password:         "new_password",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualUserRequest := new(toggl.User)
		if err := json.Unmarshal(requestBody, actualUserRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualUserRequest, expectedUserRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedUserRequest, actualUserRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateUser(context.Background(), expectedUserRequest)
}

func TestResetAPIToken(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx context.Context
		}
		out struct {
			apiToken string
			err      error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/reset_api_token_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				apiToken string
				err      error
			}{
				apiToken: "1234567890abcdefghijklmnopqrstuv",
				err:      nil,
			},
		},
		{
			name:             "401 Unauthorized",
			httpStatus:       http.StatusUnauthorized,
			testdataFilePath: "testdata/users/reset_api_token_401_unauthorized.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				apiToken string
				err      error
			}{
				apiToken: "",
				err: &toggl.TogglError{
					Message: "",
					Code:    401,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/users/reset_api_token_403_forbidden.json",
			in: struct {
				ctx context.Context
			}{
				ctx: context.Background(),
			},
			out: struct {
				apiToken string
				err      error
			}{
				apiToken: "",
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/reset_api_token_200_ok.json",
			in: struct {
				ctx context.Context
			}{
				ctx: nil,
			},
			out: struct {
				apiToken string
				err      error
			}{
				apiToken: "",
				err:      toggl.ErrContextNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualAPIToken, err := client.ResetAPIToken(c.in.ctx)
			if !reflect.DeepEqual(actualAPIToken, c.out.apiToken) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.apiToken, actualAPIToken)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestSignUp(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx  context.Context
			user *toggl.User
		}
		out struct {
			user *toggl.User
			err  error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/signup_200_ok.json",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx: context.Background(),
				user: &toggl.User{
					Email:    "test.user@toggl.com",
					Password: "password",
				},
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: &toggl.User{
					Id:                    1234567,
					DefaultWid:            9876543,
					Email:                 "test.user@toggl.com",
					Fullname:              "Test User",
					JQueryTimeofdayFormat: "",
					JQueryDateFormat:      "",
					TimeofdayFormat:       "",
					DateFormat:            "",
					StoreStartAndStopTime: false,
					BeginningOfWeek:       0,
					SidebarPiechart:       false,
					NewBlogPost: struct {
						Title string `json:"title"`
						URL   string `json:"url"`
					}{},
					SendProductEmails:      false,
					SendWeeklyReport:       false,
					SendTimerNotifications: false,
					OpenidEnabled:          false,
					Timezone:               "Etc/UTC",
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/users/signup_400_bad_request.txt",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx: context.Background(),
				user: &toggl.User{
					Email:    "test.user@toggl.com",
					Password: "password",
				},
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err: &toggl.TogglError{
					Message: "User with this email already exists\n",
					Code:    400,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/signup_200_ok.json",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx: nil,
				user: &toggl.User{
					Email:    "test.user@toggl.com",
					Password: "password",
				},
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err:  toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without user",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/users/signup_200_ok.json",
			in: struct {
				ctx  context.Context
				user *toggl.User
			}{
				ctx:  context.Background(),
				user: nil,
			},
			out: struct {
				user *toggl.User
				err  error
			}{
				user: nil,
				err:  toggl.ErrUserNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualUser, err := client.SignUp(c.in.ctx, c.in.user)
			if !reflect.DeepEqual(actualUser, c.out.user) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.user, actualUser)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, togglError)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out.err, err)
				}
			}
		})
	}
}

func TestSignUpConvertParamsToRequestBody(t *testing.T) {
	expectedSignUpRequest := &toggl.User{
		Email:       "test.user@toggl.com",
		Password:    "password",
		Timezone:    "Etc/UTC",
		CreatedWith: "toggl-go",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualSignUpRequest := new(toggl.User)
		if err := json.Unmarshal(requestBody, actualSignUpRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualSignUpRequest, expectedSignUpRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedSignUpRequest, actualSignUpRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.SignUp(context.Background(), expectedSignUpRequest)
}
