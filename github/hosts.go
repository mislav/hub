package github

import (
	"fmt"
	"os"
	"strings"

	"github.com/github/hub/git"
)

var (
	GitHubHostEnv = os.Getenv("GITHUB_HOST")
)

type GitHubHosts []string

func (h GitHubHosts) Include(host string) bool {
	for _, hh := range h {
		if hh == host {
			return true
		}
	}

	return false
}

func knownGitHubHosts() (hosts GitHubHosts) {
	ghHosts, _ := git.Config("gh.host")
	for _, ghHost := range strings.Split(ghHosts, "\n") {
		hosts = append(hosts, ghHost)
	}

	defaultHost := DefaultGitHubHost()
	hosts = append(hosts, defaultHost)
	hosts = append(hosts, fmt.Sprintf("ssh.%s", defaultHost))

	return
}

func DefaultGitHubHost() string {
	defaultHost := GitHubHostEnv
	if defaultHost == "" {
		defaultHost = GitHubHost
	}

	return defaultHost
}
