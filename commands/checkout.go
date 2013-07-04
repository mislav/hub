package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
)

var cmdCheckout = &Command{
	Run:          checkout,
	GitExtension: true,
	Usage:        "checkout PULLREQ-URL [BRANCH]",
	Short:        "Switch the active branch to another branch",
	Long: `Checks out the head of the pull request as a local branch, to allow for
reviewing, rebasing and otherwise cleaning up the commits in the pull
request before merging. The name of the local branch can explicitly be
set with BRANCH.
`,
}

/**
  $ gh checkout https://github.com/jingweno/gh/pull/73
  > git remote add -f -t feature git://github:com/foo/gh.git
  > git checkout --track -B foo-feature foo/feature

  $ gh checkout https://github.com/jingweno/gh/pull/73 custom-branch-name
**/
func checkout(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		err := transformCheckoutArgs(args)
		utils.Check(err)
	}
}

func transformCheckoutArgs(args *Args) error {
	id := parsePullRequestId(args.FirstParam())
	if id != "" {
		pullRequest, err := fetchPullRequest(id)
		if err != nil {
			return err
		}

		user := pullRequest.User.Login
		branch := pullRequest.Head.Ref

		remoteExists, err := checkIfRemoteExists(user)
		if err != nil {
			return err
		}

		args.RemoveParam(0) // Remove the pull request URL
		if remoteExists {
			updateExistingRemote(args, user, branch)
		} else {
			sshURL, err := convertToGitURL(pullRequest)
			if err != nil {
				return err
			}

			addRmote(args, user, branch, sshURL)
		}

		var newBranchName string
		if args.ParamsSize() > 0 {
			newBranchName = args.RemoveParam(0)
		} else {
			newBranchName = fmt.Sprintf("%s-%s", user, branch)
		}
		trackedBranch := fmt.Sprintf("%s/%s", user, branch)

		args.AppendParams("--track", "-B", newBranchName, trackedBranch)
	}

	return nil
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

func addRmote(args *Args, user, branch, sshURL string) {
	args.Before("git", "remote", "add", "-f", "-t", branch, user, sshURL)
}
