package toggl

const (
	projectUsersEndpoint string = "/api/v8/project_users"
)

// ProjectUser represents properties of project user.
type ProjectUser struct {
	Pid     int    `json:"pid"`
	Uid     int    `json:"uid"`
	Wid     int    `json:"wid"`
	Manager bool   `json:"manager"`
	At      string `json:"at"`
}
