package toggl

import (
	"context"
	"strconv"
	"time"
)

const (
	dashboardEndpoint string = "/api/v8/dashboard"
)

// Dashboard represents properties of two objects dashboard request returns.
type Dashboard struct {
	Activity []struct {
		UserId      int       `json:"user_id"`
		ProjectId   int       `json:"project_id"`
		Duration    int       `json:"duration"`
		Description string    `json:"description"`
		Stop        time.Time `json:"stop"`
		Tid         int       `json:"tid"`
	} `json:"activity"`
	MostActiveUser []struct {
		UserId   int `json:"user_id"`
		Duration int `json:"duration"`
	} `json:"most_active_user"`
}

// GetDashboard gets a dashboard.
func (c *Client) GetDashboard(ctx context.Context, workspace *Workspace) (*Dashboard, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	dashboard := new(Dashboard)
	endpoint := dashboardEndpoint + "/" + strconv.Itoa(workspace.Id)
	if err := c.httpGet(ctx, c.buildURL(endpoint), dashboard); err != nil {
		return nil, err
	}
	return dashboard, nil
}
