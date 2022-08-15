package track

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	workspacesPath string = "api/v9/workspaces"
)

// Workspace represents the properties of a workspace.
type Workspace struct {
	ID                          *int           `json:"id,omitempty"`
	OrganizationID              *int           `json:"organization_id,omitempty"`
	Name                        *string        `json:"name,omitempty"`
	Profile                     *int           `json:"profile,omitempty"`
	Premium                     *bool          `json:"premium,omitempty"`
	BusinessWs                  *bool          `json:"business_ws,omitempty"`
	Admin                       *bool          `json:"admin,omitempty"`
	SuspendedAt                 *time.Time     `json:"suspended_at,omitempty"`
	ServerDeletedAt             *time.Time     `json:"server_deleted_at,omitempty"`
	DefaultHourlyRate           *int           `json:"default_hourly_rate,omitempty"`
	RateLastUpdated             *string        `json:"rate_last_updated,omitempty"`
	DefaultCurrency             *string        `json:"default_currency,omitempty"`
	OnlyAdminsMayCreateProjects *bool          `json:"only_admins_may_create_projects,omitempty"`
	OnlyAdminsMayCreateTags     *bool          `json:"only_admins_may_create_tags,omitempty"`
	OnlyAdminsSeeBillableRates  *bool          `json:"only_admins_see_billable_rates,omitempty"`
	OnlyAdminsSeeTeamDashboard  *bool          `json:"only_admins_see_team_dashboard,omitempty"`
	ProjectsBillableByDefault   *bool          `json:"projects_billable_by_default,omitempty"`
	ReportsCollapse             *bool          `json:"reports_collapse,omitempty"`
	Rounding                    *int           `json:"rounding,omitempty"`
	RoundingMinutes             *int           `json:"rounding_minutes,omitempty"`
	APIToken                    *string        `json:"api_token,omitempty"`
	At                          *time.Time     `json:"at,omitempty"`
	LogoURL                     *string        `json:"logo_url,omitempty"`
	IcalURL                     *string        `json:"ical_url,omitempty"`
	IcalEnabled                 *bool          `json:"ical_enabled,omitempty"`
	CsvUpload                   *CsvUpload     `json:"csv_upload,omitempty"`
	Subscription                *Subscription  `json:"subscription,omitempty"`
	TeConstraints               *TeConstraints `json:"te_constraints,omitempty"`
}

type CsvUpload struct {
	At    *string `json:"at,omitempty"`
	LogID *int    `json:"log_id,omitempty"`
}

type Subscription struct {
	AutoRenew          *bool           `json:"auto_renew,omitempty"`
	CardDetails        *CardDetails    `json:"card_details,omitempty"`
	CompanyID          *int            `json:"company_id,omitempty"`
	ContactDetail      *ContactDetail  `json:"contact_detail,omitempty"`
	CreatedAt          *time.Time      `json:"created_at,omitempty"`
	Currency           *string         `json:"currency,omitempty"`
	CustomerID         *int            `json:"customer_id,omitempty"`
	DeletedAt          *time.Time      `json:"deleted_at,omitempty"`
	LastPricingPlanID  *int            `json:"last_pricing_plan_id,omitempty"`
	OrganizationID     *int            `json:"organization_id,omitempty"`
	PaymentDetails     *PaymentDetails `json:"payment_details,omitempty"`
	PricingPlanID      *int            `json:"pricing_plan_id,omitempty"`
	RenewalAt          *time.Time      `json:"renewal_at,omitempty"`
	SubscriptionID     *int            `json:"subscription_id,omitempty"`
	SubscriptionPeriod *Period         `json:"subscription_period,omitempty"`
	WorkspaceID        *int            `json:"workspace_id,omitempty"`
}

type CardDetails struct{}

type ContactDetail struct{}

type PaymentDetails struct{}

type Period struct{}

type TeConstraints struct {
	DescriptionPresent          *bool `json:"description_present,omitempty"`
	ProjectPresent              *bool `json:"project_present,omitempty"`
	TagPresent                  *bool `json:"tag_present,omitempty"`
	TaskPresent                 *bool `json:"task_present,omitempty"`
	TimeEntryConstraintsEnabled *bool `json:"time_entry_constraints_enabled,omitempty"`
}

