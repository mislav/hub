package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"github.com/jingweno/octokat"
	"regexp"
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
	if !args.IsEmpty() {
		err := transformCheckoutArgs(args)
		utils.Fatal(err)
	}
}

func transformCheckoutArgs(args *Args) error {
	id := parsePullRequestId(args.First())
	if id != "" {
		url := args.Remove(0)
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

		if remoteExists {
			updateExistingRemote(args, user, branch)
		} else {
			isSSH := pullRequest.Head.Repo.Private
			sshURL, err := convertPullRequestURLToGitURL(url, user, isSSH)
			if err != nil {
				return err
			}

			addRmote(args, user, branch, sshURL)
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

func fetchPullRequest(id string) (*octokat.PullRequest, error) {
	gh := github.New()
	pullRequest, err := gh.PullRequest(id)
	if err != nil {
		return nil, err
	}

	if pullRequest.Head.Repo.ID == 0 {
		user := pullRequest.User.Login
		return nil, fmt.Errorf("%s's fork is not available anymore", user)
	}

	return pullRequest, nil
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

func convertPullRequestURLToGitURL(pullRequestURL, user string, isSSH bool) (string, error) {
	project, err := github.ParseProjectFromURL(pullRequestURL)
	if err != nil {
		return "", err
	}

	return project.GitURL("", user, isSSH), nil
}
