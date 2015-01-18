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
	Usage:        "checkout PULLREQ-URL|BRANCH-URL [BRANCH]",
	Short:        "Switch the active branch to another branch",
	Long: `Checks out the head of the pull request or branch as a local branch, to allow for
reviewing, rebasing and otherwise cleaning up the commits in the pull
request or branch before merging. The name of the local branch can explicitly be
set with BRANCH.
`,
}

func init() {
	CmdRunner.Use(cmdCheckout)
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

	treeURLRegex := regexp.MustCompile("^tree/(.+)")
	pullURLRegex := regexp.MustCompile("^pull/(\\d+)")
	projectPath := url.ProjectPath()
	
	user := ""
	branch := ""
	isPrivateRepo := false

	if treeURLRegex.MatchString(projectPath) {
		tree := treeURLRegex.FindStringSubmatch(projectPath)[1]
		branch = tree
		user = url.User()
	} else if !pullURLRegex.MatchString(projectPath) {
		id := pullURLRegex.FindStringSubmatch(projectPath)[1]
		gh := github.NewClient(url.Project.Host)
		pullRequest, err := gh.PullRequest(url.Project, id)
		if err != nil {
			return err
		}
		
		user, branch = parseUserBranchFromPR(pullRequest)
		if pullRequest.Head.Repo == nil {
			return fmt.Errorf("Error: %s's fork is not available anymore", user)
		}

		isPrivateRepo = pullRequest.Head.Repo.Private
	} else {
		// not a valid PR URL
		return nil

	}

	
	if idx := args.IndexOfParam(newBranchName); idx >= 0 {
		args.RemoveParam(idx)
	}

	
	if newBranchName == "" {
		newBranchName = fmt.Sprintf("%s-%s", user, branch)
	}

	repo, err := github.LocalRepo()
	utils.Check(err)

	_, err = repo.RemoteByName(user)
	if err == nil {
		args.Before("git", "remote", "set-branches", "--add", user, branch)
		remoteURL := fmt.Sprintf("+refs/heads/%s:refs/remotes/%s/%s", branch, user, branch)
		args.Before("git", "fetch", user, remoteURL)
	} else {
		u := url.Project.GitURL("", user, isPrivateRepo)
		args.Before("git", "remote", "add", "-f", "-t", branch, user, u)
	}

	idx := args.IndexOfParam(checkoutURL)
	args.RemoveParam(idx)
	args.InsertParam(idx, "--track", "-B", newBranchName, fmt.Sprintf("%s/%s", user, branch))

	return nil
}
