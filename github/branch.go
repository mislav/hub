package github

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"regexp"
	"strings"
)

type Branch string

func (b Branch) ShortName() string {
	reg := regexp.MustCompile("^refs/(remotes/)?.+?/")
	return reg.ReplaceAllString(string(b), "")
}

func (b Branch) Upstream() (u Branch, err error) {
	name, err := git.SymbolicFullName(fmt.Sprintf("%s@{upstream}", b.ShortName()))
	if err != nil {
		return
	}

	u = Branch(name)

	return
}

func (b Branch) RemoteName() string {
	reg := regexp.MustCompile("^refs/remotes/([^/]+)")
	if reg.MatchString(string(b)) {
		return reg.FindStringSubmatch(string(b))[1]
	}

	return ""
}

func (b Branch) IsRemote() bool {
	return strings.HasPrefix(string(b), "refs/remotes")
}
