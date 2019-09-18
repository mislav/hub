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
	cachedHosts   []string
)

type GithubHostError struct {
	url *url.URL
}

func (e *GithubHostError) Error() string {
	return fmt.Sprintf("Invalid GitHub URL: %s", e.url)
}

func getKnownHost(host string) (foundHost string, err error) {
	for _, knownHost := range knownGitHubHosts() {
		// origin url may include ssh alias : instead of github.com we can have github.com-username
		// -username is used to resolve proper openssh key, host github.com-username does not exist
		// attempts to resolve will fail. -username part must be discarded
		if (host == knownHost || strings.HasPrefix(host, knownHost + "-")) {
			foundHost = knownHost
			return 
		}
	}
	err = fmt.Errorf("Not a known host")			
	return
}

func knownGitHubHosts() []string {
	if cachedHosts != nil {
		return cachedHosts
	}

	hosts := []string{}
	defaultHost := DefaultGitHubHost()
	hosts = append(hosts, defaultHost)
	hosts = append(hosts, "ssh.github.com")

	ghHosts, _ := git.ConfigAll("hub.host")
	for _, ghHost := range ghHosts {
		ghHost = strings.TrimSpace(ghHost)
		if ghHost != "" {
			hosts = append(hosts, ghHost)
		}
	}

	cachedHosts = hosts
	return hosts
}

func DefaultGitHubHost() string {
	defaultHost := GitHubHostEnv
	if defaultHost == "" {
		defaultHost = GitHubHost
	}

	return defaultHost
}
