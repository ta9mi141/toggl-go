package toggl

import "time"

// User represents properties of user.
// Some properties not listed in the documentation are also included.
type User struct {
	ID                    int       `json:"id"`
	APIToken              string    `json:"api_token"`
	DefaultWID            int       `json:"default_wid"`
	Email                 string    `json:"email"`
	Fullname              string    `json:"fullname"`
	JqueryTimeofdayFormat string    `json:"jquery_timeofday_format"`
	JqueryDateFormat      string    `json:"jquery_date_format"`
	TimeofdayFormat       string    `json:"timeofday_format"`
	DateFormat            string    `json:"date_format"`
	StoreStartAndStopTime bool      `json:"store_start_and_stop_time"`
	BeginningOfWeek       int       `json:"beginning_of_week"`
	Language              string    `json:"language"`
	ImageURL              string    `json:"image_url"`
	SidebarPiechart       bool      `json:"sidebar_piechart"`
	At                    time.Time `json:"at"`
	CreatedAt             time.Time `json:"created_at"`
	Retention             int       `json:"retention"`
	RecordTimeline        bool      `json:"record_timeline"`
	RenderTimeline        bool      `json:"render_timeline"`
	TimelineEnabled       bool      `json:"timeline_enabled"`
	TimelineExperiment    bool      `json:"timeline_experiment"`
	ShouldUpgrade         bool      `json:"should_upgrade"`
	Timezone              string    `json:"timezone"`
	OpenIDEnabled         bool      `json:"openid_enabled"`
	SendProductEmails     bool      `json:"send_product_emails"`
	SendWeeklyReport      bool      `json:"send_weekly_report"`
	SendTimeNotifications bool      `json:"send_timer_notifications"`
	Invitation            struct{}  `json:"invitation"`
	DurationFormat        string    `json:"duration_format"`
}
