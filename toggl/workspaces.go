package toggl

import "errors"

const (
	workspacesEndpoint string = "/api/v8/workspaces"
)

var (
	// ErrWorkspaceNotFound is returned when the provided workspace is nil.
	ErrWorkspaceNotFound = errors.New("The provided workspace must be non-nil")
)

// Workspace represents properties of workspace.
type Workspace struct {
	Id                          int     `json:"id"`
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
