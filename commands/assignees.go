package commands

import (
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdAssignees = &Command{
  Run: assignees,
  Usage: "assignees",
  Long: `Show all potential assignees for the current project.`,
}

func init() {
	CmdRunner.Use(cmdAssignees)
}

func assignees(command *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	if args.Noop {
		ui.Printf("Would request list of issues for %s\n", project)
	} else {
    assignees, err := gh.FetchAssignees(project)
    utils.Check(err)

    for _, assignee := range assignees {
      ui.Println(assignee.Login)
    }
  }

	args.NoForward()
}
