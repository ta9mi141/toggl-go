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

type ReportsError struct {
	Error struct {
		Message    string `json:"message"`
		Tip        string `json:"tip"`
		StatusCode int    `json:"code"`
	} `json:"error"`
}

// TODO: implement Error interface
func (e ReportsError) Error() string {}

type Option func(c *client)

// Configurable baseURL makes client testable
func BaseURL(rawurl string) option {
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

func (c *client) GetDetailed(detailedReport interface{}) error {
	err := c.get(c.buildURL(detailedEndpoint), report)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetSummary(summaryReport interface{}) error {
	err := c.get(c.buildURL(summaryEndpoint), report)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetWeekly(weeklyReport interface{}) error {
	err := c.get(c.buildURL(weeklyEndpoint), report)
	if err != nil {
		return err
	}
	return nil
}

// TODO: encode parameters
// TODO: define buildURL

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
