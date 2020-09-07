package toggl

import (
	"context"
	"errors"
)

const (
	workspaceUsersEndpoint string = "/api/v8/workspace_users"
)

// WorkspaceUser represents properties of workspace user.
type WorkspaceUser struct {
	Id        int    `json:"id"`
	Uid       int    `json:"uid"`
	Admin     bool   `json:"admin"`
	Active    bool   `json:"active"`
	InviteURL string `json:"invite_url"`
}

var (
	// ErrWorkspaceUserNotFound is returned when the provided wprkspace user is nil.
	ErrWorkspaceUserNotFound = errors.New("The provided workspace user must be non-nil")
)

type rawWorkspaceUserData struct {
	WorkspaceUser WorkspaceUser `json:"data"`
	Notifications []string      `json:"notifications"`
}

func (c *Client) InviteUsersToWorkspace(ctx context.Context, users []*User) ([]*WorkspaceUser, error) {
	return nil, nil // TODO
}

func (c *Client) UpdateWorkspaceUser(ctx context.Context, workspaceUser *WorkspaceUser) (*WorkspaceUser, error) {
	return nil, nil // TODO
}

func (c *Client) DeleteWorkspaceUser(ctx context.Context, workspaceUser *WorkspaceUser) error {
	return nil // TODO
}

func (c *Client) GetWorkspaceUsersAsWorkspaceUser(ctx context.Context, workspace *Workspace) ([]*WorkspaceUser, error) {
	return nil, nil // TODO
}
