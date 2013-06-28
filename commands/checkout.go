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
func checkout(command *Command, args *Args) {
	var err error
	if !args.IsEmpty() {
		err = transformCheckoutArgs(args)
		utils.Fatal(err)
	}
}

func transformCheckoutArgs(args *Args) error {
	id := parsePullRequestId(args.First())
	if id != "" {
		url := args.Remove(0)
		gh := github.New()
		pullRequest, err := gh.PullRequest(id)
		if err != nil {
			return err
		}

		user := pullRequest.User.Login
		branch := pullRequest.Head.Ref
		if pullRequest.Head.Repo.ID == 0 {
			return fmt.Errorf("%s's fork is not available anymore", user)
		}

		remoteExists, err := checkIfRemoteExists(user)
		if err != nil {
			return err
		}

		if remoteExists {
			updateExistingRemote(args, user, branch)
		} else {
			err = addRmote(args, user, branch, url, pullRequest.Head.Repo.Private)
			if err != nil {
				return err
			}
		}

		var newBranchName string
		if args.Size() > 0 {
			newBranchName = args.Remove(0)
		} else {
			newBranchName = fmt.Sprintf("%s-%s", user, branch)
		}
		trackedBranch := fmt.Sprintf("%s/%s", user, branch)

		args.Append("--track", "-B", newBranchName, trackedBranch)

		return nil
	}

	return nil
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

func updateExistingRemote(args *Args, user, branch string) {
	args.Before("git", "remote", "set-branches", "--add", user, branch)
	remoteURL := fmt.Sprintf("+refs/heads/%s:refs/remotes/%s/%s", branch, user, branch)
	args.Before("git", "fetch", user, remoteURL)
}

func addRmote(args *Args, user, branch, url string, isPrivate bool) error {
	project, err := github.ParseProjectFromURL(url)
	if err != nil {
		return err
	}

	sshURL := project.GitURL("", user, isPrivate)
	args.Before("git", "remote", "add", "-f", "-t", branch, user, sshURL)

	return nil
}
