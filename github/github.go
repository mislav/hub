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

	return client.PullRequest(gh.repo(), id, nil)
}

func (gh *GitHub) CreatePullRequest(base, head, title, body string) (string, error) {
	client := gh.client()
	params := octokat.PullRequestParams{Base: base, Head: head, Title: title, Body: body}
	options := octokat.Options{Params: params}
	pullRequest, err := client.CreatePullRequest(gh.repo(), &options)
	if err != nil {
		return "", err
	}

	return pullRequest.HTMLURL, nil
}

func (gh *GitHub) CreatePullRequestForIssue(base, head, issue string) (string, error) {
	client := gh.client()
	params := octokat.PullRequestForIssueParams{Base: base, Head: head, Issue: issue}
	options := octokat.Options{Params: params}
	pullRequest, err := client.CreatePullRequest(gh.repo(), &options)
	if err != nil {
		return "", err
	}

	return pullRequest.HTMLURL, nil
}

func (gh *GitHub) Repository(project Project) (*octokat.Repository, error) {
	client := gh.client()
	repo := octokat.Repo{Name: project.Name, UserName: project.Owner}

	return client.Repository(repo, nil)
}

// TODO: detach GitHub from Project
func (gh *GitHub) IsRepositoryExist(project Project) bool {
	repo, err := gh.Repository(project)

	return err == nil && repo != nil
}

func (gh *GitHub) CreateRepository(project Project, description, homepage string, isPrivate bool) (*octokat.Repository, error) {
	params := octokat.RepositoryParams{Name: project.Name, Description: description, Homepage: homepage, Private: isPrivate}
	var org string
	if project.Owner != gh.Config.FetchUser() {
		org = project.Owner
	}

	client := gh.client()
	options := octokat.Options{Params: params}
	return client.CreateRepository(org, &options)
}

func (gh *GitHub) Releases() ([]octokat.Release, error) {
	client := gh.client()
	releases, err := client.Releases(gh.repo(), nil)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (gh *GitHub) CiStatus(sha string) (*octokat.Status, error) {
	client := gh.client()
	statuses, err := client.Statuses(gh.repo(), sha, nil)
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
	r := octokat.Repo{Name: name, UserName: config.User}
	repo, err = client.Repository(r, nil)
	if repo != nil && err == nil {
		msg := fmt.Sprintf("Error creating fork: %s exists on %s", repo.FullName, GitHubHost)
		err = errors.New(msg)
		return
	}

	repo, err = client.Fork(octokat.Repo{Name: name, UserName: owner}, nil)

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
	return octokat.Repo{Name: project.Name, UserName: project.Owner}
}

func findOrCreateToken(user, password, twoFactorCode string) (string, error) {
	client := octokat.NewClient().WithLogin(user, password)
	options := &octokat.Options{}
	if twoFactorCode != "" {
		headers := octokat.Headers{"X-GitHub-OTP": twoFactorCode}
		options.Headers = headers
	}

	auths, err := client.Authorizations(options)
	if err != nil {
		return "", err
	}

	var token string
	for _, auth := range auths {
		if auth.NoteURL == OAuthAppURL {
			token = auth.Token
			break
		}
	}

	if token == "" {
		authParam := octokat.AuthorizationParams{}
		authParam.Scopes = append(authParam.Scopes, "repo")
		authParam.Note = "gh"
		authParam.NoteURL = OAuthAppURL
		options.Params = authParam

		auth, err := client.CreateAuthorization(options)
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

func (gh *GitHub) Issues() ([]octokat.Issue, error) {
	client := gh.client()
	issues, err := client.Issues(gh.repo(), nil)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
