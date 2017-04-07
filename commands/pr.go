package commands

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var (
	cmdPr = &Command{
		Run:   printHelp,
		Usage: "pr checkout <PULLREQ-NUMBER> [<BRANCH>]",
		Long: `Check out the head of a pull request as a local branch.

## Examples:
	$ hub pr checkout 73
	> git fetch origin pull/73/head:jingweno-feature
	> git checkout jingweno-feature

## See also:

hub-merge(1), hub(1), hub-checkout(1)
	`,
	}

	cmdCheckoutPr = &Command{
		Key: "checkout",
		Run: checkoutPr,
	}
)

func init() {
	cmdPr.Use(cmdCheckoutPr)
	CmdRunner.Use(cmdPr)
}

func printHelp(command *Command, args *Args) {
	fmt.Print(command.HelpText())
	os.Exit(0)
}

/**
 * Add a log messsage to /tmp/hub.log
 *
 * FIXME: Should be removed before PRing this.
 */
func johan(message string) {
	f, err := os.OpenFile("/tmp/hub.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(message + "\n"); err != nil {
		panic(err)
	}
}

func checkoutPr(command *Command, args *Args) {
	johan("")

	if args.ParamsSize() < 1 || args.ParamsSize() > 2 {
		utils.Check(fmt.Errorf("Error: Expected one or two arguments, got %d", args.ParamsSize()))
	}

	prNumberString := args.GetParam(0)
	_, err := strconv.Atoi(prNumberString)
	utils.Check(err)

	// Figure out the PR URL
	localRepo, err := github.LocalRepo()
	utils.Check(err)
	baseProject, err := localRepo.MainProject()
	utils.Check(err)
	host, err := github.CurrentConfig().PromptForHost(baseProject.Host)
	utils.Check(err)
	client := github.NewClientWithHost(host)
	pr, err := client.PullRequest(baseProject, prNumberString)
	utils.Check(err)

	// Args here are: "git pr 77" or "git pr 77 new-branch-name"
	if args.ParamsSize() == 1 {
		args.Replace(args.Executable, "checkout", pr.HtmlUrl)
	} else {
		args.Replace(args.Executable, "checkout", pr.HtmlUrl, args.GetParam(1))
	}

	// Call into the checkout code which already provides the functionality we're
	// after
	err = transformPrArgs(args)
	johan("adam")
	utils.Check(err)
	johan("frukten")
}

func transformPrArgs(args *Args) error {
	// This function initially copied from checkout.go:transformCheckoutArgs()
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

	baseRemote, err := repo.RemoteForRepo(pullRequest.Base.Repo)
	if err != nil {
		return err
	}

	var headRemote *github.Remote
	if pullRequest.IsSameRepo() {
		headRemote = baseRemote
	} else if pullRequest.Head.Repo != nil {
		headRemote, _ = repo.RemoteForRepo(pullRequest.Head.Repo)
	}

	var newArgs []string

	if headRemote != nil {
		if newBranchName == "" {
			newBranchName = pullRequest.Head.Ref
		}
		remoteBranch := fmt.Sprintf("%s/%s", headRemote.Name, pullRequest.Head.Ref)
		refSpec := fmt.Sprintf("+refs/heads/%s:refs/remotes/%s", pullRequest.Head.Ref, remoteBranch)
		if git.HasFile("refs", "heads", newBranchName) {
			newArgs = append(newArgs, newBranchName)
			args.After("git", "merge", "--ff-only", fmt.Sprintf("refs/remotes/%s", remoteBranch))
		} else {
			newArgs = append(newArgs, "-b", newBranchName, "--track", remoteBranch)
		}
		args.Before("git", "fetch", headRemote.Name, refSpec)
	} else {
		if newBranchName == "" {
			if pullRequest.Head.Repo == nil {
				newBranchName = fmt.Sprintf("pr-%s", id)
			} else {
				newBranchName = fmt.Sprintf("%s-%s", pullRequest.Head.Repo.Owner.Login, pullRequest.Head.Ref)
			}
		}
		newArgs = append(newArgs, newBranchName)

		ref := fmt.Sprintf("refs/pull/%s/head", id)
		args.Before("git", "fetch", baseRemote.Name, fmt.Sprintf("%s:%s", ref, newBranchName))

		remote := baseRemote.Name
		mergeRef := ref
		if pullRequest.MaintainerCanModify && pullRequest.Head.Repo != nil {
			project, projectErr := github.NewProjectFromRepo(pullRequest.Head.Repo)
			if projectErr != nil {
				return projectErr
			}

			remote = project.GitURL("", "", true)
			mergeRef = fmt.Sprintf("refs/heads/%s", pullRequest.Head.Ref)
		}
		args.Before("git", "config", fmt.Sprintf("branch.%s.remote", newBranchName), remote)
		args.Before("git", "config", fmt.Sprintf("branch.%s.merge", newBranchName), mergeRef)
	}
	replaceCheckoutParam(args, checkoutURL, newArgs...)

	johan("adam: " + args.ToCmd().String())
	return nil
}
