package track

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
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

// GetTimeEntriesQuery represents the additional parameters of GetTimeEntries.
type GetTimeEntriesQuery struct {
	Before    *string `url:"before,omitempty"`
	Since     *int    `url:"since,omitempty"`
	StartDate *string `url:"start_date,omitempty"`
	EndDate   *string `url:"end_date,omitempty"`
}

// GetTimeEntries lists latest time entries.
func (c *Client) GetTimeEntries(ctx context.Context, query *GetTimeEntriesQuery) ([]*TimeEntry, error) {
	var timeEntries []*TimeEntry
	apiSpecificPath := path.Join(mePath, "time_entries")
	if err := c.httpGet(ctx, apiSpecificPath, query, &timeEntries); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return timeEntries, nil
}

// GetCurrentTimeEntry loads running time entry for user id.
func (c *Client) GetCurrentTimeEntry(ctx context.Context) (*TimeEntry, error) {
	var timeEntry *TimeEntry
	apiSpecificPath := path.Join(mePath, "time_entries/current")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &timeEntry); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return timeEntry, nil
}

// CreateTimeEntryRequestBody represents a request body of CreateTimeEntry.
type CreateTimeEntryRequestBody struct {
	Billable     *bool      `json:"billable,omitempty"`
	CreatedWith  *string    `json:"created_with,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Duration     *int       `json:"duration,omitempty"`
	Duronly      *bool      `json:"duronly,omitempty"`
	PID          *int       `json:"pid,omitempty"`
	PostedFields []*string  `json:"postedFields,omitempty"`
	ProjectID    *int       `json:"project_id,omitempty"`
	Start        *time.Time `json:"start,omitempty"`
	StartDate    *string    `json:"start_date,omitempty"`
	Stop         *time.Time `json:"stop,omitempty"`
	TagAction    *string    `json:"tag_action,omitempty"`
	TagIDs       []*int     `json:"tag_ids,omitempty"`
	Tags         []*string  `json:"tags,omitempty"`
	TaskID       *int       `json:"task_id,omitempty"`
	TID          *int       `json:"tid,omitempty"`
	UID          *int       `json:"uid,omitempty"`
	UserID       *int       `json:"user_id,omitempty"`
	WID          *int       `json:"wid,omitempty"`
	WorkspaceID  *int       `json:"workspace_id,omitempty"`
}

// CreateTimeEntry creates a new workspace time entry.
func (c *Client) CreateTimeEntry(ctx context.Context, workspaceID int, reqBody *CreateTimeEntryRequestBody) (*TimeEntry, error) {
	var timeEntry *TimeEntry
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "time_entries")
	if err := c.httpPost(ctx, apiSpecificPath, reqBody, &timeEntry); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return timeEntry, nil
}

// UpdateTimeEntryRequestBody represents a request body of UpdateTimeEntry.
type UpdateTimeEntryRequestBody struct {
	Billable     *bool      `json:"billable,omitempty"`
	CreatedWith  *string    `json:"created_with,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Duration     *int       `json:"duration,omitempty"`
	Duronly      *bool      `json:"duronly,omitempty"`
	PID          *int       `json:"pid,omitempty"`
	PostedFields []*string  `json:"postedFields,omitempty"`
	ProjectID    *int       `json:"project_id,omitempty"`
	Start        *time.Time `json:"start,omitempty"`
	StartDate    *string    `json:"start_date,omitempty"`
	Stop         *time.Time `json:"stop,omitempty"`
	TagAction    *string    `json:"tag_action,omitempty"`
	TagIDs       []*int     `json:"tag_ids,omitempty"`
	Tags         []*string  `json:"tags,omitempty"`
	TaskID       *int       `json:"task_id,omitempty"`
	TID          *int       `json:"tid,omitempty"`
	UID          *int       `json:"uid,omitempty"`
	UserID       *int       `json:"user_id,omitempty"`
	WID          *int       `json:"wid,omitempty"`
	WorkspaceID  *int       `json:"workspace_id,omitempty"`
}

// UpdateTimeEntry updates a workspace time entry.
func (c *Client) UpdateTimeEntry(ctx context.Context, workspaceID, timeEntryID int, reqBody *UpdateTimeEntryRequestBody) (*TimeEntry, error) {
	return nil, nil
}
