package github

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/github/hub/version"
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

type Client struct {
	Host *Host
}

func (client *Client) FetchPullRequests(project *Project, filterParams map[string]interface{}, limit int, filter func(*PullRequest) bool) (pulls []PullRequest, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	path := fmt.Sprintf("repos/%s/%s/pulls?per_page=%d", project.Owner, project.Name, perPage(limit, 100))
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

	pulls = []PullRequest{}
	var res *simpleResponse

	for path != "" {
		res, err = api.Get(path)
		if err = checkStatus(200, "fetching pull requests", res, err); err != nil {
			return
		}
		path = res.Link("next")

		pullsPage := []PullRequest{}
		if err = res.Unmarshal(&pullsPage); err != nil {
			return
		}
		for _, pr := range pullsPage {
			if filter == nil || filter(&pr) {
				pulls = append(pulls, pr)
				if limit > 0 && len(pulls) == limit {
					path = ""
					break
				}
			}
		}
	}

	return
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
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.GetFile(fmt.Sprintf("repos/%s/%s/pulls/%s", project.Owner, project.Name, id), patchMediaType)
	if err = checkStatus(200, "getting pull request patch", res, err); err != nil {
		return
	}

	return res.Body, nil
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
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.GetFile(fmt.Sprintf("repos/%s/%s/commits/%s", project.Owner, project.Name, sha), patchMediaType)
	if err = checkStatus(200, "getting commit patch", res, err); err != nil {
		return
	}

	return res.Body, nil
}

type Gist struct {
	Files map[string]GistFile `json:"files"`
}
type GistFile struct {
	RawUrl string `json:"raw_url"`
}

func (client *Client) GistPatch(id string) (patch io.ReadCloser, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.Get(fmt.Sprintf("gists/%s", id))
	if err = checkStatus(200, "getting gist patch", res, err); err != nil {
		return
	}

	gist := Gist{}
	if err = res.Unmarshal(&gist); err != nil {
		return
	}
	rawUrl := ""
	for _, file := range gist.Files {
		rawUrl = file.RawUrl
		break
	}

	res, err = api.GetFile(rawUrl, textMediaType)
	if err = checkStatus(200, "getting gist patch", res, err); err != nil {
		return
	}

	return res.Body, nil
}

func (client *Client) Repository(project *Project) (repo *Repository, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.Get(fmt.Sprintf("repos/%s/%s", project.Owner, project.Name))
	if err = checkStatus(200, "getting commit patch", res, err); err != nil {
		return
	}

	repo = &Repository{}
	err = res.Unmarshal(&repo)
	return
}

func (client *Client) IsRepositoryExist(project *Project) bool {
	repo, err := client.Repository(project)

	return err == nil && repo != nil
}

func (client *Client) CreateRepository(project *Project, description, homepage string, isPrivate bool) (repo *Repository, err error) {
	repoURL := "user/repos"
	if project.Owner != client.Host.User {
		repoURL = fmt.Sprintf("orgs/%s/repos", project.Owner)
	}

	params := map[string]interface{}{
		"name":        project.Name,
		"description": description,
		"homepage":    homepage,
		"private":     isPrivate,
	}

	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.PostJSON(repoURL, params)
	if err = checkStatus(201, "creating repository", res, err); err != nil {
		return
	}

	repo = &Repository{}
	err = res.Unmarshal(repo)
	return
}

func (client *Client) DeleteRepository(project *Project) error {
	api, err := client.simpleApi()
	if err != nil {
		return err
	}

	repoURL := fmt.Sprintf("repos/%s/%s", project.Owner, project.Name)
	res, err := api.Delete(repoURL)
	return checkStatus(204, "deleting repository", res, err)
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
	CreatedAt       time.Time      `json:"created_at"`
	PublishedAt     time.Time      `json:"published_at"`
}

type ReleaseAsset struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	DownloadUrl string `json:"browser_download_url"`
	ApiUrl      string `json:"url"`
}

func (client *Client) FetchReleases(project *Project, limit int, filter func(*Release) bool) (releases []Release, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	path := fmt.Sprintf("repos/%s/%s/releases?per_page=%d", project.Owner, project.Name, perPage(limit, 100))

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
		for _, release := range releasesPage {
			if filter == nil || filter(&release) {
				releases = append(releases, release)
				if limit > 0 && len(releases) == limit {
					path = ""
					break
				}
			}
		}
	}

	return
}

func (client *Client) FetchRelease(project *Project, tagName string) (*Release, error) {
	releases, err := client.FetchReleases(project, 100, func(release *Release) bool {
		return release.TagName == tagName
	})

	if err == nil {
		if len(releases) < 1 {
			return nil, fmt.Errorf("Unable to find release with tag name `%s'", tagName)
		} else {
			return &releases[0], nil
		}
	} else {
		return nil, err
	}
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
	Name          string                 `json:"name"`
	FullName      string                 `json:"full_name"`
	Parent        *Repository            `json:"parent"`
	Owner         *User                  `json:"owner"`
	Private       bool                   `json:"private"`
	HasWiki       bool                   `json:"has_wiki"`
	Permissions   *RepositoryPermissions `json:"permissions"`
	HtmlUrl       string                 `json:"html_url"`
	DefaultBranch string                 `json:"default_branch"`
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
	Number int    `json:"number"`
	State  string `json:"state"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	User   *User  `json:"user"`

	PullRequest *PullRequest     `json:"pull_request"`
	Head        *PullRequestSpec `json:"head"`
	Base        *PullRequestSpec `json:"base"`

	MaintainerCanModify bool `json:"maintainer_can_modify"`

	Comments  int          `json:"comments"`
	Labels    []IssueLabel `json:"labels"`
	Assignees []User       `json:"assignees"`
	Milestone *Milestone   `json:"milestone"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	ApiUrl  string `json:"url"`
	HtmlUrl string `json:"html_url"`
}

