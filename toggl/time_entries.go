package toggl

import (
	"context"
	"errors"
	"strconv"
	"time"
)

const (
	timeEntriesEndpoint string = "/api/v8/time_entries"
)

var (
	// ErrTimeEntryNotFound is returned when the provided time entry is nil.
	ErrTimeEntryNotFound = errors.New("the provided time entry must be non-nil")
)

// TimeEntry represents properties of time entry.
type TimeEntry struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	WID         int       `json:"wid"`
	PID         int       `json:"pid"`
	TID         int       `json:"tid"`
	Start       time.Time `json:"start"`
	Stop        time.Time `json:"stop"`
	Duration    int       `json:"duration"`
	CreatedWith string    `json:"created_with"`
	Tags        []string  `json:"tags"`
	Duronly     bool      `json:"duronly"`
	At          time.Time `json:"at"`
}

type rawTimeEntryData struct {
	TimeEntry TimeEntry `json:"data"`
}

// CreateTimeEntry creates a time entry.
func (c *Client) CreateTimeEntry(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	if timeEntry == nil {
		return nil, ErrTimeEntryNotFound
	}
	rawTimeEntryData := new(rawTimeEntryData)
	if err := c.httpPost(ctx, c.buildURL(timeEntriesEndpoint), timeEntry, rawTimeEntryData); err != nil {
		return nil, err
	}
	return &rawTimeEntryData.TimeEntry, nil
}

// UpdateTimeEntry updates a time entry.
func (c *Client) UpdateTimeEntry(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	if timeEntry == nil {
		return nil, ErrTimeEntryNotFound
	}
	rawTimeEntryData := new(rawTimeEntryData)
	endpoint := timeEntriesEndpoint + "/" + strconv.Itoa(timeEntry.ID)
	if err := c.httpPut(ctx, c.buildURL(endpoint), timeEntry, rawTimeEntryData); err != nil {
		return nil, err
	}
	return &rawTimeEntryData.TimeEntry, nil
}

// DeleteTimeEntry deletes a time entry.
func (c *Client) DeleteTimeEntry(ctx context.Context, timeEntry *TimeEntry) error {
	if timeEntry == nil {
		return ErrTimeEntryNotFound
	}
	endpoint := timeEntriesEndpoint + "/" + strconv.Itoa(timeEntry.ID)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}

// GetTimeEntry gets time entry details.
func (c *Client) GetTimeEntry(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	if timeEntry == nil {
		return nil, ErrTimeEntryNotFound
	}
	rawTimeEntryData := new(rawTimeEntryData)
	endpoint := timeEntriesEndpoint + "/" + strconv.Itoa(timeEntry.ID)
	if err := c.httpGet(ctx, c.buildURL(endpoint), rawTimeEntryData); err != nil {
		return nil, err
	}
	return &rawTimeEntryData.TimeEntry, nil
}

// GetRunningTimeEntry gets running time entry.
func (c *Client) GetRunningTimeEntry(ctx context.Context) (*TimeEntry, error) {
	rawTimeEntryData := new(rawTimeEntryData)
	endpoint := timeEntriesEndpoint + "/" + "current"
	if err := c.httpGet(ctx, c.buildURL(endpoint), rawTimeEntryData); err != nil {
		return nil, err
	}
	return &rawTimeEntryData.TimeEntry, nil
}

// Start starts a time entry.
func (c *Client) Start(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	if timeEntry == nil {
		return nil, ErrTimeEntryNotFound
	}
	rawTimeEntryData := new(rawTimeEntryData)
	endpoint := timeEntriesEndpoint + "/" + "start"
	if err := c.httpPost(ctx, c.buildURL(endpoint), timeEntry, rawTimeEntryData); err != nil {
		return nil, err
	}
	return &rawTimeEntryData.TimeEntry, nil
}

// Stop stops a time entry.
func (c *Client) Stop(ctx context.Context, timeEntry *TimeEntry) (*TimeEntry, error) {
	if timeEntry == nil {
		return nil, ErrTimeEntryNotFound
	}
	rawTimeEntryData := new(rawTimeEntryData)
	endpoint := timeEntriesEndpoint + "/" + strconv.Itoa(timeEntry.ID) + "/" + "stop"
	if err := c.httpPut(ctx, c.buildURL(endpoint), timeEntry, rawTimeEntryData); err != nil {
		return nil, err
	}
	return &rawTimeEntryData.TimeEntry, nil
}
