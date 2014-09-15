package octokit

import (
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var (
	StatusesURL = Hyperlink("repos/{owner}/{repo}/statuses/{ref}")
)

// Create a StatusesService with the base url.URL
func (c *Client) Statuses(url *url.URL) (statuses *StatusesService) {
	statuses = &StatusesService{client: c, URL: url}
	return
}

type StatusesService struct {
	client *Client
	URL    *url.URL
}

func (s *StatusesService) All() (statuses []Status, result *Result) {
	result = s.client.get(s.URL, &statuses)
	return
}

type Status struct {
	*hypermedia.HALResource

	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	State       string    `json:"state,omitempty"`
	TargetURL   string    `json:"target_url,omitempty"`
	Description string    `json:"description,omitempty"`
	ID          int       `json:"id,omitempty"`
	URL         string    `json:"url,omitempty"`
	Creator     User      `json:"creator,omitempty"`
}
