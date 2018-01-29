package octokit

import (
	"github.com/jingweno/go-sawyer/hypermedia"
	"time"
)

// MilestonesURL is a URL template for accessing repository Milestones
//
// https://developer.github.com/v3/issues/milestones/
var MilestonesURL = Hyperlink("repos/{owner}/{repo}/milestones")

// Milestones creates an MilestonesService with a base url
func (c *Client) Milestones() (milestones *MilestonesService) {
	milestones = &MilestonesService{client: c}
	return
}

// MilestonesService is a service providing access to milestones for a repository
type MilestonesService struct {
	client *Client
}

// All gets a list of all milestones for a repository
//
// https://developer.github.com/v3/issues/milestones/#list-milestones-for-a-repository
func (m *MilestonesService) All(uri *Hyperlink, uriParams M) (milestones []Milestone, result *Result) {
	url, err := ExpandWithDefault(uri, &MilestonesURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = m.client.get(url, &milestones)
	return
}

// One gets a milestone for a repository
//
// https://developer.github.com/v3/issues/milestones/#get-a-single-milestone
func (l *MilestonesService) One(uri *Hyperlink, uriParams M) (milestone *Milestone, result *Result) {
	url, err := ExpandWithDefault(uri, &MilestonesURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = l.client.get(url, &milestone)
	return
}

// Creates a milestone on an repository
//
// https://developer.github.com/v3/issues/milestones/#create-a-milestone
func (m *MilestonesService) Create(uri *Hyperlink, uriParams M, requestParams interface{}) (milestone *Milestone, result *Result) {
	url, err := ExpandWithDefault(uri, &MilestonesURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = m.client.post(url, requestParams, &milestone)
	return
}

// Deletes a milestone from a repository
//
// https://developer.github.com/v3/issues/milestones/#delete-a-milestone
func (m *MilestonesService) Delete(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &MilestonesURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = m.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

type Milestone struct {
	*hypermedia.HALResource

	URL          string     `json:"url,omitempty"`
	HTMLURL      string     `json:"html_url,omitempty"`
	LabelsURL    string     `json:"labels_url",omitempty"`
	Number       int        `json:"number,omitempty"`
	ID           int        `json:"id,omitempty"`
	State        string     `json:"state,omitempty"`
	Title        string     `json:"title,omitempty"`
	Description  string     `json:"description,omitempty"`
	Creator      User       `json:"creator,omitempty"`
	OpenIssues   int        `json:"open_issues,omitempty"`
	ClosedIssues int        `json:"closed_issues,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	DueOn        *time.Time `json:"due_on,omitempty"`
}
