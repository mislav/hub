package commands

import (
	"fmt"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdFork = &Command{
	Run:   fork,
	Usage: "fork [--no-remote]",
	Long: `Fork the current project on GitHub and add a git remote for it.

## Options:
	--no-remote
		Skip adding a git remote for the fork.

## Examples:
		$ hub fork
		[ repo forked on GitHub ]
		> git remote add -f USER git@github.com:USER/REPO.git

## See also:

hub-clone(1), hub(1)
`,
}

var flagForkNoRemote bool

func init() {
	cmdFork.Flag.BoolVar(&flagForkNoRemote, "no-remote", false, "")

	CmdRunner.Use(cmdFork)
}

func fork(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	if err != nil {
		utils.Check(fmt.Errorf("Error: repository under 'origin' remote is not a GitHub project"))
	}

	config := github.CurrentConfig()
	host, err := config.PromptForHost(project.Host)
	if err != nil {
		utils.Check(github.FormatError("forking repository", err))
	}

	originRemote, err := localRepo.OriginRemote()
	if err != nil {
		utils.Check(fmt.Errorf("Error creating fork: %s", err))
	}

	forkProject := github.NewProject(host.User, project.Name, project.Host)
	newRemoteName := forkProject.Owner

	client := github.NewClient(project.Host)
	existingRepo, err := client.Repository(forkProject)
	if err == nil {
		var parentURL *github.URL
		if parent := existingRepo.Parent; parent != nil {
			parentURL, _ = github.ParseURL(parent.HTMLURL)
		}
		if parentURL == nil || !project.SameAs(parentURL.Project) {
			err = fmt.Errorf("Error creating fork: %s already exists on %s",
				forkProject, forkProject.Host)
			utils.Check(err)
		}
	} else {
		if !args.Noop {
			newRepo, err := client.ForkRepository(project)
			utils.Check(err)
			forkProject.Owner = newRepo.Owner.Login
			forkProject.Name = newRepo.Name
		}
	}

	args.NoForward()
	if !flagForkNoRemote {
		originURL := originRemote.URL.String()
		url := forkProject.GitURL("", "", true)
		args.Before("git", "remote", "add", "-f", newRemoteName, originURL)
		args.Before("git", "remote", "set-url", newRemoteName, url)

		args.AfterFn(func() error {
			ui.Printf("new remote: %s\n", newRemoteName)
			return nil
		})
	}
}
