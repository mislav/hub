package commands

import (
	"fmt"
	"regexp"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
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

func init() {
	CmdRunner.Use(cmdMerge)
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
	words := args.Words()
	if len(words) == 0 {
		return nil
	}

	mergeURL := words[0]
	url, err := github.ParseURL(mergeURL)
	if err != nil {
		return nil
	}

	pullURLRegex := regexp.MustCompile("^pull/(\\d+)")
	projectPath := url.ProjectPath()
	if !pullURLRegex.MatchString(projectPath) {
		return nil
	}

	id := pullURLRegex.FindStringSubmatch(projectPath)[1]
	gh := github.NewClient(url.Project.Host)
	pullRequest, err := gh.PullRequest(url.Project, id)
	if err != nil {
		return err
	}

	branch := pullRequest.Head.Ref
	headRepo := pullRequest.Head.Repo
	if headRepo == nil {
		return fmt.Errorf("Error: that fork is not available anymore")
	}

	u := url.GitURL(headRepo.Name, headRepo.Owner.Login, headRepo.Private)
	mergeHead := fmt.Sprintf("%s/%s", headRepo.Owner.Login, branch)
	ref := fmt.Sprintf("+refs/heads/%s:refs/remotes/%s", branch, mergeHead)
	args.Before("git", "fetch", u, ref)

	// Remove pull request URL
	idx := args.IndexOfParam(mergeURL)
	args.RemoveParam(idx)

	mergeMsg := fmt.Sprintf("Merge pull request #%v from %s\n\n%s", id, mergeHead, pullRequest.Title)
	args.AppendParams(mergeHead, "-m", mergeMsg)

	if args.IndexOfParam("--ff-only") == -1 && args.IndexOfParam("--squash") == -1 && args.IndexOfParam("--ff") == -1 {
		i := args.IndexOfParam("-m")
		args.InsertParam(i, "--no-ff")
	}

	return nil
}
