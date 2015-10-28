package github

import (
	"fmt"
	"net/url"

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
	remotesMap, err := git.Remotes()
	if err != nil {
		err = fmt.Errorf("Can't load git remote")
		return
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
