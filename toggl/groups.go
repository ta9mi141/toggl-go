package toggl

import (
	"context"
	"errors"
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

// CreateGroup creates a group.
func (c *Client) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	return nil, nil
}

// UpdateGroup updates a group.
func (c *Client) UpdateGroup(ctx context.Context, group *Group) (*Group, error) {
	return nil, nil
}

// DeleteGroup deletes a group.
func (c *Client) DeleteGroup(ctx context.Context, group *Group) error {
	return nil
}
