package github

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/octokit/go-octokit/octokit"
)

const (
	GitHubHost    string = "github.com"
	GitHubApiHost string = "api.github.com"
	UserAgent     string = "Hub"
	OAuthAppURL   string = "http://owenou.com/gh"
)

func NewClient(host string) *Client {
	c := CurrentConfigs().PromptFor(host)
	return &Client{c}
}

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
	Credential *Credential
}

func (client *Client) PullRequest(project *Project, id string) (pr *octokit.PullRequest, err error) {
	url, err := octokit.PullRequestsURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "number": id})
	if err != nil {
		return
	}

	pr, result := client.octokit().PullRequests(client.requestURL(url)).One()
	if result.HasError() {
		err = formatError("getting pull request", result)
		return
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
		err = formatError("creating pull request", result)
		if e := warnExistenceOfRepo(project, result); e != nil {
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

	params := octokit.PullRequestForIssueParams{Base: base, Head: head, Issue: issue}
	pr, result := client.octokit().PullRequests(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = formatError("creating pull request", result)
		if e := warnExistenceOfRepo(project, result); e != nil {
			err = fmt.Errorf("%s\n%s", err, e)
		}

		return
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
		err = formatError("getting repository", result)
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
	if project.Owner != client.Credential.User {
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
		err = formatError("creating repository", result)
		return
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
		err = formatError("getting release", result)
		return
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
		err = formatError("creating release", result)
		return
	}

	return
}

func (client *Client) UploadReleaseAsset(uploadUrl *url.URL, asset *os.File, contentType string) (err error) {
	c := client.octokit()
	fileInfo, err := asset.Stat()
	if err != nil {
		return
	}

	result := c.Uploads(uploadUrl).UploadAsset(asset, contentType, fileInfo.Size())
	if result.HasError() {
		err = formatError("uploading asset", result)
		return
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
		err = formatError("getting CI status", result)
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
		err = formatError("forking repository", result)
		return
	}

	return
}

func (client *Client) Issues(project *Project) (issues []octokit.Issue, err error) {
	url, err := octokit.RepoIssuesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	issues, result := client.octokit().Issues(client.requestURL(url)).All()
	if result.HasError() {
		err = formatError("getting issues", result)
		return
	}

	return
}

func (client *Client) CreateIssue(project *Project, title, body string, labels []string) (issue *octokit.Issue, err error) {
	url, err := octokit.RepoIssuesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name})
	if err != nil {
		return
	}

	params := octokit.IssueParams{
		Title:  title,
		Body:   body,
		Labels: labels,
	}

	issue, result := client.octokit().Issues(client.requestURL(url)).Create(params)
	if result.HasError() {
		err = formatError("creating issue", result)
		return
	}

	return
}

func (client *Client) GhLatestTagName() (tagName string, err error) {
	url, err := octokit.ReleasesURL.Expand(octokit.M{"owner": "jingweno", "repo": "gh"})
	if err != nil {
		return
	}

	c := octokit.NewClientWith(client.apiEndpoint(), UserAgent, nil, nil)
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

func (client *Client) FindOrCreateToken(user, password, twoFactorCode string) (token string, err error) {
	url, e := octokit.AuthorizationsURL.Expand(nil)
	if e != nil {
		err = &ClientError{e}
		return
	}

	basicAuth := octokit.BasicAuth{Login: user, Password: password, OneTimePassword: twoFactorCode}
	c := octokit.NewClientWith(client.apiEndpoint(), UserAgent, basicAuth, nil)
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

// An implementation of http.ProxyFromEnvironment that isn't broken
func proxyFromEnvironment(req *http.Request) (*url.URL, error) {
	proxy := os.Getenv("http_proxy")
	if proxy == "" {
		proxy = os.Getenv("HTTP_PROXY")
	}
	if proxy == "" {
		return nil, nil
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil || !strings.HasPrefix(proxyURL.Scheme, "http") {
		if proxyURL, err := url.Parse("http://" + proxy); err == nil {
			return proxyURL, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
	}
	return proxyURL, nil
}

func (client *Client) octokit() (c *octokit.Client) {
	tokenAuth := octokit.TokenAuth{AccessToken: client.Credential.AccessToken}
	tr := &http.Transport{Proxy: proxyFromEnvironment}
	httpClient := &http.Client{Transport: tr}
	c = octokit.NewClientWith(client.apiEndpoint(), UserAgent, tokenAuth, httpClient)

	return
}

func (client *Client) requestURL(u *url.URL) (uu *url.URL) {
	uu = u
	if client.Credential != nil && client.Credential.Host != GitHubHost {
		uu, _ = url.Parse(fmt.Sprintf("/api/v3/%s", u.Path))
	}

	return
}

func (client *Client) apiEndpoint() string {
	host := os.Getenv("GH_API_HOST")
	if host == "" && client.Credential != nil {
		host = client.Credential.Host
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

func formatError(action string, result *octokit.Result) error {
	if e, ok := result.Err.(*octokit.ResponseError); ok {
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

		return fmt.Errorf(errStr)
	}

	return result.Err
}

func warnExistenceOfRepo(project *Project, result *octokit.Result) (err error) {
	if e, ok := result.Err.(*octokit.ResponseError); ok && e.Response.StatusCode == 404 {
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
