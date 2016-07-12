package github

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/github/hub/git"
)

var (
	GitHubHostEnv = os.Getenv("GITHUB_HOST")
)

type GitHubHosts []string

type GithubHostError struct {
	url *url.URL
}

func (e *GithubHostError) Error() string {
	return fmt.Sprintf("Invalid GitHub URL: %s", e.url)
}

func (h GitHubHosts) Include(host string) bool {
	for _, hh := range h {
		if hh == host {
			return true
		}
	}

	return false
}

func knownGitHubHosts() (hosts GitHubHosts) {
	defaultHost := DefaultGitHubHost()
	hosts = append(hosts, defaultHost)
	hosts = append(hosts, "ssh."+GitHubHost)

	ghHosts, _ := git.ConfigAll("hub.host")
	for _, ghHost := range ghHosts {
		ghHost = strings.TrimSpace(ghHost)
		if ghHost != "" {
			hosts = append(hosts, ghHost)
		}
	}

	return
}

func DefaultGitHubHost() string {
	defaultHost := GitHubHostEnv
	if defaultHost == "" {
		defaultHost = GitHubHost
	}

	return defaultHost
}
