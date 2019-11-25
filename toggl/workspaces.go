package toggl

const (
	workspacesEndpoint string = "/api/v8/workspaces"
)

// Workspace represents properties of workspace.
type Workspace struct {
	Name                        string  `json:"name"`
	Premium                     bool    `json:"premium"`
	Admin                       bool    `json:"admin"`
	DefaultHourlyRate           float64 `json:"default_hourly_rate"`
	DefaultCurrency             string  `json:"default_currency"`
	OnlyAdminsMayCreateProjects bool    `json:"only_admins_may_create_projects"`
	OnlyAdminsSeeBillableRates  bool    `json:"only_admins_see_billable_rates"`
	Rounding                    int     `json:"rounding"`
	RoundingMinutes             int     `json:"rounding_minutes"`
	At                          string  `json:"at"`
	LogoURL                     string  `json:"logo_url"`
}
