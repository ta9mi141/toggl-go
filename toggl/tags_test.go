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

	"github.com/it-akumi/toggl-go/toggl"
)

func TestCreateTag(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx context.Context
			tag *toggl.Tag
		}
		out struct {
			tag *toggl.Tag
			err error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/create_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: &toggl.Tag{
					Id:   1234567,
					Wid:  1234567,
					Name: "toggl-go",
				},
				err: nil,
			},
		},
		{
			name:             "400 Bad Request",
			httpStatus:       http.StatusBadRequest,
			testdataFilePath: "testdata/tags/create_400_bad_request.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: nil,
				err: &toggl.TogglError{
					Message: "Tag already exists: toggl-go\n",
					Code:    400,
				},
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/tags/create_403_forbidden.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/create_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: nil,
				tag: &toggl.Tag{
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: nil,
				err: toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without tag",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/create_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: nil,
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: nil,
				err: toggl.ErrTagNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTag, err := client.CreateTag(c.in.ctx, c.in.tag)
			if !reflect.DeepEqual(actualTag, c.out.tag) {
				t.Errorf("\ngot : %+#v\nwant: %+#v\n", actualTag, c.out.tag)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\ngot : %#+v\nwant: %#+v\n", togglError, c.out.err)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\ngot : %#+v\nwant: %#+v\n", err, c.out.err)
				}
			}
		})
	}
}

func TestCreateTagConvertParamsToRequestBody(t *testing.T) {
	expectedTagRequest := &toggl.Tag{
		Wid:  1234567,
		Name: "toggl-go",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		actualTagRequest := new(toggl.Tag)
		if err := json.Unmarshal(requestBody, actualTagRequest); err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(actualTagRequest, expectedTagRequest) {
			t.Errorf("\ngot: %+#v\nwant: %+#v\n", actualTagRequest, expectedTagRequest)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.CreateTag(context.Background(), expectedTagRequest)
}

func TestUpdateTag(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx context.Context
			tag *toggl.Tag
		}
		out struct {
			tag *toggl.Tag
			err error
		}
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/update_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Id:   1234567,
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: &toggl.Tag{
					Id:   1234567,
					Wid:  1234567,
					Name: "toggl-go",
				},
				err: nil,
			},
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/tags/update_403_forbidden.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Id:   1234567,
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: nil,
				err: &toggl.TogglError{
					Message: "",
					Code:    403,
				},
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/tags/update_404_not_found.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Id:   1234567,
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: nil,
				err: &toggl.TogglError{
					Message: "null",
					Code:    404,
				},
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/update_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: nil,
				tag: &toggl.Tag{
					Id:   1234567,
					Wid:  1234567,
					Name: "toggl-go",
				},
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: nil,
				err: toggl.ErrContextNotFound,
			},
		},
		{
			name:             "Without tag",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/update_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: nil,
			},
			out: struct {
				tag *toggl.Tag
				err error
			}{
				tag: nil,
				err: toggl.ErrTagNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			actualTag, err := client.UpdateTag(c.in.ctx, c.in.tag)
			if !reflect.DeepEqual(actualTag, c.out.tag) {
				t.Errorf("\ngot : %+#v\nwant: %+#v\n", actualTag, c.out.tag)
			}

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out.err) {
					t.Errorf("\ngot : %#+v\nwant: %#+v\n", togglError, c.out.err)
				}
			} else {
				if !errors.Is(err, c.out.err) {
					t.Errorf("\ngot : %#+v\nwant: %#+v\n", err, c.out.err)
				}
			}
		})
	}
}

func TestUpdateTagUseURLIncludingTagId(t *testing.T) {
	tagId := 1234567
	expectedRequestURI := "/api/v8/tags/" + strconv.Itoa(tagId)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\ngot: %+#v\nwant: %+#v\n", actualRequestURI, expectedRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_, _ = client.UpdateTag(context.Background(), &toggl.Tag{
		Id:   tagId,
		Wid:  1234567,
		Name: "toggl-go",
	})
}

func TestDeleteTag(t *testing.T) {
	cases := []struct {
		name             string
		httpStatus       int
		testdataFilePath string
		in               struct {
			ctx context.Context
			tag *toggl.Tag
		}
		out error
	}{
		{
			name:             "200 OK",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/delete_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Id: 1234567,
				},
			},
			out: nil,
		},
		{
			name:             "403 Forbidden",
			httpStatus:       http.StatusForbidden,
			testdataFilePath: "testdata/tags/delete_403_forbidden.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Id: 1234567,
				},
			},
			out: &toggl.TogglError{
				Message: "",
				Code:    403,
			},
		},
		{
			name:             "404 Not Found",
			httpStatus:       http.StatusNotFound,
			testdataFilePath: "testdata/tags/delete_404_not_found.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: &toggl.Tag{
					Id: 1234567,
				},
			},
			out: &toggl.TogglError{
				Message: "null",
				Code:    404,
			},
		},
		{
			name:             "Without context",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/delete_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: nil,
				tag: &toggl.Tag{
					Id: 1234567,
				},
			},
			out: toggl.ErrContextNotFound,
		},
		{
			name:             "Without tag",
			httpStatus:       http.StatusOK,
			testdataFilePath: "testdata/tags/delete_200_ok.json",
			in: struct {
				ctx context.Context
				tag *toggl.Tag
			}{
				ctx: context.Background(),
				tag: nil,
			},
			out: toggl.ErrTagNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockServer := setupMockServer(t, c.httpStatus, c.testdataFilePath)
			defer mockServer.Close()

			client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
			err := client.DeleteTag(c.in.ctx, c.in.tag)

			var togglError toggl.Error
			if errors.As(err, &togglError) {
				if !reflect.DeepEqual(togglError, c.out) {
					t.Errorf("\ngot : %#+v\nwant: %#+v\n", togglError, c.out)
				}
			} else {
				if !errors.Is(err, c.out) {
					t.Errorf("\ngot : %#+v\nwant: %#+v\n", err, c.out)
				}
			}
		})
	}
}

func TestDeleteTagUseURLIncludingTagId(t *testing.T) {
	tagId := 1234567
	expectedRequestURI := "/api/v8/tags/" + strconv.Itoa(tagId)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestURI := r.URL.RequestURI()
		if actualRequestURI != expectedRequestURI {
			t.Errorf("\ngot: %+#v\nwant: %+#v\n", actualRequestURI, expectedRequestURI)
		}
	}))

	client := toggl.NewClient(toggl.APIToken(apiToken), baseURL(mockServer.URL))
	_ = client.DeleteTag(context.Background(), &toggl.Tag{
		Id: tagId,
	})
}
