package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
	"github.com/github/hub/ui"
)

var cmdDelete = &Command{
	Run:   delete,
	Usage: "delete [-y] [<ORGANIZATION>/]<NAME>",
	Long: `Delete an existing repository on GitHub.

## Options:

	-y
		Assume yes, force deletion of repository without asking.

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

var (
	flagDeleteAssumeYes bool
)

func init() {
	cmdDelete.Flag.BoolVarP(&flagDeleteAssumeYes, "--yes", "y", false, "YES")

	CmdRunner.Use(cmdDelete)
}

func delete(command *Command, args *Args) {
	var repoName string
	if args.IsParamsEmpty() {
		ui.Errorln("Expecting name of reposiroty to delete")
		args.NoForward()
		return
	} else {
		reg := regexp.MustCompile("^[^-]")
		if !reg.MatchString(args.FirstParam()) {
			err := fmt.Errorf("invalid argument: %s", args.FirstParam())
			utils.Check(err)
		}
		repoName = args.FirstParam()
	}

	config := github.CurrentConfig()
	host, err := config.DefaultHost()
	if err != nil {
		utils.Check(github.FormatError("deleting repository", err))
	}

	owner := host.User
	if strings.Contains(repoName, "/") {
		split := strings.SplitN(repoName, "/", 2)
		owner = split[0]
		repoName = split[1]
	}

	project := github.NewProject(owner, repoName, host.Host)
	gh := github.NewClient(project.Host)
	
	if !flagDeleteAssumeYes {
		fmt.Printf("Really delete repository '%s'(y/N)?", repoName)
		var s string
		_, err = fmt.Scan(&s)
		if err != nil {
			fmt.Println(err)
			args.NoForward()
			return
		}
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		if s != "y" {
			fmt.Println("Abort: not deleting anything.")
			args.NoForward()
			return
		}
	}

	err = gh.DeleteRepository(project)
	utils.Check(err)

	fmt.Printf("Deleted repository %v\n", repoName)

	args.NoForward()
}
