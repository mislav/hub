package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/octokat"
	"os"
	"regexp"
)

func isDir(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func parsePullRequestId(rawurl string) (id string) {
	url, err := github.ParseURL(rawurl)
	if err != nil {
		return
	}

	pullURLRegex := regexp.MustCompile("^pull/(\\d+)")
	projectPath := url.ProjectPath()
	if pullURLRegex.MatchString(projectPath) {
		id = pullURLRegex.FindStringSubmatch(projectPath)[1]
	}

	return
}

func fetchPullRequest(id string) (*octokat.PullRequest, error) {
	gh := github.New()
	pullRequest, err := gh.PullRequest(id)
	if err != nil {
		return nil, err
	}

	if pullRequest.Head.Repo.ID == 0 {
		user := pullRequest.User.Login
		return nil, fmt.Errorf("%s's fork is not available anymore", user)
	}

	return pullRequest, nil
}

func convertToGitURL(pullRequest *octokat.PullRequest) (string, error) {
	pullRequestURL := pullRequest.HTMLURL
	user := pullRequest.User.Login
	isSSH := pullRequest.Head.Repo.Private

	url, err := github.ParseURL(pullRequestURL)
	if err != nil {
		return "", err
	}

	return url.GitURL("", user, isSSH), nil
}

func parseRepoNameOwner(nameWithOwner string) (owner, name string, match bool) {
	ownerRe := fmt.Sprintf("^(%s)$", OwnerRe)
	ownerRegexp := regexp.MustCompile(ownerRe)
	if ownerRegexp.MatchString(nameWithOwner) {
		owner = ownerRegexp.FindStringSubmatch(nameWithOwner)[1]
		match = true
		return
	}

	nameWithOwnerRe := fmt.Sprintf("^(%s)\\/(%s)$", OwnerRe, NameRe)
	nameWithOwnerRegexp := regexp.MustCompile(nameWithOwnerRe)
	if nameWithOwnerRegexp.MatchString(nameWithOwner) {
		result := nameWithOwnerRegexp.FindStringSubmatch(nameWithOwner)
		owner = result[1]
		name = result[2]
		match = true
	}

	return
}
