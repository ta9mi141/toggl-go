package toggl

const (
	clientsEndpoint string = "/api/v8/clients"
)

// TogglClient represents properties of client.
type TogglClient struct {
	Name  string `json:"name"`
	Wid   int    `json:"wid"`
	Notes string `json:"notes"`
	At    string `json:"at"`
}
