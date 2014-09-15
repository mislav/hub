package octokit

import (
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var (
	RepoIssuesURL = Hyperlink("repos/{owner}/{repo}/issues{/number}")
)

// Create a IssuesService with the base url.URL
func (c *Client) Issues(url *url.URL) (issues *IssuesService) {
	issues = &IssuesService{client: c, URL: url}
	return
}

type IssuesService struct {
	client *Client
	URL    *url.URL
}

func (i *IssuesService) One() (issue *Issue, result *Result) {
	result = i.client.get(i.URL, &issue)
	return
}

func (i *IssuesService) All() (issues []Issue, result *Result) {
	result = i.client.get(i.URL, &issues)
	return
}

func (i *IssuesService) Create(params interface{}) (issue *Issue, result *Result) {
	result = i.client.post(i.URL, params, &issue)
	return
}

func (i *IssuesService) Update(params interface{}) (issue *Issue, result *Result) {
	result = i.client.patch(i.URL, params, &issue)
	return
}

type Issue struct {
	*hypermedia.HALResource

	URL     string `json:"url,omitempty,omitempty"`
	HTMLURL string `json:"html_url,omitempty,omitempty"`
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

type IssueParams struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	State     string   `json:"state,omitempty"`
	Milestone uint64   `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}
