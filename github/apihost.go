package github

import (
	"net/url"
	"os"
)

type apiHost struct {
	Host string
}

func (ah *apiHost) String() string {
	host := os.Getenv("HUB_TEST_HOST")
	if host == "" && ah.Host != "" {
		host = ah.Host
	}

	if host == GitHubHost {
		host = GitHubApiHost
	}

	return absolute(host)
}

func absolute(endpoint string) string {
	u, _ := url.Parse(endpoint)
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	return u.String()
}
