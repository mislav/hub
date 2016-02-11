package octokit

import (
	"io"
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// URLs for accessing specific endpoints on the gists API. The most general
// is GistsURL, used to list, get, create, edit and delete a gist. GistsUserURL,
// GistsPublicURL and GistsStarredURL are used for listing certain sets of gists.
// GistsRevisionURL can be used to access a specific revision of a gist.
// GistsCommitsURL can be used to access all the commits on a specific gist.
// GistsStarURL is used for starring and unstarring gists, and checking if a gist
// is starred. GistsForksURL is used to fork or delete a fork of a gist, and to
// list a gist's forks.
//
// https://developer.github.com/v3/gists
var (
	GistsUserURL     = Hyperlink("users/{username}/gists")
	GistsURL         = Hyperlink("gists{/gist_id}")
	GistsPublicURL   = Hyperlink("gists/public")
	GistsStarredURL  = Hyperlink("gists/starred")
	GistsRevisionURL = Hyperlink("gists/{gist_id}/{commit_sha}")
	GistsCommitsURL  = Hyperlink("gists/{gist_id}/commits")
	GistsStarURL     = Hyperlink("gists/{gist_id}/star")
	GistsForksURL    = Hyperlink("gists/{gist_id}/forks")
)

// Gists creates a GistsService to be used with any proper URL
//
// https://developer.github.com/v3/gists/
func (c *Client) Gists() (gists *GistsService) {
	gists = &GistsService{client: c}
	return
}

// GistsService is a service providing access to gists from a particular url
type GistsService struct {
	client *Client
}

// All gets a list of all gists associated with the url of the service
//
// https://developer.github.com/v3/gists/#list-gists
func (g *GistsService) All(uri *Hyperlink, uriParams M) (gists []Gist, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = g.client.get(url, &gists)
	return
}

// One gets a specific gist based on the url of the service
//
// https://developer.github.com/v3/gists/#get-a-single-gist
// https://developer.github.com/v3/gists/#get-a-specific-revision-of-a-gist
func (g *GistsService) One(uri *Hyperlink, uriParams M) (gist *Gist, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = g.client.get(url, &gist)
	return
}

// Raw gets the raw contents of first file in a specific gist
//
// https://developer.github.com/v3/gists/#truncation
func (g *GistsService) Raw(uri *Hyperlink, uriParams M) (body io.ReadCloser, result *Result) {
	var rawURL *url.URL

	gist, result := g.One(uri, uriParams)
	for _, file := range gist.Files {
		rawURL, _ = url.Parse(file.RawURL)
		break
	}

	body, result = g.client.getBody(rawURL, textMediaType)
	return
}

// Create posts a new gist based on parameters in a Gist struct to
// the specified URL
//
// https://developer.github.com/v3/gists/#create-a-gist
func (g *GistsService) Create(uri *Hyperlink, uriParams M, requestParams interface{}) (gist *Gist, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = g.client.post(url, requestParams, &gist)
	return
}

// Update modifies a specific gist based on the url of the service
//
// https://developer.github.com/v3/gists/#edit-a-gist
func (g *GistsService) Update(uri *Hyperlink, uriParams M, requestParams interface{}) (gist *Gist, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = g.client.patch(url, requestParams, &gist)
	return
}

// Commits gets a list of all commits to the given gist
//
// https://developer.github.com/v3/gists/#list-gist-commits
func (g *GistsService) Commits(uri *Hyperlink, uriParams M) (gistCommits []GistCommit, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsCommitsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = g.client.get(url, &gistCommits)
	return
}

// Star stars a gist
//
// https://developer.github.com/v3/gists/#star-a-gist
func (g *GistsService) Star(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsStarURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = g.client.put(url, nil, nil)
	success = (!result.HasError() && result.Response.StatusCode == 204)
	return
}

// Unstar unstars a gist
//
// https://developer.github.com/v3/gists/#unstar-a-gist
func (g *GistsService) Unstar(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsStarURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = g.client.delete(url, nil, nil)
	success = (!result.HasError() && result.Response.StatusCode == 204)
	return
}

// CheckStar checks if a gist is starred
//
// https://developer.github.com/v3/gists/#check-if-a-gist-is-starred
func (g *GistsService) CheckStar(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsStarURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = g.client.get(url, nil)
	success = (!result.HasError() && result.Response.StatusCode == 204)
	return
}

// Fork forks a gist
//
// https://developer.github.com/v3/gists/#fork-a-gist
func (g *GistsService) Fork(uri *Hyperlink, uriParams M) (gist *Gist, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsForksURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = g.client.post(url, nil, &gist)
	return
}

// ListForks lists all the forks of a gist
//
// https://developer.github.com/v3/gists/#list-gist-forks
func (g *GistsService) ListForks(uri *Hyperlink, uriParams M) (gistForks []GistFork, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsForksURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = g.client.get(url, &gistForks)
	return
}

// Delete deletes a gist by its id
//
// https://developer.github.com/v3/gists/#delete-a-gist
func (g *GistsService) Delete(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &GistsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = g.client.delete(url, nil, nil)
	success = (!result.HasError() && result.Response.StatusCode == 204)
	return
}

// GistFile is a representation of the file stored in a gist
type GistFile struct {
	*hypermedia.HALResource

	FileName  string `json:"filename,omitempty"`
	Type      string `json:"type,omitempty"`
	Language  string `json:"language,omitempty"`
	RawURL    string `json:"raw_url,omitempty"`
	Size      int    `json:"size,omitempty"`
	Truncated bool   `json:"truncated,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Gist is a representation of a gist on github, a standalone file that acts as a
// sole element of its own repository
type Gist struct {
	*hypermedia.HALResource

	ID          string               `json:"id,omitempty"`
	Comments    float64              `json:"comments,omitempty"`
	CommentsURL string               `json:"comments_url,omitempty"`
	CommitsURL  string               `json:"commits_url,omitempty"`
	CreatedAt   *time.Time           `json:"created_at,omitempty"`
	Description string               `json:"description,omitempty"`
	Files       map[string]*GistFile `json:"files,omitempty"`
	Forks       []GistFork           `json:"forks,omitempty"`
	ForksURL    Hyperlink            `json:"forks_url,omitempty"`
	GitPullURL  Hyperlink            `json:"git_pull_url,omitempty"`
	History     []GistCommit         `json:"history,omitempty"`
	GitPushURL  Hyperlink            `json:"git_push_url,omitempty"`
	HtmlURL     Hyperlink            `json:"html_url,omitempty"`
	Owner       *User                `json:"owner,omitempty"`
	Public      bool                 `json:"public,omitempty"`
	UpdatedAt   *time.Time           `json:"updated_at,omitempty"`
	URL         string               `json:"url,omitempty"`
	User        *User                `json:"user,omitempty"`
}

//GistFork represents information about a fork of a gist
type GistFork struct {
	*hypermedia.HALResource

	ID        string     `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	URL       string     `json:"url,omitempty"`
	User      *User      `json:"user,omitempty"`
}

// GistCommit is a representation of one of the commits to a gist
type GistCommit struct {
	*hypermedia.HALResource

	ChangeStatus *GistChangeStatus `json:"change_status,omitempty"`
	CommittedAt  *time.Time        `json:"committed_at,omitempty"`
	URL          string            `json:"url,omitempty"`
	User         *User             `json:"user,omitempty"`
	Version      string            `json:"version,omitempty"`
}

// GistChangeStatus represents all changes on a given Gist
type GistChangeStatus struct {
	Additions int `json:"additions,omitempty"`
	Deletions int `json:"deletions,omitempty"`
	Total     int `json:"total,omitempty"`
}
