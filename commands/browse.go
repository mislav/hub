package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
)

var cmdBrowse = &Command{
	Run:   browse,
	Usage: "browse [[USER]/REPOSITORY] [SUBPAGE]",
	Short: "Open a GitHub page in the default browser",
	Long: `Open repository's GitHub page in the system's default web browser using
open(1) or the BROWSER env variable. If the repository isn't specified,
browse opens the page of the repository found in the current directory.
If SUBPAGE is specified, the browser will open on the specified subpage:
one of "wiki", "commits", "issues" or other (the default is "tree").
`,
}

func browse(cmd *Command, args []string) {
	var subpage string
	project := github.CurrentProject()
	ownerWithName := project.OwnerWithName()

	if len(args) == 1 {
		ownerWithName = utils.CatPaths(project.Owner, args[0])
	} else if len(args) == 2 {
		if args[0] != "--" {
			ownerWithName = args[0]
		}
		subpage = args[1]
	} else if len(args) > 2 {
		cmd.PrintUsage()
	}

	url := project.WebUrl(ownerWithName, subpage)

	fmt.Println(url)
}
