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
	Host  string
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
		if err != nil {
			utils.Check(fmt.Errorf("Aborted: not currently on any branch."))
		}
		head = Branch(headBranch).ShortName()
	}

	return &Repo{base, head, p}
}

func (p *Project) LocalRepo() *Repo {
	return p.LocalRepoWith("", "")
}

func CurrentProject() *Project {
	remote, err := git.OriginRemote()
	utils.Check(err)

	owner, name := parseOwnerAndName(remote.URL.String())

	return &Project{Name: name, Owner: owner}
}

func NewProjectFromURL(url *url.URL) (p *Project, err error) {
	if !KnownHosts().Include(url.Host) {
		err = fmt.Errorf("Invalid GitHub URL: %s", url)
		return
	}

	parts := strings.SplitN(url.Path, "/", 4)
	if len(parts) < 2 {
		err = fmt.Errorf("Invalid GitHub URL: %s", url)
		return
	}

	name := strings.TrimRight(parts[2], ".git")
	p = &Project{Name: name, Owner: parts[1], Host: url.Host}

	return
}

func NewProject(owner, name, host string) *Project {
	if strings.Contains(owner, "/") {
		result := strings.SplitN(owner, "/", 2)
		owner = result[0]
		if name == "" {
			name = result[1]
		}
	} else if strings.Contains(name, "/") {
		result := strings.SplitN(name, "/", 2)
		if owner == "" {
			owner = result[0]
		}
		name = result[1]
	}

	if host == "" {
		host = GitHubHost
	}

	if owner == "" {
		c := CurrentConfigs().PromptFor(host)
		owner = c.User
	}

	if name == "" {
		name, _ = utils.DirName()
	}

	return &Project{Name: name, Owner: owner, Host: host}
}

// Deprecated:
// NewProjectFromOwnerAndName creates a new Project from a specified owner and name
//
// If the owner or the name string is in the format of OWNER/NAME, it's split and used as the owner and name of the Project.
// If the owner string is empty, the current user is used as the name of the Project.
// If the name string is empty, the current dir name is used as the name of the Project.
func NewProjectFromOwnerAndName(owner, name string) *Project {
	if strings.Contains(owner, "/") {
		result := strings.SplitN(owner, "/", 2)
		owner = result[0]
		if name == "" {
			name = result[1]
		}
	} else if strings.Contains(name, "/") {
		result := strings.SplitN(name, "/", 2)
		if owner == "" {
			owner = result[0]
		}
		name = result[1]
	}

	if owner == "" {
		owner = CurrentConfig().FetchUser()
	}

	if name == "" {
		name, _ = utils.DirName()
	}

	return &Project{Name: name, Owner: owner}
}

func parseOwnerAndName(remote string) (owner string, name string) {
	url, err := mustMatchGitHubURL(remote)
	utils.Check(err)

	return url[1], url[2]
}

func MatchURL(url string) []string {
	httpRegex := regexp.MustCompile("https://github\\.com/(.+)/(.+?)(\\.git|$)")
	if httpRegex.MatchString(url) {
		return httpRegex.FindStringSubmatch(url)
	}

	readOnlyRegex := regexp.MustCompile("git://.*github\\.com/(.+)/(.+?)(\\.git|$)")
	if readOnlyRegex.MatchString(url) {
		return readOnlyRegex.FindStringSubmatch(url)
	}

	sshRegex := regexp.MustCompile("git@github\\.com:(.+)/(.+?)(\\.git|$)")
	if sshRegex.MatchString(url) {
		return sshRegex.FindStringSubmatch(url)
	}

	return nil
}

func mustMatchGitHubURL(url string) ([]string, error) {
	githubURL := MatchURL(url)
	if githubURL == nil {
		return nil, errors.New("The origin remote doesn't point to a GitHub repository: " + url)
	}

	return githubURL, nil
}
