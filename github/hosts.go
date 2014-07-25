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
	defaultHost := DefaultGitHubHost()
	hosts = append(hosts, defaultHost)
	hosts = append(hosts, fmt.Sprintf("ssh.%s", defaultHost))

	ghHosts, _ := git.Config("hub.host")
	for _, ghHost := range strings.Split(ghHosts, "\n") {
		ghHost = strings.TrimSpace(ghHost)
		if ghHosts != "" {
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
