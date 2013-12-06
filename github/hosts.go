package github

import (
	"fmt"
	"os"
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
	host := os.Getenv("GITHUB_HOST")
	var mainHost string
	if host != "" {
		mainHost = host
	} else {
		mainHost = GitHubHost
	}

	hosts = append(hosts, mainHost)
	hosts = append(hosts, fmt.Sprintf("ssh.%s", mainHost))

	return
}
