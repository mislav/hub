package git

import (
	"errors"
	"regexp"
)

type GitRemote struct {
	Name string
	URL  string
}

func Remotes() ([]*GitRemote, error) {
	r := regexp.MustCompile("(.+)\t(.+github.com.+) \\(push\\)")
	output, err := execGitCmd("remote", "-v")
	if err != nil {
		return nil, errors.New("Can't load git remote")
	}

	remotes := make([]*GitRemote, 0)
	for _, o := range output {
		if r.MatchString(o) {
			match := r.FindStringSubmatch(o)
			remotes = append(remotes, &GitRemote{Name: match[1], URL: match[2]})
		}
	}

	if len(remotes) == 0 {
		return nil, errors.New("Can't find git remote (push)")
	}

	return remotes, nil
}

func OriginRemote() (*GitRemote, error) {
	remotes, err := Remotes()
	if err != nil {
		return nil, err
	}

	for _, r := range remotes {
		if r.Name == "origin" {
			return r, nil
		}
	}

	return nil, errors.New("Can't find git remote orign (push)")
}
