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

type PullRequest struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
	IssueUrl string `json:"issue_url"`
}

func (c *Client) CreatePullRequest(repo Repository, params PullRequestParams) (*PullRequest, error) {
	path := fmt.Sprintf("repos/%s/%s/pulls", repo.UserName, repo.Name)
	body, err := c.postWithParams(path, nil, params)
	if err != nil {
		return nil, err
	}

	var pr PullRequest
	err = jsonUnmarshal(body, &pr)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}
