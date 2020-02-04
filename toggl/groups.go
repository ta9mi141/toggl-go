package toggl

import (
	"context"
	"errors"
	"strconv"
	"time"
)

const (
	groupsEndpoint string = "/api/v8/groups"
)

var (
	// ErrGroupNotFound is returned when the provided group is nil.
	ErrGroupNotFound = errors.New("The provided group must be non-nil")
)

// Group represents properties of group.
type Group struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Wid  int       `json:"wid"`
	At   time.Time `json:"at"`
}

type rawGroupData struct {
	Group Group `json:"data"`
}

// CreateGroup creates a group.
func (c *Client) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	if group == nil {
		return nil, ErrGroupNotFound
	}
	rawGroupData := new(rawGroupData)
	if err := c.httpPost(ctx, c.buildURL(groupsEndpoint), group, rawGroupData); err != nil {
		return nil, err
	}
	return &rawGroupData.Group, nil
}

// UpdateGroup updates a group.
func (c *Client) UpdateGroup(ctx context.Context, group *Group) (*Group, error) {
	if group == nil {
		return nil, ErrGroupNotFound
	}
	rawGroupData := new(rawGroupData)
	endpoint := groupsEndpoint + "/" + strconv.Itoa(group.Id)
	if err := c.httpPut(ctx, c.buildURL(endpoint), group, rawGroupData); err != nil {
		return nil, err
	}
	return &rawGroupData.Group, nil
}

// DeleteGroup deletes a group.
func (c *Client) DeleteGroup(ctx context.Context, group *Group) error {
	if group == nil {
		return ErrGroupNotFound
	}
	endpoint := groupsEndpoint + "/" + strconv.Itoa(group.Id)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}
