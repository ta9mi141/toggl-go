package toggl

import (
	"context"
	"errors"
	"time"
)

const (
	timeEntriesEndpoint string = "/api/v8/time_entries"
)

var (
	// ErrTimeEntryNotFound is returned when the provided time entry is nil.
	ErrTimeEntryNotFound = errors.New("The provided time entry must be non-nil")
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

// CreateTimeEntry creates a time entry.
func (c *Client) CreateTimeEntry(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	return nil, nil
}

// UpdateTimeEntry updates a time entry.
func (c *Client) UpdateTimeEntry(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	return nil, nil
}

// DeleteTimeEntry deletes a time entry.
func (c *Client) DeleteTimeEntry(ctx context.Context, timeEntry *TimeEntry) error {
	return nil
}

// GetTimeEntry gets time entry details.
func (c *Client) GetTimeEntry(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	return nil, nil
}

// GetRunningTimeEntry gets running time entry.
func (c *Client) GetRunningTimeEntry(ctx context.Context) (*TimeEntry, error) {
	return nil, nil
}

// Start starts a time entry.
func (c *Client) Start(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	return nil, nil
}

// Stop stops a time entry.
func (c *Client) Stop(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	return nil, nil
}
