package github

import (
	"encoding/base64"
	"fmt"
	"github.com/jingweno/gh/config"
	"net/http"
	"strings"
)

const (
	GitHubApiUrl  string = "https://" + GitHubApiHost
	GitHubApiHost string = "api.github.com"
	OAuthAppUrl   string = "http://owenou.com/gh"
)

type GitHub struct {
	httpClient    *http.Client
	authorization string
	project       *Project
	config        *config.Config
}

func (gh *GitHub) updateToken(token string) {
	gh.config.Token = token
	config.Save(gh.config)
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

func (gh *GitHub) CreatePullRequest(params PullRequestParams) (*PullRequestResponse, error) {
	return createPullRequest(gh, params)
}

func (gh *GitHub) ListStatuses(ref string) ([]Status, error) {
	return listStatuses(gh, ref)
}

func hashAuth(u, p string) string {
	var a = fmt.Sprintf("%s:%s", u, p)
	return base64.StdEncoding.EncodeToString([]byte(a))
}

func New() *GitHub {
	project := CurrentProject()
	c, _ := config.Load()
	c.FetchUser()

	gh := GitHub{&http.Client{}, "", project, &c}
	if c.Token != "" {
		gh.updateTokenAuth(c.Token)
	}

	return &gh
}
