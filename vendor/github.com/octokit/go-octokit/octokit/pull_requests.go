package octokit

import (
	"io"
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var (
	PullRequestsURL = Hyperlink("repos/{owner}/{repo}/pulls{/number}")
)

// Create a PullRequestsService with the base url.URL
func (c *Client) PullRequests(url *url.URL) (pullRequests *PullRequestsService) {
	pullRequests = &PullRequestsService{client: c, URL: url}
	return
}

type PullRequestsService struct {
	client *Client
	URL    *url.URL
}

func (p *PullRequestsService) One() (pull *PullRequest, result *Result) {
	result = p.client.get(p.URL, &pull)
	return
}

func (p *PullRequestsService) Create(params interface{}) (pull *PullRequest, result *Result) {
	result = p.client.post(p.URL, params, &pull)
	return
}

func (p *PullRequestsService) All() (pulls []PullRequest, result *Result) {
	result = p.client.get(p.URL, &pulls)
	return
}

func (p *PullRequestsService) Diff() (diff io.ReadCloser, result *Result) {
	return p.client.getBody(p.URL, diffMediaType)
}

func (p *PullRequestsService) Patch() (patch io.ReadCloser, result *Result) {
	return p.client.getBody(p.URL, patchMediaType)
}

type PullRequest struct {
	*hypermedia.HALResource

	URL               string            `json:"url,omitempty"`
	ID                int               `json:"id,omitempty"`
	HTMLURL           string            `json:"html_url,omitempty"`
	DiffURL           string            `json:"diff_url,omitempty"`
	PatchURL          string            `json:"patch_url,omitempty"`
	IssueURL          string            `json:"issue_url,omitempty"`
	Title             string            `json:"title,omitempty"`
	Number            int               `json:"number,omitempty"`
	State             string            `json:"state,omitempty"`
	User              User              `json:"user,omitempty"`
	Body              string            `json:"body,omitempty"`
	CreatedAt         time.Time         `json:"created_at,omitempty"`
	UpdatedAt         time.Time         `json:"updated_at,omitempty"`
	ClosedAt          *time.Time        `json:"closed_at,omitempty"`
	MergedAt          *time.Time        `json:"merged_at,omitempty"`
	MergeCommitSha    string            `json:"merge_commit_sha,omitempty"`
	Assignee          *User             `json:"assignee,omitempty"`
	CommitsURL        string            `json:"commits_url,omitempty"`
	ReviewCommentsURL string            `json:"review_comments_url,omitempty"`
	ReviewCommentURL  string            `json:"review_comment_url,omitempty"`
	CommentsURL       string            `json:"comments_url,omitempty"`
	Head              PullRequestCommit `json:"head,omitempty"`
	Base              PullRequestCommit `json:"base,omitempty"`
	Merged            bool              `json:"merged,omitempty"`
	MergedBy          User              `json:"merged_by,omitempty"`
	Comments          int               `json:"comments,omitempty"`
	ReviewComments    int               `json:"review_comments,omitempty"`
	Commits           int               `json:"commits,omitempty"`
	Additions         int               `json:"additions,omitempty"`
	Deletions         int               `json:"deletions,omitempty"`
	ChangedFiles      int               `json:"changed_files,omitempty"`
}

type PullRequestCommit struct {
	Label string      `json:"label,omitempty"`
	Ref   string      `json:"ref,omitempty"`
	Sha   string      `json:"sha,omitempty"`
	User  User        `json:"user,omitempty"`
	Repo  *Repository `json:"repo,omitempty"`
}

type PullRequestParams struct {
	Base  string `json:"base,omitempty"`
	Head  string `json:"head,omitempty"`
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type PullRequestForIssueParams struct {
	Base  string `json:"base,omitempty"`
	Head  string `json:"head,omitempty"`
	Issue string `json:"issue,omitempty"`
}
