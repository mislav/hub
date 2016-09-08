package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdFetch = &Command{
	Run:          fetch,
	GitExtension: true,
	Usage:        "fetch <USER>[,<USER2>...]",
	Long: `Add missing remotes prior to performing git fetch.

## Examples:
		$ hub fetch --multiple jingweno mislav
		> git remote add jingweno git://github.com/jingweno/REPO.git
		> git remote add jingweno git://github.com/mislav/REPO.git
		> git fetch jingweno
		> git fetch mislav

## See also:

hub-remote(1), hub(1), git-fetch(1)
`,
}

func init() {
	CmdRunner.Use(cmdFetch)
}

func fetch(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		err := tranformFetchArgs(args)
		utils.Check(err)
	}
}

func tranformFetchArgs(args *Args) error {
	names := parseRemoteNames(args)

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	currentProject, currentProjectErr := localRepo.CurrentProject()

	projects := make(map[*github.Project]bool)
	ownerRegexp := regexp.MustCompile(fmt.Sprintf("^%s$", OwnerRe))
	for _, name := range names {
		if ownerRegexp.MatchString(name) && !isCloneable(name) {
			_, err := localRepo.RemoteByName(name)
			if err != nil {
				utils.Check(currentProjectErr)
				project := github.NewProject(name, currentProject.Name, "")
				gh := github.NewClient(project.Host)
				repo, err := gh.Repository(project)
				if err != nil {
					continue
				}

				projects[project] = repo.Private
			}
		}
	}

	for project, private := range projects {
		args.Before("git", "remote", "add", project.Owner, project.GitURL("", "", private))
	}

	return nil
}

func parseRemoteNames(args *Args) (names []string) {
	words := args.Words()
	if i := args.IndexOfParam("--multiple"); i != -1 {
		if args.ParamsSize() > 1 {
			names = words
		}
	} else if len(words) > 0 {
		remoteName := words[0]
		commaPattern := fmt.Sprintf("^%s(,%s)+$", OwnerRe, OwnerRe)
		remoteNameRegexp := regexp.MustCompile(commaPattern)
		if remoteNameRegexp.MatchString(remoteName) {
			i := args.IndexOfParam(remoteName)
			args.RemoveParam(i)
			names = strings.Split(remoteName, ",")
			args.InsertParam(i, names...)
			args.InsertParam(i, "--multiple")
		} else {
			names = append(names, remoteName)
		}
	}

	return
}
