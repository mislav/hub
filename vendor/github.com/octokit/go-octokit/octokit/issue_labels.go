package octokit

// IssueLabelsURL is a URL template for accessing issue labels
//
// https://developer.github.com/v3/issues/labels/
var IssueLabelsURL = Hyperlink("repos/{owner}/{repo}/issues/{number}/labels{/name}")

// IssueLabels creates an IssueLabelsService with a base url
func (c *Client) IssueLabels() (issueLabels *IssueLabelsService) {
	issueLabels = &IssueLabelsService{client: c}
	return
}

// IssueLabelsService is a service providing access to labels for an issue
type IssueLabelsService struct {
	client *Client
}

// Adds labels to an issue
//
// https://developer.github.com/v3/issues/labels/#add-labels-to-an-issue
func (l *IssueLabelsService) Add(uri *Hyperlink, uriParams M, labelsToAdd []string) (labels []Label, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueLabelsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = l.client.post(url, labelsToAdd, &labels)
	return
}

// All gets a list of all labels for an issue
//
// https://developer.github.com/v3/issues/labels/#list-labels-on-an-issue
func (l *IssueLabelsService) All(uri *Hyperlink, uriParams M) (labels []Label, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueLabelsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = l.client.get(url, &labels)
	return
}

// Removes a label from an issue
//
// https://developer.github.com/v3/issues/labels/#remove-a-label-from-an-issue
func (l *IssueLabelsService) Remove(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueLabelsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = l.client.delete(url, nil, nil)

	success = (result.Response.StatusCode == 204)
	return
}

// Removes all labels from an issue
//
// https://developer.github.com/v3/issues/labels/#remove-all-labels-from-an-issue
func (l *IssueLabelsService) RemoveAll(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueLabelsURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = l.client.delete(url, nil, nil)

	success = (result.Response.StatusCode == 204)
	return
}

// Replace all labels for an issue
//
// https://developer.github.com/v3/issues/labels/#replace-all-labels-for-an-issue
func (l *IssueLabelsService) ReplaceAll(uri *Hyperlink, uriParams M, newLabels []string) (labels []Label, result *Result) {
	url, err := ExpandWithDefault(uri, &IssueLabelsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = l.client.put(url, newLabels, &labels)
	return
}