type PullRequest Issue

type PullRequestSpec struct {
	Label string      `json:"label"`
	Ref   string      `json:"ref"`
	Sha   string      `json:"sha"`
	Repo  *Repository `json:"repo"`
}

func (pr *PullRequest) IsSameRepo() bool {
	return pr.Head != nil && pr.Head.Repo != nil &&
		pr.Head.Repo.Name == pr.Base.Repo.Name &&
		pr.Head.Repo.Owner.Login == pr.Base.Repo.Owner.Login
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

func (client *Client) FetchIssues(project *Project, filterParams map[string]interface{}, limit int, filter func(*Issue) bool) (issues []Issue, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	path := fmt.Sprintf("repos/%s/%s/issues?per_page=%d", project.Owner, project.Name, perPage(limit, 100))
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
		for _, issue := range issuesPage {
			if filter == nil || filter(&issue) {
				issues = append(issues, issue)
				if limit > 0 && len(issues) == limit {
					path = ""
					break
				}
			}
		}
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

func (client *Client) FetchLabels(project *Project) (labels []IssueLabel, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	path := fmt.Sprintf("repos/%s/%s/labels?per_page=100", project.Owner, project.Name)

	labels = []IssueLabel{}
	var res *simpleResponse

	for path != "" {
		res, err = api.Get(path)
		if err = checkStatus(200, "fetching labels", res, err); err != nil {
			return
		}
		path = res.Link("next")

		labelsPage := []IssueLabel{}
		if err = res.Unmarshal(&labelsPage); err != nil {
			return
		}
		labels = append(labels, labelsPage...)
	}

	return
}

func (client *Client) FetchMilestones(project *Project) (milestones []Milestone, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	path := fmt.Sprintf("repos/%s/%s/milestones?per_page=100", project.Owner, project.Name)

	milestones = []Milestone{}
	var res *simpleResponse

	for path != "" {
		res, err = api.Get(path)
		if err = checkStatus(200, "fetching milestones", res, err); err != nil {
			return
		}
		path = res.Link("next")

		milestonesPage := []Milestone{}
		if err = res.Unmarshal(&milestonesPage); err != nil {
			return
		}
		milestones = append(milestones, milestonesPage...)
	}

	return
}

func (client *Client) CurrentUser() (user *User, err error) {
	api, err := client.simpleApi()
	if err != nil {
		return
	}

	res, err := api.Get("user")
	if err = checkStatus(200, "getting current user", res, err); err != nil {
		return
	}

	user = &User{}
	err = res.Unmarshal(user)
	return
}

type AuthorizationEntry struct {
	Token string `json:"token"`
}

func isToken(api *simpleClient, password string) bool {
	api.PrepareRequest = func(req *http.Request) {
		req.Header.Set("Authorization", "token "+password)
	}

	res, _ := api.Get("user")
	if res != nil && res.StatusCode == 200 {
		return true
	}
	return false
}

func (client *Client) FindOrCreateToken(user, password, twoFactorCode string) (token string, err error) {
	api := client.apiClient()

	if len(password) >= 40 && isToken(api, password) {
		return password, nil
	}

	params := map[string]interface{}{
		"scopes":   []string{"repo"},
		"note_url": OAuthAppURL,
	}

	api.PrepareRequest = func(req *http.Request) {
		req.SetBasicAuth(user, password)
		if twoFactorCode != "" {
			req.Header.Set("X-GitHub-OTP", twoFactorCode)
		}
	}

	count := 1
	maxTries := 9
	for {
		params["note"], err = authTokenNote(count)
		if err != nil {
			return
		}

		res, postErr := api.PostJSON("authorizations", params)
		if postErr != nil {
			err = postErr
			break
		}

		if res.StatusCode == 201 {
			auth := &AuthorizationEntry{}
			if err = res.Unmarshal(auth); err != nil {
				return
			}
			token = auth.Token
			break
		} else if res.StatusCode == 422 && count < maxTries {
			count++
		} else {
			errInfo, e := res.ErrorInfo()
			if e == nil {
				err = errInfo
			} else {
				err = e
			}
			return
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

func (client *Client) simpleApi() (c *simpleClient, err error) {
	err = client.ensureAccessToken()
	if err != nil {
		return
	}

	c = client.apiClient()
	c.PrepareRequest = func(req *http.Request) {
		req.Header.Set("Authorization", "token "+client.Host.AccessToken)
	}
	return
}

func (client *Client) apiClient() *simpleClient {
	httpClient := newHttpClient(os.Getenv("HUB_TEST_HOST"), os.Getenv("HUB_VERBOSE") != "")
	apiRoot := client.absolute(normalizeHost(client.Host.Host))
	if client.Host != nil && client.Host.Host != GitHubHost {
		apiRoot.Path = "/api/v3/"
	}

	return &simpleClient{
		httpClient: httpClient,
		rootUrl:    apiRoot,
	}
}

func (client *Client) absolute(host string) *url.URL {
	u, _ := url.Parse("https://" + host + "/")
	if client.Host != nil && client.Host.Protocol != "" {
		u.Scheme = client.Host.Protocol
	}
	return u
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

func perPage(limit, max int) int {
	if limit > 0 {
		limit = limit + (limit / 2)
		if limit < max {
			return limit
		}
	}
	return max
}
