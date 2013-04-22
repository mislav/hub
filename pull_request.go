package main

import (
	"log"
)

var cmdPullRequest = &Command{
	Run:   pullRequest,
	Usage: "pull-request [-f] [TITLE|-i ISSUE] [-b BASE] [-h HEAD]",
	Short: "Open a pull request on GitHub",
	Long: `Opens a pull request on GitHub for the project that the "origin" remote
points to. The default head of the pull request is the current branch.
Both base and head of the pull request can be explicitly given in one of
the following formats: "branch", "owner:branch", "owner/repo:branch".
This command will abort operation if it detects that the current topic
branch has local commits that are not yet pushed to its upstream branch
on the remote. To skip this check, use -f.

If TITLE is omitted, a text editor will open in which title and body of
the pull request can be entered in the same manner as git commit message.

If instead of normal TITLE an issue number is given with -i, the pull
request will be attached to an existing GitHub issue. Alternatively, instead
of title you can paste a full URL to an issue on GitHub.
`,
}

var flagPullRequestBase, flagPullRequestHead string

func init() {
	// TODO: delay calculation of owner and current branch until being used
	cmdPullRequest.Flag.StringVar(&flagPullRequestBase, "b", git.Owner()+":master", "BASE")
	cmdPullRequest.Flag.StringVar(&flagPullRequestHead, "h", git.Owner()+":"+git.CurrentBranch(), "HEAD")
}

func pullRequest(cmd *Command, args []string) {
	params := PullRequestParams{"title", "body", flagPullRequestBase, flagPullRequestHead}
	err := gh.CreatePullRequest(git.Owner(), git.Repo(), params)
	if err != nil {
		log.Fatal(err)
	}
}
