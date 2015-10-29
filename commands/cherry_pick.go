package commands

import (
	"regexp"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdCherryPick = &Command{
	Run:          cherryPick,
	GitExtension: true,
	Usage:        "cherry-pick GITHUB-REF",
	Short:        "Apply the changes introduced by some existing commits",
	Long: `Cherry-pick a commit from a fork using either full URL to the commit
or GitHub-flavored Markdown notation, which is user@sha. If the remote
doesn't yet exist, it will be added. A git-fetch(1) user is issued
prior to the cherry-pick attempt.
`,
}

func init() {
	CmdRunner.Use(cmdCherryPick)
}

/*
  $ gh cherry-pick https://github.com/jingweno/gh/commit/a319d88#comments
  > git remote add -f --no-tags jingweno git://github.com/jingweno/gh.git
  > git cherry-pick a319d88

  $ gh cherry-pick jingweno@a319d88
  > git remote add -f --no-tags jingweno git://github.com/jingweno/gh.git
  > git cherry-pick a319d88

  $ gh cherry-pick jingweno@SHA
  > git fetch jingweno
  > git cherry-pick SHA
*/
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
	url, err := github.ParseURL(ref)
	if err == nil {
		commitRegex := regexp.MustCompile("^commit\\/([a-f0-9]{7,40})")
		projectPath := url.ProjectPath()
		if commitRegex.MatchString(projectPath) {
			sha = commitRegex.FindStringSubmatch(projectPath)[1]
			project = url.Project

			return
		}
	}

	ownerWithShaRegexp := regexp.MustCompile("^([a-zA-Z0-9][a-zA-Z0-9-]*)@([a-f0-9]{7,40})$")
	if ownerWithShaRegexp.MatchString(ref) {
		matches := ownerWithShaRegexp.FindStringSubmatch(ref)
		sha = matches[2]
		localRepo, err := github.LocalRepo()
		utils.Check(err)

		project, err = localRepo.CurrentProject()
		utils.Check(err)
		project.Owner = matches[1]
	}

	return
}
