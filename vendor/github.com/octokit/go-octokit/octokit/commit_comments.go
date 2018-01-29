package octokit

import (
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// RepoCommentsURL is a template for comments linked to a specific repository
// CommitCommentsURL is a template for comments linked to a specific commit
var (
	RepoCommentsURL   = Hyperlink("/repos/{owner}/{repo}/comments{/id}")
	CommitCommentsURL = Hyperlink("/repos/{owner}/{repo}/commits/{sha}/comments")
)

// Create a CommitCommentsService
func (c *Client) CommitComments() (k *CommitCommentsService) {
	k = &CommitCommentsService{client: c}
	return
}

// A service to return comments for commits
type CommitCommentsService struct {
	client *Client
}

// All commit comments for a single commit
//
// https://developer.github.com/v3/repos/comments/#list-comments-for-a-single-commit
func (c *CommitCommentsService) All(uri *Hyperlink, uriParams M) (comments []CommitComment, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.get(url, &comments)
	return
}

// One commit comment by id 
//
// https://developer.github.com/v3/repos/comments/#get-a-single-commit-comment
func (c *CommitCommentsService) One(uri *Hyperlink, uriParams M) (comment *CommitComment, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.get(url, &comment)
	return
}

// Creates a comment on a commit
//
// https://developer.github.com/v3/repos/comments/#create-a-commit-comment
func (c *CommitCommentsService) Create(uri *Hyperlink, uriParams M, requestParams interface{}) (comment *CommitComment, result *Result) {
	url, err := ExpandWithDefault(uri, &CommitCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.post(url, requestParams, &comment)
	return
}

// Updates a comment on a commit
//
// https://developer.github.com/v3/repos/comments/#update-a-commit-comment
func (c *CommitCommentsService) Update(uri *Hyperlink, uriParams M, requestParams interface{}) (comment *CommitComment, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.patch(url, requestParams, &comment)
	return
}

// Deletes a comment on a commit
//
// https://developer.github.com/v3/repos/comments/#delete-a-commit-comment
func (c *CommitCommentsService) Delete(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoCommentsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = c.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

type CommitComment struct {
	*hypermedia.HALResource

	ID        int        `json:"id,omitempty"`
	URL       string     `json:"url,omitempty"`
	User      User       `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Body      string     `json:"body,omitempty"`
	HTMLURL   string     `json:"html_url,omitempty"`
	Position  int        `json:"position,omitempty"`
	Line      int        `json:"line,omitempty"`
	Path      string     `json:"path,omitempty"`
	CommitID  string     `json:"commit_id,omitempty"`
}
