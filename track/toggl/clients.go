package toggl

import (
	"context"
	"time"
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
	return nil, nil
}
