package toggl

import "time"

const (
	timeEntriesEndpoint string = "/api/v8/time_entries"
)

// TimeEntry represents properties of time entry.
type TimeEntry struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Wid         int       `json:"wid"`
	Pid         int       `json:"pid"`
	Tid         int       `json:"tid"`
	Start       time.Time `json:"start"`
	Stop        time.Time `json:"stop"`
	Duration    int       `json:"duration"`
	CreatedWith string    `json:"created_with"`
	Tags        []string  `json:"tags"`
	Duronly     bool      `json:"duronly"`
	At          time.Time `json:"at"`
}
