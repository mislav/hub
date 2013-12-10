package github

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"os"
	"strings"
)

type Hosts []string

func (h Hosts) Include(host string) bool {
	for _, hh := range h {
		if hh == host {
			return true
		}
	}

	return false
}

func KnownHosts() (hosts Hosts) {
	ghHosts, _ := git.Config("gh.host")
	for _, ghHost := range strings.Split(ghHosts, "\n") {
		hosts = append(hosts, ghHost)
	}

	defaultHost := DefaultHost()
	hosts = append(hosts, defaultHost)
	hosts = append(hosts, fmt.Sprintf("ssh.%s", defaultHost))

	return
}

func DefaultHost() string {
	defaultHost := os.Getenv("GITHUB_HOST")
	if defaultHost == "" {
		defaultHost = GitHubHost
	}

	return defaultHost
}
