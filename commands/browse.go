package commands

import (
	"github.com/jingweno/gh/cmd"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
)

var cmdBrowse = &Command{
	Run:   browse,
	Usage: "browse [-u USER] [-r REPOSITORY] [-p SUBPAGE]",
	Short: "Open a GitHub page in the default browser",
	Long: `Open repository's GitHub page in the system's default web browser using
open(1) or the BROWSER env variable. If the repository isn't specified,
browse opens the page of the repository found in the current directory.
If SUBPAGE is specified, the browser will open on the specified subpage:
one of "wiki", "commits", "issues" or other (the default is "tree").
`,
}

var flagBrowseUser, flagBrowseRepo, flagBrowseSubpage string

func init() {
	cmdBrowse.Flag.StringVar(&flagBrowseUser, "u", "", "USER")
	cmdBrowse.Flag.StringVar(&flagBrowseRepo, "r", "", "REPOSITORY")
	cmdBrowse.Flag.StringVar(&flagBrowseSubpage, "p", "", "SUBPAGE")
}

func browse(command *Command, args []string) {
	project := github.CurrentProject()
	if flagBrowseUser == "" {
		flagBrowseUser = project.Owner
	}
	if flagBrowseRepo == "" {
		flagBrowseRepo = project.Name
	}

	ownerWithName := utils.ConcatPaths(flagBrowseUser, flagBrowseRepo)

	url := project.WebUrl(ownerWithName, flagBrowseSubpage)
	browser, err := utils.BrowserLauncher()
	utils.Check(err)

	browser = append(browser, url)
	c := cmd.NewWithArray(browser)
	err = c.Exec()
	utils.Check(err)
}
