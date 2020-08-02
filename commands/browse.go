package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/utils"
)

var cmdBrowse = &Command{
	Run:   browse,
	Usage: "browse [-uc] [[<USER>/]<REPOSITORY>|--] [<SUBPAGE>]",
	Long: `Open a GitHub repository in a web browser.

## Options:
	-u, --url
		Print the URL instead of opening it.

	-c, --copy
		Put the URL in clipboard instead of opening it.

	[<USER>/]<REPOSITORY>
		Defaults to repository in the current working directory.

	<SUBPAGE>
		One of "wiki", "commits", "issues", or other (default: "tree").

## Examples:
		$ hub browse
		> open https://github.com/REPO

		$ hub browse -- issues
		> open https://github.com/REPO/issues

		$ hub browse jingweno/gh
		> open https://github.com/jingweno/gh

		$ hub browse gh wiki
		> open https://github.com/USER/gh/wiki

## See also:

hub-compare(1), hub(1)
`,
}

func init() {
	CmdRunner.Use(cmdBrowse)
}

func browse(command *Command, args *Args) {
	var (
		dest    string
		subpage string
		path    string
		project *github.Project
		branch  *github.Branch
		err     error
	)

	if !args.IsParamsEmpty() {
		dest = args.RemoveParam(0)
	}

	if !args.IsParamsEmpty() {
		subpage = args.RemoveParam(0)
	}

	if args.Terminator {
		subpage = dest
		dest = ""
	}

	localRepo, _ := github.LocalRepo()
	if dest != "" {
		project = github.NewProject("", dest, "")
		branch = localRepo.MasterBranch()
	} else if subpage != "" && subpage != "commits" && subpage != "tree" && subpage != "blob" && subpage != "settings" {
		project, err = localRepo.MainProject()
		branch = localRepo.MasterBranch()
		utils.Check(err)
	} else {
		currentBranch, err := localRepo.CurrentBranch()
		if err != nil {
			currentBranch = localRepo.MasterBranch()
		}

		var owner string
		mainProject, err := localRepo.MainProject()
		if err == nil {
			host, err := github.CurrentConfig().PromptForHost(mainProject.Host)
			if err != nil {
				utils.Check(github.FormatError("in browse", err))
			} else {
				owner = host.User
			}
		}

		branch, project, _ = localRepo.RemoteBranchAndProject(owner, currentBranch.IsMaster())
		if branch == nil {
			branch = localRepo.MasterBranch()
		}
	}

	if project == nil {
		utils.Check(command.UsageError(""))
	}

	if subpage == "commits" {
		path = fmt.Sprintf("commits/%s", branchInURL(branch))
	} else if subpage == "tree" || subpage == "" {
		if !branch.IsMaster() {
			path = fmt.Sprintf("tree/%s", branchInURL(branch))
		}
	} else {
		path = subpage
	}

	pageURL := project.WebURL("", "", path)

	args.NoForward()
	flagBrowseURLPrint := args.Flag.Bool("--url")
	flagBrowseURLCopy := args.Flag.Bool("--copy")
	printBrowseOrCopy(args, pageURL, !flagBrowseURLPrint && !flagBrowseURLCopy, flagBrowseURLCopy)
}

func branchInURL(branch *github.Branch) string {
	parts := strings.Split(branch.ShortName(), "/")
	newPath := make([]string, len(parts))
	for i, s := range parts {
		newPath[i] = url.QueryEscape(s)
	}
	return strings.Join(newPath, "/")
}
