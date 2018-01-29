package octokit

import (
	"io"
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

const (
	MergeStateClean    = "clean"
	MergeStateUnstable = "unstable"
	MergeStateDirty    = "dirty"
	MergeStateUnknown  = "unknown"
)

// PullRequestsURL is a template for accessing pull requests in a particular
// repository for a particular owner that can be expanded to a full address.
//
// https://developer.github.com/v3/pulls/
var PullRequestsURL = Hyperlink("repos/{owner}/{repo}/pulls{/number}")

// PullRequests creates a PullRequestsService with a base url
//
// https://developer.github.com/v3/pulls/
func (c *Client) PullRequests(url *url.URL) (pullRequests *PullRequestsService) {
	pullRequests = &PullRequestsService{client: c, URL: url}
	return
}

// PullRequestService is a service providing access to pull requests from
// a particular url
type PullRequestsService struct {
	client *Client
	URL    *url.URL
}

// One gets a specific pull request based on the url of the service
//
// https://developer.github.com/v3/pulls/#get-a-single-pull-request
func (p *PullRequestsService) One() (pull *PullRequest, result *Result) {
	result = p.client.get(p.URL, &pull)
	return
}

// Create posts a new pull request based on the parameters given to the
// pull request service url
//
// https://developer.github.com/v3/pulls/#create-a-pull-request
func (p *PullRequestsService) Create(params interface{}) (pull *PullRequest, result *Result) {
	result = p.client.post(p.URL, params, &pull)
	return
}

// All gets a list of all pull requests associated with the url of the service
//
// https://developer.github.com/v3/pulls/#list-pull-requests
func (p *PullRequestsService) All() (pulls []PullRequest, result *Result) {
	result = p.client.get(p.URL, &pulls)
	return
}

// Diff gets the difference of the data in the specific pull request to the branch
// that the pull request is out to be merged with
//
// https://developer.github.com/v3/pulls/#get-a-single-pull-request
func (p *PullRequestsService) Diff() (diff io.ReadCloser, result *Result) {
	return p.client.getBody(p.URL, diffMediaType)
}

// Patch gets all the patches made to the specific pull request associated with
// the url given to the service
//
// https://developer.github.com/v3/pulls/#update-a-pull-request
func (p *PullRequestsService) Patch() (patch io.ReadCloser, result *Result) {
	return p.client.getBody(p.URL, patchMediaType)
}

// PullRequest represents a pull request on GitHub and all associated parameters
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
	Mergeable         *bool             `json:"mergeable,omitempty"`
	MergeableState    string            `json:"mergeable_state,omitempty"`
}

// PullRequestCommit represents one of the commits associated with the given
// pull request - the head being pulled in or the base being pulled into
type PullRequestCommit struct {
	Label string      `json:"label,omitempty"`
	Ref   string      `json:"ref,omitempty"`
	Sha   string      `json:"sha,omitempty"`
	User  User        `json:"user,omitempty"`
	Repo  *Repository `json:"repo,omitempty"`
}

// PullRequestParams represent the set of parameters used to create a new
// pull request
type PullRequestParams struct {
	Base     string `json:"base,omitempty"`
	Head     string `json:"head,omitempty"`
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	Assignee string `json:"assignee,omitempty"`
}

// PullRequestForIssueParams represent the set of parameters used to
// create a new pull request for a specific issue
type PullRequestForIssueParams struct {
	Base  string `json:"base,omitempty"`
	Head  string `json:"head,omitempty"`
	Issue string `json:"issue,omitempty"`
}
