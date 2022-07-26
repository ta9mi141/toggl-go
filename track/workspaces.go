package track

import (
	"context"
	"time"
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
	return nil, nil
}
