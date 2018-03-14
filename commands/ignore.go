package commands

import (
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdIgnore = &Command{
		Run:   ignore,
		Usage: "ignore [-l <LANGUAGE>]",
		Long: `Show available templates from the GitHub .gitignore repository

## Options:
	-l, --language <LANGUAGE>
		Show .gitignore template for <LANGUAGE>
`,
	}

	flagIgnoreLanguage string
)

func init() {
	cmdIgnore.Flag.StringVarP(&flagIgnoreLanguage, "language", "l", "", "LANGUAGE")
	CmdRunner.Use(cmdIgnore)
}

func ignore(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	var gitignore *github.Gitignore
	gitignore, err = gh.Gitignore(flagIgnoreLanguage)
	utils.Check(err)

	ui.Println(gitignore.Source)

	args.NoForward()
}
