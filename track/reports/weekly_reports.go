package reports

import (
	"context"
	"time"
)

// WeeklyReport represents the properties of a weekly report.
type WeeklyReport []struct {
	UserID    *int   `json:"user_id,omitempty"`
	ProjectID *int   `json:"project_id,omitempty"`
	Seconds   []*int `json:"seconds,omitempty"`
}

// SearchWeeklyReportRequestBody represents a request body of SearchWeeklyReport.
type SearchWeeklyReportRequestBody struct {
	Billable           *bool      `json:"billable,omitempty"`
	ClientIDs          []*int     `json:"client_ids,omitempty"`
	Description        *string    `json:"description,omitempty"`
	EndDate            *string    `json:"end_date,omitempty"`
	GroupIDs           []*int     `json:"group_ids,omitempty"`
	MaxDurationSeconds *int       `json:"max_duration_seconds,omitempty"`
	MinDurationSeconds *int       `json:"min_duration_seconds,omitempty"`
	PostedFields       []*string  `json:"postedFields,omitempty"`
	ProjectIDs         []*int     `json:"project_ids,omitempty"`
	Rounding           *int       `json:"rounding,omitempty"`
	RoundingMinutes    *int       `json:"rounding_minutes,omitempty"`
	StartTime          *time.Time `json:"startTime,omitempty"`
	StartDate          *string    `json:"start_date,omitempty"`
	TagIDs             []*int     `json:"tag_ids,omitempty"`
	TaskIDs            []*int     `json:"task_ids,omitempty"`
	UserIDs            []*int     `json:"user_ids,omitempty"`
}

// SearchWeeklyReport returns time entries for weekly report.
func (c *APIClient) SearchWeeklyReport(ctx context.Context, workspaceID int, reqBody *SearchWeeklyReportRequestBody) (*WeeklyReport, error) {
	return nil, nil
}
