package commands

import (
//"github.com/jingweno/gh/github"
//"github.com/jingweno/gh/utils"
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
*/
func fetch(command *Command, args *Args) {
}
