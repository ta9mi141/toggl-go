package toggl

import (
	"context"
	"time"
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
	return nil, nil
}
