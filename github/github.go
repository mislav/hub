package github

import (
	"encoding/base64"
	"fmt"
	"github.com/jingweno/gh/octokit"
	"github.com/jingweno/gh/utils"
	"net/http"
	"strings"
)

const (
	GitHubHost    string = "github.com"
	GitHubApiUrl  string = "https://" + GitHubApiHost
	GitHubApiHost string = "api.github.com"
	OAuthAppUrl   string = "http://owenou.com/gh"
)

type GitHub struct {
	httpClient    *http.Client
	authorization string
	Project       *Project
	config        *Config
}

func (gh *GitHub) updateToken(token string) {
	gh.config.Token = token
	saveConfig(gh.config)
}

func (gh *GitHub) updateTokenAuth(token string) {
	gh.authorization = fmt.Sprintf("token %s", token)
}

func (gh *GitHub) updateBasicAuth(user, pass string) {
	gh.authorization = fmt.Sprintf("Basic %s", hashAuth(user, pass))
}

func (gh *GitHub) setAuth(auth string) {
	gh.authorization = auth
}

func (gh *GitHub) isBasicAuth() bool {
	return strings.HasPrefix(gh.authorization, "Basic")
}

func (gh *GitHub) CreatePullRequest(base, head, title, body string) (string, error) {
	client := gh.client()
	params := octokit.PullRequestParams{base, head, title, body}
	pullRequest, err := client.CreatePullRequest(gh.repo(), params)
	if err != nil {
		return "", err
	}

	return pullRequest.HtmlUrl, nil
}

func (gh *GitHub) CreatePullRequestForIssue(base, head, issue string) (string, error) {
	client := gh.client()
	params := octokit.PullRequestForIssueParams{base, head, issue}
	pullRequest, err := client.CreatePullRequestForIssue(gh.repo(), params)
	if err != nil {
		return "", err
	}

	return pullRequest.HtmlUrl, nil
}

func (gh *GitHub) CIStatus(sha string) (*octokit.Status, error) {
	client := gh.client()
	statuses, err := client.Statuses(gh.repo(), sha)
	if err != nil {
		return nil, err
	}

	if len(statuses) == 0 {
		return nil, nil
	} else {
		return &statuses[0], nil
	}
}

func (gh *GitHub) repo() octokit.Repository {
	project := gh.Project
	return octokit.Repository{project.Name, project.Owner}
}

func findOrCreateToken(user, password string) (string, error) {
	client := octokit.NewClientWithPassword(user, password)
	auths, err := client.Authorizations()
	if err != nil {
		return "", err
	}

	var token string
	for _, auth := range auths {
		if auth.NoteUrl == OAuthAppUrl {
			token = auth.Token
			break
		}
	}

	if token == "" {
		authParam := octokit.AuthorizationParams{}
		authParam.Scopes = append(authParam.Scopes, "repo")
		authParam.Note = "gh"
		authParam.NoteUrl = OAuthAppUrl

		auth, err := client.CreatedAuthorization(authParam)
		if err != nil {
			return "", err
		}

		token = auth.Token
	}

	return token, nil
}

func (gh *GitHub) client() *octokit.Client {
	config := gh.config
	if config.User == "" {
		config.FetchUser()
	}

	if config.Token == "" {
		password := config.FetchPassword()
		token, err := findOrCreateToken(config.User, password)
		utils.Check(err)

		config.Token = token
		err = saveConfig(config)
		utils.Check(err)
	}

	return octokit.NewClientWithToken(config.Token)
}

func hashAuth(u, p string) string {
	var a = fmt.Sprintf("%s:%s", u, p)
	return base64.StdEncoding.EncodeToString([]byte(a))
}

func New() *GitHub {
	project := CurrentProject()
	c, _ := loadConfig()
	c.FetchUser()

	gh := GitHub{&http.Client{}, "", project, &c}
	if c.Token != "" {
		gh.updateTokenAuth(c.Token)
	}

	return &gh
}
