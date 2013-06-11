package github

import (
	"errors"
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"regexp"
)

type Project struct {
	Name  string
	Owner string
}

func (p *Project) WebURL(name, owner, path string) string {
	if owner == "" {
		owner = p.Owner
	}
	if name == "" {
		name = p.Name
	}

	url := fmt.Sprintf("https://%s", utils.ConcatPaths(GitHubHost, owner, name))
	if path != "" {
		url = utils.ConcatPaths(url, path)
	}

	return url
}

func (p *Project) LocalRepoWith(base, head string) *Repo {
	if base == "" {
		base = "master"
	}
	if head == "" {
		head, _ = git.Head()
	}

	return &Repo{base, head, p}
}

func (p *Project) LocalRepo() *Repo {
	return p.LocalRepoWith("", "")
}

func CurrentProject() *Project {
	owner, name := parseOwnerAndName()

	return &Project{name, owner}
}

func parseOwnerAndName() (name, remote string) {
	remote, err := git.Remote()
	utils.Check(err)

	url, err := mustMatchGitHubURL(remote)
	utils.Check(err)

	return url[1], url[2]
}

func mustMatchGitHubURL(url string) ([]string, error) {
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
