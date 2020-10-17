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
	Id       int       `json:"id"`
	Pid      int       `json:"pid"`
	Uid      int       `json:"uid"`
	Wid      int       `json:"wid"`
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
	ErrProjectUserNotFound = errors.New("The provided project user must be non-nil")
)

// CreateProjectUsers creates project user.
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

// UpdateProjectUsers updates project user.
func (c *Client) UpdateProjectUser(ctx context.Context, projectUser *ProjectUser) (*ProjectUser, error) {
	if projectUser == nil {
		return nil, ErrProjectUserNotFound
	}
	rawProjectUserData := new(rawProjectUserData)
	endpoint := projectUsersEndpoint + "/" + strconv.Itoa(projectUser.Id)
	if err := c.httpPut(ctx, c.buildURL(endpoint), projectUser, rawProjectUserData); err != nil {
		return nil, err
	}
	return &rawProjectUserData.ProjectUser, nil
}

// DeleteProjectUsers deletes project users.
func (c *Client) DeleteProjectUsers(ctx context.Context, projectUsers []*ProjectUser) error {
	if len(projectUsers) == 0 {
		return ErrProjectUserNotFound
	}

	var projectUserIds []int
	for _, projectUser := range projectUsers {
		if projectUser == nil {
			return ErrProjectUserNotFound
		}
		projectUserIds = append(projectUserIds, projectUser.Id)
	}

	endpoint := projectUsersEndpoint + "/" + arrayToString(projectUserIds, ",")
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
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.Id) + "/project_users"
	if err := c.httpGet(ctx, c.buildURL(endpoint), &projectUsers); err != nil {
		return nil, err
	}
	return projectUsers, nil
}
