package github

import (
	"github.com/jingweno/gh/utils"
	"github.com/jingweno/octokat"
)

const (
	GitHubHost  string = "github.com"
	OAuthAppUrl string = "http://owenou.com/gh"
)

type GitHub struct {
	Project *Project
	config  *Config
}

func (gh *GitHub) CreatePullRequest(base, head, title, body string) (string, error) {
	client := gh.client()
	params := octokat.PullRequestParams{base, head, title, body}
	pullRequest, err := client.CreatePullRequest(gh.repo(), params)
	if err != nil {
		return "", err
	}

	return pullRequest.HtmlUrl, nil
}

func (gh *GitHub) CreatePullRequestForIssue(base, head, issue string) (string, error) {
	client := gh.client()
	params := octokat.PullRequestForIssueParams{base, head, issue}
	pullRequest, err := client.CreatePullRequestForIssue(gh.repo(), params)
	if err != nil {
		return "", err
	}

	return pullRequest.HtmlUrl, nil
}

func (gh *GitHub) CiStatus(sha string) (*octokat.Status, error) {
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

func (gh *GitHub) repo() octokat.Repo {
	project := gh.Project
	return octokat.Repo{project.Name, project.Owner}
}

func findOrCreateToken(user, password string) (string, error) {
	client := octokat.NewClient().WithLogin(user, password)
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
		authParam := octokat.AuthorizationParams{}
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

func (gh *GitHub) client() *octokat.Client {
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

	return octokat.NewClient().WithToken(config.Token)
}

func New() *GitHub {
	project := CurrentProject()
	c, _ := loadConfig()
	c.FetchUser()

	return &GitHub{project, &c}
}
