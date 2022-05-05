package toggl

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	projectsPath string = "api/v8/projects"
)

type Project struct {
	ID            *int       `json:"id,omitempty"`
	WID           *int       `json:"wid,omitempty"`
	CID           *int       `json:"cid,omitempty"`
	Name          *string    `json:"Name,omitempty"`
	Billable      *bool      `json:"billable,omitempty"`
	IsPrivate     *bool      `json:"is_private,omitempty"`
	Active        *bool      `json:"active,omitempty"`
	Template      *bool      `json:"template,omitempty"`
	At            *time.Time `json:"at,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	Color         *string    `json:"color,omitempty"`
	AutoEstimates *bool      `json:"auto_estimates,omitempty"`
	ActualHours   *int       `json:"actual_hours,omitempty"`
	HexColor      *string    `json:"hex_color,omitempty"`
}

type projectResponse struct {
	Project Project `json:"data"`
}

// GetProject gets the project.
func (c *Client) GetProject(ctx context.Context, id int) (*Project, error) {
	response := new(projectResponse)
	apiSpecificPath := path.Join(projectsPath, strconv.Itoa(id))
	if err := c.httpGet(ctx, apiSpecificPath, response); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return &response.Project, nil
}
