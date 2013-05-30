package github

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type PullRequestParams struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Base  string `json:"base"`
	Head  string `json:"head"`
}

type PullRequestResponse struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
	IssueUrl string `json:"issue_url"`
}

func createPullRequest(gh *GitHub, params PullRequestParams) (*PullRequestResponse, error) {
	project := gh.project
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(b)
	url := fmt.Sprintf("/repos/%s/%s/pulls", project.Owner, project.Name)
	response, err := httpPost(gh, url, nil, buffer)
	if err != nil {
		return nil, err
	}

	var pullRequestResponse PullRequestResponse
	err = unmarshalBody(response, &pullRequestResponse)
	if err != nil {
		return nil, err
	}

	return &pullRequestResponse, nil
}
