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

func (b *Branch) ShortName() string {
	reg := regexp.MustCompile("^refs/(remotes/)?.+?/")
	return reg.ReplaceAllString(b.Name, "")
}

func (b *Branch) LongName() string {
	reg := regexp.MustCompile("^refs/(remotes/)?")
	return reg.ReplaceAllString(b.Name, "")
}

func (b *Branch) PushTarget(owner string, preferUpstream bool) (branch *Branch) {
	var err error
	pushDefault, _ := git.Config("push.default")
	if pushDefault == "upstream" || pushDefault == "tracking" {
		branch, err = b.Upstream()
		if err != nil {
			return
		}
	} else {
		shortName := b.ShortName()
		remotes := b.Repo.remotesForPublish(owner)

		var remotesInOrder []Remote
		if preferUpstream {
			// reverse the remote lookup order
			// see OriginNamesInLookupOrder
			for i := len(remotes) - 1; i >= 0; i-- {
				remotesInOrder = append(remotesInOrder, remotes[i])
			}
		} else {
			remotesInOrder = remotes
		}

		for _, remote := range remotesInOrder {
			if git.HasFile("refs", "remotes", remote.Name, shortName) {
				name := fmt.Sprintf("refs/remotes/%s/%s", remote.Name, shortName)
				branch = &Branch{b.Repo, name}
				break
			}
		}
	}

	return
}

func (b *Branch) RemoteName() string {
	reg := regexp.MustCompile("^refs/remotes/([^/]+)")
	if reg.MatchString(b.Name) {
		return reg.FindStringSubmatch(b.Name)[1]
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
