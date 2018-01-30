package commands

import (
	"fmt"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdLabel = &Command{
		Key:   "labels",
		Run:   listLabels,
		Usage: "issue labels",
		Long:  "List the labels available in this repository.",
	}
)

func listLabels(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	if args.Noop {
		ui.Printf("Would request list of labels for %s\n", project)
		return
	}

	labels, err := gh.FetchLabels(project)
	utils.Check(err)

	for _, label := range labels {
		ui.Printf(formatLabel(label, true))
	}

	args.NoForward()
}

func formatLabel(label github.IssueLabel, colorize bool) string {
	format := "%l%n"
	if colorize {
		format = "%c%n"
	}

	color, err := utils.NewColor(label.Color)
	if err != nil {
		utils.Check(err)
	}

	placeholders := map[string]string{
		"l": label.Name,
		"c": fmt.Sprintf("\033[38;5;%d;48;2;%d;%d;%dm %s \033[m",
			getSuitableTextColor(color), color.Red, color.Green, color.Blue, label.Name),
	}

	return ui.Expand(format, placeholders, colorize)
}

func getSuitableTextColor(color *utils.Color) int {
	if color.Brightness() < 0.65 {
		return 15 // white text
	}
	return 16 // black text
}
