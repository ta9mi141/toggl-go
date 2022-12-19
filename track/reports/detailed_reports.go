package reports

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// DetailedReport represents the properties of a detailed report.
type DetailedReport []struct {
	UserID                *int         `json:"user_id,omitempty"`
	Username              *string      `json:"username,omitempty"`
	ProjectID             *int         `json:"project_id,omitempty"`
	TaskID                *int         `json:"task_id,omitempty"`
	Billable              *bool        `json:"billable,omitempty"`
	Description           *string      `json:"description,omitempty"`
	TagIDs                []*int       `json:"tag_ids,omitempty"`
	BillableAmountInCents *int         `json:"billable_amount_in_cents,omitempty"`
	HourlyRateInCents     *int         `json:"hourly_rate_in_cents,omitempty"`
	Currency              *string      `json:"currency,omitempty"`
	TimeEntries           []*timeEntry `json:"time_entries,omitempty"`
	RowNumber             *int         `json:"row_number,omitempty"`
}

type timeEntry struct {
	ID      *int       `json:"id,omitempty"`
	Seconds *int       `json:"seconds,omitempty"`
	Start   *time.Time `json:"start,omitempty"`
	Stop    *time.Time `json:"stop,omitempty"`
	At      *time.Time `json:"at,omitempty"`
}

// SearchDetailedReportRequestBody represents a request body of SearchDetailedReport.
type SearchDetailedReportRequestBody struct {
	Billable           *bool      `json:"billable,omitempty"`
	ClientIDs          []*int     `json:"client_ids,omitempty"`
	Description        *string    `json:"description,omitempty"`
	EndDate            *string    `json:"end_date,omitempty"`
	FirstID            *int       `json:"first_id,omitempty"`
	FirstRowNumber     *int       `json:"first_row_number,omitempty"`
	FirstTimestamp     *int       `json:"first_timestamp,omitempty"`
	GroupIDs           []*int     `json:"group_ids,omitempty"`
	Grouped            *bool      `json:"grouped,omitempty"`
	HideAmounts        *bool      `json:"hide_amounts,omitempty"`
	MaxDurationSeconds *int       `json:"max_duration_seconds,omitempty"`
	MinDurationSeconds *int       `json:"min_duration_seconds,omitempty"`
	OrderBy            *string    `json:"order_by,omitempty"`
	OrderDir           *string    `json:"order_dir,omitempty"`
	PostedFields       []*string  `json:"postedFields,omitempty"`
	ProjectIDs         []*int     `json:"project_ids,omitempty"`
	Rounding           *int       `json:"rounding,omitempty"`
	RoundingMinutes    *int       `json:"rounding_minutes,omitempty"`
	StartTime          *time.Time `json:"startTime,omitempty"`
	StartDate          *string    `json:"start_date,omitempty"`
	TagIDs             []*int     `json:"tag_ids,omitempty"`
	TaskIDs            []*int     `json:"task_ids,omitempty"`
	TimeEntryIDs       []*int     `json:"time_entry_ids,omitempty"`
	UserIDs            []*int     `json:"user_ids,omitempty"`
}

// SearchDetailedReport returns time entries for detailed report.
func (c *APIClient) SearchDetailedReport(ctx context.Context, workspaceID int, reqBody *SearchDetailedReportRequestBody) (*DetailedReport, error) {
	var detailedReport *DetailedReport
	apiSpecificPath := path.Join(reportsPath, strconv.Itoa(workspaceID), "search/time_entries")
	if err := c.httpPost(ctx, apiSpecificPath, reqBody, &detailedReport); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return detailedReport, nil
}
