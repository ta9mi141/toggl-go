package toggl

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Project represents the properties of a project.
type Project struct {
	ID                  *int                 `json:"id,omitempty"`
	WorkspaceID         *int                 `json:"workspace_id,omitempty"`
	ClientID            *int                 `json:"client_id,omitempty"`
	Name                *string              `json:"name,omitempty"`
	IsPrivate           *bool                `json:"is_private,omitempty"`
	Active              *bool                `json:"active,omitempty"`
	At                  *time.Time           `json:"at,omitempty"`
	CreatedAt           *time.Time           `json:"created_at,omitempty"`
	ServerDeletedAt     *time.Time           `json:"server_deleted_at,omitempty"`
	Color               *string              `json:"color,omitempty"`
	Billable            *bool                `json:"billable,omitempty"`
	Template            *bool                `json:"template,omitempty"`
	AutoEstimates       *bool                `json:"auto_estimates,omitempty"`
	EstimatedHours      *int                 `json:"estimated_hours,omitempty"`
	Rate                *int                 `json:"rate,omitempty"`
	RateLastUpdated     *string              `json:"rate_last_updated,omitempty"`
	Currency            *string              `json:"currency,omitempty"`
	Recurring           *bool                `json:"recurring,omitempty"`
	RecurringParameters *recurringParameters `json:"recurring_parameters,omitempty"`
	CurrentPeriod       *currentPeriod       `json:"current_period,omitempty"`
	FixedFee            *int                 `json:"fixed_fee,omitempty"`
	ActualHours         *int                 `json:"actual_hours,omitempty"`
	WID                 *int                 `json:"wid,omitempty"`
	CID                 *int                 `json:"cid,omitempty"`
	ForeignID           *string              `json:"foreign_id,omitempty"`
	FirstTimeEntry      *string              `json:"first_time_entry,omitempty"`
}

type recurringParameters struct {
	Items []*recurringParameter `json:"items,omitempty"`
}

type recurringParameter struct {
	CustomPeriod       *int    `json:"custom_period,omitempty"`
	EstimatedSeconds   *int    `json:"estimated_seconds,omitempty"`
	ParameterStartDate *string `json:"parameter_start_date,omitempty"`
	ParameterEndDate   *string `json:"parameter_end_date,omitempty"`
	Period             *string `json:"period,omitempty"`
	ProjectStartDate   *string `json:"project_start_date,omitempty"`
}

type currentPeriod struct {
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
}

// GetProjectsQuery represents the additional parameters of GetProjects.
// Currently user_ids, client_ids, and group_ids are not supported.
type GetProjectsQuery struct {
	Active        *bool   `url:"active,omitempty"`
	Since         *int    `url:"since,omitempty"`
	Billable      *bool   `url:"billable,omitempty"`
	Name          *string `url:"name,omitempty"`
	Page          *int    `url:"page,omitempty"`
	SortField     *string `url:"sort_field,omitempty"`
	SortOrder     *string `url:"sort_order,omitempty"`
	OnlyTemplates *bool   `url:"only_templates,omitempty"`
}

// GetProjects gets projects for given workspace.
func (c *Client) GetProjects(ctx context.Context, workspaceID int, query *GetProjectsQuery) ([]*Project, error) {
	var projects []*Project
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "projects")
	if err := c.httpGet(ctx, apiSpecificPath, query, &projects); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return projects, nil
}

// GetProjectQuery represents the additional parameters of GetProject.
type GetProjectQuery struct {
	WithFirstTimeEntry *bool `url:"with_first_time_entry,omitempty"`
}

// GetProject gets project for given workspace.
func (c *Client) GetProject(ctx context.Context, workspaceID, projectID int, query *GetProjectQuery) (*Project, error) {
	var project *Project
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "projects", strconv.Itoa(projectID))
	if err := c.httpGet(ctx, apiSpecificPath, query, &project); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return project, nil
}

// CreateProjectRequestBody represents a request body of CreateProject.
type CreateProjectRequestBody struct {
	Active              *bool                `json:"active,omitempty"`
	AutoEstimates       *bool                `json:"auto_estimates,omitempty"`
	Billable            *bool                `json:"billable,omitempty"`
	CID                 *int                 `json:"cid,omitempty"`
	ClientID            *int                 `json:"client_id,omitempty"`
	ClientName          *string              `json:"client_name,omitempty"`
	Color               *string              `json:"color,omitempty"`
	Currency            *string              `json:"currency,omitempty"`
	EstimatedHours      *int                 `json:"estimated_hours,omitempty"`
	FixedFee            *int                 `json:"fixed_fee,omitempty"`
	ForeignID           *string              `json:"foreign_id,omitempty"`
	IsPrivate           *bool                `json:"is_private,omitempty"`
	Name                *string              `json:"name,omitempty"`
	PostedFields        []*string            `json:"postedFields,omitempty"`
	Rate                *int                 `json:"rate,omitempty"`
	RateChangeMode      *string              `json:"rate_change_mode,omitempty"`
	Recurring           *bool                `json:"recurring,omitempty"`
	RecurringParameters *recurringParameters `json:"recurring_parameters,omitempty"`
	Template            *bool                `json:"template,omitempty"`
	TemplateID          *int                 `json:"template_id,omitempty"`
}

// CreateProject creates project for given workspace.
func (c *Client) CreateProject(ctx context.Context, workspaceID int, reqBody *CreateProjectRequestBody) (*Project, error) {
	var project *Project
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "projects")
	if err := c.httpPost(ctx, apiSpecificPath, reqBody, &project); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return project, nil
}

// DeleteProject deletes project for given workspace.
func (c *Client) DeleteProject(ctx context.Context, workspaceID, projectID int) error {
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "projects", strconv.Itoa(projectID))
	if err := c.httpDelete(ctx, apiSpecificPath); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
