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
func (c *APIClient) GetTags(ctx context.Context, workspaceID int) ([]*Tag, error) {
	var tags []*Tag
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "tags")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &tags); err != nil {
		return nil, errors.Wrap(err, "failed to get tags")
	}
	return tags, nil
}

// CreateTagRequestBody represents a request body of CreateTag.
type CreateTagRequestBody struct {
	Name        *string `json:"name,omitempty"`
	WorkspaceID *int    `json:"workspace_id,omitempty"`
}

// CreateTag creates workspace tags.
func (c *APIClient) CreateTag(ctx context.Context, workspaceID int, reqBody *CreateTagRequestBody) (*Tag, error) {
	var tag *Tag
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "tags")
	if err := c.httpPost(ctx, apiSpecificPath, reqBody, &tag); err != nil {
		return nil, errors.Wrap(err, "failed to create tag")
	}
	return tag, nil
}

// UpdateTagRequestBody represents a request body of UpdateTag.
type UpdateTagRequestBody struct {
	Name        *string `json:"name,omitempty"`
	WorkspaceID *int    `json:"workspace_id,omitempty"`
}

// UpdateTag updates workspace tags.
func (c *APIClient) UpdateTag(ctx context.Context, workspaceID, tagID int, reqBody *UpdateTagRequestBody) (*Tag, error) {
	var tag *Tag
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "tags", strconv.Itoa(tagID))
	if err := c.httpPut(ctx, apiSpecificPath, reqBody, &tag); err != nil {
		return nil, errors.Wrap(err, "failed to update tag")
	}
	return tag, nil
}

// DeleteTag deletes workspace tags.
func (c *APIClient) DeleteTag(ctx context.Context, workspaceID, tagID int) error {
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "tags", strconv.Itoa(tagID))
	if err := c.httpDelete(ctx, apiSpecificPath); err != nil {
		return errors.Wrap(err, "failed to delete tag")
	}
	return nil
}
