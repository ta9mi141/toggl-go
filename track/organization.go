package track

import "time"

// Organization represents the properties of an organization.
type Organization struct {
	ID                      *int       `json:"id,omitempty"`
	Name                    *string    `json:"name,omitempty"`
	PricingPlanID           *int       `json:"pricing_plan_id,omitempty"`
	CreatedAt               *time.Time `json:"created_at,omitempty"`
	At                      *time.Time `json:"at,omitempty"`
	ServerDeletedAt         *time.Time `json:"server_deleted_at,omitempty"`
	IsMultiWorkspaceEnabled *bool      `json:"is_multi_workspace_enabled,omitempty"`
	SuspendedAt             *string    `json:"suspended_at,omitempty"`
	UserCount               *int       `json:"user_count,omitempty"`
	TrialInfo               *TrialInfo `json:"trial_info,omitempty"`
	IsChargify              *bool      `json:"is_chargify,omitempty"`
	MaxWorkspaces           *int       `json:"max_workspaces,omitempty"`
	Admin                   *bool      `json:"admin,omitempty"`
	Owner                   *bool      `json:"owner,omitempty"`
}

type TrialInfo struct {
	Trial             *bool   `json:"trial,omitempty"`
	TrialAvailable    *bool   `json:"trial_available,omitempty"`
	TrialEndDate      *string `json:"trial_end_date,omitempty"`
	NextPaymentDate   *string `json:"next_payment_date,omitempty"`
	LastPricingPlanID *int    `json:"last_pricing_plan_id,omitempty"`
}
