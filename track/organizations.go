package track

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	organizationsPath string = "api/v9/organizations"
)

// Organization represents the properties of an organization.
type Organization struct {
	ID                      *int       `json:"id,omitempty"`
	Name                    *string    `json:"name,omitempty"`
	PricingPlanID           *int       `json:"pricing_plan_id,omitempty"`
	CreatedAt               *time.Time `json:"created_at,omitempty"`
	At                      *time.Time `json:"at,omitempty"`
	ServerDeletedAt         *time.Time `json:"server_deleted_at,omitempty"`
	IsMultiWorkspaceEnabled *bool      `json:"is_multi_workspace_enabled,omitempty"`
	SuspendedAt             *time.Time `json:"suspended_at,omitempty"`
	UserCount               *int       `json:"user_count,omitempty"`
	TrialInfo               *TrialInfo `json:"trial_info,omitempty"`
	IsChargify              *bool      `json:"is_chargify,omitempty"`
	IsUnified               *bool      `json:"is_unified,omitempty"`
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

// GetOrganization returns organization name and current pricing plan.
func (c *Client) GetOrganization(ctx context.Context, organizationID int) (*Organization, error) {
	var organization *Organization
	apiSpecificPath := path.Join(organizationsPath, strconv.Itoa(organizationID))
	if err := c.httpGet(ctx, apiSpecificPath, nil, &organization); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return organization, nil
}
