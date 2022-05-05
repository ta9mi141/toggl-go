package toggl

import (
	"context"
	"time"
)

type Project struct {
	ID            *int       `json:"id,omitempty"`
	WID           *int       `json:"wid",omitempty`
	CID           *int       `json:"cid,omitempty"`
	Name          *string    `json:"Name,omitempty"`
	Billable      *bool      `json:"billable",omitempty`
	IsPrivate     *bool      `json:"is_private",omitempty`
	Active        *bool      `json:"active",omitempty`
	Template      *bool      `json:"template",omitempty`
	At            *time.Time `json:"at,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	Color         *string    `json:"color,omitempty"`
	AutoEstimates *bool      `json:"auto_estimates",omitempty`
	ActualHours   *int       `json:"actual_hours,omitempty"`
	HexColor      *string    `json:"hex_color,omitempty"`
}

const (
	projectsPath string = "api/v8/projects"
)

// GetProject gets the project.
func (c *Client) GetProject(ctx context.Context, id int) (*Project, error) {
	return nil, nil
}
