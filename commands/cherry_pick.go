package commands

import (
	"fmt"
	"regexp"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdCherryPick = &Command{
	Run:          cherryPick,
	GitExtension: true,
	Usage: `
cherry-pick <COMMIT-URL>
cherry-pick <USER>@<SHA>
`,
	Long: `Cherry-pick a commit from a fork on GitHub.

## Examples:
		$ hub cherry-pick https://github.com/jingweno/gh/commit/a319d88#comments
		> git remote add -f --no-tags jingweno git://github.com/jingweno/gh.git
		> git cherry-pick a319d88

		$ hub cherry-pick jingweno@a319d88

## See also:

hub-am(1), hub(1), git-cherry-pick(1)
`,
}

func init() {
	CmdRunner.Use(cmdCherryPick)
}

func cherryPick(command *Command, args *Args) {
	if args.IndexOfParam("-m") == -1 && args.IndexOfParam("--mainline") == -1 {
		transformCherryPickArgs(args)
	}
}

func transformCherryPickArgs(args *Args) {
	if args.IsParamsEmpty() {
		return
	}

	ref := args.LastParam()
	project, sha := parseCherryPickProjectAndSha(ref)
	if project != nil {
		args.ReplaceParam(args.IndexOfParam(ref), sha)

		remote := gitRemoteForProject(project)
		if remote != nil {
			args.Before("git", "fetch", remote.Name)
		} else {
			args.Before("git", "remote", "add", "-f", "--no-tags", project.Owner, project.GitURL("", "", false))
		}
	}
}

func parseCherryPickProjectAndSha(ref string) (project *github.Project, sha string) {
	shaRe := "[a-f0-9]{7,40}"

	var mainProject *github.Project
	localRepo, mainProjectErr := github.LocalRepo()
	if mainProjectErr == nil {
		mainProject, mainProjectErr = localRepo.MainProject()
	}

	url, err := github.ParseURL(ref)
	if err == nil {
		projectPath := url.ProjectPath()

		commitRegex := regexp.MustCompile(fmt.Sprintf("^commit/(%s)", shaRe))
		if matches := commitRegex.FindStringSubmatch(projectPath); len(matches) > 0 {
			sha = matches[1]
			project = url.Project
			return
		}
	}

	ownerWithShaRegexp := regexp.MustCompile(fmt.Sprintf("^(%s)@(%s)$", OwnerRe, shaRe))
	if matches := ownerWithShaRegexp.FindStringSubmatch(ref); len(matches) > 0 {
		utils.Check(mainProjectErr)
		project = mainProject
		project.Owner = matches[1]
		sha = matches[2]
	}

	return
}
