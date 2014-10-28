package github

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/octokit/go-octokit/octokit"
)

const (
	GitHubHost    string = "github.com"
	GitHubApiHost string = "api.github.com"
	UserAgent     string = "Hub"
	OAuthAppName  string = "hub"
	OAuthAppURL   string = "http://hub.github.com/"
)

func NewClient(h string) *Client {
	return NewClientWithHost(&Host{Host: h})
}

func NewClientWithHost(host *Host) *Client {
	return &Client{host}
}

type AuthError struct {
	error
}

func (e *AuthError) Error() string {
	return e.error.Error()
}

func (e *AuthError) Is2FAError() bool {
	re, ok := e.error.(*octokit.ResponseError)
	return ok && re.Type == octokit.ErrorOneTimePasswordRequired
}

type Client struct {
	Host *Host
}

func (client *Client) PullRequest(project *Project, id string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "number": id})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting pull request", err)
		return
	}

	pr, result := api.PullRequests(client.requestURL(url)).One()
	if result.HasError() {
		err = FormatError("getting pull request", result.Err)
		return
	}

	return
}

func (client *Client) PullRequestPatch(project *Project, id string) (patch io.ReadCloser, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "number": id})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting pull request", err)
		return
	}

	patch, result := api.PullRequests(client.requestURL(url)).Patch()
	if result.HasError() {
		err = FormatError("getting pull request", result.Err)
		return
	}

	return
}

func (client *Client) CreatePullRequest(project *Project, base, head, title, body string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("creating pull request", err)
		return
	}

	params := octokit.PullRequestParams{Base: base, Head: head, Title: title, Body: body}
	pr, result := api.PullRequests(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = FormatError("creating pull request", result.Err)
		if e := warnExistenceOfRepo(project, result.Err); e != nil {
			err = fmt.Errorf("%s\n%s", err, e)
		}

		return
	}

	return
}

func (client *Client) CreatePullRequestForIssue(project *Project, base, head, issue string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("creating pull request", err)
		return
	}

	params := octokit.PullRequestForIssueParams{Base: base, Head: head, Issue: issue}
	pr, result := api.PullRequests(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = FormatError("creating pull request", result.Err)
		if e := warnExistenceOfRepo(project, result.Err); e != nil {
			err = fmt.Errorf("%s\n%s", err, e)
		}

		return
	}

	return
}

func (client *Client) CommitPatch(project *Project, sha string) (patch io.ReadCloser, err error) {
	url, err := octokit.CommitsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "sha": sha})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting pull request", err)
		return
	}

	patch, result := api.Commits(client.requestURL(url)).Patch()
	if result.HasError() {
		err = FormatError("getting pull request", result.Err)
		return
	}

	return
}

func (client *Client) GistPatch(id string) (patch io.ReadCloser, err error) {
	url, err := octokit.GistsURL.Expand(octokit.M{"gist_id": id})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting pull request", err)
		return
	}

	patch, result := api.Gists(client.requestURL(url)).Raw()
	if result.HasError() {
		err = FormatError("getting pull request", result.Err)
		return
	}

	return
}

func (client *Client) Repository(project *Project) (repo *octokit.Repository, err error) {
	url, err := octokit.RepositoryURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting repository", err)
		return
	}

	repo, result := api.Repositories(client.requestURL(url)).One()
	if result.HasError() {
		err = FormatError("getting repository", result.Err)
		return
	}

	return
}

func (client *Client) IsRepositoryExist(project *Project) bool {
	repo, err := client.Repository(project)

	return err == nil && repo != nil
}

func (client *Client) CreateRepository(project *Project, description, homepage string, isPrivate bool) (repo *octokit.Repository, err error) {
	var repoURL octokit.Hyperlink
	if project.Owner != client.Host.User {
		repoURL = octokit.OrgRepositoriesURL
	} else {
		repoURL = octokit.UserRepositoriesURL
	}

	url, err := repoURL.Expand(octokit.M{"org": project.Owner})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("creating repository", err)
		return
	}

	params := octokit.Repository{
		Name:        project.Name,
		Description: description,
		Homepage:    homepage,
		Private:     isPrivate,
	}
	repo, result := api.Repositories(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = FormatError("creating repository", result.Err)
		return
	}

	return
}

