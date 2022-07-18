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

// Me represents properties of an user.
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
	me := new(Me)
	if err := c.httpGet(ctx, mePath, me); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return me, nil
}

// PutMeRequestBody represents a request body of PutMe.
type PutMeRequestBody struct {
	BeginningOfWeek    *int    `json:"beginning_of_week,omitempty"`
	CountryID          *int    `json:"country_id,omitempty"`
	CurrentPassword    *string `json:"current_password,omitempty"`
	DefaultWorkspaceID *int    `json:"default_workspace_id,omitempty"`
	Email              *string `json:"email,omitempty"`
	Fullname           *string `json:"fullname,omitempty"`
	Password           *string `json:"password,omitempty"`
	Timezone           *string `json:"timezone,omitempty"`
}

// PutMe updates details for the current user.
func (c *Client) PutMe(ctx context.Context, reqBody *PutMeRequestBody) (*Me, error) {
	me := new(Me)
	if err := c.httpPut(ctx, mePath, reqBody, me); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return me, nil
}

// GetMyOrganizations gets all organizations a given user is part of.
func (c *Client) GetMyOrganizations(ctx context.Context) ([]*Organization, error) {
	var organizations []*Organization
	apiSpecificPath := path.Join(mePath, "organizations")
	if err := c.httpGet(ctx, apiSpecificPath, &organizations); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return organizations, nil
}
