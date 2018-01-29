package octokit

// Team is a representation of a team on GitHub, containing
// all identifying information related to the specific team.

import "github.com/jingweno/go-sawyer/hypermedia"

// Hyperlinks to the various team locations on github.
// OrganizationTeamsURL is the template for teams within a particular organization.
// TeamURL is a template for a particular team.
// TeamMembersURL is a template for members on a particular team.
//
// https://developer.github.com/v3/repos/
var (
	// Unlike TeamReposURL, the docs for the teams API do _not_ list `page` or
	// `per_page` querystring params.
	OrganizationTeamsURL = Hyperlink("/orgs/{org}/teams")
	TeamURL              = Hyperlink("/teams/{id}")
	TeamMembersURL       = Hyperlink("/teams/{id}/members")
	TeamMembershipURL    = Hyperlink("/teams/{id}/memberships/{username}")
	TeamRepositoriesURL  = Hyperlink("/teams/{id}/repos")
	TeamRepositoryURL    = Hyperlink("/teams/{id}/repos/{owner}/{repo}")
	CurrentUserTeams     = Hyperlink("/user/teams") // pass to GetTeams()
)

// Teams returns a TeamsService with a base url
//
// https://developer.github.com/v3/orgs/teams/
func (c *Client) Teams() (team *TeamsService) {
	team = &TeamsService{client: c}
	return
}

// TeamsService for getting as well as updating team information
type TeamsService struct {
	client *Client
}

// One returns a single Team for a given URL.
func (t *TeamsService) One(uri *Hyperlink, uriParams M) (
	team Team, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamURL, uriParams)
	if err != nil {
		return Team{}, &Result{Err: err}
	}
	result = t.client.get(url, &team)
	return
}

// All returns a slice of Teams for a given URL.
func (t *TeamsService) All(uri *Hyperlink, uriParams M) (
	teams []Team, result *Result) {
	url, err := ExpandWithDefault(uri, &OrganizationTeamsURL, uriParams)
	if err != nil {
		return []Team(nil), &Result{Err: err}
	}
	result = t.client.get(url, &teams)
	return
}

// Create a team
//
// https://developer.github.com/v3/orgs/teams/#create-team
func (t *TeamsService) Create(uri *Hyperlink, input TeamParams,
	uriParams M) (team Team, result *Result) {
	url, err := ExpandWithDefault(uri, &OrganizationTeamsURL, uriParams)
	if err != nil {
		return Team{}, &Result{Err: err}
	}
	result = t.client.post(url, input, &team)
	return
}

// Update specified team's information
//
// https://developer.github.com/v3/orgs/teams/#edit-team
func (t *TeamsService) Update(uri *Hyperlink, input TeamParams,
	uriParams M) (team Team, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamURL, uriParams)
	if err != nil {
		return Team{}, &Result{Err: err}
	}
	result = t.client.patch(url, input, &team)
	return
}

// Delete specified teams
//
// https://developer.github.com/v3/orgs/teams/#delete-team
func (t *TeamsService) Delete(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}
	result = t.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

// GetTeams gets the teams that belong to a particular organization
//
// https://developer.github.com/v3/orgs/teams/#list-teams
func (o *OrganizationService) GetTeams(uri *Hyperlink, uriParams M) (
	teams []Team, result *Result) {
	return o.client.Teams().All(&OrganizationTeamsURL, uriParams)
}

// Get the specified team's information
//
// https://developer.github.com/v3/orgs/teams/#get-team
func (t *TeamsService) Get(uri *Hyperlink, uriParams M) (
	team Team, result *Result) {
	return t.One(&TeamURL, uriParams)
}

// GetMembers reutrns the members (Users) in the specified team
//
// https://developer.github.com/v3/orgs/teams/#list-team-members
func (t *TeamsService) GetMembers(uri *Hyperlink, uriParams M) (
	members []User, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamMembersURL, uriParams)
	if err != nil {
		return []User(nil), &Result{Err: err}
	}
	return t.client.Users(url).All()
}

// (deprecated) https://developer.github.com/v3/orgs/teams/#get-team-member
// (deprecated) https://developer.github.com/v3/orgs/teams/#add-team-member
// (deprecated) https://developer.github.com/v3/orgs/teams/#remove-team-member

