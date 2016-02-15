package octokit

import (
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var IssueCommentsURL = Hyperlink("/repos/{owner}/{repo}/issues{/number}/comments{/id}")

// Create a IssueCommentsService
func (c *Client) IssueComments() (k *IssueCommentsService) {
	k = &IssueCommentsService{client: c}
	return
}

// A service to return comments for issues
type IssueCommentsService struct {
	client *Client
}

// Get a list of all issue comments
//
// https://developer.github.com/v3/issues/comments/#list-comments-on-an-issue
func (c *IssueCommentsService) All(uri *Hyperlink, uriParams M) (comments []IssueComment, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.get(url, &comments)
	return
}

// Get a single comment by id
//
// https://developer.github.com/v3/issues/comments/#get-a-single-comment
func (c *IssueCommentsService) One(uri *Hyperlink, uriParams M) (comment *IssueComment, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.get(url, &comment)
	return
}

// Creates a comment on an issue
//
// https://developer.github.com/v3/issues/comments/#create-a-comment
func (c *IssueCommentsService) Create(uri *Hyperlink, uriParams M, requestParams interface{}) (comment *IssueComment, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.post(url, requestParams, &comment)
	return
}

// Updates a comment on an issue
//
// https://developer.github.com/v3/issues/comments/#edit-a-comment
func (c *IssueCommentsService) Update(uri *Hyperlink, uriParams M, requestParams interface{}) (comment *IssueComment, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.patch(url, requestParams, &comment)
	return
}

// Deletes a comment on an issue
//
// https://developer.github.com/v3/issues/comments/#delete-a-comment
func (c *IssueCommentsService) Delete(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueCommentsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = c.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

type IssueComment struct {
	*hypermedia.HALResource

	ID        int        `json:"id,omitempty"`
	URL       string     `json:"url,omitempty"`
	User      User       `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Body      string     `json:"body,omitempty"`
	HTMLURL   string     `json:"html_url,omitempty"`
	IssueURL  string     `json:"issue_url",omitempty"`
}
