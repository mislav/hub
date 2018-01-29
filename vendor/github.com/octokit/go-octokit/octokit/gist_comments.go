package octokit

import (
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var GistCommentsURL = Hyperlink("/gists/{gist_id}/comments{/id}")

// Create a GistCommentsService
func (c *Client) GistComments() (k *GistCommentsService) {
	k = &GistCommentsService{client: c}
	return
}

// A service to return comments for gists
type GistCommentsService struct {
	client *Client
}

// Get a list of all gist comments
//
// https://developer.github.com/v3/gists/comments/#list-comments-on-a-gist
func (c *GistCommentsService) All(uri *Hyperlink, uriParams M) (comments []GistComment, result *Result) {
	url, err := ExpandWithDefault(uri, &GistCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.get(url, &comments)
	return
}

// Get a single comment by id
//
// https://developer.github.com/v3/gists/comments/#get-a-single-comment
func (c *GistCommentsService) One(uri *Hyperlink, uriParams M) (comment *GistComment, result *Result) {
	url, err := ExpandWithDefault(uri, &GistCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.get(url, &comment)
	return
}

// Creates a comment on a gist
//
// https://developer.github.com/v3/gists/comments/#create-a-comment
func (c *GistCommentsService) Create(uri *Hyperlink, uriParams M, requestParams interface{}) (comment *GistComment, result *Result) {
	url, err := ExpandWithDefault(uri, &GistCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.post(url, requestParams, &comment)
	return
}

// Updates a comment on a gist
//
// https://developer.github.com/v3/gists/comments/#edit-a-comment
func (c *GistCommentsService) Update(uri *Hyperlink, uriParams M, requestParams interface{}) (comment *GistComment, result *Result) {
	url, err := ExpandWithDefault(uri, &GistCommentsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = c.client.patch(url, requestParams, &comment)
	return
}

// Deletes a comment on a gist
//
// https://developer.github.com/v3/gists/comments/#delete-a-comment
func (c *GistCommentsService) Delete(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &GistCommentsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = c.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

type GistComment struct {
	*hypermedia.HALResource

	ID        int        `json:"id,omitempty"`
	URL       string     `json:"url,omitempty"`
	User      User       `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Body      string     `json:"body,omitempty"`
}
