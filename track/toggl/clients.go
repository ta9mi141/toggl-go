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
	ForeignID       *string    `json:"foreign_id,omitempty"`
	ServerDeletedAt *time.Time `json:"server_deleted_at,omitempty"`
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
