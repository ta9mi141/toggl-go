package toggl

import (
	"context"
	"errors"
	"time"
)

const (
	workspacesEndpoint string = "/api/v8/workspaces"
)

var (
	// ErrWorkspaceNotFound is returned when the provided workspace is nil.
	ErrWorkspaceNotFound = errors.New("The provided workspace must be non-nil")
)

// Workspace represents properties of workspace.
type Workspace struct {
	Id                          int       `json:"id"`
	Name                        string    `json:"name"`
	Premium                     bool      `json:"premium"`
	Admin                       bool      `json:"admin"`
	DefaultHourlyRate           float64   `json:"default_hourly_rate"`
	DefaultCurrency             string    `json:"default_currency"`
	OnlyAdminsMayCreateProjects bool      `json:"only_admins_may_create_projects"`
	OnlyAdminsSeeBillableRates  bool      `json:"only_admins_see_billable_rates"`
	Rounding                    int       `json:"rounding"`
	RoundingMinutes             int       `json:"rounding_minutes"`
	At                          time.Time `json:"at"`
	LogoURL                     string    `json:"logo_url"`
}

type rawWorkspaceData struct {
	Workspace Workspace `json:"data"`
}

// GetWorkspaces gets data about all the workspaces where the token owner belongs to.
func (c *Client) GetWorkspaces(ctx context.Context) ([]*Workspace, error) {
	return nil, nil
}

// GetWorkspace gets data about the single workspace.
func (c *Client) GetWorkspace(ctx context.Context, workspace *Workspace) (*Workspace, error) {
	return nil, nil
}

// UpdateWorkspace updates the workspace.
func (c *Client) UpdateWorkspace(ctx context.Context, workspace *Workspace) (*Workspace, error) {
	return nil, nil
}

// GetWorkspaceUsers gets workspace users.
func (c *Client) GetWorkspaceUsers(ctx context.Context, workspace *Workspace) ([]*User, error) {
	return nil, nil
}

// GetWorkspaceClients gets workspace clients.
func (c *Client) GetWorkspaceClients(ctx context.Context, workspace *Workspace) ([]*Client, error) {
	return nil, nil
}

// GetWorkspaceGroups gets workspace groups.
func (c *Client) GetWorkspaceGroups(ctx context.Context, workspace *Workspace) ([]*Group, error) {
	return nil, nil
}

// GetWorkspaceProjects gets workspace projects.
func (c *Client) GetWorkspaceProjects(ctx context.Context, workspace *Workspace, params ...QueryString) ([]*Project, error) {
	return nil, nil
}

// GetWorkspaceTags gets workspace tags.
func (c *Client) GetWorkspaceTags(ctx context.Context, workspace *Workspace) ([]*Tag, error) {
	return nil, nil
}
