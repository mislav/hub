package git

import (
	"errors"
	"net/url"
	"regexp"
)

type Remote struct {
	Name string
	URL  *url.URL
}

func Remotes() (remotes []Remote, err error) {
	re := regexp.MustCompile(`(.+)\s+(.+)\s+\((push|fetch)\)`)

	output, err := execGitCmd("remote", "-v")
	if err != nil {
		err = errors.New("Can't load git remote")
		return
	}

	remotesMap := make(map[string]string)
	for _, o := range output {
		if re.MatchString(o) {
			match := re.FindStringSubmatch(o)
			remotesMap[match[1]] = match[2]
		}
	}

	for k, v := range remotesMap {
		url, e := ParseURL(v)
		if e != nil {
			err = e
			return
		}

		remotes = append(remotes, Remote{Name: k, URL: url})
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

	return nil, errors.New("Can't find git remote origin")
}
