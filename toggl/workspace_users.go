package toggl

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

const (
	workspaceUsersEndpoint string = "/api/v8/workspace_users"
)

// WorkspaceUser represents properties of workspace user.
type WorkspaceUser struct {
	ID        int    `json:"id"`
	UID       int    `json:"uid"`
	Admin     bool   `json:"admin"`
	Active    bool   `json:"active"`
	InviteURL string `json:"invite_url"`
}

var (
	// ErrWorkspaceUserNotFound is returned when the provided wprkspace user is nil.
	ErrWorkspaceUserNotFound = errors.New("the provided workspace user must be non-nil")
)

type rawInvitedUsersData struct {
	WorkspaceUsers []*WorkspaceUser `json:"data"`
	Notifications  []string         `json:"notifications"`
}

type rawWorkspaceUserData struct {
	WorkspaceUser WorkspaceUser `json:"data"`
}

func (c *Client) InviteUsersToWorkspace(ctx context.Context, workspace *Workspace, users []*User) ([]*WorkspaceUser, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	emails := []string{}
	for _, user := range users {
		if user == nil {
			return nil, ErrUserNotFound
		}
		emails = append(emails, user.Email)
	}

	invitedUsers := struct {
		Emails []string `json:"emails"`
	}{
		Emails: emails,
	}
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.ID) + "/invite"
	rawInvitedUsersData := new(rawInvitedUsersData)

	if err := c.httpPost(ctx, c.buildURL(endpoint), invitedUsers, rawInvitedUsersData); err != nil {
		return nil, err
	}

	if rawInvitedUsersData.Notifications == nil {
		return rawInvitedUsersData.WorkspaceUsers, nil
	} else {
		return rawInvitedUsersData.WorkspaceUsers, errors.New(strings.Join(rawInvitedUsersData.Notifications, "\n"))
	}
}

func (c *Client) UpdateWorkspaceUser(ctx context.Context, workspaceUser *WorkspaceUser) (*WorkspaceUser, error) {
	if workspaceUser == nil {
		return nil, ErrWorkspaceUserNotFound
	}
	endpoint := workspaceUsersEndpoint + "/" + strconv.Itoa(workspaceUser.ID)
	rawWorkspaceUserData := new(rawWorkspaceUserData)
	if err := c.httpPut(ctx, c.buildURL(endpoint), workspaceUser, rawWorkspaceUserData); err != nil {
		return nil, err
	}
	return &rawWorkspaceUserData.WorkspaceUser, nil
}

func (c *Client) DeleteWorkspaceUser(ctx context.Context, workspaceUser *WorkspaceUser) error {
	if workspaceUser == nil {
		return ErrWorkspaceUserNotFound
	}
	endpoint := workspaceUsersEndpoint + "/" + strconv.Itoa(workspaceUser.ID)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetWorkspaceUsersAsWorkspaceUser(ctx context.Context, workspace *Workspace) ([]*WorkspaceUser, error) {
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	endpoint := workspacesEndpoint + "/" + strconv.Itoa(workspace.ID) + "/workspace_users"
	var workspaceUsers []*WorkspaceUser
	if err := c.httpGet(ctx, c.buildURL(endpoint), &workspaceUsers); err != nil {
		return nil, err
	}
	return workspaceUsers, nil
}
