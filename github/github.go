package github

import (
	"errors"
	"fmt"
	"github.com/jingweno/octokat"
)

const (
	GitHubHost  string = "github.com"
	OAuthAppURL string = "http://owenou.com/gh"
)

type GitHub struct {
	Project *Project
	Config  *Config
}

func (gh *GitHub) PullRequest(id string) (*octokat.PullRequest, error) {
	client := gh.client()

	return client.PullRequest(gh.repo(), id)
}

func (gh *GitHub) CreatePullRequest(base, head, title, body string) (string, error) {
	client := gh.client()
	params := octokat.PullRequestParams{base, head, title, body}
	pullRequest, err := client.CreatePullRequest(gh.repo(), params)
	if err != nil {
		return "", err
	}

	return pullRequest.HTMLURL, nil
}

func (gh *GitHub) Repository(project Project) (*octokat.Repository, error) {
	client := gh.client()

	return client.Repository(octokat.Repo{project.Name, project.Owner})
}

// TODO: detach GitHub from Project
func (gh *GitHub) IsRepositoryExist(project Project) bool {
	repo, err := gh.Repository(project)

	return err == nil && repo != nil
}

func (gh *GitHub) CreateRepository(project Project, description, homepage string, isPrivate bool) (*octokat.Repository, error) {
	params := octokat.Params{"description": description, "homepage": homepage, "private": isPrivate}
	if project.Owner != gh.Config.FetchUser() {
		params.Put("organization", project.Owner)
	}

	client := gh.client()

	return client.CreateRepository(project.Name, &params)
}

func (gh *GitHub) CreatePullRequestForIssue(base, head, issue string) (string, error) {
	client := gh.client()
	params := octokat.PullRequestForIssueParams{base, head, issue}
	pullRequest, err := client.CreatePullRequestForIssue(gh.repo(), params)
	if err != nil {
		return "", err
	}

	return pullRequest.HTMLURL, nil
}

func (gh *GitHub) CiStatus(sha string) (*octokat.Status, error) {
	client := gh.client()
	statuses, err := client.Statuses(gh.repo(), sha)
	if err != nil {
		return nil, err
	}

	if len(statuses) == 0 {
		return nil, nil
	}

	return &statuses[0], nil
}

func (gh *GitHub) ForkRepository(name, owner string, noRemote bool) (repo *octokat.Repository, err error) {
	client := gh.client()
	config := gh.Config
	repo, err = client.Repository(octokat.Repo{name, config.User})
	if repo != nil && err == nil {
		msg := fmt.Sprintf("Error creating fork: %s exists on %s", repo.FullName, GitHubHost)
		err = errors.New(msg)
		return
	}

	repo, err = client.Fork(octokat.Repo{name, owner}, nil)

	return
}

func (gh *GitHub) ExpandRemoteUrl(owner, name string, isSSH bool) (url string) {
	project := gh.Project
	if owner == "origin" {
		config := gh.Config
		owner = config.FetchUser()
	}

	return project.GitURL(name, owner, isSSH)
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
		if auth.NoteUrl == OAuthAppURL {
			token = auth.Token
			break
		}
	}

	if token == "" {
		authParam := octokat.AuthorizationParams{}
		authParam.Scopes = append(authParam.Scopes, "repo")
		authParam.Note = "gh"
		authParam.NoteUrl = OAuthAppURL

		auth, err := client.CreatedAuthorization(authParam)
		if err != nil {
			return "", err
		}

		token = auth.Token
	}

	return token, nil
}

func (gh *GitHub) client() *octokat.Client {
	config := gh.Config
	config.FetchCredentials()

	return octokat.NewClient().WithToken(config.Token)
}

func New() *GitHub {
	project := CurrentProject()
	c := CurrentConfig()
	c.FetchUser()

	return &GitHub{project, c}
}

// TODO: detach project from GitHub
func NewWithoutProject() *GitHub {
	c := CurrentConfig()
	c.FetchUser()

	return &GitHub{nil, c}
}
