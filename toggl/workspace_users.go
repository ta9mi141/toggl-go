package toggl

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
