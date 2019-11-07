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

## See also:

hub-am(1), hub(1), git-cherry-pick(1)
`,
}

var (
	shaReStr       = "[a-f0-9]{7,40}"
	commitRegex    = regexp.MustCompile(fmt.Sprintf("^commit/(%s)", shaReStr))
	pullRegex      = regexp.MustCompile(fmt.Sprintf(`^pull/(\d+)/commits/(%s)`, shaReStr))
	ownerWithShaRe = regexp.MustCompile(fmt.Sprintf("^(%s)@(%s)$", OwnerReStr, shaReStr))
)

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

	var project *github.Project
	var sha, refspec string

	var mainProject *github.Project
	localRepo, mainProjectErr := github.LocalRepo()
	if mainProjectErr == nil {
		mainProject, mainProjectErr = localRepo.MainProject()
	}

	ref := args.LastParam()
	if url, err := github.ParseURL(ref); err == nil {
		projectPath := url.ProjectPath()
		if matches := commitRegex.FindStringSubmatch(projectPath); len(matches) > 0 {
			sha = matches[1]
			project = url.Project
		} else if matches := pullRegex.FindStringSubmatch(projectPath); len(matches) > 0 {
			pullId := matches[1]
			sha = matches[2]
			utils.Check(mainProjectErr)
			project = mainProject
			refspec = fmt.Sprintf("refs/pull/%s/head", pullId)
		}
	} else {
		if matches := ownerWithShaRe.FindStringSubmatch(ref); len(matches) > 0 {
			utils.Check(mainProjectErr)
			project = mainProject
			project.Owner = matches[1]
			sha = matches[2]
		}
	}

	if project != nil {
		args.ReplaceParam(args.IndexOfParam(ref), sha)

		tmpName := "_hub-cherry-pick"
		remoteName := tmpName

		if remote, err := localRepo.RemoteForProject(project); err == nil {
			remoteName = remote.Name
		} else {
			args.Before("git", "remote", "add", remoteName, project.GitURL("", "", false))
		}

		fetchArgs := []string{"git", "fetch", "-q", "--no-tags", remoteName}
		if refspec != "" {
			fetchArgs = append(fetchArgs, refspec)
		}
		args.Before(fetchArgs...)

		if remoteName == tmpName {
			args.Before("git", "remote", "rm", remoteName)
		}
	}
}
