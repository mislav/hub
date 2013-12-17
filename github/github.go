package github

import (
	"fmt"
	"github.com/jingweno/go-octokit/octokit"
	"net/url"
	"os"
)

const (
	GitHubHost    string = "github.com"
	GitHubApiHost string = "api.github.com"
	OAuthAppURL   string = "http://owenou.com/gh"
)

type GitHub struct {
	Credentials *Credentials
}

func (gh *GitHub) PullRequest(project *Project, id string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "number": id})
	if err != nil {
		return
	}

	client := gh.octokit()
	pr, result := client.PullRequests(gh.requestURL(url)).One()
	if result.HasError() {
		err = result.Err
	}

	return
}

func (gh *GitHub) CreatePullRequest(project *Project, base, head, title, body string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	client := gh.octokit()
	params := octokit.PullRequestParams{Base: base, Head: head, Title: title, Body: body}
	pr, result := client.PullRequests(gh.requestURL(url)).Create(params)
	if result.HasError() {
		err = result.Err
	}

	return
}

func (gh *GitHub) CreatePullRequestForIssue(project *Project, base, head, issue string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	client := gh.octokit()
	params := octokit.PullRequestForIssueParams{Base: base, Head: head, Issue: issue}
	pr, result := client.PullRequests(gh.requestURL(url)).Create(params)
	if result.HasError() {
		err = result.Err
	}

	return
}

func (gh *GitHub) Repository(project *Project) (repo *octokit.Repository, err error) {
	url, err := octokit.RepositoryURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	client := gh.octokit()
	repo, result := client.Repositories(gh.requestURL(url)).One()
	if result.HasError() {
		err = result.Err
	}

	return
}

func (gh *GitHub) IsRepositoryExist(project *Project) bool {
	repo, err := gh.Repository(project)

	return err == nil && repo != nil
}

func (gh *GitHub) CreateRepository(project *Project, description, homepage string, isPrivate bool) (repo *octokit.Repository, err error) {
	var repoURL octokit.Hyperlink
	if project.Owner != gh.Credentials.User {
		repoURL = octokit.OrgRepositoriesURL
	} else {
		repoURL = octokit.UserRepositoriesURL
	}

	url, err := repoURL.Expand(octokit.M{"org": project.Owner})
	if err != nil {
		return
	}

	client := gh.octokit()
	params := octokit.Repository{
		Name:        project.Name,
		Description: description,
		Homepage:    homepage,
		Private:     isPrivate,
	}
	repo, result := client.Repositories(gh.requestURL(url)).Create(params)
	if result.HasError() {
		err = result.Err
	}

	return
}

func (gh *GitHub) Releases(project *Project) (releases []octokit.Release, err error) {
	url, err := octokit.ReleasesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	client := gh.octokit()
	releases, result := client.Releases(gh.requestURL(url)).All()
	if result.HasError() {
		err = result.Err
		return
	}

	return
}

func (gh *GitHub) CIStatus(project *Project, sha string) (status *octokit.Status, err error) {
	url, err := octokit.StatusesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "ref": sha})
	if err != nil {
		return
	}

	client := gh.octokit()
	statuses, result := client.Statuses(gh.requestURL(url)).All()
	if result.HasError() {
		err = result.Err
		return
	}

	if len(statuses) > 0 {
		status = &statuses[0]
	}

	return
}

func (gh *GitHub) ForkRepository(project *Project) (repo *octokit.Repository, err error) {
	url, err := octokit.ForksURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	client := gh.octokit()
	repo, result := client.Repositories(gh.requestURL(url)).Create(nil)
	if result.HasError() {
		err = result.Err
	}

	return
}

func (gh *GitHub) Issues(project *Project) (issues []octokit.Issue, err error) {
	url, err := octokit.RepoIssuesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	client := gh.octokit()
	issues, result := client.Issues(gh.requestURL(url)).All()
	if result.HasError() {
		err = result.Err
		return
	}

	return
}

func (gh *GitHub) FindOrCreateToken(user, password, twoFactorCode string) (token string, err error) {
	url, err := octokit.AuthorizationsURL.Expand(nil)
	if err != nil {
		return
	}

	basicAuth := octokit.BasicAuth{Login: user, Password: password, OneTimePassword: twoFactorCode}
	client := octokit.NewClientWith(gh.apiEndpoint(), nil, basicAuth)
	authsService := client.Authorizations(gh.requestURL(url))

	auths, result := authsService.All()
	if result.HasError() {
		err = result.Err
		return
	}

	for _, auth := range auths {
		if auth.NoteURL == OAuthAppURL {
			token = auth.Token
			break
		}
	}

	if token == "" {
		authParam := octokit.AuthorizationParams{}
		authParam.Scopes = append(authParam.Scopes, "repo")
		authParam.Note = "gh"
		authParam.NoteURL = OAuthAppURL

		auth, result := authsService.Create(authParam)
		if result.HasError() {
			err = result.Err
			return
		}

		token = auth.Token
	}

	return
}

func (gh *GitHub) octokit() (c *octokit.Client) {
	tokenAuth := octokit.TokenAuth{AccessToken: gh.Credentials.AccessToken}
	c = octokit.NewClientWith(gh.apiEndpoint(), nil, tokenAuth)

	return
}

func (gh *GitHub) requestURL(u *url.URL) (uu *url.URL) {
	uu = u
	if gh.Credentials.Host != GitHubHost {
		uu, _ = url.Parse(fmt.Sprintf("/api/v3/%s", u.Path))
	}

	return
}

func (gh *GitHub) apiEndpoint() string {
	host := os.Getenv("GH_API_HOST")
	if host == "" {
		host = gh.Credentials.Host
	}

	if host == GitHubHost {
		host = GitHubApiHost
	}

	return absolute(host)
}

func absolute(endpoint string) string {
	u, _ := url.Parse(endpoint)
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	return u.String()
}

func NewClient(host string) *GitHub {
	c := CurrentConfigs().PromptFor(host)
	return &GitHub{Credentials: c}
}
