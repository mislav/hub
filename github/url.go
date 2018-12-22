package github

import (
	"net/url"
	"strings"
)

type URL struct {
	url.URL
	*Project
}

func (url URL) ProjectPath() (projectPath string) {
	split := strings.SplitN(url.Path, "/", 4)
	if len(split) > 3 {
		projectPath = split[3]
	}

	return
}

func ParseURL(rawurl string) (*URL, error) {
	repoUrl, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	project, err := NewProjectFromURL(repoUrl)
	if err != nil {
		return nil, err
	}

	return &URL{Project: project, URL: *repoUrl}, nil
}
