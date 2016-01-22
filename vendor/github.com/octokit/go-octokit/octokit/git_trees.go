package octokit

import (
	"net/url"
)

var GitTreesURL = Hyperlink("repos/{owner}/{repo}/git/trees/{sha}{?recursive}")

func (c *Client) GitTrees(url *url.URL) (trees *GitTreesService) {
	trees = &GitTreesService{client: c, URL: url}
	return
}

type GitTreesService struct {
	client *Client
	URL    *url.URL
}

// Get a Git Tree
func (c *GitTreesService) One() (tree *GitTree, result *Result) {
	result = c.client.get(c.URL, &tree)
	return
}

type GitTree struct {
	Sha       string         `json:"sha,omitempty"`
	Tree      []GitTreeEntry `json:"tree,omitempty"`
	Truncated bool           `json:"truncated,omitempty"`
	URL       string         `json:"url,omitempty"`
}

type GitTreeEntry struct {
	Mode string `json:"mode,omitempty"`
	Path string `json:"path,omitempty"`
	Sha  string `json:"sha,omitempty"`
	Size int    `json:"size,omitempty"`
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}
