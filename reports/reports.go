package reports

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	basicAuthPassword string = "api_token" // Defined in Toggl Reports API
	defaultBaseURL    string = "https://toggl.com"
	detailedEndpoint  string = "/reports/api/v2/details"
	summaryEndpoint   string = "/reports/api/v2/summary"
	weeklyEndpoint    string = "/reports/api/v2/weekly"
)

type client struct {
	client   *http.Client
	header   http.Header
	apiToken string
	url      *url.URL
}

type StandardRequestParameters struct {
	UserAgent   string
	WorkSpaceId string
	Since       time.Time
	Until       time.Time
}

func (params *StandardRequestParameters) values() url.Values {
	values := url.Values{}

	// user_agent and workspace_id are required.
	values.Add("user_agent", params.UserAgent)
	values.Add("workspace_id", params.WorkSpaceId)
	// since and until must be ISO 8601 date (YYYY-MM-DD) format
	if !params.Since.IsZero() {
		values.Add("since", params.Since.Format("2006-01-02"))
	}
	if !params.Until.IsZero() {
		values.Add("until", params.Until.Format("2006-01-02"))
	}

	return values
}

type urlEncoder interface {
	urlEncode() string
}

type DetailedRequestParameters struct {
	*StandardRequestParameters
	Page int
}

func (params *DetailedRequestParameters) urlEncode() string {
	values := params.StandardRequestParameters.values()

	if params.Page != 0 {
		values.Add("page", fmt.Sprint(params.Page))
	}

	return values.Encode()
}

type SummaryRequestParameters struct {
	*StandardRequestParameters
	Grouping            string
	Subgrouping         string
	SubgroupingIds      bool
	GroupedTimeEntryIds bool
}

func (params *SummaryRequestParameters) urlEncode() string {
	values := params.StandardRequestParameters.values()

	if params.Grouping != "" {
		values.Add("grouping", params.Grouping)
	}
	if params.Subgrouping != "" {
		values.Add("subgrouping", params.Subgrouping)
	}
	if params.GroupedTimeEntryIds == true {
		values.Add("grouped_time_entry_ids", "true")
	}
	if params.SubgroupingIds == true {
		values.Add("subgrouping_ids", "true")
	}

	return values.Encode()
}

type WeeklyRequestParameters struct {
	*StandardRequestParameters
	Grouping  string
	Calculate string
}

func (params *WeeklyRequestParameters) urlEncode() string {
	values := params.StandardRequestParameters.values()

	if params.Grouping != "" {
		values.Add("grouping", params.Grouping)
	}
	if params.Calculate != "" {
		values.Add("calculate", params.Calculate)
	}

	return values.Encode()
}

type ReportsError struct {
	Err struct {
		Message    string `json:"message"`
		Tip        string `json:"tip"`
		StatusCode int    `json:"code"`
	} `json:"error"`
}

func (e ReportsError) Error() string {
	return fmt.Sprintf(
		"HTTP Status: %d\n%s\n\n%s\n",
		e.Err.StatusCode,
		e.Err.Message,
		e.Err.Tip,
	)
}

type Option func(c *client)

// Configurable baseURL makes client testable
func BaseURL(rawurl string) Option {
	return func(c *client) {
		url, _ := url.Parse(rawurl)
		c.url = url
	}
}

func NewClient(apiToken string, options ...Option) *client {
	url, _ := url.Parse(defaultBaseURL)
	newClient := &client{
		client:   &http.Client{},
		header:   make(http.Header),
		apiToken: apiToken,
		url:      url,
	}
	for _, option := range options {
		option(newClient)
	}
	newClient.header.Set("Content-type", "application/json")
	return newClient
}

func (c *client) GetDetailed(params *DetailedRequestParameters, detailedReport interface{}) error {
	err := c.get(c.buildURL(detailedEndpoint, params), detailedReport)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetSummary(params *SummaryRequestParameters, summaryReport interface{}) error {
	err := c.get(c.buildURL(summaryEndpoint, params), summaryReport)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetWeekly(params *WeeklyRequestParameters, weeklyReport interface{}) error {
	err := c.get(c.buildURL(weeklyEndpoint, params), weeklyReport)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) buildURL(endpoint string, params urlEncoder) string {
	c.url.Path = endpoint
	return c.url.String() + "?" + params.urlEncode()
}

func (c *client) get(url string, report interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.apiToken, basicAuthPassword)
	resp, err := checkResponse(c.client.Do(req))
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
	if resp.StatusCode <= 199 || 300 <= resp.StatusCode {
		var reportsError = new(ReportsError)
		if err := decodeJSON(resp, reportsError); err != nil {
			return nil, err
		}
		return nil, reportsError
	}
	return resp, nil
}

func decodeJSON(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
