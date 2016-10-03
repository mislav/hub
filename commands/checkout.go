package commands

import (
	"fmt"
	"regexp"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdCheckout = &Command{
	Run:          checkout,
	GitExtension: true,
	Usage:        "checkout <PULLREQ-URL> [<BRANCH>]",
	Long: `Check out the head of a pull request as a local branch.

## Examples:
		$ hub checkout https://github.com/jingweno/gh/pull/73
		> git remote add -f --no-tags -t feature git://github:com/jingweno/gh.git
		> git checkout --track -B jingweno-feature jingweno/feature

## See also:

hub-merge(1), hub-am(1), hub(1), git-checkout(1)
`,
}

func init() {
	CmdRunner.Use(cmdCheckout)
}

func checkout(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		err := transformCheckoutArgs(args)
		utils.Check(err)
	}
}

func transformCheckoutArgs(args *Args) error {
	words := args.Words()

	if len(words) == 0 {
		return nil
	}

	checkoutURL := words[0]
	var newBranchName string
	if len(words) > 1 {
		newBranchName = words[1]
	}

	url, err := github.ParseURL(checkoutURL)
	if err != nil {
		// not a valid GitHub URL
		return nil
	}

	pullURLRegex := regexp.MustCompile("^pull/(\\d+)")
	projectPath := url.ProjectPath()
	if !pullURLRegex.MatchString(projectPath) {
		// not a valid PR URL
		return nil
	}

	err = sanitizeCheckoutFlags(args)
	if err != nil {
		return err
	}

	id := pullURLRegex.FindStringSubmatch(projectPath)[1]
	gh := github.NewClient(url.Project.Host)
	pullRequest, err := gh.PullRequest(url.Project, id)
	if err != nil {
		return err
	}

	if idx := args.IndexOfParam(newBranchName); idx >= 0 {
		args.RemoveParam(idx)
	}

	repo, err := github.LocalRepo()
	if err != nil {
		return err
	}

	remote, err := repo.RemoteForRepo(pullRequest.Base.Repo)
	if err != nil {
		return err
	}

	var refSpec string
	var newArgs []string

	if pullRequest.IsSameRepo() {
		if newBranchName == "" {
			newBranchName = pullRequest.Head.Ref
		}
		remoteBranch := fmt.Sprintf("%s/%s", remote.Name, pullRequest.Head.Ref)
		refSpec = fmt.Sprintf("+refs/heads/%s:refs/remotes/%s", pullRequest.Head.Ref, remoteBranch)
		newArgs = append(newArgs, "-b", newBranchName, "--track", remoteBranch)
	} else {
		if newBranchName == "" {
			if pullRequest.Head.Repo == nil {
				newBranchName = fmt.Sprintf("pr-%s", id)
			} else {
				newBranchName = fmt.Sprintf("%s-%s", pullRequest.Head.Repo.Owner.Login, pullRequest.Head.Ref)
			}
		}
		refSpec = fmt.Sprintf("pull/%s/head:%s", id, newBranchName)
		newArgs = append(newArgs, newBranchName)
	}

	args.Before("git", "fetch", remote.Name, refSpec)
	replaceCheckoutParam(args, checkoutURL, newArgs...)
	return nil
}

func sanitizeCheckoutFlags(args *Args) error {
	if i := args.IndexOfParam("-b"); i != -1 {
		return fmt.Errorf("Unsupported flag -b when checking out pull request")
	}

	if i := args.IndexOfParam("--orphan"); i != -1 {
		return fmt.Errorf("Unsupported flag --orphan when checking out pull request")
	}

	return nil
}

func replaceCheckoutParam(args *Args, checkoutURL string, replacement ...string) {
	idx := args.IndexOfParam(checkoutURL)
	args.RemoveParam(idx)
	args.InsertParam(idx, replacement...)
}
