package toggl

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	workspacesPath string = "api/v8/workspaces"
)

// Workspace represents properties of workspace.
// Some properties not listed in the documentation are also included.
type Workspace struct {
	ID                          *int       `json:"id,omitempty"`
	Name                        *string    `json:"name,omitempty"`
	Profile                     *int       `json:"profile,omitempty"`
	Premium                     *bool      `json:"premium,omitempty"`
	Admin                       *bool      `json:"admin,omitempty"`
	DefaultHourlyRate           *float64   `json:"default_hourly_rate,omitempty"`
	DefaultCurrency             *string    `json:"default_currency,omitempty"`
	OnlyAdminsMayCreateProjects *bool      `json:"only_admins_may_create_projects,omitempty"`
	OnlyAdminsSeeBillableRates  *bool      `json:"only_admins_see_billable_rates,omitempty"`
	OnlyAdminsSeeTeamDashboard  *bool      `json:"only_admins_see_team_dashboard,omitempty"`
	ProjectsBillableByDefault   *bool      `json:"projects_billable_by_default,omitempty"`
	Rounding                    *int       `json:"rounding,omitempty"`
	RoundingMinutes             *int       `json:"rounding_minutes,omitempty"`
	APIToken                    *string    `json:"api_token,omitempty"`
	At                          *time.Time `json:"at,omitempty"`
	IcalEnabled                 *bool      `json:"ical_enabled,omitempty"`
	LogoURL                     *string    `json:"logo_url,omitempty"`
	IcalURL                     *string    `json:"ical_url,omitempty"`
}

type workspaceRequest struct {
	Workspace Workspace `json:"workspace"`
}

type workspaceResponse struct {
	Workspace Workspace `json:"data"`
}

// GetWorkspaces gets all the workspaces where the token owner belongs to.
func (c *Client) GetWorkspaces(ctx context.Context) ([]*Workspace, error) {
	var workspaces []*Workspace
	if err := c.httpGet(ctx, workspacesPath, &workspaces); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return workspaces, nil
}

// GetWorkspace gets the single workspace.
func (c *Client) GetWorkspace(ctx context.Context, id int) (*Workspace, error) {
	response := new(workspaceResponse)
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(id))
	if err := c.httpGet(ctx, apiSpecificPath, response); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return &response.Workspace, nil
}

// GetWorkspaceUsers gets the workspace users.
func (c *Client) GetWorkspaceUsers(ctx context.Context, id int) ([]*User, error) {
	var users []*User
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(id), "users")
	if err := c.httpGet(ctx, apiSpecificPath, &users); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return users, nil
}

// GetWorkspaceProjectsParameter is the additional parameter of GetWorkspaceProjects.
type GetWorkspaceProjectsParameter requestParameter

// Active determines whether the request filters projects by their state.
func Active(active string) GetWorkspaceProjectsParameter {
	return activeParameter(active)
}

type activeParameter string

func (p activeParameter) apply() {}

// ActualHours determines whether the request gets the completed hours per project.
func ActualHours(actualHours bool) GetWorkspaceProjectsParameter {
	return actualHoursParameter(actualHours)
}

type actualHoursParameter bool

func (p actualHoursParameter) apply() {}

// OnlyTemplates determines whether the request gets only project templates.
func OnlyTemplates(onlyTemplates bool) GetWorkspaceProjectsParameter {
	return onlyTemplatesParameter(onlyTemplates)
}

type onlyTemplatesParameter bool

func (p onlyTemplatesParameter) apply() {}

// GetWorkspaceProjects gets the workspace projects.
func (c *Client) GetWorkspaceProjects(ctx context.Context, id int, params ...GetWorkspaceProjectsParameter) ([]*Project, error) {
	var projects []*Project
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(id), "projects")
	if err := c.httpGet(ctx, apiSpecificPath, &projects); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return projects, nil
}

// UpdateWorkspace updates the workspace.
func (c *Client) UpdateWorkspace(ctx context.Context, id int, workspace *Workspace) (*Workspace, error) {
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(id))
	request := &workspaceRequest{Workspace: *workspace}
	response := new(workspaceResponse)

	if err := c.httpPut(ctx, apiSpecificPath, request, response); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return &response.Workspace, nil
}
