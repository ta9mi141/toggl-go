package toggl

import (
	"context"
	"errors"
	"time"
)

const (
	usersEndpoint         string = "/api/v8/me"
	resetAPITokenEndpoint string = "/api/v8/reset_token"
	signUpEndpoint        string = "/api/v8/signups"
)

var (
	// ErrUserNotFound is returned when the provided user is nil.
	ErrUserNotFound = errors.New("The provided user must be non-nil")
)

// User represents properties of user including ones only used to update or sign up.
type User struct {
	Id                    int       `json:"id"`
	APIToken              string    `json:"api_token"`
	DefaultWid            int       `json:"default_wid"`
	Email                 string    `json:"email"`
	Fullname              string    `json:"fullname"`
	JQueryTimeofdayFormat string    `json:"jquery_timeofday_format"`
	JQueryDateFormat      string    `json:"jquery_date_format"`
	TimeofdayFormat       string    `json:"timeofday_format"`
	DateFormat            string    `json:"date_format"`
	StoreStartAndStopTime bool      `json:"store_start_and_stop_time"`
	BeginningOfWeek       int       `json:"beginning_of_week"`
	Language              string    `json:"language"`
	ImageUrl              string    `json:"image_url"`
	SidebarPiechart       bool      `json:"sidebar_piechart"`
	At                    time.Time `json:"at"`
	NewBlogPost           struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"new_blog_post"`
	SendProductEmails      bool           `json:"send_product_emails"`
	SendWeeklyReport       bool           `json:"send_weekly_report"`
	SendTimerNotifications bool           `json:"send_timer_notifications"`
	OpenidEnabled          bool           `json:"openid_enabled"`
	Timezone               string         `json:"timezone"`
	TimeEntries            []*TimeEntry   `json:"time_entries"`
	Projects               []*Project     `json:"projects"`
	Tags                   []*Tag         `json:"tags"`
	Workspaces             []*Workspace   `json:"workspaces"`
	Clients                []*TogglClient `json:"clients"`
	Password               string         `json:"password"`
	CurrentPassword        string         `json:"current_password"`
	CreatedWith            string         `json:"created_with"`
}

type rawUserData struct {
	User User `json:"data"`
}

// GetUser gets current user data.
func (c *Client) GetUser(ctx context.Context, params ...QueryString) (*User, error) {
	rawUserData := new(rawUserData)
	if err := c.httpGet(ctx, c.buildURL(usersEndpoint, params...), rawUserData); err != nil {
		return nil, err
	}
	return &rawUserData.User, nil
}

// UpdateUser updates user data.
func (c *Client) UpdateUser(ctx context.Context, user *User) (*User, error) {
	if user == nil {
		return nil, ErrUserNotFound
	}
	rawUserData := new(rawUserData)
	if err := c.httpPut(ctx, c.buildURL(usersEndpoint), user, rawUserData); err != nil {
		return nil, err
	}
	return &rawUserData.User, nil
}

// ResetAPIToken resets API token and returns the new API token.
func (c *Client) ResetAPIToken(ctx context.Context) (string, error) {
	var newAPIToken string
	if err := c.httpPost(ctx, c.buildURL(resetAPITokenEndpoint), nil, &newAPIToken); err != nil {
		return "", err
	}
	return newAPIToken, nil
}

// SignUp creates new user.
func (c *Client) SignUp(ctx context.Context, user *User) (*User, error) {
	if user == nil {
		return nil, ErrUserNotFound
	}
	rawUserData := new(rawUserData)
	if err := c.httpPost(ctx, c.buildURL(signUpEndpoint), user, rawUserData); err != nil {
		return nil, err
	}
	return &rawUserData.User, nil
}
