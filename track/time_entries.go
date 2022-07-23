package track

import (
	"context"
	"time"
)

// TimeEntry represents the properties of a time entry.
type TimeEntry struct {
	ID              *int       `json:"id,omitempty"`
	WorkspaceID     *int       `json:"workspace_id,omitempty"`
	ProjectID       *int       `json:"project_id,omitempty"`
	TaskID          *int       `json:"task_id,omitempty"`
	Billable        *bool      `json:"billable,omitempty"`
	Start           *time.Time `json:"start,omitempty"`
	Stop            *time.Time `json:"stop,omitempty"`
	Duration        *int       `json:"duration,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Tags            []*string  `json:"tags,omitempty"`
	TagIDs          []*int     `json:"tag_ids,omitempty"`
	Duronly         *bool      `json:"duronly,omitempty"`
	At              *time.Time `json:"at,omitempty"`
	ServerDeletedAt *string    `json:"server_deleted_at,omitempty"`
	UserID          *int       `json:"user_id,omitempty"`
	UID             *int       `json:"uid,omitempty"`
	WID             *int       `json:"wid,omitempty"`
	PID             *int       `json:"pid,omitempty"`
	TID             *int       `json:"tid,omitempty"`
}

// GetTimeEntriesQueries represents the additional parameters of GetTimeEntries.
type GetTimeEntriesQueries struct {
	Before    *string `url:"before,omitempty"`
	Since     *int    `url:"since,omitempty"`
	StartDate *string `url:"start_date,omitempty"`
	EndDate   *string `url:"end_date,omitempty"`
}

// GetTimeEntries returns latest time entries.
func (c *Client) GetTimeEntries(ctx context.Context, queries *GetTimeEntriesQueries) ([]*TimeEntry, error) {
	return nil, nil
}
