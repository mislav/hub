package github

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/utils"
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

	ownerWithName := fmt.Sprintf("%s/%s", owner, name)
	if strings.Contains(ownerWithName, ".wiki") {
		ownerWithName = strings.TrimSuffix(ownerWithName, ".wiki")
		if path != "wiki" {
			if strings.HasPrefix(path, "commits") {
				path = "_history"
			} else if path != "" {
				path = fmt.Sprintf("_%s", path)
			}

			if path != "" {
				path = utils.ConcatPaths("wiki", path)
			} else {
				path = "wiki"
			}
		}
	}

	url := fmt.Sprintf("https://%s", utils.ConcatPaths(p.Host, ownerWithName))
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

	host := rawHost(p.Host)

	if useHttpProtocol() {
		url = fmt.Sprintf("https://%s/%s/%s.git", host, owner, name)
	} else if isSSH {
		url = fmt.Sprintf("git@%s:%s/%s.git", host, owner, name)
	} else {
		url = fmt.Sprintf("git://%s/%s/%s.git", host, owner, name)
	}

	return url
}

// Remove the scheme from host when the host url is absolute.
func rawHost(host string) string {
	u, err := url.Parse(host)
	utils.Check(err)

	if u.IsAbs() {
		return u.Host
	} else {
		return u.Path
	}
}

func useHttpProtocol() bool {
	https := os.Getenv("HUB_PROTOCOL")
	if https == "" {
		https, _ = git.Config("hub.protocol")
	}

	return https == "https"
}

func NewProjectFromURL(url *url.URL) (p *Project, err error) {
	if !knownGitHubHosts().Include(url.Host) {
		err = fmt.Errorf("Invalid GitHub URL: %s", url)
		return
	}

	parts := strings.SplitN(url.Path, "/", 4)
	if len(parts) <= 2 {
		err = fmt.Errorf("Invalid GitHub URL: %s", url)
		return
	}

	name := strings.TrimSuffix(parts[2], ".git")
	p = NewProject(parts[1], name, url.Host)

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
		host = DefaultGitHubHost()
	}

	if owner == "" {
		h, e := CurrentConfigs().PromptForHost(host)
		utils.Check(e)
		owner = h.User
	}

	if name == "" {
		name, _ = utils.DirName()
	}

	return &Project{Name: name, Owner: owner, Host: host}
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
