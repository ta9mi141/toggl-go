package toggl

const (
	tagsEndpoint string = "/api/v8/tags"
)

// Tag represents properties of tag.
type Tag struct {
	Name string `json:"name"`
	Wid  int    `json:"wid"`
}
