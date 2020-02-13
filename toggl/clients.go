package toggl

import (
	"context"
	"errors"
	"time"
)

const (
	clientsEndpoint string = "/api/v8/clients"
)

var (
	// ErrTogglClinetNotFound is returned when the provided toggl client is nil.
	ErrTogglClientNotFound = errors.New("The provided toggl client must be non-nil")
)

// TogglClient represents properties of client.
type TogglClient struct {
	Id    int       `json:"id"`
	Name  string    `json:"name"`
	Wid   int       `json:"wid"`
	Notes string    `json:"notes"`
	At    time.Time `json:"at"`
}

func (c *Client) GetTogglClient(ctx context.Context, togglClient *TogglClient) (*TogglClient, error) {
	return nil, nil
}

func (c *Client) GetTogglClients(ctx context.Context) ([]*TogglClient, error) {
	return nil, nil
}

func (c *Client) GetTogglClientProjects(ctx context.Context, togglClient *TogglClient, active ...string) ([]*Project, error) {
	return nil, nil
}

func (c *Client) CreateTogglClient(ctx context.Context, togglClient *TogglClient) (*TogglClient, error) {
	return nil, nil
}

func (c *Client) UpdateTogglClient(ctx context.Context, togglClient *TogglClient) (*TogglClient, error) {
	return nil, nil
}

func (c *Client) DeleteTogglClient(ctx context.Context, togglClient *TogglClient) error {
	return nil
}
