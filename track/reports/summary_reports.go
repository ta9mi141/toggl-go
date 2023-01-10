package reports

import (
	"context"
	"time"
)

// SummaryReport represents the properties of a summary report.
type SummaryReport struct {
	Groups []*group `json:"groups,omitempty"`
}

type group struct {
	ID        *int        `json:"id,omitempty"`
	SubGroups []*subGroup `json:"sub_groups,omitempty"`
}

type subGroup struct {
	ID      *int    `json:"id,omitempty"`
	Title   *string `json:"title,omitempty"`
	Seconds *int    `json:"seconds,omitempty"`
}

// SearchSummaryReportRequestBody represents a request body of SearchSummaryReport.
type SearchSummaryReportRequestBody struct {
	Audit               *audit     `json:"audit,omitempty"`
	Billable            *bool      `json:"billable,omitempty"`
	ClientIDs           []*int     `json:"client_ids,omitempty"`
	Description         *string    `json:"description,omitempty"`
	EndDate             *string    `json:"end_date,omitempty"`
	GroupIDs            []*int     `json:"group_ids,omitempty"`
	Grouping            *string    `json:"grouping,omitempty"`
	IncludeTimeEntryIDs *bool      `json:"include_time_entry_ids,omitempty"`
	MaxDurationSeconds  *int       `json:"max_duration_seconds,omitempty"`
	MinDurationSeconds  *int       `json:"min_duration_seconds,omitempty"`
	PostedFields        []*string  `json:"postedFields,omitempty"`
	ProjectIDs          []*int     `json:"project_ids,omitempty"`
	Rounding            *int       `json:"rounding,omitempty"`
	RoundingMinutes     *int       `json:"rounding_minutes,omitempty"`
	StartTime           *time.Time `json:"startTime,omitempty"`
	StartDate           *string    `json:"start_date,omitempty"`
	SubGrouping         *string    `json:"sub_grouping,omitempty"`
	TagIDs              []*int     `json:"tag_ids,omitempty"`
	TaskIDs             []*int     `json:"task_ids,omitempty"`
	UserIDs             []*int     `json:"user_ids,omitempty"`
}

type audit struct {
	GroupFilter       *groupFilter `json:"group_filter,omitempty"`
	ShowEmptyGroups   *bool        `json:"show_empty_groups,omitempty"`
	ShowTrackedGroups *bool        `json:"show_tracked_groups,omitempty"`
}

type groupFilter struct {
	Currency           *string `json:"currency,omitempty"`
	MaxAmountCents     *int    `json:"max_amount_cents,omitempty"`
	MaxDurationSeconds *int    `json:"max_duration_seconds,omitempty"`
	MinAmountCents     *int    `json:"min_amount_cents,omitempty"`
	MinDurationSeconds *int    `json:"min_duration_seconds,omitempty"`
}

// SearchSummaryReport returns time entries for summary report.
func (c *APIClient) SearchSummaryReport(ctx context.Context, workspaceID int, reqBody *SearchSummaryReportRequestBody) (*SummaryReport, error) {
	return nil, nil
}
