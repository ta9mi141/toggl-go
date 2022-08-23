package toggl

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Tag represents the properties of a tag.
type Tag struct {
	ID          *int       `json:"id,omitempty"`
	WorkspaceID *int       `json:"workspace_id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	At          *time.Time `json:"at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// GetTags lists workspace tags.
func (c *Client) GetTags(ctx context.Context, workspaceID int) ([]*Tag, error) {
	var tags []*Tag
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "tags")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &tags); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return tags, nil
}
