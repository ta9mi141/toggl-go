package toggl

import (
	"context"
	"errors"
	"strconv"
	"time"
)

const (
	projectUsersEndpoint string = "/api/v8/project_users"
)

// ProjectUser represents properties of project user including additional fields.
type ProjectUser struct {
	ID       int       `json:"id"`
	PID      int       `json:"pid"`
	UID      int       `json:"uid"`
	WID      int       `json:"wid"`
	Manager  bool      `json:"manager"`
	At       time.Time `json:"at"`
	Fullname string    `json:"fullname"`
	Fields   string    `json:"fields"`
}

type rawProjectUserData struct {
	ProjectUser ProjectUser `json:"data"`
}

var (
	// ErrProjectUserNotFound is returned when the provided project user is nil.
	ErrProjectUserNotFound = errors.New("the provided project user must be non-nil")
)

// CreateProjectUser creates a project user.
func (c *Client) CreateProjectUser(ctx context.Context, projectUser *ProjectUser) (*ProjectUser, error) {
	if projectUser == nil {
		return nil, ErrProjectUserNotFound
	}
	rawProjectUserData := new(rawProjectUserData)
	if err := c.httpPost(ctx, c.buildURL(projectUsersEndpoint), projectUser, rawProjectUserData); err != nil {
		return nil, err
	}
	return &rawProjectUserData.ProjectUser, nil
}

// UpdateProjectUser updates a project user.
func (c *Client) UpdateProjectUser(ctx context.Context, projectUser *ProjectUser) (*ProjectUser, error) {
	if projectUser == nil {
		return nil, ErrProjectUserNotFound
	}
	rawProjectUserData := new(rawProjectUserData)
	endpoint := projectUsersEndpoint + "/" + strconv.Itoa(projectUser.ID)
	if err := c.httpPut(ctx, c.buildURL(endpoint), projectUser, rawProjectUserData); err != nil {
		return nil, err
	}
	return &rawProjectUserData.ProjectUser, nil
}

// DeleteProjectUser deletes a project user.
func (c *Client) DeleteProjectUser(ctx context.Context, projectUser *ProjectUser) error {
	if projectUser == nil {
		return ErrProjectUserNotFound
	}
	endpoint := projectUsersEndpoint + "/" + strconv.Itoa(projectUser.ID)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}

// GetProjectUsersInWorkspace get project users in a workspace.
func (c *Client) GetProjectUsersInWorkspace(ctx context.Context, workspace *Workspace) ([]*ProjectUser, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	var projectUsers []*ProjectUser
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.ID) + "/project_users"
	if err := c.httpGet(ctx, c.buildURL(endpoint), &projectUsers); err != nil {
		return nil, err
	}
	return projectUsers, nil
}
