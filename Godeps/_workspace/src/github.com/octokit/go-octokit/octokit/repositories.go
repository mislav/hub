package octokit

import (
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var (
	RepositoryURL       = Hyperlink("repos/{owner}/{repo}")
	ForksURL            = Hyperlink("repos/{owner}/{repo}/forks")
	UserRepositoriesURL = Hyperlink("user/repos")
	OrgRepositoriesURL  = Hyperlink("orgs/{org}/repos")
)

// Create a RepositoriesService with the base url.URL
func (c *Client) Repositories(url *url.URL) (repos *RepositoriesService) {
	repos = &RepositoriesService{client: c, URL: url}
	return
}

type RepositoriesService struct {
	client *Client
	URL    *url.URL
}

func (r *RepositoriesService) One() (repo *Repository, result *Result) {
	result = r.client.get(r.URL, &repo)
	return
}

func (r *RepositoriesService) All() (repos []Repository, result *Result) {
	result = r.client.get(r.URL, &repos)
	return
}

func (r *RepositoriesService) Create(params interface{}) (repo *Repository, result *Result) {
	result = r.client.post(r.URL, params, &repo)
	return
}

type Repository struct {
	*hypermedia.HALResource

	ID            int           `json:"id,omitempty"`
	Owner         User          `json:"owner,omitempty"`
	Name          string        `json:"name,omitempty"`
	FullName      string        `json:"full_name,omitempty"`
	Description   string        `json:"description,omitempty"`
	Private       bool          `json:"private"`
	Fork          bool          `json:"fork,omitempty"`
	URL           string        `json:"url,omitempty"`
	HTMLURL       string        `json:"html_url,omitempty"`
	CloneURL      string        `json:"clone_url,omitempty"`
	GitURL        string        `json:"git_url,omitempty"`
	SSHURL        string        `json:"ssh_url,omitempty"`
	SVNURL        string        `json:"svn_url,omitempty"`
	MirrorURL     string        `json:"mirror_url,omitempty"`
	Homepage      string        `json:"homepage,omitempty"`
	Language      string        `json:"language,omitempty"`
	Forks         int           `json:"forks,omitempty"`
	ForksCount    int           `json:"forks_count,omitempty"`
	Watchers      int           `json:"watchers,omitempty"`
	WatchersCount int           `json:"watchers_count,omitempty"`
	Size          int           `json:"size,omitempty"`
	MasterBranch  string        `json:"master_branch,omitempty"`
	OpenIssues    int           `json:"open_issues,omitempty"`
	PushedAt      time.Time     `json:"pushed_at,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
	Permissions   Permissions   `json:"permissions,omitempty"`
	Organization  *Organization `json:"organization,omitempty"`
	Parent        *Repository   `json:"parent,omitempty"`
	Source        *Repository   `json:"source,omitempty"`
	HasIssues     bool          `json:"has_issues,omitempty"`
	HasWiki       bool          `json:"has_wiki,omitempty"`
	HasDownloads  bool          `json:"has_downloads,omitempty"`
}

type Permissions struct {
	Admin bool
	Push  bool
	Pull  bool
}
