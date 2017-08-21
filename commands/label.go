package commands

import (
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)


var (
	cmdLabel = &Command{
		Run: listLabels,
		Usage: `
label
`,
		Long: `Manage GitHub labels for the current repository.

## Commands:

With no arguments, show a list of open issues.
`,
	}
)

func init() {
	CmdRunner.Use(cmdLabel)
}

func listLabels(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

    if args.Noop {
		ui.Printf("Would request list of labels for %s\n", project)
	} else {
		labels, err := gh.FetchLabels(project)
		utils.Check(err)

		for _, label := range labels {
			ui.Printf(formatLabel(label))
		}
	}

	args.NoForward()
}

func formatLabel(label github.IssueLabel) string {
    format := "%l%n"

	placeholders := map[string]string{
		"l":  label.Name,
	}

	return ui.Expand(format, placeholders, false)
}
