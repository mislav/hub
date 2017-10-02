package github

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/github/hub/version"
	"github.com/octokit/go-octokit/octokit"
)

const (
	GitHubHost    string = "github.com"
	GitHubApiHost string = "api.github.com"
	OAuthAppURL   string = "http://hub.github.com/"
)

var UserAgent = "Hub " + version.Version

func NewClient(h string) *Client {
	return NewClientWithHost(&Host{Host: h})
}

func NewClientWithHost(host *Host) *Client {
	return &Client{host}
}

type AuthError struct {
	Err error
}

func (e *AuthError) Error() string {
	return e.Err.Error()
}

func (e *AuthError) IsRequired2FACodeError() bool {
	re, ok := e.Err.(*octokit.ResponseError)
	return ok && re.Type == octokit.ErrorOneTimePasswordRequired
}

func (e *AuthError) IsDuplicatedTokenError() bool {
	re, ok := e.Err.(*octokit.ResponseError)
	return ok && re.Type == octokit.ErrorUnprocessableEntity
}

type Client struct {
	Host *Host
}

func (client *Client) PullRequest(project *Project, id string) (pr *PullRequest, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.Get(fmt.Sprintf("repos/%s/%s/pulls/%s", project.Owner, project.Name, id))
	if err = checkStatus(200, "getting pull request", res, err); err != nil {
		return
	}

	pr = &PullRequest{}
	err = res.Unmarshal(pr)

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

type PullRequest struct {
	ApiUrl              string           `json:"url"`
	Number              int              `json:"number"`
	HtmlUrl             string           `json:"html_url"`
	Title               string           `json:"title"`
	MaintainerCanModify bool             `json:"maintainer_can_modify"`
	Head                *PullRequestSpec `json:"head"`
	Base                *PullRequestSpec `json:"base"`
}

type PullRequestSpec struct {
	Label string      `json:"label"`
	Ref   string      `json:"ref"`
	Sha   string      `json:"sha"`
	Repo  *Repository `json:"repo"`
}

func (pr *PullRequest) IsSameRepo() bool {
	return pr.Head.Repo != nil &&
		pr.Head.Repo.Name == pr.Base.Repo.Name &&
		pr.Head.Repo.Owner.Login == pr.Base.Repo.Owner.Login
}

func (client *Client) CreatePullRequest(project *Project, params map[string]interface{}) (pr *PullRequest, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.PostJSON(fmt.Sprintf("repos/%s/%s/pulls", project.Owner, project.Name), params)
	if err = checkStatus(201, "creating pull request", res, err); err != nil {
		if res != nil && res.StatusCode == 404 {
			projectUrl := strings.SplitN(project.WebURL("", "", ""), "://", 2)[1]
			err = fmt.Errorf("%s\nAre you sure that %s exists?", err, projectUrl)
		}
		return
	}

	pr = &PullRequest{}
	err = res.Unmarshal(pr)

	return
}

func (client *Client) RequestReview(project *Project, prNumber int, params map[string]interface{}) (err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.PostReview(fmt.Sprintf("repos/%s/%s/pulls/%d/requested_reviewers", project.Owner, project.Name, prNumber), params)
	if err = checkStatus(201, "requesting reviewer", res, err); err != nil {
		return
	}

	res.Body.Close()
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

type Release struct {
	Name            string         `json:"name"`
	TagName         string         `json:"tag_name"`
	TargetCommitish string         `json:"target_commitish"`
	Body            string         `json:"body"`
	Draft           bool           `json:"draft"`
	Prerelease      bool           `json:"prerelease"`
	Assets          []ReleaseAsset `json:"assets"`
	TarballUrl      string         `json:"tarball_url"`
	ZipballUrl      string         `json:"zipball_url"`
	HtmlUrl         string         `json:"html_url"`
	UploadUrl       string         `json:"upload_url"`
	ApiUrl          string         `json:"url"`
}

type ReleaseAsset struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	DownloadUrl string `json:"browser_download_url"`
	ApiUrl      string `json:"url"`
}

func (client *Client) FetchReleases(project *Project) (releases []Release, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	path := fmt.Sprintf("repos/%s/%s/releases?per_page=100", project.Owner, project.Name)

	releases = []Release{}
	var res *simpleResponse

	for path != "" {
		res, err = api.Get(path)
		if err = checkStatus(200, "fetching releases", res, err); err != nil {
			return
		}
		path = res.Link("next")

		releasesPage := []Release{}
		if err = res.Unmarshal(&releasesPage); err != nil {
			return
		}
		releases = append(releases, releasesPage...)
	}

	return
}

func (client *Client) FetchRelease(project *Project, tagName string) (foundRelease *Release, err error) {
	releases, err := client.FetchReleases(project)
	if err != nil {
		return
	}

	for _, release := range releases {
		if release.TagName == tagName {
			foundRelease = &release
			break
		}
	}

	if foundRelease == nil {
		err = fmt.Errorf("Unable to find release with tag name `%s'", tagName)
	}
	return
}

func (client *Client) CreateRelease(project *Project, releaseParams *Release) (release *Release, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.PostJSON(fmt.Sprintf("repos/%s/%s/releases", project.Owner, project.Name), releaseParams)
	if err = checkStatus(201, "creating release", res, err); err != nil {
		return
	}

	release = &Release{}
	err = res.Unmarshal(release)
	return
}

func (client *Client) EditRelease(release *Release, releaseParams map[string]interface{}) (updatedRelease *Release, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.PatchJSON(release.ApiUrl, releaseParams)
	if err = checkStatus(200, "editing release", res, err); err != nil {
		return
	}

	updatedRelease = &Release{}
	err = res.Unmarshal(updatedRelease)
	return
}

func (client *Client) DeleteRelease(release *Release) (err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.Delete(release.ApiUrl)
	if err = checkStatus(204, "deleting release", res, err); err != nil {
		return
	}

	return
}

func (client *Client) UploadReleaseAsset(release *Release, filename, label string) (asset *ReleaseAsset, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	parts := strings.SplitN(release.UploadUrl, "{", 2)
	uploadUrl := parts[0]
	uploadUrl += "?name=" + url.QueryEscape(filepath.Base(filename))
	if label != "" {
		uploadUrl += "&label=" + url.QueryEscape(label)
	}

	res, err := api.PostFile(uploadUrl, filename)
	if err = checkStatus(201, "uploading release asset", res, err); err != nil {
		return
	}

	asset = &ReleaseAsset{}
	err = res.Unmarshal(asset)
	return
}

func (client *Client) DeleteReleaseAsset(asset *ReleaseAsset) (err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.Delete(asset.ApiUrl)
	err = checkStatus(204, "deleting release asset", res, err)

	return
}

func (client *Client) DownloadReleaseAsset(url string) (asset io.ReadCloser, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	resp, err := api.GetFile(url, "application/octet-stream")
	if err = checkStatus(200, "downloading asset", resp, err); err != nil {
		return
	}

	return resp.Body, err
}

type CIStatusResponse struct {
	State    string     `json:"state"`
	Statuses []CIStatus `json:"statuses"`
}

type CIStatus struct {
	State     string `json:"state"`
	Context   string `json:"context"`
	TargetUrl string `json:"target_url"`
}

func (client *Client) FetchCIStatus(project *Project, sha string) (status *CIStatusResponse, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.Get(fmt.Sprintf("repos/%s/%s/commits/%s/status", project.Owner, project.Name, sha))
	if err = checkStatus(200, "fetching statuses", res, err); err != nil {
		return
	}

	status = &CIStatusResponse{}
	err = res.Unmarshal(status)

	return
}

type Repository struct {
	Name        string                 `json:"name"`
	Parent      *Repository            `json:"parent"`
	Owner       *User                  `json:"owner"`
	Private     bool                   `json:"private"`
	Permissions *RepositoryPermissions `json:"permissions"`
	HtmlUrl     string                 `json:"html_url"`
}

type RepositoryPermissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

func (client *Client) ForkRepository(project *Project, params map[string]interface{}) (repo *Repository, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.PostJSON(fmt.Sprintf("repos/%s/%s/forks", project.Owner, project.Name), params)
	if err = checkStatus(202, "creating fork", res, err); err != nil {
		return
	}

	repo = &Repository{}
	err = res.Unmarshal(repo)

	return
}

type Issue struct {
	Number      int          `json:"number"`
	State       string       `json:"state"`
	Title       string       `json:"title"`
	Body        string       `json:"body"`
	User        *User        `json:"user"`
	Assignees   []User       `json:"assignees"`
	Labels      []IssueLabel `json:"labels"`
	PullRequest *PullRequest `json:"pull_request"`
	HtmlUrl     string       `json:"html_url"`
	Comments    int          `json:"comments"`
	Milestone   *Milestone   `json:"milestone"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type IssueLabel struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type User struct {
	Login string `json:"login"`
}

type Milestone struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func (client *Client) FetchIssues(project *Project, filterParams map[string]interface{}) (issues []Issue, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	path := fmt.Sprintf("repos/%s/%s/issues?per_page=100", project.Owner, project.Name)
	if filterParams != nil {
		query := url.Values{}
		for key, value := range filterParams {
			switch v := value.(type) {
			case string:
				query.Add(key, v)
			}
		}
		path += "&" + query.Encode()
	}

	issues = []Issue{}
	var res *simpleResponse

	for path != "" {
		res, err = api.Get(path)
		if err = checkStatus(200, "fetching issues", res, err); err != nil {
			return
		}
		path = res.Link("next")

		issuesPage := []Issue{}
		if err = res.Unmarshal(&issuesPage); err != nil {
			return
		}
		issues = append(issues, issuesPage...)
	}

	return
}

func (client *Client) CreateIssue(project *Project, params interface{}) (issue *Issue, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.PostJSON(fmt.Sprintf("repos/%s/%s/issues", project.Owner, project.Name), params)
	if err = checkStatus(201, "creating issue", res, err); err != nil {
		return
	}

	issue = &Issue{}
	err = res.Unmarshal(issue)
	return
}

func (client *Client) UpdateIssue(project *Project, issueNumber int, params map[string]interface{}) (err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.PatchJSON(fmt.Sprintf("repos/%s/%s/issues/%d", project.Owner, project.Name, issueNumber), params)
	if err = checkStatus(200, "updating issue", res, err); err != nil {
		return
	}

	res.Body.Close()
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

	basicAuth := octokit.BasicAuth{
		Login:           user,
		Password:        password,
		OneTimePassword: twoFactorCode,
	}
	c := client.newOctokitClient(basicAuth)
	authsService := c.Authorizations(client.requestURL(authUrl))

	authParam := octokit.AuthorizationParams{
		Scopes:  []string{"repo"},
		NoteURL: OAuthAppURL,
	}

	count := 1
	for {
		note, e := authTokenNote(count)
		if e != nil {
			err = e
			return
		}

		authParam.Note = note
		auth, result := authsService.Create(authParam)
		if !result.HasError() {
			token = auth.Token
			break
		}

		authErr := &AuthError{result.Err}
		if authErr.IsDuplicatedTokenError() {
			if count >= 9 {
				err = authErr
				break
			} else {
				count++
				continue
			}
		} else {
			err = authErr
			break
		}
	}

	return
}

func (client *Client) ensureAccessToken() (err error) {
	if client.Host.AccessToken == "" {
		host, err := CurrentConfig().PromptForHost(client.Host.Host)
		if err == nil {
			client.Host = host
		}
	}
	return
}

func (client *Client) api() (c *octokit.Client, err error) {
	err = client.ensureAccessToken()
	if err != nil {
		return
	}

	tokenAuth := octokit.TokenAuth{AccessToken: client.Host.AccessToken}
	c = client.newOctokitClient(tokenAuth)

	return
}

func (client *Client) simpleApi() (c *simpleClient, err error) {
	err = client.ensureAccessToken()
	if err != nil {
		return
	}

	httpClient := newHttpClient(os.Getenv("HUB_TEST_HOST"), os.Getenv("HUB_VERBOSE") != "")
	apiRoot := client.requestURL(client.absolute(normalizeHost(client.Host.Host)))

	c = &simpleClient{
		httpClient:  httpClient,
		rootUrl:     apiRoot,
		accessToken: client.Host.AccessToken,
	}
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

func (client *Client) absolute(host string) *url.URL {
	u, _ := url.Parse("https://" + host + "/")
	if client.Host != nil && client.Host.Protocol != "" {
		u.Scheme = client.Host.Protocol
	}
	return u
}

func (client *Client) requestURL(base *url.URL) *url.URL {
	if client.Host != nil && client.Host.Host != GitHubHost {
		newUrl, _ := url.Parse(base.String())
		basePath := base.Path
		if !strings.HasPrefix(basePath, "/") {
			basePath = "/" + basePath
		}
		newUrl.Path = "/api/v3" + basePath
		return newUrl
	} else {
		return base
	}
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

func checkStatus(expectedStatus int, action string, response *simpleResponse, err error) error {
	if err != nil {
		return fmt.Errorf("Error %s: %s", action, err.Error())
	} else if response.StatusCode != expectedStatus {
		errInfo, err := response.ErrorInfo()
		if err == nil {
			return FormatError(action, errInfo)
		} else {
			return fmt.Errorf("Error %s: %s (HTTP %d)", action, err.Error(), response.StatusCode)
		}
	} else {
		return nil
	}
}

func FormatError(action string, err error) (ee error) {
	switch e := err.(type) {
	default:
		ee = err
	case *AuthError:
		return FormatError(action, e.Err)
	case *octokit.ResponseError:
		info := &errorInfo{
			Message:  e.Message,
			Response: e.Response,
			Errors:   []fieldError{},
		}
		for _, err := range e.Errors {
			info.Errors = append(info.Errors, fieldError{
				Field:   err.Field,
				Message: err.Message,
				Code:    err.Code,
			})
		}
		return FormatError(action, info)
	case *errorInfo:
		statusCode := e.Response.StatusCode
		var reason string
		if s := strings.SplitN(e.Response.Status, " ", 2); len(s) >= 2 {
			reason = strings.TrimSpace(s[1])
		}

		errStr := fmt.Sprintf("Error %s: %s (HTTP %d)", action, reason, statusCode)

		var errorSentences []string
		for _, err := range e.Errors {
			switch err.Code {
			case "custom":
				errorSentences = append(errorSentences, err.Message)
			case "missing_field":
				errorSentences = append(errorSentences, fmt.Sprintf("Missing field: \"%s\"", err.Field))
			case "already_exists":
				errorSentences = append(errorSentences, fmt.Sprintf("Duplicate value for \"%s\"", err.Field))
			case "invalid":
				errorSentences = append(errorSentences, fmt.Sprintf("Invalid value for \"%s\"", err.Field))
			case "unauthorized":
				errorSentences = append(errorSentences, fmt.Sprintf("Not allowed to change field \"%s\"", err.Field))
			}
		}

		var errorMessage string
		if len(errorSentences) > 0 {
			errorMessage = strings.Join(errorSentences, "\n")
		} else {
			errorMessage = e.Message
		}

		if errorMessage != "" {
			errStr = fmt.Sprintf("%s\n%s", errStr, errorMessage)
		}

		ee = fmt.Errorf(errStr)
	}

	return
}

func authTokenNote(num int) (string, error) {
	n := os.Getenv("USER")

	if n == "" {
		n = os.Getenv("USERNAME")
	}

	if n == "" {
		whoami := exec.Command("whoami")
		whoamiOut, err := whoami.Output()
		if err != nil {
			return "", err
		}
		n = strings.TrimSpace(string(whoamiOut))
	}

	h, err := os.Hostname()
	if err != nil {
		return "", err
	}

	if num > 1 {
		return fmt.Sprintf("hub for %s@%s %d", n, h, num), nil
	}

	return fmt.Sprintf("hub for %s@%s", n, h), nil
}
