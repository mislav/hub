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

	url, err := mustMatchGitHubUrl(remote)
	check(err)

	return url[1], url[2]
}

func mustMatchGitHubUrl(url string) ([]string, error) {
	httpRegex := regexp.MustCompile("https://github.com/(.+)/(.+).git")
	if httpRegex.MatchString(url) {
		return httpRegex.FindStringSubmatch(url), nil
	}

	readOnlyRegex := regexp.MustCompile("git://github.com/(.+)/(.+).git")
	if readOnlyRegex.MatchString(url) {
		return readOnlyRegex.FindStringSubmatch(url), nil
	}

	sshRegex := regexp.MustCompile("git@github.com:(.+)/(.+).git")
	if sshRegex.MatchString(url) {
		return sshRegex.FindStringSubmatch(url), nil
	}

	return nil, errors.New("The origin remote doesn't point to a GitHub repository: " + url)
}
