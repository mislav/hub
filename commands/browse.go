package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
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
	cmdBrowse.Flag.BoolVar(&flagBrowseURLOnly, "u", false, "URL only")
	cmdBrowse.Flag.StringVar(&flagBrowseProject, "p", "", "PROJECT")
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
	)
	localRepo := github.LocalRepo()
	if flagBrowseProject != "" {
		// gh browse -p jingweno/gh
		// gh browse -p gh
		project = github.NewProject("", flagBrowseProject, "")
	} else {
		// gh browse
		p, err := localRepo.CurrentProject()
		utils.Check(err)
		project = p

		currentBranch, err := localRepo.CurrentBranch()
		utils.Check(err)

		upstream, err := currentBranch.Upstream()
		if err == nil {
			branch = upstream
		}
	}

	if project == nil {
		err := fmt.Errorf(command.FormattedUsage())
		utils.Check(err)
	}

	master, err := localRepo.MasterBranch()
	utils.Check(err)
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
		if !reflect.DeepEqual(branch, master) {
			subpage = fmt.Sprintf("tree/%s", branchInURL(branch))
		}
	}

	url := project.WebURL("", "", subpage)
	launcher, err := utils.BrowserLauncher()
	utils.Check(err)

	if flagBrowseURLOnly {
		args.Replace("echo", url)
	} else {
		args.Replace(launcher[0], "", launcher[1:]...)
		args.AppendParams(url)
	}
}

func branchInURL(branch *github.Branch) string {
	return strings.Replace(branch.ShortName(), ".", "/", -1)
}