func (client *Client) Releases(project *Project) (releases []octokit.Release, err error) {
	url, err := octokit.ReleasesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting release", err)
		return
	}

	releases, result := api.Releases(client.requestURL(url)).All()
	if result.HasError() {
		err = FormatError("getting release", result.Err)
		return
	}

	return
}

func (client *Client) CreateRelease(project *Project, params octokit.ReleaseParams) (release *octokit.Release, err error) {
	url, err := octokit.ReleasesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("creating release", err)
		return
	}

	release, result := api.Releases(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = FormatError("creating release", result.Err)
		return
	}

	return
}

func (client *Client) UploadReleaseAsset(uploadUrl *url.URL, asset *os.File, contentType string) (err error) {
	fileInfo, err := asset.Stat()
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("uploading asset", err)
		return
	}

	result := api.Uploads(uploadUrl).UploadAsset(asset, contentType, fileInfo.Size())
	if result.HasError() {
		err = FormatError("uploading asset", result.Err)
		return
	}

	return
}

func (client *Client) CIStatus(project *Project, sha string) (status *octokit.Status, err error) {
	url, err := octokit.StatusesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "ref": sha})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting CI status", err)
		return
	}

	statuses, result := api.Statuses(client.requestURL(url)).All()
	if result.HasError() {
		err = FormatError("getting CI status", result.Err)
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

	api, err := client.api()
	if err != nil {
		err = FormatError("creating fork", err)
		return
	}

	repo, result := api.Repositories(client.requestURL(url)).Create(nil)
	if result.HasError() {
		err = FormatError("creating fork", result.Err)
		return
	}

	return
}

func (client *Client) Issues(project *Project) (issues []octokit.Issue, err error) {
	url, err := octokit.RepoIssuesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting issues", err)
		return
	}

	issues, result := api.Issues(client.requestURL(url)).All()
	if result.HasError() {
		err = FormatError("getting issues", result.Err)
		return
	}

	return
}

func (client *Client) CreateIssue(project *Project, title, body string, labels []string) (issue *octokit.Issue, err error) {
	url, err := octokit.RepoIssuesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("creating issues", err)
		return
	}

	params := octokit.IssueParams{
		Title:  title,
		Body:   body,
		Labels: labels,
	}
	issue, result := api.Issues(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = FormatError("creating issue", result.Err)
		return
	}

	return
}

func (client *Client) GhLatestTagName() (tagName string, err error) {
	url, err := octokit.ReleasesURL.Expand(octokit.M{"owner": "jingweno", "repo": "gh"})
	if err != nil {
		return
	}

	c := client.newOctokitClient(nil)
	releases, result := c.Releases(client.requestURL(url)).All()
	if result.HasError() {
		err = fmt.Errorf("Error getting gh release: %s", result.Err)
		return
	}

	if len(releases) == 0 {
		err = fmt.Errorf("No gh release is available")
		return
	}

	tagName = releases[0].TagName

	return
}

func (client *Client) CurrentUser() (user *octokit.User, err error) {
	url, err := octokit.CurrentUserURL.Expand(nil)
	if err != nil {
		return
	}

	api, err := client.api()
	if err != nil {
		err = FormatError("getting current user", err)
		return
	}

	user, result := api.Users(client.requestURL(url)).One()
	if result.HasError() {
		err = FormatError("getting current user", result.Err)
		return
	}

	return
}

