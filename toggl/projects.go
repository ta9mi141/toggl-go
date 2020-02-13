package toggl

import "time"

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
	Color      int       `json:"color"`
	CreatedAt  string    `json:"created_at"`
}
