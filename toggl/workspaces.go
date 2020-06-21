package toggl

import (
	"context"
	"errors"
	"strconv"
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
	var workspaces []*Workspace
	if err := c.httpGet(ctx, c.buildURL(workspacesEndpoint), &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil
}

// GetWorkspace gets data about the single workspace.
func (c *Client) GetWorkspace(ctx context.Context, workspace *Workspace) (*Workspace, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	rawWorkspaceData := new(rawWorkspaceData)
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.Id)
	if err := c.httpGet(ctx, c.buildURL(endpoint), rawWorkspaceData); err != nil {
		return nil, err
	}
	return &rawWorkspaceData.Workspace, nil
}

// UpdateWorkspace updates the workspace.
func (c *Client) UpdateWorkspace(ctx context.Context, workspace *Workspace) (*Workspace, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	rawWorkspaceData := new(rawWorkspaceData)
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.Id)
	if err := c.httpPut(ctx, c.buildURL(endpoint), workspace, rawWorkspaceData); err != nil {
		return nil, err
	}
	return &rawWorkspaceData.Workspace, nil
}

// GetWorkspaceUsers gets workspace users.
func (c *Client) GetWorkspaceUsers(ctx context.Context, workspace *Workspace) ([]*User, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	var workspaceUsers []*User
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.Id) + "/users"
	if err := c.httpGet(ctx, c.buildURL(endpoint), &workspaceUsers); err != nil {
		return nil, err
	}
	return workspaceUsers, nil
}

// GetWorkspaceClients gets workspace clients.
func (c *Client) GetWorkspaceClients(ctx context.Context, workspace *Workspace) ([]*TogglClient, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	var workspaceClients []*TogglClient
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.Id) + "/clients"
	if err := c.httpGet(ctx, c.buildURL(endpoint), &workspaceClients); err != nil {
		return nil, err
	}
	return workspaceClients, nil
}

// GetWorkspaceGroups gets workspace groups.
func (c *Client) GetWorkspaceGroups(ctx context.Context, workspace *Workspace) ([]*Group, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	var workspaceGroups []*Group
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.Id) + "/groups"
	if err := c.httpGet(ctx, c.buildURL(endpoint), &workspaceGroups); err != nil {
		return nil, err
	}
	return workspaceGroups, nil
}

// GetWorkspaceProjects gets workspace projects.
func (c *Client) GetWorkspaceProjects(ctx context.Context, workspace *Workspace, params ...QueryString) ([]*Project, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	var workspaceProjects []*Project
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.Id) + "/projects"
	if err := c.httpGet(ctx, c.buildURL(endpoint, params...), &workspaceProjects); err != nil {
		return nil, err
	}
	return workspaceProjects, nil
}

// GetWorkspaceTags gets workspace tags.
func (c *Client) GetWorkspaceTags(ctx context.Context, workspace *Workspace) ([]*Tag, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	var workspaceTags []*Tag
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.Id) + "/tags"
	if err := c.httpGet(ctx, c.buildURL(endpoint), &workspaceTags); err != nil {
		return nil, err
	}
	return workspaceTags, nil
}
