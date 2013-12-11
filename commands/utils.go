package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"github.com/jingweno/go-octokit/octokit"
	"os"
	"regexp"
	"strings"
)

func isDir(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func parseUserBranchFromPR(pullRequest *octokit.PullRequest) (user string, branch string) {
	userBranch := strings.SplitN(pullRequest.Head.Label, ":", 2)
	user = userBranch[0]
	if len(userBranch) > 1 {
		branch = userBranch[1]
	} else {
		branch = pullRequest.Head.Ref
	}

	return
}

func parseRepoNameOwner(nameWithOwner string) (owner, name string, match bool) {
	ownerRe := fmt.Sprintf("^(%s)$", OwnerRe)
	ownerRegexp := regexp.MustCompile(ownerRe)
	if ownerRegexp.MatchString(nameWithOwner) {
		owner = ownerRegexp.FindStringSubmatch(nameWithOwner)[1]
		match = true
		return
	}

	nameWithOwnerRe := fmt.Sprintf("^(%s)\\/(%s)$", OwnerRe, NameRe)
	nameWithOwnerRegexp := regexp.MustCompile(nameWithOwnerRe)
	if nameWithOwnerRegexp.MatchString(nameWithOwner) {
		result := nameWithOwnerRegexp.FindStringSubmatch(nameWithOwner)
		owner = result[1]
		name = result[2]
		match = true
	}

	return
}

func hasGitRemote(name string) bool {
	remotes, err := git.Remotes()
	utils.Check(err)
	for _, remote := range remotes {
		if remote.Name == name {
			return true
		}
	}

	return false
}
