package toggl

const (
	dashboardEndpoint string = "/api/v8/dashboard"
)

// Dashboard represents properties of two objects dashboard request returns.
type Dashboard struct {
	Activity struct {
		UserId      int    `json:"user_id"`
		ProjectId   int    `json:"project_id"`
		Duration    int    `json:"duration"`
		Description string `json:"description"`
		Stop        string `json:"stop"`
		Tid         int    `json:"tid"`
	} `json:"activity"`
	MostActiveUser struct {
		UserId   int `json:"user_id"`
		Duration int `json:"duration"`
	} `json:"most_active_user"`
}
