package commands

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
)

var cmdDelete = &Command{
	Run:   deleteRepo,
	Usage: "delete [-y] [<ORGANIZATION>/]<NAME>",
	Long: `Delete an existing repository on GitHub.

## Options:

	-y, --yes
		Skip the confirmation prompt and immediately delete the repository.

	[<ORGANIZATION>/]<NAME>
		The name for the repository on GitHub.

## Examples:
		$ hub delete recipes
		[ personal repo deleted on GitHub ]

		$ hub delete sinatra/recipes
		[ repo deleted in GitHub organization ]

## See also:

hub-init(1), hub(1)
`,
}

func init() {
	CmdRunner.Use(cmdDelete)
}

func deleteRepo(command *Command, args *Args) {
	var repoName string
	if !args.IsParamsEmpty() {
		repoName = args.FirstParam()
	}

	re := regexp.MustCompile(NameWithOwnerRe)
	if !re.MatchString(repoName) {
		utils.Check(command.UsageError(""))
	}

	config := github.CurrentConfig()
	host, err := config.DefaultHost()
	if err != nil {
		utils.Check(github.FormatError("deleting repository", err))
	}

	owner := host.User
	if strings.Contains(repoName, "/") {
		split := strings.SplitN(repoName, "/", 2)
		owner, repoName = split[0], split[1]
	}

	project := github.NewProject(owner, repoName, host.Host)
	gh := github.NewClient(project.Host)

	if !args.Flag.Bool("--yes") {
		ui.Printf("Really delete repository '%s' (yes/N)? ", project)
		answer := ""
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			answer = strings.TrimSpace(scanner.Text())
		}
		utils.Check(scanner.Err())
		if answer != "yes" {
			utils.Check(fmt.Errorf("Please type 'yes' for confirmation."))
		}
	}

	if args.Noop {
		ui.Printf("Would delete repository '%s'.\n", project)
	} else {
		err = gh.DeleteRepository(project)
		if err != nil && strings.Contains(err.Error(), "HTTP 403") {
			ui.Errorf("Please edit the token used for hub at https://%s/settings/tokens\n", project.Host)
			ui.Errorln("and verify that the `delete_repo` scope is enabled.")
		}
		utils.Check(err)
		ui.Printf("Deleted repository '%s'.\n", project)
	}

	args.NoForward()
}
