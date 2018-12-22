package commands

import (
	"github.com/github/hub/ui"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var (
	flagQualified bool
)

var (
	cmdInternalListRepos = &Command{
		Run: internalListRepos,
		Usage: "internal-list-repos [-q] [--qualified] [<OWNER>]",
		Long: `List GitHub repos for a given owner.

This command is for hub's internal use only. It may change or go away 
at any time. Do not code against it.

## Options:
    <OWNER>
        The user or organization to list repos for. Defaults to your own
        GitHub username.

    -q, --qualified
        Include the "<USER>/" prefix on all repos listed

## See also:

hub(1)
`,
	}
)

func init() {
	cmdInternalListRepos.Flag.BoolVarP(&flagQualified, "qualified", "q", false, "QUALIFY")

	CmdRunner.Use(cmdInternalListRepos)
}

func internalListRepos(cmd *Command, args *Args) {
	gh := github.NewClient(github.GitHubHost)

	words := args.Words()
	var ownerName string
	if len(words) > 0 {
		ownerName = words[0]
	}

	repos, err := gh.FetchRepositories(ownerName)
	utils.Check(err)

	for _, repo := range repos {
		if flagQualified {
			ui.Printf("%s/%s\n", ownerName, repo.Name)
		} else {
			ui.Printf("%s\n", repo.Name)
		}
	}

	args.NoForward()
}