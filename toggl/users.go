package toggl

const (
	usersEndpoint         string = "/api/v8/me"
	resetApiTokenEndpoint string = "/api/v8/reset_token"
	signupEndpoint        string = "/api/v8/signups"
)

// User represents properties of user.
type User struct {
	ApiToken              string `json:"api_token"`
	DefaultWid            int    `json:"default_wid"`
	Email                 string `json:"email"`
	Fullname              string `json:"fullname"`
	JQueryTimeofdayFormat string `json:"jquery_timeofday_format"`
	JQueryDateFormat      string `json:"jquery_date_format"`
	TimeofdayFormat       string `json:"timeofday_format"`
	DateFormat            string `json:"date_format"`
	StoreStartAndStopTime bool   `json:"store_start_and_stop_time"`
	BeginningOfWeek       int    `json:"beginning_of_week"`
	Language              string `json:"language"`
	ImageUrl              string `json:"image_url"`
	SidebarPiechart       bool   `json:"sidebar_piechart"`
	At                    string `json:"at"`
	NewBlogPost           struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"new_blog_post"`
	SendProductEmails      bool   `json:"send_product_emails"`
	SendWeeklyReport       bool   `json:"send_weekly_report"`
	SendTimerNotifications bool   `json:"send_timer_notifications"`
	OpenidEnabled          bool   `json:"openid_enabled"`
	Timezone               string `json:"timezone"`
}
