package toggl

import (
	"context"
	"path"
	"time"

	"github.com/pkg/errors"
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
func (c *APIClient) GetMe(ctx context.Context) (*Me, error) {
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
func (c *APIClient) UpdateMe(ctx context.Context, reqBody *UpdateMeRequestBody) (*Me, error) {
	var me *Me
	if err := c.httpPut(ctx, mePath, reqBody, &me); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return me, nil
}

// GetMyOrganizations gets all organizations a given user is part of.
func (c *APIClient) GetMyOrganizations(ctx context.Context) ([]*Organization, error) {
	var organizations []*Organization
	apiSpecificPath := path.Join(mePath, "organizations")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &organizations); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return organizations, nil
}

// GetMyProjectsQuery represents the additional parameters of GetMyProjects.
type GetMyProjectsQuery struct {
	IncludeArchived *string `url:"include_archived,omitempty"`
}

// GetMyProjects gets projects.
func (c *APIClient) GetMyProjects(ctx context.Context, query *GetMyProjectsQuery) ([]*Project, error) {
	var projects []*Project
	apiSpecificPath := path.Join(mePath, "projects")
	if err := c.httpGet(ctx, apiSpecificPath, query, &projects); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return projects, nil
}

// GetMyProjectsPaginatedQuery represents the additional parameters of GetMyProjectsPaginated.
type GetMyProjectsPaginatedQuery struct {
	StartProjectID *int `url:"start_project_id,omitempty"`
}

// GetMyProjectsPaginated gets paginated projects.
func (c *APIClient) GetMyProjectsPaginated(ctx context.Context, query *GetMyProjectsPaginatedQuery) ([]*Project, error) {
	var projects []*Project
	apiSpecificPath := path.Join(mePath, "projects/paginated")
	if err := c.httpGet(ctx, apiSpecificPath, query, &projects); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return projects, nil
}

// GetMyTags returns tags for the current user.
func (c *APIClient) GetMyTags(ctx context.Context) ([]*Tag, error) {
	var tags []*Tag
	apiSpecificPath := path.Join(mePath, "tags")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &tags); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return tags, nil
}

// GetMyClients gets clients.
func (c *APIClient) GetMyClients(ctx context.Context) ([]*Client, error) {
	var clients []*Client
	apiSpecificPath := path.Join(mePath, "clients")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &clients); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return clients, nil
}
