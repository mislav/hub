package github

import (
	"fmt"
	"github.com/jingweno/gh/utils"
	"github.com/jingweno/go-octokit/octokit"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	GitHubHost    string = "github.com"
	GitHubApiHost string = "api.github.com"
	OAuthAppURL   string = "http://owenou.com/gh"
)

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
		err = result.Err
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
		err = result.Err
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
		err = result.Err
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
		err = result.Err
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
		err = result.Err
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
		err = result.Err
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
		err = result.Err
		return
	}

	return
}

func (client *Client) UploadReleaseAsset(release *octokit.Release, asset *os.File, fi os.FileInfo) (err error) {
	uploadUrl, err := octokit.Hyperlink(release.UploadURL).Expand(octokit.M{"name": fi.Name()})
	utils.Check(err)

	c := client.octokit()

	content, err := ioutil.ReadAll(asset)
	utils.Check(err)

	contentType := http.DetectContentType(content)

	fmt.Printf("-- Uploading %s to %s\n", contentType, uploadUrl.String())
	request, err := http.NewRequest("POST", uploadUrl.String(), asset)
	utils.Check(err)

	request.Header.Add("Content-Type", contentType)

	if c.AuthMethod != nil {
		request.Header.Add("Authorization", c.AuthMethod.String())
	}

	if basicAuth, ok := c.AuthMethod.(octokit.BasicAuth); ok && basicAuth.OneTimePassword != "" {
		request.Header.Add("X-GitHub-OTP", basicAuth.OneTimePassword)
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	utils.Check(err)

	if response.Status != "201" {
		return fmt.Errorf("Error uploading the release asset, status %s", response.Status)
	}
	return nil
}

func (client *Client) CIStatus(project *Project, sha string) (status *octokit.Status, err error) {
	url, err := octokit.StatusesURL.Expand(octokit.M{"owner": project.Owner, "repo": project.Name, "ref": sha})
	if err != nil {
		return
	}

	statuses, result := client.octokit().Statuses(client.requestURL(url)).All()
	if result.HasError() {
		err = result.Err
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
		err = result.Err
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
		err = result.Err
		return
	}

	return
}

func (client *Client) FindOrCreateToken(user, password, twoFactorCode string) (token string, err error) {
	url, err := octokit.AuthorizationsURL.Expand(nil)
	if err != nil {
		return
	}

	basicAuth := octokit.BasicAuth{Login: user, Password: password, OneTimePassword: twoFactorCode}
	c := octokit.NewClientWith(client.apiEndpoint(), nil, basicAuth)
	authsService := c.Authorizations(client.requestURL(url))

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
