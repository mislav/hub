package octokit

import (
	"github.com/jingweno/go-sawyer/hypermedia"
	"time"
)

// RepoIssuesURL is a template for accessing issues in a particular
// repository for a particular owner that can be expanded to a full address.
//
// https://developer.github.com/v3/issues/
var RepoIssuesURL = Hyperlink("repos/{owner}/{repo}/issues{/number}{?filter,state,labels,sort}")

// Issues creates an IssuesService with a base url
func (c *Client) Issues() (issues *IssuesService) {
	issues = &IssuesService{client: c}
	return
}

// IssuesService is a service providing access to issues from a particular url
type IssuesService struct {
	client *Client
}

// One gets a specific issue based on the url of the service
//
// https://developer.github.com/v3/issues/#get-a-single-issue
func (i *IssuesService) One(uri *Hyperlink, uriParams M) (issue *Issue, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoIssuesURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = i.client.get(url, &issue)
	return
}

// All gets a list of all issues associated with the url of the service
//
// https://developer.github.com/v3/issues/#list-issues-for-a-repository
func (i *IssuesService) All(uri *Hyperlink, uriParams M) (issues []Issue, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoIssuesURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = i.client.get(url, &issues)
	return
}

// Create posts a new issue with particular parameters to the issues service url
//
// https://developer.github.com/v3/issues/#create-an-issue
func (i *IssuesService) Create(uri *Hyperlink, uriParams M, requestParams interface{}) (issue *Issue, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoIssuesURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = i.client.post(url, requestParams, &issue)
	return
}

// Update modifies a specific issue given parameters on the service url
//
// https://developer.github.com/v3/issues/#edit-an-issue
func (i *IssuesService) Update(uri *Hyperlink, uriParams M, requestParams interface{}) (issue *Issue, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoIssuesURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = i.client.patch(url, requestParams, &issue)
	return
}

// Issue represents an issue on GitHub with all associated fields
type Issue struct {
	*hypermedia.HALResource

	URL     string `json:"url,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
	Number  int    `json:"number,omitempty"`
	State   string `json:"state,omitempty"`
	Title   string `json:"title,omitempty"`
	Body    string `json:"body,omitempty"`
	User    User   `json:"user,omitempty"`
	Labels  []struct {
		URL   string `json:"url,omitempty"`
		Name  string `json:"name,omitempty"`
		Color string `json:"color,omitempty"`
	}
	Assignee  User `json:"assignee,omitempty"`
	Milestone struct {
		URL          string     `json:"url,omitempty"`
		Number       int        `json:"number,omitempty"`
		State        string     `json:"state,omitempty"`
		Title        string     `json:"title,omitempty"`
		Description  string     `json:"description,omitempty"`
		Creator      User       `json:"creator,omitempty"`
		OpenIssues   int        `json:"open_issues,omitempty"`
		ClosedIssues int        `json:"closed_issues,omitempty"`
		CreatedAt    time.Time  `json:"created_at,omitempty"`
		DueOn        *time.Time `json:"due_on,omitempty"`
	}
	Comments    int `json:"comments,omitempty"`
	PullRequest struct {
		HTMLURL  string `json:"html_url,omitempty"`
		DiffURL  string `json:"diff_url,omitempty"`
		PatchURL string `json:"patch_url,omitempty"`
	} `json:"pull_request,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
}

// IssueParams represents the struture used to create or update an Issue
type IssueParams struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	State     string   `json:"state,omitempty"`
	Milestone uint64   `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}
