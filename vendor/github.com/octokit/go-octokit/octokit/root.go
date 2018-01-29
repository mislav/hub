package octokit

import (
	"net/url"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// RootURL is simply the root GitHub address. Accessing this address provides all
// other accessible templates and addresses as hypermedia relations.
//
// https://api.github.com/
var RootURL = Hyperlink("")

// Rel fetches and expands the given name in the the Hyperlink map m
func (c *Client) Rel(name string, m map[string]interface{}) (*url.URL, error) {
	if c.rootRels == nil || len(c.rootRels) == 0 {
		u, _ := url.Parse("/")
		root, res := c.Root(u).One()
		if res.HasError() {
			return nil, res
		}
		c.rootRels = root.Rels()
	}

	return c.rootRels.Rel(name, m)
}

// Root creates a RootService with a base url
//
// https://api.github.com/
func (c *Client) Root(url *url.URL) (root *RootService) {
	root = &RootService{client: c, URL: url}
	return
}

// RootService is a representation of a simple service to access hyperlinks
// to all the other accessible URLs
type RootService struct {
	client *Client
	URL    *url.URL
}

// One accesses the root URI templates and assigns them to result
//
// https://developer.github.com/v3/#root-endpoint
func (r *RootService) One() (root *Root, result *Result) {
	root = &Root{HALResource: &hypermedia.HALResource{}}
	result = r.client.get(r.URL, &root)
	if root != nil {
		// Cached hyperlinks
		root.PullsURL = hypermedia.Hyperlink(PullRequestsURL)
	}

	return
}

// Root represents the base with hyperlinks in template form to all API calls
//
// https://api.github.com/
type Root struct {
	*hypermedia.HALResource

	UserSearchURL               hypermedia.Hyperlink `rel:"user_search" json:"user_search_url,omitempty"`
	UserRepositoriesURL         hypermedia.Hyperlink `rel:"user_repositories" json:"user_repositories_url,omitempty"`
	UserOrganizationsURL        hypermedia.Hyperlink `rel:"user_organizations" json:"user_organizations_url,omitempty"`
	UserURL                     hypermedia.Hyperlink `rel:"user" json:"user_url,omitempty"`
	TeamURL                     hypermedia.Hyperlink `rel:"team" json:"team_url,omitempty"`
	StarredGistsURL             hypermedia.Hyperlink `rel:"starred_gists" json:"starred_gists_url,omitempty"`
	StarredURL                  hypermedia.Hyperlink `rel:"starred" json:"starred_url,omitempty"`
	CurrentUserRepositoriesURL  hypermedia.Hyperlink `rel:"current_user_repositories" json:"current_user_repositories_url,omitempty"`
	RepositorySearchURL         hypermedia.Hyperlink `rel:"repository_search" json:"repository_search_url,omitempty"`
	RepositoryURL               hypermedia.Hyperlink `rel:"repository" json:"repository_url,omitempty"`
	RateLimitURL                hypermedia.Hyperlink `rel:"rate_limit" json:"rate_limit_url,omitempty"`
	GistsURL                    hypermedia.Hyperlink `rel:"gists" json:"gists_url,omitempty"`
	FollowingURL                hypermedia.Hyperlink `rel:"following" json:"following_url,omitempty"`
	FeedsURL                    hypermedia.Hyperlink `rel:"feeds" json:"feeds_url,omitempty"`
	EventsURL                   hypermedia.Hyperlink `rel:"events" json:"events_url,omitempty"`
	EmojisURL                   hypermedia.Hyperlink `rel:"emojis" json:"emojis_url,omitempty"`
	EmailsURL                   hypermedia.Hyperlink `rel:"emails" json:"emails_url,omitempty"`
	AuthorizationsURL           hypermedia.Hyperlink `rel:"authorizations" json:"authorizations_url,omitempty"`
	CurrentUserURL              hypermedia.Hyperlink `rel:"current_user" json:"current_user_url,omitempty"`
	HubURL                      hypermedia.Hyperlink `rel:"hub" json:"hub_url,omitempty"`
	IssueSearchURL              hypermedia.Hyperlink `rel:"issue_search" json:"issue_search_url,omitempty"`
	IssuesURL                   hypermedia.Hyperlink `rel:"issues" json:"issues_url,omitempty"`
	KeysURL                     hypermedia.Hyperlink `rel:"keys" json:"keys_url,omitempty"`
	NotificationsURL            hypermedia.Hyperlink `rel:"notifications" json:"notifications_url,omitempty"`
	OrganizationRepositoriesURL hypermedia.Hyperlink `rel:"organization_repositories" json:"organization_repositories_url,omitempty"`
	OrganizationURL             hypermedia.Hyperlink `rel:"organization" json:"organization_url,omitempty"`
	PublicGistsURL              hypermedia.Hyperlink `rel:"public_gists" json:"public_gists_url,omitempty"`
	PullsURL                    hypermedia.Hyperlink `rel:"pulls" json:"-"`
	rels                        hypermedia.Relations `json:"-"`
}

// Rels gets the link relations from the HALResource's Links field.
func (r *Root) Rels() hypermedia.Relations {
	if r.rels == nil || len(r.rels) == 0 {
		r.rels = hypermedia.HyperFieldDecoder(r)
		for key, hyperlink := range r.HALResource.Rels() {
			r.rels[key] = hyperlink
		}
	}
	return r.rels
}
