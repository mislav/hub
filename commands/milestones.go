package commands

import (
	"fmt"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var (
	cmdMilestone = &Command{
		// TODO: Add usage information, eventually get rid of `Key: "milestone",`
		Key: "milestone",
		Run: listMilestones,
	}

	cmdCreateMilestone = &Command{
		Key: "create",
		Run: createMilestone,
	}

	// TODO: Add update, delete, get subcommands
)

func init() {
	cmdMilestone.Use(cmdCreateMilestone)
	CmdRunner.Use(cmdMilestone)
}

func listMilestones(cmd *Command, args *Args) {
	return
}

func createMilestone(cmd *Command, args *Args) {
	// TODO: Add flags for title, description, state, due_on date
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.CurrentProject()
	utils.Check(err)
	gh := github.NewClient(project.Host)

	messageBuilder := &github.MessageBuilder{
		Filename: "MILESTONE_EDITMSG",
		Title:    "milestone",
	}
	messageBuilder.AddCommentedSection(fmt.Sprintf(`Creating milestone for %s

Write a description for this milestone. The first block of
text is the title and the rest is the description `, project))
	messageBuilder.Edit = true
	title, body, err := messageBuilder.Extract()
	messageBuilder.Cleanup()
	utils.Check(err)

	if title == "" {
		utils.Check(fmt.Errorf("Aborting milestone creation due to empty milestone title"))
	}
	params := &github.MilestoneParams{
		Title:       title,
		Description: body,
		State:       "open",
		DueOn:       "",
	}
	var milestone *github.Milestone
	args.NoForward()
	milestone, err = gh.CreateMilestone(project, params)
	utils.Check(err)

	printBrowseOrCopy(args, milestone.HtmlUrl, false, false)
}
