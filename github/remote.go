package github

import (
	"fmt"
	"github.com/github/hub/git"
	"net/url"
	"regexp"
	"strings"
)

type Remote struct {
	Name string
	URL  *url.URL
}

func (remote *Remote) Project() (*Project, error) {
	return NewProjectFromURL(remote.URL)
}

func Remotes() (remotes []Remote, err error) {
	re := regexp.MustCompile(`(.+)\s+(.+)\s+\((push|fetch)\)`)

	rs, err := git.Remotes()
	if err != nil {
		err = fmt.Errorf("Can't load git remote")
		return
	}

	remotesMap := make(map[string]string)
	for _, r := range rs {
		if re.MatchString(r) {
			match := re.FindStringSubmatch(r)
			name := strings.TrimSpace(match[1])
			url := strings.TrimSpace(match[2])
			remotesMap[name] = url
		}
	}

	for n, u := range remotesMap {
		url, e := git.ParseURL(u)
		if e != nil {
			err = e
			return
		}

		remotes = append(remotes, Remote{Name: n, URL: url})
	}

	return
}

func OriginRemote() (*Remote, error) {
	remotes, err := Remotes()
	if err != nil {
		return nil, err
	}

	for _, r := range remotes {
		if r.Name == "origin" {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("Can't find git remote origin")
}
