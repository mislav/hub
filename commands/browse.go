package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdBrowse = &Command{
	Run:   browse,
	Usage: "browse [-u] [[<USER>/]<REPOSITORY>|--] [SUBPAGE]",
	Short: "Open a GitHub page in the default browser",
	Long: `Open repository's GitHub page in the system's default web browser using
"open(1)" or the "BROWSER" env variable. If the repository isn't
specified, "browse" opens the page of the repository found in the current
directory. If SUBPAGE is specified, the browser will open on the specified
subpage: one of "wiki", "commits", "issues" or other (the default is
"tree"). With "-u", outputs the URL rather than opening the browser.
`,
}

var (
	flagBrowseURLOnly bool
)

func init() {
	cmdBrowse.Flag.BoolVarP(&flagBrowseURLOnly, "url-only", "u", false, "URL")

	CmdRunner.Use(cmdBrowse)
}

/*
  $ gh browse
  > open https://github.com/CURRENT_REPO

  $ gh browse -- issues
  > open https://github.com/CURRENT_REPO/issues

  $ gh browse jingweno/gh
  > open https://github.com/jingweno/gh

  $ gh browse gh
  > open https://github.com/YOUR_LOGIN/gh

  $ gh browse gh wiki
  > open https://github.com/YOUR_LOGIN/gh/wiki
*/
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

		branch, project, _ = localRepo.RemoteBranchAndProject("", currentBranch.IsMaster())
		if branch == nil {
			branch = localRepo.MasterBranch()
		}
	}

	if project == nil {
		err := fmt.Errorf(command.FormattedUsage())
		utils.Check(err)
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

	pageUrl := project.WebURL("", "", path)
	launcher, err := utils.BrowserLauncher()
	utils.Check(err)

	if flagBrowseURLOnly {
		args.Replace("echo", pageUrl)
	} else {
		args.Replace(launcher[0], "", launcher[1:]...)
		args.AppendParams(pageUrl)
	}
}

func branchInURL(branch *github.Branch) string {
	parts := strings.Split(branch.ShortName(), "/")
	newPath := make([]string, len(parts))
	for i, s := range parts {
		newPath[i] = url.QueryEscape(s)
	}
	return strings.Join(newPath, "/")
}
