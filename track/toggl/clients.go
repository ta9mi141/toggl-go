package toggl

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Client represents the properties of a client.
type Client struct {
	ID              *int       `json:"id,omitempty"`
	WID             *int       `json:"wid,omitempty"`
	Name            *string    `json:"name,omitempty"`
	At              *time.Time `json:"at,omitempty"`
	ServerDeletedAt *time.Time `json:"server_deleted_at,omitempty"`
	Archived        *bool      `json:"archived,omitempty"`
}

// GetClients lists clients from workspace.
func (c *APIClient) GetClients(ctx context.Context, workspaceID int) ([]*Client, error) {
	var clients []*Client
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "clients")
	if err := c.httpGet(ctx, apiSpecificPath, nil, &clients); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return clients, nil
}

// GetClient loads client from workspace.
func (c *APIClient) GetClient(ctx context.Context, workspaceID, clientID int) (*Client, error) {
	var client *Client
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "clients", strconv.Itoa(clientID))
	if err := c.httpGet(ctx, apiSpecificPath, nil, &client); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return client, nil
}
