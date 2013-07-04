package commands

import (
	"fmt"
	"github.com/jingweno/gh/utils"
	"github.com/jingweno/octokat"
)

var cmdMerge = &Command{
	Run:          merge,
	GitExtension: true,
	Usage:        "merge PULLREQ-URL",
	Short:        "Join two or more development histories (branches) together",
	Long: `Merge the pull request with a commit message that includes the pull request
ID and title, similar to the GitHub Merge Button.
`,
}

/*
  $ gh merge https://github.com/jingweno/gh/pull/73
  > git fetch git://github.com/jingweno/gh.git +refs/heads/feature:refs/remotes/jingweno/feature
  > git merge jingweno/feature --no-ff -m 'Merge pull request #73 from jingweno/feature...'
*/
func merge(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		err := transformMergeArgs(args)
		utils.Check(err)
	}
}

func transformMergeArgs(args *Args) error {
	id := parsePullRequestId(args.FirstParam())
	if id != "" {
		pullRequest, err := fetchPullRequest(id)
		if err != nil {
			return err
		}

		err = fetchAndMerge(args, pullRequest)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchAndMerge(args *Args, pullRequest *octokat.PullRequest) error {
	user := pullRequest.User.Login
	branch := pullRequest.Head.Ref
	url, err := convertToGitURL(pullRequest)
	if err != nil {
		return err
	}

	args.RemoveParam(0) // Remove the pull request URL

	mergeHead := fmt.Sprintf("%s/%s", user, branch)
	ref := fmt.Sprintf("+refs/heads/%s:refs/remotes/%s", branch, mergeHead)
	args.Before("git", "fetch", url, ref)
	mergeMsg := fmt.Sprintf("Merge pull request #%v from %s\n\n%s", pullRequest.Number, mergeHead, pullRequest.Title)
	args.AppendParams(mergeHead, "--no-ff", "-m", mergeMsg)

	return nil
}
