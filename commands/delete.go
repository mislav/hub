package commands

import (
	"os"
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdDelete = &Command{
	Run:   delete,
	Usage: "delete [[<ORGANIZATION>/]<NAME>]",
	Long: `Delete an existing repository on GitHub.

## Options:

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
)

func init() {
	CmdRunner.Use(cmdDelete)
}

func delete(command *Command, args *Args) {
	var repoName string
	if args.IsParamsEmpty() {
		dirName, err := git.WorkdirName()
		utils.Check(err)
		repoName = github.SanitizeProjectName(dirName)
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

	if ! gh.IsRepositoryExist(project) {
		fmt.Println("No such repository")
		args.NoForward()
		return
	}
	
	if os.Getenv("HUB_UNSAFE_DELETE") == "" {
		fmt.Printf("Repository '%s' exists. Really delete it (y/N)?", repoName)
		var s string	
		_, err = fmt.Scan(&s)
		if err != nil {
			fmt.Println(err);
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
		
		fmt.Println("Please write the name of the repository again (this operation can not be undone!): ")
		_, err = fmt.Scan(&s)
		if err != nil {
			fmt.Println(err);
			args.NoForward()
			return
		}
		s = strings.TrimSpace(s)
		if s != repoName {
			fmt.Println("Names don't match.. bailing out.. no deletion")
			args.NoForward()
			return
		}
	}
	
	err = gh.DeleteRepository(project)
	utils.Check(err)
	
	fmt.Printf("Deleted repository %v\n", repoName)
	
	args.NoForward()
}
