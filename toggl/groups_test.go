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

func TestCreateGroup(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx   context.Context
			group *toggl.Group
		}
		out struct {
			group *toggl.Group
			err   error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/create_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: &toggl.Group{
					Id:   123456,
					Wid:  1234567,
					Name: "toggl-go",
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/groups/create_400_bad_request.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: nil,
				err: &toggl.TogglError{
					Message: "Name has already been taken\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/groups/create_403_forbidden.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/create_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: nil,
				group: &toggl.Group{
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: nil,
				err:   toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without group",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/create_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx:   context.Background(),
				group: nil,
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: nil,
				err:   toggl.ErrGroupNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualGroup, err := client.CreateGroup(c.in.ctx, c.in.group)
			if !reflect.DeepEqual(actualGroup, c.out.group) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.group, actualGroup)
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

func TestCreateGroupConvertParamsToRequestBody(t *testing.T) {
	expectedGroupRequest := &toggl.Group{
		Wid:  1234567,
		Name: "toggl-go",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualGroupRequest := new(toggl.Group)
		if err := json.Unmarshal(requestBody, actualGroupRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualGroupRequest, expectedGroupRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedGroupRequest, actualGroupRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.CreateGroup(context.Background(), expectedGroupRequest)
}

func TestUpdateGroup(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx   context.Context
			group *toggl.Group
		}
		out struct {
			group *toggl.Group
			err   error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/update_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Id:   123456,
					Name: "toggl-go",
				},
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: &toggl.Group{
					Id:   123456,
					Wid:  1234567,
					Name: "toggl-go",
					At:   time.Date(2020, time.February, 2, 6, 40, 53, 0, time.FixedZone("", 0)),
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/groups/update_400_bad_request.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Id:   1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: nil,
				err: &toggl.TogglError{
					Message: "Invalid group ID\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/groups/update_403_forbidden.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Id:   123456,
					Name: "toggl-go",
				},
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/update_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: nil,
				group: &toggl.Group{
					Id:   123456,
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: nil,
				err:   toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without group",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/update_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx:   context.Background(),
				group: nil,
			},
			out: struct {
				group *toggl.Group
				err   error
			}{
				group: nil,
				err:   toggl.ErrGroupNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualGroup, err := client.UpdateGroup(c.in.ctx, c.in.group)
			if !reflect.DeepEqual(actualGroup, c.out.group) {
				t.Errorf("\nwant: %+#v\ngot : %+#v\n", c.out.group, actualGroup)
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

func TestUpdateGroupConvertParamsToRequestBody(t *testing.T) {
	expectedGroupRequest := &toggl.Group{
		Wid:  1234567,
		Name: "updated",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualGroupRequest := new(toggl.Group)
		if err := json.Unmarshal(requestBody, actualGroupRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualGroupRequest, expectedGroupRequest) {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedGroupRequest, actualGroupRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateGroup(context.Background(), expectedGroupRequest)
}

func TestUpdateGroupUseURLIncludingGroupId(t *testing.T) {
	groupId := 123456
	expectedRequestURI := "/api/v8/groups/" + strconv.Itoa(groupId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateGroup(context.Background(), &toggl.Group{
		Id:   groupId,
		Wid:  1234567,
		Name: "toggl-go",
	})
}

func TestDeleteGroup(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx   context.Context
			group *toggl.Group
		}
		out error
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/delete_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Id: 123456,
				},
			},
			out: nil,
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/groups/delete_400_bad_request.html",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Id: 123456,
				},
			},
			out: &toggl.TogglError{
				Message: "<html>\n<head><title>400 Bad Request</title></head>\n<body>\n<center><h1>400 Bad Request</h1></center>\n<hr><center>nginx</center>\n</body>\n</html>\n",
				Code:    400,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/groups/delete_403_forbidden.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: context.Background(),
				group: &toggl.Group{
					Id: 123456,
				},
			},
			out: &toggl.TogglError{
				Message: "",
				Code:    403,
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/delete_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx: nil,
				group: &toggl.Group{
					Id: 123456,
				},
			},
			out: toggl.ErrContextNotFound,
		},
		{
			name:             "Without group",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/groups/delete_200_ok.json",
			in: struct {
				ctx   context.Context
				group *toggl.Group
			}{
				ctx:   context.Background(),
				group: nil,
			},
			out: toggl.ErrGroupNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			err := client.DeleteGroup(c.in.ctx, c.in.group)

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

func TestDeleteGroupUseURLIncludingGroupId(t *testing.T) {
	groupId := 123456
	expectedRequestURI := "/api/v8/groups/" + strconv.Itoa(groupId) + "?"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\nwant: %+#v\ngot : %+#v\n", expectedRequestURI, actualRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_ = client.DeleteGroup(context.Background(), &toggl.Group{
		Id: groupId,
	})
}
