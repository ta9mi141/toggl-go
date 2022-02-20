package toggl

import "time"

// User represents properties of user.
// Some properties not listed in the documentation are also included.
type User struct {
	ID                    *int        `json:"id,omitempty"`
	APIToken              *string     `json:"api_token,omitempty"`
	DefaultWID            *int        `json:"default_wid,omitempty"`
	Email                 *string     `json:"email,omitempty"`
	Fullname              *string     `json:"fullname,omitempty"`
	JqueryTimeofdayFormat *string     `json:"jquery_timeofday_format,omitempty"`
	JqueryDateFormat      *string     `json:"jquery_date_format,omitempty"`
	TimeofdayFormat       *string     `json:"timeofday_format,omitempty"`
	DateFormat            *string     `json:"date_format,omitempty"`
	StoreStartAndStopTime *bool       `json:"store_start_and_stop_time,omitempty"`
	BeginningOfWeek       *int        `json:"beginning_of_week,omitempty"`
	Language              *string     `json:"language,omitempty"`
	ImageURL              *string     `json:"image_url,omitempty"`
	SidebarPiechart       *bool       `json:"sidebar_piechart,omitempty"`
	At                    *time.Time  `json:"at,omitempty"`
	CreatedAt             *time.Time  `json:"created_at,omitempty"`
	Retention             *int        `json:"retention,omitempty"`
	RecordTimeline        *bool       `json:"record_timeline,omitempty"`
	RenderTimeline        *bool       `json:"render_timeline,omitempty"`
	TimelineEnabled       *bool       `json:"timeline_enabled,omitempty"`
	TimelineExperiment    *bool       `json:"timeline_experiment,omitempty"`
	ShouldUpgrade         *bool       `json:"should_upgrade,omitempty"`
	Timezone              *string     `json:"timezone,omitempty"`
	OpenIDEnabled         *bool       `json:"openid_enabled,omitempty"`
	SendProductEmails     *bool       `json:"send_product_emails,omitempty"`
	SendWeeklyReport      *bool       `json:"send_weekly_report,omitempty"`
	SendTimeNotifications *bool       `json:"send_timer_notifications,omitempty"`
	Invitation            *Invitation `json:"invitation,omitempty"`
	DurationFormat        *string     `json:"duration_format,omitempty"`
}

// Invitation represents properties of invitation.
// Some properties not listed in the documentation are also included.
type Invitation struct {
	ID            *int    `json:"id,omitempty"`
	Code          *string `json:"code,omitempty"`
	SenderName    *string `json:"sender_name,omitempty"`
	SenderEmail   *string `json:"sender_email,omitempty"`
	WorkspaceName *string `json:"workspace_name,omitempty"`
}