// GetWorkspace gets information of single workspace.
func (c *Client) GetWorkspace(ctx context.Context, workspaceID int) (*Workspace, error) {
	var workspace *Workspace
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID))
	if err := c.httpGet(ctx, apiSpecificPath, nil, &workspace); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return workspace, nil
}

// WorkspaceUser represents the properties of a user who belong to a workspace.
type WorkspaceUser struct {
	ID                *int       `json:"id,omitempty"`
	UserID            *int       `json:"user_id,omitempty"`
	WorkspaceID       *int       `json:"workspace_id,omitempty"`
	Admin             *bool      `json:"admin,omitempty"`
	OrganizationAdmin *bool      `json:"organization_admin,omitempty"`
	WorkspaceAdmin    *bool      `json:"workspace_admin,omitempty"`
	Active            *bool      `json:"active,omitempty"`
	Email             *string    `json:"email,omitempty"`
	Timezone          *string    `json:"timezone,omitempty"`
	Inactive          *bool      `json:"inactive,omitempty"`
	At                *time.Time `json:"at,omitempty"`
	Name              *string    `json:"name,omitempty"`
	Rate              *int       `json:"rate,omitempty"`
	RateLastUpdated   *string    `json:"rate_last_updated,omitempty"`
	LabourCost        *int       `json:"labour_cost,omitempty"`
	InviteURL         *string    `json:"invite_url,omitempty"`
	InvitationCode    *string    `json:"invitation_code,omitempty"`
	AvatarFileName    *string    `json:"avatar_file_name,omitempty"`
	GroupIDs          *groupIDs  `json:"group_ids,omitempty"`
	IsDirect          *bool      `json:"is_direct,omitempty"`
}

type groupIDs struct {
}

// GetWorkspaceUsers returns any users who belong to the workspace directly or through at least one group.
func (c *Client) GetWorkspaceUsers(ctx context.Context, organizationID, workspaceID int) ([]*WorkspaceUser, error) {
	var workspaceUsers []*WorkspaceUser
	apiSpecificPath := path.Join(organizationsPath, strconv.Itoa(organizationID), "workspaces", strconv.Itoa(workspaceID))
	if err := c.httpGet(ctx, apiSpecificPath, nil, &workspaceUsers); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return workspaceUsers, nil
}

// UpdateWorkspaceRequestBody represents a request body of UpdateWorkspace.
type UpdateWorkspaceRequestBody struct {
	Admins                      []*int  `json:"admins,omitempty"`
	DefaultCurrency             *string `json:"default_currency,omitempty"`
	DefaultHourlyRate           *int    `json:"default_hourly_rate,omitempty"`
	InitialPricingPlan          *int    `json:"initial_pricing_plan,omitempty"`
	Name                        *string `json:"name,omitempty"`
	OnlyAdminsMayCreateProjects *bool   `json:"only_admins_may_create_projects,omitempty"`
	OnlyAdminsMayCreateTags     *bool   `json:"only_admins_may_create_tags,omitempty"`
	OnlyAdminsSeeBillableRates  *bool   `json:"only_admins_see_billable_rates,omitempty"`
	OnlyAdminsSeeTeamDashboard  *bool   `json:"only_admins_see_team_dashboard,omitempty"`
	OrganizationID              *int    `json:"organization_id,omitempty"`
	ProjectsBillableByDefault   *bool   `json:"projects_billable_by_default,omitempty"`
	RateChangeMode              *string `json:"rate_change_mode,omitempty"`
	ReportsCollapse             *bool   `json:"reports_collapse,omitempty"`
	Rounding                    *int    `json:"rounding,omitempty"`
	RoundingMinutes             *int    `json:"rounding_minutes,omitempty"`
}

// UpdateWorkspace updates a specific workspace.
func (c *Client) UpdateWorkspace(ctx context.Context, workspaceID int, reqBody *UpdateWorkspaceRequestBody) (*Workspace, error) {
	return nil, nil
}
