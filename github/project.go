package github

import (
	"errors"
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"net/url"
	"regexp"
	"strings"
)

type Project struct {
	Name  string
	Owner string
}

func (p Project) String() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Name)
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

func (p *Project) GitURL(name, owner string, isSSH bool) (url string) {
	if name == "" {
		name = p.Name
	}
	if owner == "" {
		owner = p.Owner
	}

	if isSSH {
		url = fmt.Sprintf("git@%s:%s/%s.git", GitHubHost, owner, name)
	} else {
		url = fmt.Sprintf("git://%s.git", utils.ConcatPaths(GitHubHost, owner, name))
	}

	return url
}

func (p *Project) LocalRepoWith(base, head string) *Repo {
	if base == "" {
		base = "master"
	}
	if head == "" {
		headBranch, err := git.Head()
		utils.Check(err)
		head = headBranch.ShortName()
	}

	return &Repo{base, head, p}
}

func (p *Project) LocalRepo() *Repo {
	return p.LocalRepoWith("", "")
}

func CurrentProject() *Project {
	remote, err := git.OriginRemote()
	utils.Check(err)

	owner, name := parseOwnerAndName(remote.URL)

	return &Project{name, owner}
}

func ParseProjectFromURL(uu string) (*Project, error) {
	u, err := url.Parse(uu)
	if err != nil {
		return nil, err
	}

	if u.Host != GitHubHost || u.Scheme != "https" {
		return nil, fmt.Errorf("Invalid GitHub URL: %s", u)
	}

	parts := strings.SplitN(u.Path, "/", 4)
	if len(parts) >= 2 {
		return &Project{Name: parts[2], Owner: parts[1]}, nil
	}

	return nil, fmt.Errorf("Invalid GitHub URL: %s", u)
}

func NewProjectFromNameAndOwner(name, owner string) Project {
	if strings.Contains(owner, "/") {
		result := strings.SplitN(owner, "/", 2)
		owner = result[0]
		name = result[1]
	} else if strings.Contains(name, "/") {
		result := strings.SplitN(owner, "/", 2)
		owner = result[0]
		name = result[1]
	}

	if owner == "" {
		owner = CurrentConfig().FetchUser()
	}

	if name == "" {
		name, _ = utils.DirName()
	}

	return Project{Name: name, Owner: owner}
}

func parseOwnerAndName(remote string) (owner string, name string) {
	url, err := mustMatchGitHubURL(remote)
	utils.Check(err)

	return url[1], url[2]
}

func mustMatchGitHubURL(url string) ([]string, error) {
	httpRegex := regexp.MustCompile("https://github\\.com/(.+)/(.+?)(\\.git|$)")
	if httpRegex.MatchString(url) {
		return httpRegex.FindStringSubmatch(url), nil
	}

	readOnlyRegex := regexp.MustCompile("git://github\\.com/(.+)/(.+?)(\\.git|$)")
	if readOnlyRegex.MatchString(url) {
		return readOnlyRegex.FindStringSubmatch(url), nil
	}

	sshRegex := regexp.MustCompile("git@github\\.com:(.+)/(.+?)(\\.git|$)")
	if sshRegex.MatchString(url) {
		return sshRegex.FindStringSubmatch(url), nil
	}

	return nil, errors.New("The origin remote doesn't point to a GitHub repository: " + url)
}
