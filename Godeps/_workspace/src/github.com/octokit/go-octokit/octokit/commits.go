package octokit

import (
	"io"
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var CommitsURL = Hyperlink("repos/{owner}/{repo}/commits{/sha}")

func (c *Client) Commits(url *url.URL) (commits *CommitsService) {
	commits = &CommitsService{client: c, URL: url}
	return
}

type CommitsService struct {
	client *Client
	URL    *url.URL
}

// Get all commits on CommitsService#URL
func (c *CommitsService) All() (commits []Commit, result *Result) {
	result = c.client.get(c.URL, &commits)
	return
}

// Get a commit based on CommitsService#URL
func (c *CommitsService) One() (commit *Commit, result *Result) {
	result = c.client.get(c.URL, &commit)
	return
}

// Get a commit patch based on CommitsService#URL
func (c *CommitsService) Patch() (patch io.ReadCloser, result *Result) {
	patch, result = c.client.getBody(c.URL, patchMediaType)
	return
}

type CommitFile struct {
	Additions   int    `json:"additions,omitempty"`
	BlobURL     string `json:"blob_url,omitempty"`
	Changes     int    `json:"changes,omitempty"`
	ContentsURL string `json:"contents_url,omitempty"`
	Deletions   int    `json:"deletions,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Patch       string `json:"patch,omitempty"`
	RawURL      string `json:"raw_url,omitempty"`
	Sha         string `json:"sha,omitempty"`
	Status      string `json:"status,omitempty"`
}

type CommitStats struct {
	Additions int `json:"additions,omitempty"`
	Deletions int `json:"deletions,omitempty"`
	Total     int `json:"total,omitempty"`
}

type CommitCommit struct {
	Author struct {
		Date  *time.Time `json:"date,omitempty"`
		Email string     `json:"email,omitempty"`
		Name  string     `json:"name,omitempty"`
	} `json:"author,omitempty"`
	CommentCount int `json:"comment_count,omitempty"`
	Committer    struct {
		Date  *time.Time `json:"date,omitempty"`
		Email string     `json:"email,omitempty"`
		Name  string     `json:"name,omitempty"`
	} `json:"committer,omitempty"`
	Message string `json:"message,omitempty"`
	Tree    struct {
		Sha string `json:"sha,omitempty"`
		URL string `json:"url,omitempty"`
	} `json:"tree,omitempty"`
	URL string `json:"url,omitempty"`
}

type Commit struct {
	*hypermedia.HALResource

	Author      *User         `json:"author,omitempty"`
	CommentsURL string        `json:"comments_url,omitempty"`
	Commit      *CommitCommit `json:"commit,omitempty"`
	Committer   *User         `json:"committer,omitempty"`
	Files       []CommitFile  `json:"files,omitempty"`
	HtmlURL     string        `json:"html_url,omitempty"`
	Parents     []Commit      `json:"parents,omitempty"`
	Sha         string        `json:"sha,omitempty"`
	Stats       CommitStats   `json:"stats,omitempty"`
	URL         string        `json:"url,omitempty"`
}
