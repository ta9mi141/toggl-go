package webhooks

import (
	"context"
	"path"

	"github.com/pkg/errors"
)

// EventFilters represents the properties of event filters.
type EventFilters struct {
	Client        []*string `json:"client,omitempty"`
	Project       []*string `json:"project,omitempty"`
	ProjectGroup  []*string `json:"project_group,omitempty"`
	ProjectUser   []*string `json:"project_user,omitempty"`
	Tag           []*string `json:"tag,omitempty"`
	Task          []*string `json:"task,omitempty"`
	TimeEntry     []*string `json:"time_entry,omitempty"`
	Workspace     []*string `json:"workspace,omitempty"`
	WorkspaceUser []*string `json:"workspace_user,omitempty"`
}

// GetEventFilters gets the list of supported event filters.
func (c *APIClient) GetEventFilters(ctx context.Context) (*EventFilters, error) {
	var eventFilters *EventFilters
	apiSpecificPath := path.Join(webhooksPath, "event_filters")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &eventFilters); err != nil {
		return nil, errors.Wrap(err, "failed to get event filters")
	}
	return eventFilters, nil
}
