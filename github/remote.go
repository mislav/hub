package github

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/github/hub/git"
)

var (
	OriginNamesInLookupOrder = []string{"upstream", "github", "origin"}
)

type Remote struct {
	Name string
	URL  *url.URL
}

func (remote *Remote) String() string {
	return remote.Name
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

	// build the remotes map
	remotesMap := make(map[string]string)
	for _, r := range rs {
		if re.MatchString(r) {
			match := re.FindStringSubmatch(r)
			name := strings.TrimSpace(match[1])
			url := strings.TrimSpace(match[2])
			remotesMap[name] = url
		}
	}

	// construct remotes in priority order
	names := OriginNamesInLookupOrder
	for _, name := range names {
		if u, ok := remotesMap[name]; ok {
			url, e := git.ParseURL(u)
			if e == nil {
				remotes = append(remotes, Remote{Name: name, URL: url})
				delete(remotesMap, name)
			}
		}
	}

	// the rest of the remotes
	for n, u := range remotesMap {
		url, e := git.ParseURL(u)
		if e == nil {
			remotes = append(remotes, Remote{Name: n, URL: url})
		}
	}

	return
}
