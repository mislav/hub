package octokit

// Organization is a representation of an organization on GitHub, containing
// all identifying information related to the specific organization.

import (
	"time"
)

var (
	OrganizationURL      = Hyperlink("/orgs/{org}")
	OrganizationReposURL = Hyperlink("/orgs/{org}/repos{?type,page,per_page,sort}")
	YourOrganizationsURL = Hyperlink("/user/orgs")
	UserOrganizationsURL = Hyperlink("/users/{username}/orgs")
)

func (c *Client) Organization() (organization *OrganizationService) {
	organization = &OrganizationService{client: c}
	return
}

// A service for getting as well as updating organization information
type OrganizationService struct {
	client *Client
}

// Get the specified organization's information
//
// https://developer.github.com/v3/orgs/#get-an-organization
func (g *OrganizationService) OrganizationGet(uri *Hyperlink, params M) (
	organization Organization, result *Result) {
	if uri == nil {
		uri = &OrganizationURL
	}
	url, err := uri.Expand(params)
	if err != nil {
		return Organization{}, &Result{Err: err}
	}
	result = g.client.get(url, &organization)
	return
}

// Update specified organization's information
//
// https://developer.github.com/v3/orgs/#edit-an-organization
func (g *OrganizationService) OrganizationUpdate(uri *Hyperlink,
	input OrganizationParams, URLParams M) (organization Organization,
	result *Result) {
	if uri == nil {
		uri = &OrganizationURL
	}
	url, err := uri.Expand(URLParams)
	if err != nil {
		return Organization{}, &Result{Err: err}
	}
	result = g.client.patch(url, input, &organization)
	return
}

// Get the list of repository information of an organization
//
// https://developer.github.com/v3/repos/#list-organization-repositories
func (g *OrganizationService) OrganizationRepos(uri *Hyperlink, params M) (
	repos []Repository, result *Result) {
	if uri == nil {
		uri = &OrganizationReposURL
	}
	url, err := uri.Expand(params)
	if err != nil {
		return nil, &Result{Err: err}
	}
	result = g.client.get(url, &repos)
	return
}

// Get information for the list of organizations the current user belongs to
//
// https://developer.github.com/v3/orgs/#list-your-organizations
func (g *OrganizationService) YourOrganizations(uri *Hyperlink, params M) (
	organizations []Organization, result *Result) {
	if uri == nil {
		uri = &YourOrganizationsURL
	}
	url, err := uri.Expand(params)
	if err != nil {
		return nil, &Result{Err: err}
	}
	result = g.client.get(url, &organizations)
	return
}

// Get the information for the list of organizations the specified user belongs to
//
// https://developer.github.com/v3/orgs/#list-user-organizations
func (g *OrganizationService) UserOrganizations(uri *Hyperlink, params M) (
	organizations []Organization, result *Result) {
	if uri == nil {
		uri = &UserOrganizationsURL
	}
	url, err := uri.Expand(params)
	if err != nil {
		return nil, &Result{Err: err}
	}
	result = g.client.get(url, &organizations)
	return
}

type Organization struct {
	Description      string     `json:"description, omitempty"`
	AvatarURL        string     `json:"avatar_url,omitempty"`
	PublicMembersURL Hyperlink  `json:"public_member_url,omitempty"`
	MembersURL       Hyperlink  `json:"members_url,omitempty"`
	EventsURL        Hyperlink  `json:"events_url,omitempty"`
	ReposURL         Hyperlink  `json:"repos_url,omitempty"`
	URL              string     `json:"url,omitempty"`
	ID               int        `json:"id,omitempty"`
	Login            string     `json:"login,omitempty"`
	Name             string     `json:"name, omitempty"`
	Company          string     `json:"company, omitempty"`
	Blog             string     `json:"blog, omitempty"`
	Location         string     `json:"location, omitempty"`
	Email            string     `json:"email, omitempty"`
	PublicRepos      int        `json:"public_repos,omitempty"`
	PublicGists      int        `json:"public_gists,omitempty"`
	Followers        int        `json:"followers,omitempty"`
	Followering      int        `json:"following,omitempty"`
	HTMLURL          string     `json:"html_url,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	Type             string     `json:"type,omitempty"`

	//Limited Access Fields
	Total_Private_Repos int    `json:"total_private_repos,omitempty"`
	Owned_Private_Repos int    `json:"owned_private_repos,omitempty"`
	Private_Gists       int    `json:"private_gists,omitempty"`
	Disk_Usage          int    `json:"disk_usage,omitempty"`
	Collaborators       int    `json:"collaborators,omitempty"`
	BillingEmail        string `json:"billing_email,omitempty"`
	Plan                Plan   `json:"plan,omitempty"`
}

type Plan struct {
	Name         string `json:"name,omitempty"`
	Space        int    `json:"space,omitempty"`
	PrivateRepos int    `json:"private_repos,omitempty"`
}

// OrganizationParams represents the struture used to create or update an Organization
type OrganizationParams struct {
	BillingEmail string `json:"billing_email,omitempty"`
	Blog         string `json:"blog,omitempty"`
	Company      string `json:"company,omitempty"`
	Email        string `json:"email,omitempty"`
	Location     string `json:"location,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
}
