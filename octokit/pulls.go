package octokit

import (
	"fmt"
)

type PullRequestParams struct {
	Base  string `json:"base"`
	Head  string `json:"head"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PullRequestForIssueParams struct {
	Base  string `json:"base"`
	Head  string `json:"head"`
	Issue string `json:"issue"`
}

type PullRequest struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
	IssueUrl string `json:"issue_url"`
}

func (c *Client) CreatePullRequest(repo Repository, params PullRequestParams) (*PullRequest, error) {
	return c.createPullRequest(repo, params)
}

func (c *Client) CreatePullRequestForIssue(repo Repository, params PullRequestForIssueParams) (*PullRequest, error) {
	return c.createPullRequest(repo, params)
}

func (c *Client) createPullRequest(repo Repository, params interface{}) (*PullRequest, error) {
	path := fmt.Sprintf("repos/%s/pulls", repo)
	var pr PullRequest
	err := c.jsonPost(path, nil, params, &pr)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}
