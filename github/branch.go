package github

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/git"
)

type Branch struct {
	Repo *GitHubRepo
	Name string
}

var (
	shortRemoteRe = regexp.MustCompile("^refs/(remotes/)?.+?/")
	longRemoteRe  = regexp.MustCompile("^refs/(remotes/)?")
	remoteRe      = regexp.MustCompile("^refs/remotes/([^/]+)")
)

func (b *Branch) ShortName() string {
	return shortRemoteRe.ReplaceAllString(b.Name, "")
}

func (b *Branch) LongName() string {
	return longRemoteRe.ReplaceAllString(b.Name, "")
}

func (b *Branch) RemoteName() string {
	if remoteRe.MatchString(b.Name) {
		return remoteRe.FindStringSubmatch(b.Name)[1]
	}

	return ""
}

func (b *Branch) Upstream() (u *Branch, err error) {
	name, err := git.SymbolicFullName(fmt.Sprintf("%s@{upstream}", b.ShortName()))
	if err != nil {
		return
	}

	u = &Branch{b.Repo, name}

	return
}

func (b *Branch) IsMaster() bool {
	masterName := b.Repo.MasterBranch().ShortName()
	return b.ShortName() == masterName
}

func (b *Branch) IsRemote() bool {
	return strings.HasPrefix(b.Name, "refs/remotes")
}
