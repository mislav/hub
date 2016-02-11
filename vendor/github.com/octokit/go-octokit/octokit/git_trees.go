package octokit

import (
	"net/url"
)

// GitTreesURL is a template for accessing git trees at a particular sha hash or branch
// of a particular repository of a particular owner. The request may be set to be
// recursive for a particular level of depth (0 is no recursion) to follow sub-trees from
// the primary repository.
//
// https://developer.github.com/v3/git/trees/
var GitTreesURL = Hyperlink("repos/{owner}/{repo}/git/trees/{sha}{?recursive}")

// GitTrees creates a GitTreesService with a base url
//
// https://developer.github.com/v3/git/trees/
func (c *Client) GitTrees(url *url.URL) (trees *GitTreesService) {
	trees = &GitTreesService{client: c, URL: url}
	return
}

// GitTreesService is a service providing access to GitTrees from a particular url
type GitTreesService struct {
	client *Client
	URL    *url.URL
}

// One gets a specific GitTree based on the url of the service. May specify
// to get the tree recursively
//
// https://developer.github.com/v3/git/trees/#get-a-tree
// https://developer.github.com/v3/git/trees/#get-a-tree-recursively
func (c *GitTreesService) One() (tree *GitTree, result *Result) {
	result = c.client.get(c.URL, &tree)
	return
}

// GitTree represents a tree on GitHub, a level in the GitHub equivalent of a
// directory structure
type GitTree struct {
	Sha       string         `json:"sha,omitempty"`
	Tree      []GitTreeEntry `json:"tree,omitempty"`
	Truncated bool           `json:"truncated,omitempty"`
	URL       string         `json:"url,omitempty"`
}

// GitTreeEntry represents an element within a GitTree on GitHub, which may either
// be another Tree or a single Blob
type GitTreeEntry struct {
	Mode string `json:"mode,omitempty"`
	Path string `json:"path,omitempty"`
	Sha  string `json:"sha,omitempty"`
	Size int    `json:"size,omitempty"`
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}