// GetMembership returns a user's membership details for a given team
//
// https://developer.github.com/v3/orgs/teams/#get-team-membership
func (t *TeamsService) GetMembership(uri *Hyperlink, uriParams M) (
	membership TeamMembership, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamMembershipURL, uriParams)
	if err != nil {
		return TeamMembership{}, &Result{Err: err}
	}
	result = t.client.get(url, &membership)
	return
}

// AddMembership adds an organization user to a given team with a specified role
//
// https://developer.github.com/v3/orgs/teams/#add-team-membership
func (t *TeamsService) AddMembership(uri *Hyperlink, uriParams M, role string) (
	membership TeamMembership, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamMembershipURL, uriParams)
	if err != nil {
		return TeamMembership{}, &Result{Err: err}
	}
	result = t.client.put(url, M{"role": role}, &membership)
	return
}

// RemoveMembership removes a user from a given Team.
//
// https://developer.github.com/v3/orgs/teams/#remove-team-membership
func (t *TeamsService) RemoveMembership(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamMembershipURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}
	result = t.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

// GetRepositories returns the repos a team has access to.
//
// https://developer.github.com/v3/orgs/teams/#list-team-repos
func (t *TeamsService) GetRepositories(uri *Hyperlink, uriParams M) (
	repos []Repository, result *Result) {
	if uri == nil {
		uri = &TeamRepositoriesURL
	}
	return t.client.Repositories().All(uri, uriParams)
}

// CheckRepository returns whether a team manages a repository.
//
// https://developer.github.com/v3/orgs/teams/#check-if-a-team-manages-a-repository
func (t *TeamsService) CheckRepository(uri *Hyperlink, uriParams M) (manages bool, repo *Repository, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamRepositoryURL, uriParams)
	if err != nil {
		result = &Result{Err: err}
		return
	}
	result = t.client.get(url, &repo)
	manages = (result.Response.StatusCode == 204 || result.Response.StatusCode == 200)
	return
}

// UpdateRepository adds (or updates) an organization's repository with the
// given team.
//
// https://developer.github.com/v3/orgs/teams/#add-or-update-team-repository
func (t *TeamsService) UpdateRepository(uri *Hyperlink, uriParams M, permission string) (
	success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamRepositoryURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}
	result = t.client.put(url, M{"permission": permission}, nil)
	success = (result.Response.StatusCode == 204)
	return
}

// RemoveRepository removes a team's access to an organization's repository.
//
// https://developer.github.com/v3/orgs/teams/#remove-team-repository
func (t *TeamsService) RemoveRepository(uri *Hyperlink, uriParams M) (
	success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &TeamRepositoryURL, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}
	result = t.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

// Team is a representation of a team on GitHub, containing all identifying
// information related to the specific team.
type Team struct {
	*hypermedia.HALResource

	Name            string    `json:"name, omitempty"`
	ID              int       `json:"id,omitempty"`
	Slug            string    `json:"slug, omitempty"`
	Description     string    `json:"description, omitempty"`
	Permission      string    `json:"permission, omitempty"`
	Privacy         string    `json:"privacy, omitempty"`
	URL             string    `json:"url,omitempty"`
	MembersURL      Hyperlink `json:"members_url,omitempty"`
	RepositoriesURL Hyperlink `json:"repositories_url,omitempty"`

	// When individuallty fetched, created, or updated...
	MembersCount int           `json:"members_count,omitempty"`
	ReposCount   int           `json:"repos_count,omitempty"`
	Organization *Organization `json:"organization,omitempty"`
}

// TeamParams represents the struture used to create or update a Team
type TeamParams struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	RepoNames   []string `json:"repo_names,omitempty"`
	Privacy     string   `json:"privacy, omitempty"`
	Permission  string   `json:"permission, omitempty"` // deprecated
}

// TeamMembership represents the membership status of a user on a Team.
type TeamMembership struct {
	URL   string `json:"url,omitempty"`
	Role  string `json:"role,omitempty"`
	State string `json:"state, omitempty"`
}
