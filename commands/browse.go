package commands

import (
	"fmt"
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
	"net/url"
	"reflect"
	"strings"
)

var cmdBrowse = &Command{
	Run:   browse,
	Usage: "browse [-u] [-p] [[<USER>/]<REPOSITORY>] [SUBPAGE]",
	Short: "Open a GitHub page in the default browser",
	Long: `Open repository's GitHub page in the system's default web browser using
"open(1)" or the "BROWSER" env variable. If the repository isn't
specified with "-p", "browse" opens the page of the repository found in the current
directory. If SUBPAGE is specified, the browser will open on the specified
subpage: one of "wiki", "commits", "issues" or other (the default is
"tree"). With "-u", outputs the URL rather than opening the browser.
`,
}

var (
	flagBrowseURLOnly bool
	flagBrowseProject string
)

func init() {
	cmdBrowse.Flag.BoolVarP(&flagBrowseURLOnly, "url-only", "u", false, "URL only")
	cmdBrowse.Flag.StringVarP(&flagBrowseProject, "project", "p", "", "PROJECT")

	CmdRunner.Use(cmdBrowse)
}

/*
  $ gh browse
  > open https://github.com/YOUR_USER/CURRENT_REPO

  $ gh browse commit/SHA
  > open https://github.com/YOUR_USER/CURRENT_REPO/commit/SHA

  $ gh browse issues
  > open https://github.com/YOUR_USER/CURRENT_REPO/issues

  $ gh browse -p jingweno/gh
  > open https://github.com/jingweno/gh

  $ gh browse -p jingweno/gh commit/SHA
  > open https://github.com/jingweno/gh/commit/SHA

  $ gh browse -p resque
  > open https://github.com/YOUR_USER/resque

  $ gh browse -p resque network
  > open https://github.com/YOUR_USER/resque/network
*/
func browse(command *Command, args *Args) {
	var (
		project *github.Project
		branch  *github.Branch
		err     error
	)
	localRepo := github.LocalRepo()
	if flagBrowseProject != "" {
		// gh browse -p jingweno/gh
		// gh browse -p gh
		project = github.NewProject("", flagBrowseProject, "")
	} else {
		// gh browse
		branch, project, err = localRepo.RemoteBranchAndProject("")
		utils.Check(err)
	}

	if project == nil {
		err := fmt.Errorf(command.FormattedUsage())
		utils.Check(err)
	}

	master := localRepo.MasterBranch()
	if branch == nil {
		branch = master
	}

	var subpage string
	if !args.IsParamsEmpty() {
		subpage = args.RemoveParam(0)
	}

	if subpage == "commits" {
		subpage = fmt.Sprintf("commits/%s", branchInURL(branch))
	} else if subpage == "tree" || subpage == "" {
		if !reflect.DeepEqual(branch, master) && branch.IsRemote() {
			subpage = fmt.Sprintf("tree/%s", branchInURL(branch))
		}
	}

	pageUrl := project.WebURL("", "", subpage)
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
	parts := strings.Split(strings.Replace(branch.ShortName(), ".", "/", -1), "/")
	newPath := make([]string, len(parts))
	for i, s := range parts {
		newPath[i] = url.QueryEscape(s)
	}
	return strings.Join(newPath, "/")
}
