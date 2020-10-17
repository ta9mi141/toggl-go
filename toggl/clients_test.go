package toggl_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/ta9mi1shi1/toggl-go/toggl"
)

func TestGetTogglClient(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			togglClient *toggl.TogglClient
		}
		out struct {
			togglClient *toggl.TogglClient
			err         error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/get_client_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: &toggl.TogglClient{
					Id:   12345678,
					Name: "Very Big Company",
					Wid:  1234567,
					At:   time.Date(2020, time.February, 10, 10, 25, 54, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/clients/get_client_400_bad_request.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err: &toggl.TogglError{
					Message: "Invalid client_id\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/clients/get_client_403_forbidden.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/get_client_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         nil,
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err:         toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without client",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/get_client_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: nil,
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err:         toggl.ErrTogglClientNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTogglClient, err := client.GetTogglClient(c.in.ctx, c.in.togglClient)
			if !reflect.DeepEqual(actualTogglClient, c.out.togglClient) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.togglClient, actualTogglClient)
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

func TestGetTogglClientUseURLIncludingClientId(t *testing.T) {
	togglClientId := 12345678
	expectedRequestURI := "/api/v8/clients/" + strconv.Itoa(togglClientId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetTogglClient(context.Background(), &toggl.TogglClient{
		Id: togglClientId,
	})
}

func TestGetTogglClients(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               context.Context
		out              struct {
			togglClients []*toggl.TogglClient
			err          error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/get_clients_200_ok.json",
			in:               context.Background(),
			out: struct {
				togglClients []*toggl.TogglClient
				err          error
			}{
				togglClients: []*toggl.TogglClient{
					{
						Id:    12349455,
						Name:  "Very Big Company",
						Wid:   777,
						Notes: "something about the client",
						At:    time.Date(2020, time.February, 10, 10, 25, 54, 0, time.FixedZone("", 0)),
					},
					{
						Id:    1239456,
						Name:  "Small startup",
						Wid:   777,
						Notes: "Really cool people",
						At:    time.Date(2019, time.March, 11, 10, 25, 54, 0, time.FixedZone("", 0)),
					},
				},
				err: nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/clients/get_clients_403_forbidden.json",
			in:               context.Background(),
			out: struct {
				togglClients []*toggl.TogglClient
				err          error
			}{
				togglClients: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/get_clients_200_ok.json",
			in:               nil,
			out: struct {
				togglClients []*toggl.TogglClient
				err          error
			}{
				togglClients: nil,
				err:          toggl.ErrContextNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTogglClients, err := client.GetTogglClients(c.in)
			if !reflect.DeepEqual(actualTogglClients, c.out.togglClients) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.togglClients, actualTogglClients)
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

func TestGetTogglClientProjects(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			togglClient *toggl.TogglClient
		}
		out struct {
			togglClientProjects []*toggl.Project
			err                 error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/get_client_projects_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: struct {
				togglClientProjects []*toggl.Project
				err                 error
			}{
				togglClientProjects: []*toggl.Project{
					{
						Id:        909,
						Wid:       777,
						Cid:       987,
						Name:      "Very lucrative project",
						IsPrivate: true,
						Active:    true,
						At:        time.Date(2013, time.March, 6, 9, 15, 18, 0, time.FixedZone("", 0)),
					},
					{
						Id:        32143,
						Wid:       777,
						Cid:       987,
						Name:      "Factory server infrastructure",
						IsPrivate: true,
						Active:    true,
						At:        time.Date(2013, time.March, 6, 9, 16, 6, 0, time.FixedZone("", 0)),
					},
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/clients/get_client_projects_400_bad_request.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: struct {
				togglClientProjects []*toggl.Project
				err                 error
			}{
				togglClientProjects: nil,
				err: &toggl.TogglError{
					Message: "Missing or invalid client ID\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/clients/get_client_projects_403_forbidden.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: struct {
				togglClientProjects []*toggl.Project
				err                 error
			}{
				togglClientProjects: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/get_client_projects_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         nil,
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: struct {
				togglClientProjects []*toggl.Project
				err                 error
			}{
				togglClientProjects: nil,
				err:                 toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without client",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/get_client_projects_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: nil,
			},
			out: struct {
				togglClientProjects []*toggl.Project
				err                 error
			}{
				togglClientProjects: nil,
				err:                 toggl.ErrTogglClientNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTogglClientProjects, err := client.GetTogglClientProjects(c.in.ctx, c.in.togglClient)
			if !reflect.DeepEqual(actualTogglClientProjects, c.out.togglClientProjects) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.togglClientProjects, actualTogglClientProjects)
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

func TestGetTogglClientProjectsUseURLIncludingClientId(t *testing.T) {
	togglClientId := 12345678
	expectedRequestURI := "/api/v8/clients/" + strconv.Itoa(togglClientId) + "/projects" + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetTogglClientProjects(context.Background(), &toggl.TogglClient{
		Id: togglClientId,
	})
}

func TestGetTogglClientProjectsUseURLIncludingQueryStrings(t *testing.T) {
	togglClientId := 12345678
	active := "both"
	expectedRequestURI := "/api/v8/clients/" + strconv.Itoa(togglClientId) + "/projects?active=" + active
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.GetTogglClientProjects(context.Background(), &toggl.TogglClient{Id: togglClientId}, toggl.Active(active))
}

func TestCreateTogglClient(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			togglClient *toggl.TogglClient
		}
		out struct {
			togglClient *toggl.TogglClient
			err         error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/create_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx: context.Background(),
				togglClient: &toggl.TogglClient{
					Name: "Very Big Company",
					Wid:  1234567,
				},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: &toggl.TogglClient{
					Id:   12345678,
					Name: "Very Big Company",
					Wid:  1234567,
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/clients/create_400_bad_request.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx: context.Background(),
				togglClient: &toggl.TogglClient{
					Name: "Very Big Company",
					Wid:  777,
				},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err: &toggl.TogglError{
					Message: "User 1234567 cannot access workspace 777\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/clients/create_403_forbidden.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx: context.Background(),
				togglClient: &toggl.TogglClient{
					Name: "Very Big Company",
					Wid:  777,
				},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/create_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx: nil,
				togglClient: &toggl.TogglClient{
					Name: "Very Big Company",
					Wid:  777,
				},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err:         toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without client",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/create_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: nil,
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err:         toggl.ErrTogglClientNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTogglClient, err := client.CreateTogglClient(c.in.ctx, c.in.togglClient)
			if !reflect.DeepEqual(actualTogglClient, c.out.togglClient) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.togglClient, actualTogglClient)
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

func TestCreateTogglClientConvertParamsToRequestBody(t *testing.T) {
	expectedTogglClientRequest := &toggl.TogglClient{
		Wid:  1234567,
		Name: "Very Big Company",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualTogglClientRequest := new(toggl.TogglClient)
		if err := json.Unmarshal(requestBody, actualTogglClientRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualTogglClientRequest, expectedTogglClientRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedTogglClientRequest, actualTogglClientRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.CreateTogglClient(context.Background(), expectedTogglClientRequest)
}

func TestUpdateTogglClient(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			togglClient *toggl.TogglClient
		}
		out struct {
			togglClient *toggl.TogglClient
			err         error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/update_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx: context.Background(),
				togglClient: &toggl.TogglClient{
					Id:    12345678,
					Name:  "Very Big Company",
					Notes: "something about the client",
				},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: &toggl.TogglClient{
					Id:    12345678,
					Name:  "Very Big Company",
					Wid:   1234567,
					Notes: "something about the client",
					At:    time.Date(2020, time.February, 10, 10, 15, 46, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/clients/update_400_bad_request.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx: context.Background(),
				togglClient: &toggl.TogglClient{
					Name:  "Very Big Company",
					Notes: "something about the client",
				},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err: &toggl.TogglError{
					Message: "Client can't be blank\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/clients/update_403_forbidden.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx: context.Background(),
				togglClient: &toggl.TogglClient{
					Id:    12345678,
					Name:  "Very Big Company",
					Notes: "something about the client",
				},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/update_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx: nil,
				togglClient: &toggl.TogglClient{
					Id:    12345678,
					Name:  "Very Big Company",
					Notes: "something about the client",
				},
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err:         toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without client",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/update_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: nil,
			},
			out: struct {
				togglClient *toggl.TogglClient
				err         error
			}{
				togglClient: nil,
				err:         toggl.ErrTogglClientNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTogglClient, err := client.UpdateTogglClient(c.in.ctx, c.in.togglClient)
			if !reflect.DeepEqual(actualTogglClient, c.out.togglClient) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.togglClient, actualTogglClient)
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

func TestUpdateTogglClientConvertParamsToRequestBody(t *testing.T) {
	expectedTogglClientRequest := &toggl.TogglClient{
		Name:  "Very Big Company",
		Notes: "something about the client",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualTogglClientRequest := new(toggl.TogglClient)
		if err := json.Unmarshal(requestBody, actualTogglClientRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualTogglClientRequest, expectedTogglClientRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedTogglClientRequest, actualTogglClientRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateTogglClient(context.Background(), expectedTogglClientRequest)
}

func TestUpdateTogglClientUseURLIncludingClientId(t *testing.T) {
	togglClientId := 12345678
	expectedRequestURI := "/api/v8/clients/" + strconv.Itoa(togglClientId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateTogglClient(context.Background(), &toggl.TogglClient{
		Id: togglClientId,
	})
}

func TestDeleteTogglClient(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx         context.Context
			togglClient *toggl.TogglClient
		}
		out error
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/delete_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: nil,
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/clients/delete_400_bad_request.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: &toggl.TogglError{
				Message: "User 1234567 cannot access client 12345678\n",
				Code:    400,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/clients/delete_403_forbidden.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: &toggl.TogglError{
				Message: "",
				Code:    403,
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/delete_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         nil,
				togglClient: &toggl.TogglClient{Id: 12345678},
			},
			out: toggl.ErrContextNotFound,
		},
		{
			name:             "Without client",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/clients/delete_200_ok.json",
			in: struct {
				ctx         context.Context
				togglClient *toggl.TogglClient
			}{
				ctx:         context.Background(),
				togglClient: nil,
			},
			out: toggl.ErrTogglClientNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			err := client.DeleteTogglClient(c.in.ctx, c.in.togglClient)

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out, togglError)
				}
			} else {
				if !errors.Is(err, c.out) {
					t.Errorf("\nwant: %#+v\ngot : %#+v\n", c.out, err)
				}
			}
		})
	}
}

func TestDeleteTogglClientUseURLIncludingClientId(t *testing.T) {
	togglClientId := 12345678
	expectedRequestURI := "/api/v8/clients/" + strconv.Itoa(togglClientId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_ = client.DeleteTogglClient(context.Background(), &toggl.TogglClient{
		Id: togglClientId,
	})
}
