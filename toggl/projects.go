package toggl

import (
	"context"
	"errors"
	"strconv"
	"time"
)

const (
	projectsEndpoint string = "/api/v8/projects"
)

// Project represents properties of project.
type Project struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Wid        int       `json:"wid"`
	Cid        int       `json:"cid"`
	Active     bool      `json:"active"`
	IsPrivate  bool      `json:"is_private"`
	Template   bool      `json:"template"`
	TemplateId int       `json:"template_id"`
	At         time.Time `json:"at"`
	Color      string    `json:"color"`
	CreatedAt  time.Time `json:"created_at"`
}

type rawProjectData struct {
	Project Project `json:"data"`
}

var (
	// ErrProjectNotFound is returned when the provided project is nil.
	ErrProjectNotFound = errors.New("The provided project must be non-nil")
)

// CreateProject creates a project.
func (c *Client) CreateProject(ctx context.Context, project *Project) (*Project, error) {
	if project == nil {
		return nil, ErrProjectNotFound
	}
	rawProjectData := new(rawProjectData)
	if err := c.httpPost(ctx, c.buildURL(projectsEndpoint), project, rawProjectData); err != nil {
		return nil, err
	}
	return &rawProjectData.Project, nil
}

// UpdateProject updates a project.
func (c *Client) UpdateProject(ctx context.Context, project *Project) (*Project, error) {
	if project == nil {
		return nil, ErrProjectNotFound
	}
	rawProjectData := new(rawProjectData)
	endpoint := projectsEndpoint + "/" + strconv.Itoa(project.Id)
	if err := c.httpPut(ctx, c.buildURL(endpoint), project, rawProjectData); err != nil {
		return nil, err
	}
	return &rawProjectData.Project, nil
}

// DeleteProject deletes a project.
func (c *Client) DeleteProject(ctx context.Context, project *Project) error {
	if project == nil {
		return ErrProjectNotFound
	}
	endpoint := projectsEndpoint + "/" + strconv.Itoa(project.Id)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}

// GetProject gets a project.
func (c *Client) GetProject(ctx context.Context, project *Project) (*Project, error) {
	if project == nil {
		return nil, ErrProjectNotFound
	}
	rawProjectData := new(rawProjectData)
	endpoint := projectsEndpoint + "/" + strconv.Itoa(project.Id)
	if err := c.httpGet(ctx, c.buildURL(endpoint), rawProjectData); err != nil {
		return nil, err
	}
	return &rawProjectData.Project, nil
}

// GetProjectUsers gets project users.
func (c *Client) GetProjectUsers(ctx context.Context, project *Project) ([]*ProjectUser, error) {
	if project == nil {
		return nil, ErrProjectNotFound
	}
	var projectUsers []*ProjectUser
	endpoint := projectsEndpoint + "/" + strconv.Itoa(project.Id) + "/project_users"
	if err := c.httpGet(ctx, c.buildURL(endpoint), &projectUsers); err != nil {
		return nil, err
	}
	return projectUsers, nil
}
