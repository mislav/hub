package octokit

import (
	"net/url"
)

var (
	EmojisURL = Hyperlink("/emojis")
)

// Create a EmojisService with the base url.URL
func (c *Client) Emojis(url *url.URL) (emojis *EmojisService) {
	emojis = &EmojisService{client: c, URL: url}
	return
}

type EmojisService struct {
	client *Client
	URL    *url.URL
}

func (s *EmojisService) All() (emojis map[string]string, result *Result) {
	result = s.client.get(s.URL, &emojis)
	return
}
