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

// OrganizationUser represents the properties of a user in an organization.
type OrganizationUser struct {
	ID             *int        `json:"id,omitempty"`
	Name           *string     `json:"name,omitempty"`
	Email          *string     `json:"email,omitempty"`
	UserID         *int        `json:"user_id,omitempty"`
	AvatarURL      *string     `json:"avatar_url,omitempty"`
	Admin          *bool       `json:"admin,omitempty"`
	Owner          *bool       `json:"owner,omitempty"`
	Joined         *bool       `json:"joined,omitempty"`
	InvitationCode *string     `json:"invitation_code,omitempty"`
	Inactive       *bool       `json:"inactive,omitempty"`
	CanEditEmail   *bool       `json:"can_edit_email,omitempty"`
	Workspaces     *workspaces `json:"workspaces,omitempty"`
	Groups         *groups     `json:"groups,omitempty"`
}

type workspaces struct {
	Items []*workspace `json:"items,omitempty"`
}

type workspace struct {
	Admin       *bool   `json:"admin,omitempty"`
	Name        *string `json:"name,omitempty"`
	WorkspaceID *int    `json:"workspace_id,omitempty"`
}

type groups struct {
	Items []*group `json:"items,omitempty"`
}

type group struct {
	GroupID *int    `json:"group_id,omitempty"`
	Name    *string `json:"name,omitempty"`
}

// GetOrganizationUsersQuery represents the additional parameters of GetOrganizationUsers.
type GetOrganizationUsersQuery struct {
	Filter       *string `url:"filter,omitempty"`
	ActiveStatus *string `url:"active_status,omitempty"`
	OnlyAdmins   *string `url:"only_admins,omitempty"`
	Groups       *string `url:"groups,omitempty"`
	Workspaces   *string `url:"workspaces,omitempty"`
	Page         *int    `url:"page,omitempty"`
	PerPage      *int    `url:"per_page,omitempty"`
	SortDir      *string `url:"sort_dir,omitempty"`
}

// GetOrganizationUsers returns list of users in an organization.
func (c *Client) GetOrganizationUsers(ctx context.Context, organizationID int, query *GetOrganizationUsersQuery) ([]*OrganizationUser, error) {
	return nil, nil
}
