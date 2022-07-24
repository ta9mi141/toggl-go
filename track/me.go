package track

import (
	"context"
	"path"
	"time"

	"github.com/pkg/errors"
)

const (
	mePath string = "api/v9/me"
)

// Me represents the properties of an user.
// Some properties not listed in the documentation are also included.
type Me struct {
	ID                 *int       `json:"id,omitempty"`
	APIToken           *string    `json:"api_token,omitempty"`
	Email              *string    `json:"email,omitempty"`
	Fullname           *string    `json:"fullname,omitempty"`
	Timezone           *string    `json:"timezone,omitempty"`
	DefaultWorkspaceID *int       `json:"default_workspace_id,omitempty"`
	BeginningOfWeek    *int       `json:"beginning_of_week,omitempty"`
	ImageURL           *string    `json:"image_url,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	OpenIDEmail        *bool      `json:"openid_email,omitempty"`
	OpenIDEnabled      *bool      `json:"openid_enabled,omitempty"`
	CountryID          *int       `json:"country_id,omitempty"`
	At                 *time.Time `json:"at,omitempty"`
	IntercomHash       *string    `json:"intercom_hash,omitempty"`
	HasPassword        *bool      `json:"has_password,omitempty"`
	Options            struct{}   `json:"options,omitempty"`
}

// GetMe returns details for the current user.
func (c *Client) GetMe(ctx context.Context) (*Me, error) {
	var me *Me
	if err := c.httpGet(ctx, mePath, nil, &me); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return me, nil
}

// UpdateMeRequestBody represents a request body of UpdateMe.
type UpdateMeRequestBody struct {
	BeginningOfWeek    *int    `json:"beginning_of_week,omitempty"`
	CountryID          *int    `json:"country_id,omitempty"`
	CurrentPassword    *string `json:"current_password,omitempty"`
	DefaultWorkspaceID *int    `json:"default_workspace_id,omitempty"`
	Email              *string `json:"email,omitempty"`
	Fullname           *string `json:"fullname,omitempty"`
	Password           *string `json:"password,omitempty"`
	Timezone           *string `json:"timezone,omitempty"`
}

// UpdateMe updates details for the current user.
func (c *Client) UpdateMe(ctx context.Context, reqBody *UpdateMeRequestBody) (*Me, error) {
	var me *Me
	if err := c.httpPut(ctx, mePath, reqBody, &me); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return me, nil
}

// GetMyOrganizations gets all organizations a given user is part of.
func (c *Client) GetMyOrganizations(ctx context.Context) ([]*Organization, error) {
	var organizations []*Organization
	apiSpecificPath := path.Join(mePath, "organizations")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &organizations); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return organizations, nil
}

// GetProjectsQuery represents the additional parameters of GetProjects.
type GetProjectsQuery struct {
	IncludeArchived *string `url:"include_archived,omitempty"`
}

// GetProjects gets projects.
func (c *Client) GetProjects(ctx context.Context, query *GetProjectsQuery) ([]*Project, error) {
	var projects []*Project
	apiSpecificPath := path.Join(mePath, "projects")
	if err := c.httpGet(ctx, apiSpecificPath, query, &projects); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return projects, nil
}

// GetProjectsPaginatedQuery represents the additional parameters of GetProjectsPaginated.
type GetProjectsPaginatedQuery struct {
	StartProjectID *int `url:"start_project_id,omitempty"`
}

// GetProjectsPaginated gets paginated projects.
func (c *Client) GetProjectsPaginated(ctx context.Context, query *GetProjectsPaginatedQuery) ([]*Project, error) {
	var projects []*Project
	apiSpecificPath := path.Join(mePath, "projects/paginated")
	if err := c.httpGet(ctx, apiSpecificPath, query, &projects); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return projects, nil
}
