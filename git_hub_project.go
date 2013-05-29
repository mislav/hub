package main

import (
	"errors"
	"regexp"
)

type GitHubProject struct {
	Name  string
	Owner string
}

func CurrentProject() *GitHubProject {
	owner, name := parseOwnerAndName()

	return &GitHubProject{name, owner}
}

func parseOwnerAndName() (name, remote string) {
	remote, err := git.Remote()
	check(err)

	url, err := mustMatchGitUrl(remote)
	check(err)

	return url[1], url[2]
}

func mustMatchGitUrl(url string) ([]string, error) {
	httpRegex := regexp.MustCompile("https://.+/(.+)/(.+).git")
	if httpRegex.MatchString(url) {
		return httpRegex.FindStringSubmatch(url), nil
	}

	readOnlyRegex := regexp.MustCompile("git://.+/(.+)/(.+).git")
	if readOnlyRegex.MatchString(url) {
		return readOnlyRegex.FindStringSubmatch(url), nil
	}

	sshRegex := regexp.MustCompile(".+:(.+)/(.+).git")
	if sshRegex.MatchString(url) {
		return sshRegex.FindStringSubmatch(url), nil
	}

	return nil, errors.New("Can't parse git owner from URL: " + url)
}
