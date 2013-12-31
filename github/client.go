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

type ClientError struct {
	error
}

func (e *ClientError) Error() string {
	return e.error.Error()
}

func (e *ClientError) Is2FAError() bool {
	re, ok := e.error.(*octokit.ResponseError)
	return ok && re.Type == octokit.ErrorOneTimePasswordRequired
}

type Client struct {
	Credentials *Credentials
}

func (client *Client) PullRequest(project *Project, id string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "number": id})
	if err != nil {
		return
	}

	pr, result := client.octokit().PullRequests(client.requestURL(url)).One()
	if result.HasError() {
		err = fmt.Errorf("Error getting pull request: %s", result.Err)
	}

	return
}

func (client *Client) CreatePullRequest(project *Project, base, head, title, body string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	params := octokit.PullRequestParams{Base: base, Head: head, Title: title, Body: body}
	pr, result := client.octokit().PullRequests(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = fmt.Errorf("Error creating pull request: %s", result.Err)
	}

	return
}

func (client *Client) CreatePullRequestForIssue(project *Project, base, head, issue string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	params := octokit.PullRequestForIssueParams{Base: base, Head: head, Issue: issue}
	pr, result := client.octokit().PullRequests(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = fmt.Errorf("Error creating pull request: %s", result.Err)
	}

	return
}

func (client *Client) Repository(project *Project) (repo *octokit.Repository, err error) {
	url, err := octokit.RepositoryURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	repo, result := client.octokit().Repositories(client.requestURL(url)).One()
	if result.HasError() {
		err = fmt.Errorf("Error getting repository: %s", result.Err)
	}

	return
}

func (client *Client) IsRepositoryExist(project *Project) bool {
	repo, err := client.Repository(project)

	return err == nil && repo != nil
}

func (client *Client) CreateRepository(project *Project, description, homepage string, isPrivate bool) (repo *octokit.Repository, err error) {
	var repoURL octokit.Hyperlink
	if project.Owner != client.Credentials.User {
		repoURL = octokit.OrgRepositoriesURL
	} else {
		repoURL = octokit.UserRepositoriesURL
	}

	url, err := repoURL.Expand(octokit.M{"org": project.Owner})
	if err != nil {
		return
	}

	params := octokit.Repository{
		Name:        project.Name,
		Description: description,
		Homepage:    homepage,
		Private:     isPrivate,
	}
	repo, result := client.octokit().Repositories(client.requestURL(url)).Create(params)
	if result.HasError() {
		if result.Response == nil || result.Response.StatusCode == 500 {
			err = fmt.Errorf("Error creating repository: Internal Server Error (HTTP 500)")
		} else {
			err = fmt.Errorf("Error creating repository: %v", result.Err)
		}
	}

	return
}

func (client *Client) Releases(project *Project) (releases []octokit.Release, err error) {
	url, err := octokit.ReleasesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	releases, result := client.octokit().Releases(client.requestURL(url)).All()
	if result.HasError() {
		err = fmt.Errorf("Error getting release: %s", result.Err)
	}

	return
}

func (client *Client) CreateRelease(project *Project, params octokit.ReleaseParams) (release *octokit.Release, err error) {
	url, err := octokit.ReleasesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	release, result := client.octokit().Releases(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = fmt.Errorf("Error creating release: %s", result.Err)
	}

	return
}

func (client *Client) UploadReleaseAsset(uploadUrl *url.URL, asset *os.File, contentType string) (err error) {
	c := client.octokit()
	result := c.Uploads(uploadUrl).UploadAsset(asset, contentType)
	if result.HasError() {
		err = fmt.Errorf("Error uploading asset: %s", result.Err)
	}
	return
}

func (client *Client) CIStatus(project *Project, sha string) (status *octokit.Status, err error) {
	url, err := octokit.StatusesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "ref": sha})
	if err != nil {
		return
	}

	statuses, result := client.octokit().Statuses(client.requestURL(url)).All()
	if result.HasError() {
		err = fmt.Errorf("Error getting CI status: %s", result.Err)
		return
	}

	if len(statuses) > 0 {
		status = &statuses[0]
	}

	return
}

func (client *Client) ForkRepository(project *Project) (repo *octokit.Repository, err error) {
	url, err := octokit.ForksURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	repo, result := client.octokit().Repositories(client.requestURL(url)).Create(nil)
	if result.HasError() {
		err = fmt.Errorf("Error forking repository: %s", result.Err)
	}

	return
}

func (client *Client) Issues(project *Project) (issues []octokit.Issue, err error) {
	var result *octokit.Result

	err = client.issuesService(project, func(service *octokit.IssuesService) error {
		issues, result = service.All()
		return resultError(result)
	})

	return
}

func (client *Client) CreateIssue(project *Project, title, body string, labels []string) (issue *octokit.Issue, err error) {
	params := octokit.IssueParams{
		Title:  title,
		Body:   body,
		Labels: labels,
	}

	var result *octokit.Result

	err = client.issuesService(project, func(service *octokit.IssuesService) error {
		issue, result = service.Create(params)
		return resultError(result)
	})

	return
}

func (client *Client) FindOrCreateToken(user, password, twoFactorCode string) (token string, err error) {
	url, e := octokit.AuthorizationsURL.Expand(nil)
	if e != nil {
		err = &ClientError{e}
		return
	}

	basicAuth := octokit.BasicAuth{Login: user, Password: password, OneTimePassword: twoFactorCode}
	c := octokit.NewClientWith(client.apiEndpoint(), nil, basicAuth)
	authsService := c.Authorizations(client.requestURL(url))

	auths, result := authsService.All()
	if result.HasError() {
		err = &ClientError{result.Err}
		return
	}

	for _, auth := range auths {
		if auth.App.URL == OAuthAppURL {
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
			err = &ClientError{result.Err}
			return
		}

		token = auth.Token
	}

	return
}

func (client *Client) octokit() (c *octokit.Client) {
	tokenAuth := octokit.TokenAuth{AccessToken: client.Credentials.AccessToken}
	c = octokit.NewClientWith(client.apiEndpoint(), nil, tokenAuth)

	return
}

func (client *Client) requestURL(u *url.URL) (uu *url.URL) {
	uu = u
	if client.Credentials.Host != GitHubHost {
		uu, _ = url.Parse(fmt.Sprintf("/api/v3/%s", u.Path))
	}

	return
}

func (client *Client) apiEndpoint() string {
	host := os.Getenv("GH_API_HOST")
	if host == "" {
		host = client.Credentials.Host
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

func NewClient(host string) *Client {
	c := CurrentConfigs().PromptFor(host)
	return &Client{Credentials: c}
}

func (client *Client) issuesService(project *Project, fn func(service *octokit.IssuesService) error) (err error) {
	url, err := octokit.RepoIssuesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	service := client.octokit().Issues(client.requestURL(url))
	return fn(service)
}

func resultError(result *octokit.Result) (err error) {
	if result != nil && result.HasError() {
		err = result.Err
	}
	return
}
