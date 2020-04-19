package toggl

import (
	"context"
	"errors"
	"strconv"
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

type rawTogglClientData struct {
	TogglClient TogglClient `json:"data"`
}

// GetTogglClient gets client details.
func (c *Client) GetTogglClient(ctx context.Context, togglClient *TogglClient) (*TogglClient, error) {
	if togglClient == nil {
		return nil, ErrTogglClientNotFound
	}
	rawTogglClientData := new(rawTogglClientData)
	endpoint := clientsEndpoint + "/" + strconv.Itoa(togglClient.Id)
	if err := c.httpGet(ctx, c.buildURL(endpoint), rawTogglClientData); err != nil {
		return nil, err
	}
	return &rawTogglClientData.TogglClient, nil
}

// GetTogglClients gets clients visible to user.
func (c *Client) GetTogglClients(ctx context.Context) ([]*TogglClient, error) {
	var togglClients []*TogglClient
	if err := c.httpGet(ctx, c.buildURL(clientsEndpoint), &togglClients); err != nil {
		return nil, err
	}
	return togglClients, nil
}

// GetTogglClientProjects gets client projects.
func (c *Client) GetTogglClientProjects(ctx context.Context, togglClient *TogglClient, params ...QueryString) ([]*Project, error) {
	if togglClient == nil {
		return nil, ErrTogglClientNotFound
	}
	var projects []*Project
	endpoint := clientsEndpoint + "/" + strconv.Itoa(togglClient.Id) + "/projects"
	if err := c.httpGet(ctx, c.buildURL(endpoint, params...), &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

// CreateTogglClient creates a client.
func (c *Client) CreateTogglClient(ctx context.Context, togglClient *TogglClient) (*TogglClient, error) {
	if togglClient == nil {
		return nil, ErrTogglClientNotFound
	}
	rawTogglClientData := new(rawTogglClientData)
	if err := c.httpPost(ctx, c.buildURL(clientsEndpoint), togglClient, rawTogglClientData); err != nil {
		return nil, err
	}
	return &rawTogglClientData.TogglClient, nil
}

// UpdateTogglClient updates a client.
func (c *Client) UpdateTogglClient(ctx context.Context, togglClient *TogglClient) (*TogglClient, error) {
	if togglClient == nil {
		return nil, ErrTogglClientNotFound
	}
	rawTogglClientData := new(rawTogglClientData)
	endpoint := clientsEndpoint + "/" + strconv.Itoa(togglClient.Id)
	if err := c.httpPut(ctx, c.buildURL(endpoint), togglClient, rawTogglClientData); err != nil {
		return nil, err
	}
	return &rawTogglClientData.TogglClient, nil
}

// DeleteTogglClient deletes a client.
func (c *Client) DeleteTogglClient(ctx context.Context, togglClient *TogglClient) error {
	if togglClient == nil {
		return ErrTogglClientNotFound
	}
	endpoint := clientsEndpoint + "/" + strconv.Itoa(togglClient.Id)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}
