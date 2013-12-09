package git

import (
	"fmt"
	"net/url"
	"regexp"
)

func ParseURL(rawurl string) (u *url.URL, err error) {
	sshGitRegexp := regexp.MustCompile(`(.+)@(.+):(.+)(\.git)?`)
	if sshGitRegexp.MatchString(rawurl) {
		match := sshGitRegexp.FindStringSubmatch(rawurl)
		user := match[1]
		host := match[2]
		path := match[3]
		ext := match[4]
		rawurl = fmt.Sprintf("ssh://%s@%s/%s%s", user, host, path, ext)
	}

	return url.Parse(rawurl)
}
