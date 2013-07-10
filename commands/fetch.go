package commands

import (
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"regexp"
	"strings"
)

var cmdFetch = &Command{
	Run:          fetch,
	GitExtension: true,
	Usage:        "fetch [USER...]",
	Short:        "Download data, tags and branches from a remote repository",
	Long: `Adds missing remote(s) with git remote add prior to fetching. New
remotes are only added if they correspond to valid forks on GitHub.
`,
}

/*
  $ gh fetch jingweno
  > git remote add jingweno git://github.com/jingweno/REPO.git
  > git fetch jingweno

  $ git fetch jingweno,foo
  > git remote add jingweno ...
  > git remote add foo ...
  > git fetch --multiple jingweno foo

  $ git fetch --multiple jingweno foo
  > git remote add jingweno ...
  > git remote add foo ...
  > git fetch --multiple jingweno foo
*/
func fetch(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		tranformFetchArgs(args)
	}
}

func tranformFetchArgs(args *Args) {
	remotes, err := git.Remotes()
	utils.Check(err)

	names := parseRemoteNames(args)
	gh := github.New()
	projects := []github.Project{}
	ownerRegexp := regexp.MustCompile(OwnerRe)
	for _, name := range names {
		if ownerRegexp.MatchString(name) && !isRemoteExist(remotes, name) {
			project := github.NewProjectFromNameAndOwner("", name)
			repo, err := gh.Repository(project)
			if err != nil {
				continue
			}

			project = github.NewProjectFromNameAndOwner("", repo.FullName)
			projects = append(projects, project)
		}
	}

	for _, project := range projects {
		var isSSH bool
		if project.Owner == gh.Config.FetchUser() {
			isSSH = true
		}
		args.Before("git", "remote", "add", project.Owner, project.GitURL("", "", isSSH))
	}
}

func parseRemoteNames(args *Args) (names []string) {
	if i := args.IndexOfParam("--multiple"); i != -1 {
		if args.ParamsSize() > 1 {
			names = args.Params[1:]
		}
	} else {
		remoteName := args.FirstParam()
		remoteNameRegexp := regexp.MustCompile("^\\w+(,\\w+)$")
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

func isRemoteExist(remotes []*git.GitRemote, name string) bool {
	for _, r := range remotes {
		if r.Name == name {
			return true
		}
	}

	return false
}
