package commands

import (
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
)

var cmdBrowse = &Command{
	Run:   browse,
	Usage: "browse [-u USER] [-r REPOSITORY] [SUBPAGE]",
	Short: "Open a GitHub page in the default browser",
	Long: `Open repository's GitHub page in the system's default web browser using
open(1) or the BROWSER env variable. If the repository isn't specified,
browse opens the page of the repository found in the current directory.
If SUBPAGE is specified, the browser will open on the specified subpage:
one of "wiki", "commits", "issues" or other (the default is "tree").
`,
}

var flagBrowseUser, flagBrowseRepo string

func init() {
	cmdBrowse.Flag.StringVar(&flagBrowseUser, "u", "", "USER")
	cmdBrowse.Flag.StringVar(&flagBrowseRepo, "r", "", "REPOSITORY")
}

/*
  $ gh browse
  > open https://github.com/YOUR_USER/CURRENT_REPO

  $ gh browse commit/SHA
  > open https://github.com/YOUR_USER/CURRENT_REPO/commit/SHA

  $ gh browse issues
  > open https://github.com/YOUR_USER/CURRENT_REPO/issues

  $ gh browse -u jingweno -r gh
  > open https://github.com/jingweno/gh

  $ gh browse -u jingweno -r gh commit/SHA
  > open https://github.com/jingweno/gh/commit/SHA

  $ gh browse -r resque
  > open https://github.com/YOUR_USER/resque

  $ gh browse -r resque network
  > open https://github.com/YOUR_USER/resque/network
*/
func browse(command *Command, args *Args) {
	subpage := "tree"
	if !args.IsParamsEmpty() {
		subpage = args.RemoveParam(0)
	}

	project := github.CurrentProject()
	if subpage == "tree" || subpage == "commits" {
		repo := project.LocalRepo()
		subpage = utils.ConcatPaths(subpage, repo.Head)
	}

	url := project.WebURL(flagBrowseRepo, flagBrowseUser, subpage)
	launcher, err := utils.BrowserLauncher()
	if err != nil {
		utils.Check(err)
	}

  args.Replace(launcher[0], "", launcher[1:]...)
	args.AppendParams(url)
}
