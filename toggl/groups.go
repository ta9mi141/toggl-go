package toggl

const (
	groupsEndpoint string = "/api/v8/groups"
)

// Group represents properties of group.
type Group struct {
	Name string `json:"name"`
	Wid  int    `json:"wid"`
	At   string `json:"at"`
}
