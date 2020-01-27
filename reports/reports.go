/*
Package reports is a library of Toggl Reports API v2 for Go programming language.

This package deals with 3 types of reports, detailed report, summary report, and weekly report.
Though each report has their own data structure of successful response, they're not defined in this package.
Users must define a structure corresponding responses of each report to use before sending request.

See API documentation for more details.
https://github.com/toggl/toggl_api_docs/blob/master/reports.md
*/
package reports

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

const (
	basicAuthPassword string = "api_token" // Defined in Toggl Reports API v2
	defaultBaseURL    string = "https://toggl.com"
)

// Client implements a basic request handling used by all of the reports.
type Client struct {
	HTTPClient *http.Client
	APIToken   string
	URL        *url.URL
	header     http.Header
}

// StandardRequestParameters represents request parameters used in all of the reports.
type StandardRequestParameters struct {
	UserAgent           string
	WorkSpaceId         string
	Since               time.Time
	Until               time.Time
	Billable            string
	ClientIds           string
	ProjectIds          string
	UserIds             string
	MembersOfGroupIds   string
	OrMembersOfGroupIds string
	TagIds              string
	TaskIds             string
	TimeEntryIds        string
	Description         string
	WithoutDescription  bool
	OrderField          string
	OrderDesc           bool
	DistinctRates       bool
	Rounding            bool
	DisplayHours        string
}

func (params *StandardRequestParameters) values() url.Values {
	values := url.Values{
		// user_agent and workspace_id are required.
		"user_agent":   []string{params.UserAgent},
		"workspace_id": []string{params.WorkSpaceId},
	}

	// since and until must be ISO 8601 date (YYYY-MM-DD) format
	if !params.Since.IsZero() {
		values.Add("since", params.Since.Format("2006-01-02"))
	}
	if !params.Until.IsZero() {
		values.Add("until", params.Until.Format("2006-01-02"))
	}
	if params.Billable != "" {
		values.Add("billable", params.Billable)
	}
	if params.ClientIds != "" {
		values.Add("client_ids", params.ClientIds)
	}
	if params.ProjectIds != "" {
		values.Add("project_ids", params.ProjectIds)
	}
	if params.UserIds != "" {
		values.Add("user_ids", params.UserIds)
	}
	if params.MembersOfGroupIds != "" {
		values.Add("members_of_group_ids", params.MembersOfGroupIds)
	}
	if params.OrMembersOfGroupIds != "" {
		values.Add("or_members_of_group_ids", params.OrMembersOfGroupIds)
	}
	if params.TagIds != "" {
		values.Add("tag_ids", params.TagIds)
	}
	if params.TaskIds != "" {
		values.Add("task_ids", params.TaskIds)
	}
	if params.TimeEntryIds != "" {
		values.Add("time_entry_ids", params.TimeEntryIds)
	}
	if params.Description != "" {
		values.Add("description", params.Description)
	}
	if params.WithoutDescription == true {
		values.Add("without_description", "true")
	}
	if params.OrderField != "" {
		values.Add("order_field", params.OrderField)
	}
	if params.OrderDesc == true {
		values.Add("order_desc", "on")
	}
	if params.DistinctRates == true {
		values.Add("distinct_rates", "on")
	}
	if params.Rounding == true {
		values.Add("rounding", "on")
	}
	if params.DisplayHours != "" {
		values.Add("display_hours", params.DisplayHours)
	}

	return values
}

type urlEncoder interface {
	urlEncode() string
}

// Error wraps error interface with status code and tip.
// Use errors.As method or type assertions with Error's StatusCode
// and Tip methods to get detailed information about the error.
type Error interface {
	error
	// StatusCode returns HTTP status code of the error
	StatusCode() int
	// Tip shows what to do in case of the error
	Tip() string
}

// ReportsError represents a response of unsuccessful request.
type ReportsError struct {
	Err struct {
		Message string `json:"message"`
		Tip     string `json:"tip"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func (e ReportsError) Error() string {
	return e.Err.Message
}

func (e ReportsError) StatusCode() int {
	return e.Err.Code
}

func (e ReportsError) Tip() string {
	return e.Err.Tip
}

const (
	tooManyRequestErrorMessage string = "Too Many Requests"
	tooManyRequestErrorTip     string = "Add delay between requests"
)

var (
	// ErrContextNotFound is returned when the provided context is nil.
	ErrContextNotFound = errors.New("The provided ctx must be non-nil")
)

// Option represents optional parameters of NewClient.
type Option func(c *Client)

// HTTPClient sets an HTTP client to use when sending requests.
// By default, http.DefaultClient is used.
func HTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// NewClient returns a pointer to a new initialized client.
func NewClient(apiToken string, options ...Option) *Client {
	url, _ := url.Parse(defaultBaseURL)
	newClient := &Client{
		HTTPClient: http.DefaultClient,
		APIToken:   apiToken,
		URL:        url,
		header:     make(http.Header),
	}
	newClient.header.Set("Content-type", "application/json")
	for _, option := range options {
		option(newClient)
	}
	return newClient
}

func (c *Client) buildURL(endpoint string, params urlEncoder) string {
	c.URL.Path = endpoint
	return c.URL.String() + "?" + params.urlEncode()
}

func (c *Client) get(ctx context.Context, url string, report interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.APIToken, basicAuthPassword)

	if ctx == nil {
		return ErrContextNotFound
	}
	req = req.WithContext(ctx)

	resp, err := checkResponse(c.HTTPClient.Do(req))
	if err != nil {
		return err
	}
	if err = decodeJSON(resp, report); err != nil {
		return err
	}
	return nil
}

func checkResponse(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}
	switch {
	case resp.StatusCode == http.StatusTooManyRequests:
		// Since the response of "429 Too Many Requests" is not in form of JSON,
		// the error must be handled individually before calling the function "decodeJSON".
		tooManyRequestsError := new(ReportsError)
		tooManyRequestsError.Err.Code = http.StatusTooManyRequests
		tooManyRequestsError.Err.Message = tooManyRequestErrorMessage
		tooManyRequestsError.Err.Tip = tooManyRequestErrorTip
		return nil, tooManyRequestsError
	case resp.StatusCode <= 199 || 300 <= resp.StatusCode:
		reportsError := new(ReportsError)
		if err := decodeJSON(resp, reportsError); err != nil {
			return nil, err
		}
		return nil, reportsError
	default:
		return resp, nil
	}
}

func decodeJSON(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
