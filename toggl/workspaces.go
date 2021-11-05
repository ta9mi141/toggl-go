package toggl

import (
	"context"
	"time"
)

// Workspace represents properties of workspace.
// Some properties not listed in the documentation are also included.
type Workspace struct {
	ID                          int       `json:"id"`
	Name                        string    `json:"name"`
	Profile                     int       `json:"profile"`
	Premium                     bool      `json:"premium"`
	Admin                       bool      `json:"admin"`
	DefaultHourlyRate           float64   `json:"default_hourly_rate"`
	DefaultCurrency             string    `json:"default_currency"`
	OnlyAdminsMayCreateProjects bool      `json:"only_admins_may_create_projects"`
	OnlyAdminsSeeBillableRates  bool      `json:"only_admins_see_billable_rates"`
	OnlyAdminsSeeTeamDashboard  bool      `json:"only_admins_see_team_dashboard"`
	ProjectsBillableByDefault   bool      `json:"projects_billable_by_default"`
	Rounding                    int       `json:"rounding"`
	RoundingMinutes             int       `json:"rounding_minutes"`
	APIToken                    string    `json:"api_token"`
	At                          time.Time `json:"at"`
	IcalEnabled                 bool      `json:"ical_enabled"`
	LogoURL                     string    `json:"logo_url,omitempty"`
	IcalURL                     string    `json:"ical_url,omitempty"`
}

// GetWorkspaces gets data about all the workspaces where the token owner belongs to.
func (c *Client) GetWorkspaces(ctx context.Context) ([]*Workspace, error) {
	var workspaces []*Workspace
	if err := c.httpGet(ctx, "workspaces", &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil
}

// GetWorkspace gets data about the single workspace.
func (c *Client) GetWorkspace(ctx context.Context, id int) (*Workspace, error) {
	return nil, nil
}
