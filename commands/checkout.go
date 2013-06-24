package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"regexp"
)

var cmdCheckout = &Command{
	Run:          checkout,
	GitExtension: true,
	Usage:        "checkout PULLREQ-URL [BRANCH]",
	Short:        "Switch the active branch to another branch",
}

/**
  $ gh checkout https://github.com/jingweno/gh/pull/73
  # > git remote add -f -t feature git://github:com/foo/gh.git
  # > git checkout --track -B foo-feature foo/feature

  $ gh checkout https://github.com/jingweno/gh/pull/73 custom-branch-name
**/
func checkout(command *Command, args []string) {
	var err error
	if len(args) > 0 {
		args, err = transformCheckoutArgs(args)
		utils.Fatal(err)
	}

	err = git.SysExec("checkout", args...)
	utils.Check(err)
}

func transformCheckoutArgs(args []string) ([]string, error) {
	id := parsePullRequestId(args[0])
	if id != "" {
		newArgs, url := removeItem(args, 0)
		gh := github.New()
		pullRequest, err := gh.PullRequest(id)
		if err != nil {
			return nil, err
		}

		user := pullRequest.User.Login
		branch := pullRequest.Head.Ref
		if pullRequest.Head.Repo.ID == 0 {
			return nil, fmt.Errorf("%s's fork is not available anymore", user)
		}

		remoteExists, err := checkIfRemoteExists(user)
		if err != nil {
			return nil, err
		}

		if remoteExists {
			err = updateExistingRemote(user, branch)
		} else {
			err = addRmote(user, branch, url, pullRequest.Head.Repo.Private)
		}
		if err != nil {
			return nil, err
		}

		var newBranchName string
		if len(newArgs) > 0 {
			newArgs, newBranchName = removeItem(newArgs, 0)
		} else {
			newBranchName = fmt.Sprintf("%s-%s", user, branch)
		}
		trackedBranch := fmt.Sprintf("%s/%s", user, branch)

		newArgs = append(newArgs, "--track", "-B", newBranchName, trackedBranch)

		return newArgs, nil
	}

	return args, nil
}

func parsePullRequestId(url string) string {
	pullURLRegex := regexp.MustCompile("https://github\\.com/.+/.+/pull/(\\d+)")
	if pullURLRegex.MatchString(url) {
		return pullURLRegex.FindStringSubmatch(url)[1]
	}

	return ""
}

func checkIfRemoteExists(remote string) (bool, error) {
	remotes, err := git.Remotes()
	if err != nil {
		return false, err
	}

	for _, r := range remotes {
		if r.Name == remote {
			return true, nil
		}
	}

	return false, nil
}

func updateExistingRemote(user, branch string) error {
	err := git.Spawn("remote", "set-branches", "--add", user, branch)
	if err != nil {
		return err
	}
	remoteURL := fmt.Sprintf("+refs/heads/%s:refs/remotes/%s/%s", branch, user, branch)

	return git.Spawn("fetch", user, remoteURL)
}

func addRmote(user, branch, url string, isPrivate bool) error {
	project, err := github.ParseProjectFromURL(url)
	if err != nil {
		return err
	}

	sshURL := project.GitURL("", user, isPrivate)

	return git.Spawn("remote", "add", "-f", "-t", branch, user, sshURL)
}
