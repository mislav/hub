package octokit

import (
	"github.com/jingweno/go-sawyer/hypermedia"
)

// RepoLabelsURL is a URL template for accessing labels for a repository.
//
// https://developer.github.com/v3/issues/labels/
var RepoLabelsURL = Hyperlink("repos/{owner}/{repo}/labels{/name}")

// Labels creates a LabelsService
func (c *Client) Labels() (labels *LabelsService) {
	labels = &LabelsService{client: c}
	return
}

// LabelsService is a service providing access to labels for a repository
type LabelsService struct {
	client *Client
}

// All gets a list of all labels for a repository
//
// https://developer.github.com/v3/issues/labels/#list-all-labels-for-this-repository
func (l *LabelsService) All(uri *Hyperlink, uriParams M) (labels []Label, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoLabelsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = l.client.get(url, &labels)
	return
}

// One gets a label for a repository
//
// https://developer.github.com/v3/issues/labels/#get-a-single-label
func (l *LabelsService) One(uri *Hyperlink, uriParams M) (label *Label, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoLabelsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = l.client.get(url, &label)
	return
}

// Create a new label for a repository
//
// https://developer.github.com/v3/issues/labels/#create-a-label
func (l *LabelsService) Create(uri *Hyperlink, uriParams M, requestParams interface{}) (label *Label, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoLabelsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = l.client.post(url, requestParams, &label)
	return
}

// Updates a label for a repository
//
// https://developer.github.com/v3/issues/labels/#update-a-label
func (l *LabelsService) Update(uri *Hyperlink, uriParams M, requestParams interface{}) (label *Label, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoLabelsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = l.client.patch(url, requestParams, &label)
	return
}

// Delete a label for a repository
//
// https://developer.github.com/v3/issues/labels/#delete-a-label
func (l *LabelsService) Delete(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &RepoLabelsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = l.client.delete(url, nil, nil)

	success = (result.Response.StatusCode == 204)
	return
}

// Label represents a label for a GitHub repository
type Label struct {
	*hypermedia.HALResource

	URL   string `json:"url,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}
