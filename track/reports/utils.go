package reports

import (
	"context"
	"path"
	"strconv"

	"github.com/pkg/errors"
)

// Project represents the properties of a filtered project.
type Project struct {
	ID       *int    `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	ClientID *int    `json:"client_id,omitempty"`
	Color    *string `json:"color,omitempty"`
	Active   *bool   `json:"active,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Billable *bool   `json:"billable,omitempty"`
}

// ListProjectsRequestBody represents a request body of ListProjects.
type ListProjectsRequestBody struct {
	ClientIDs  []*int  `json:"client_ids,omitempty"`
	Currency   *string `json:"currency,omitempty"`
	IDs        []*int  `json:"ids,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	IsBillable *bool   `json:"is_billable,omitempty"`
	IsPrivate  *bool   `json:"is_private,omitempty"`
	Name       *string `json:"name,omitempty"`
	Start      *int    `json:"start,omitempty"`
}

// ListProjects returns filtered projects from a workspace.
func (c *APIClient) ListProjects(ctx context.Context, workspaceID int, reqBody *ListProjectsRequestBody) ([]*Project, error) {
	var projects []*Project
	apiSpecificPath := path.Join(reportsPath, strconv.Itoa(workspaceID), "filters/projects")
	if err := c.httpPost(ctx, apiSpecificPath, reqBody, &projects); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return projects, nil
}
