package webhooks

import (
	"context"
	"errors"
	"net/http"
	"path"
	"reflect"
	"testing"

	"github.com/ta9mi141/toggl-go/track"
	"github.com/ta9mi141/toggl-go/track/internal"
)

func TestGetEventFilters(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			statusCode   int
			testdataFile string
		}
		out struct {
			eventFilters *EventFilters
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
				testdataFile: "testdata/event_filters/get_event_filters_200_ok.json",
			},
			out: struct {
				eventFilters *EventFilters
				err          error
			}{
				eventFilters: &EventFilters{
					Client:        []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
					Project:       []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
					ProjectGroup:  []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
					ProjectUser:   []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
					Tag:           []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
					Task:          []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
					TimeEntry:     []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
					Workspace:     []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
					WorkspaceUser: []*string{track.Ptr("created"), track.Ptr("updated"), track.Ptr("deleted")},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiSpecificPath := path.Join(webhooksPath, "event_filters")
			mockServer := internal.NewMockServer(t, apiSpecificPath, tt.in.statusCode, tt.in.testdataFile)
			defer mockServer.Close()

			apiClient := NewAPIClient(internal.APIToken, withBaseURL(mockServer.URL))
			eventFilters, err := apiClient.GetEventFilters(context.Background())

			if !reflect.DeepEqual(eventFilters, tt.out.eventFilters) {
				internal.Errorf(t, eventFilters, tt.out.eventFilters)
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