func (client *Client) FindOrCreateToken(user, password, twoFactorCode string) (token string, err error) {
	authUrl, e := octokit.AuthorizationsURL.Expand(nil)
	if e != nil {
		err = &AuthError{e}
		return
	}

	basicAuth := octokit.BasicAuth{Login: user, Password: password, OneTimePassword: twoFactorCode}
	c := client.newOctokitClient(basicAuth)
	authsService := c.Authorizations(client.requestURL(authUrl))

	if twoFactorCode != "" {
		// dummy request to trigger a 2FA SMS since a HTTP GET won't do it
		authsService.Create(nil)
	}

	auths, result := authsService.All()
	if result.HasError() {
		err = &AuthError{result.Err}
		return
	}

	var moreAuths []octokit.Authorization
	for result.NextPage != nil {
		authUrl, e := result.NextPage.Expand(nil)
		if e != nil {
			return "", e
		}
		authUrl, _ = url.Parse(authUrl.RequestURI())

		as := c.Authorizations(authUrl)
		moreAuths, result = as.All()
		if result.HasError() {
			err = &AuthError{result.Err}
			return
		}

		auths = append(auths, moreAuths...)
	}

	for _, auth := range auths {
		if auth.Note == OAuthAppName || auth.NoteURL == OAuthAppURL {
			token = auth.Token
			break
		}
	}

	if token == "" {
		authParam := octokit.AuthorizationParams{}
		authParam.Scopes = append(authParam.Scopes, "repo")
		authParam.Note = OAuthAppName
		authParam.NoteURL = OAuthAppURL

		auth, result := authsService.Create(authParam)
		if result.HasError() {
			err = &AuthError{result.Err}
			return
		}

		token = auth.Token
	}

	return
}

func (client *Client) api() (c *octokit.Client, err error) {
	if client.Host.AccessToken == "" {
		host, e := CurrentConfig().PromptForHost(client.Host.Host)
		if e != nil {
			err = e
			return
		}
		client.Host = host
	}

	tokenAuth := octokit.TokenAuth{AccessToken: client.Host.AccessToken}
	c = client.newOctokitClient(tokenAuth)

	return
}

func (client *Client) newOctokitClient(auth octokit.AuthMethod) *octokit.Client {
	var host string
	if client.Host != nil {
		host = client.Host.Host
	}
	host = normalizeHost(host)
	apiHostURL := client.absolute(host)

	httpClient := newHttpClient(os.Getenv("HUB_TEST_HOST"), os.Getenv("HUB_VERBOSE") != "")
	c := octokit.NewClientWith(apiHostURL.String(), UserAgent, auth, httpClient)

	return c
}

func (client *Client) absolute(endpoint string) *url.URL {
	u, _ := url.Parse(endpoint)
	if u.Scheme == "" && client.Host != nil {
		u.Scheme = client.Host.Protocol
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	return u
}

func (client *Client) requestURL(u *url.URL) (uu *url.URL) {
	uu = u
	if client.Host != nil && client.Host.Host != GitHubHost {
		uu, _ = url.Parse(fmt.Sprintf("/api/v3/%s", u.Path))
	}

	return
}

func normalizeHost(host string) string {
	host = strings.ToLower(host)
	if host == "" {
		host = GitHubHost
	}

	if host == GitHubHost {
		host = GitHubApiHost
	}

	return host
}

func FormatError(action string, err error) (ee error) {
	switch e := err.(type) {
	default:
		ee = err
	case *octokit.ResponseError:
		statusCode := e.Response.StatusCode
		var reason string
		if s := strings.SplitN(e.Response.Status, " ", 2); len(s) >= 2 {
			reason = strings.TrimSpace(s[1])
		}

		errStr := fmt.Sprintf("Error %s: %s (HTTP %d)", action, reason, statusCode)

		var messages []string
		if statusCode == 422 {
			if e.Message != "" {
				messages = append(messages, e.Message)
			}

			if len(e.Errors) > 0 {
				for _, e := range e.Errors {
					messages = append(messages, e.Error())
				}
			}
		}

		if len(messages) > 0 {
			errStr = fmt.Sprintf("%s\n%s", errStr, strings.Join(messages, "\n"))
		}

		ee = fmt.Errorf(errStr)
	case *AuthError:
		errStr := fmt.Sprintf("Error %s: Unauthorized (HTTP 401)", action)
		ee = fmt.Errorf(errStr)
	}

	return
}

func warnExistenceOfRepo(project *Project, ee error) (err error) {
	if e, ok := ee.(*octokit.ResponseError); ok && e.Response.StatusCode == 404 {
		var url string
		if s := strings.SplitN(project.WebURL("", "", ""), "://", 2); len(s) >= 2 {
			url = s[1]
		}
		if url != "" {
			err = fmt.Errorf("Are you sure that %s exists?", url)
		}
	}

	return
}
