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
		return nil, errors.Wrap(err, "failed to get clients")
	}
	return clients, nil
}

// GetClient loads client from workspace.
func (c *APIClient) GetClient(ctx context.Context, workspaceID, clientID int) (*Client, error) {
	var client *Client
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "clients", strconv.Itoa(clientID))
	if err := c.httpGet(ctx, apiSpecificPath, nil, &client); err != nil {
		return nil, errors.Wrap(err, "failed to get client")
	}
	return client, nil
}

// CreateClientRequestBody represents a request body of CreateClient.
type CreateClientRequestBody struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	WID  *int    `json:"wid,omitempty"`
}

// CreateClient creates workspace client.
func (c *APIClient) CreateClient(ctx context.Context, workspaceID int, reqBody *CreateClientRequestBody) (*Client, error) {
	var client *Client
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "clients")
	if err := c.httpPost(ctx, apiSpecificPath, reqBody, &client); err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}
	return client, nil
}

// UpdateClientRequestBody represents a request body of UpdateClient.
type UpdateClientRequestBody struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	WID  *int    `json:"wid,omitempty"`
}

// UpdateClient updates workspace client.
func (c *APIClient) UpdateClient(ctx context.Context, workspaceID, clientID int, reqBody *UpdateClientRequestBody) (*Client, error) {
	var client *Client
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "clients", strconv.Itoa(clientID))
	if err := c.httpPut(ctx, apiSpecificPath, reqBody, &client); err != nil {
		return nil, errors.Wrap(err, "failed to update client")
	}
	return client, nil
}

// DeleteClient deletes workspace client.
func (c *APIClient) DeleteClient(ctx context.Context, workspaceID, clientID int) error {
	apiSpecificPath := path.Join(workspacesPath, strconv.Itoa(workspaceID), "clients", strconv.Itoa(clientID))
	if err := c.httpDelete(ctx, apiSpecificPath); err != nil {
		return errors.Wrap(err, "failed to delete client")
	}
	return nil
}
