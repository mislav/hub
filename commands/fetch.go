package commands

import (
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
		err := tranformFetchArgs(args)
		utils.Check(err)
	}
}

func tranformFetchArgs(args *Args) error {
	names := parseRemoteNames(args)
	localRepo := github.LocalRepo()

	gh := github.New()
	projects := make(map[*github.Project]bool)
	ownerRegexp := regexp.MustCompile(OwnerRe)
	for _, name := range names {
		if ownerRegexp.MatchString(name) {
			_, err := localRepo.RemotesByName(name)
			if err != nil {
				project := github.NewProject(name, "", "")
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
