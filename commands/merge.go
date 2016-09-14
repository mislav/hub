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
	Usage:        "merge <PULLREQ-URL>",
	Long: `Merge a pull request with a message like the GitHub Merge Button.

## Examples:
		$ hub merge https://github.com/jingweno/gh/pull/73
		> git fetch origin refs/pull/73/head
		> git merge FETCH_HEAD --no-ff -m "Merge pull request #73 from jingweno/feature..."

## See also:

hub-checkout(1), hub(1), git-merge(1)
`,
}

func init() {
	CmdRunner.Use(cmdMerge)
}

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

	repo, err := github.LocalRepo()
	if err != nil {
		return err
	}

	remote, err := repo.RemoteForRepo(pullRequest.Base.Repo)
	if err != nil {
		return err
	}

	branch := pullRequest.Head.Ref
	headRepo := pullRequest.Head.Repo
	if headRepo == nil {
		return fmt.Errorf("Error: that fork is not available anymore")
	}

	args.Before("git", "fetch", remote.Name, fmt.Sprintf("refs/pull/%s/head", id))

	// Remove pull request URL
	idx := args.IndexOfParam(mergeURL)
	args.RemoveParam(idx)

	mergeMsg := fmt.Sprintf("Merge pull request #%s from %s/%s\n\n%s", id, headRepo.Owner.Login, branch, pullRequest.Title)
	args.AppendParams("FETCH_HEAD", "-m", mergeMsg)

	if args.IndexOfParam("--ff-only") == -1 && args.IndexOfParam("--squash") == -1 && args.IndexOfParam("--ff") == -1 {
		i := args.IndexOfParam("-m")
		args.InsertParam(i, "--no-ff")
	}

	return nil
}
