package toggl

const (
	timeEntriesEndpoint string = "/api/v8/time_entries"
)

// TimeEntry represents properties of time entry.
type TimeEntry struct {
	Description string   `json:"description"`
	Wid         int      `json:"wid"`
	Pid         int      `json:"pid"`
	Tid         int      `json:"tid"`
	Start       string   `json:"start"`
	Stop        string   `json:"stop"`
	Duration    int      `json:"duration"`
	CreatedWith string   `json:"created_with"`
	Tags        []string `json:"tags"`
	Duronly     bool     `json:"duronly"`
	At          string   `json:"at"`
}
