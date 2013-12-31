package commands

var cmdSubmodule = &Command{
	Run:          submodule,
	GitExtension: true,
	Usage:        "submodule add [-p] OPTIONS [USER/]REPOSITORY DIRECTORY",
	Short:        "Initialize, update or inspect submodules",
	Long: `Submodule repository "git://github.com/USER/REPOSITORY.git" into
DIRECTORY as  with  git-submodule(1).  When  USER/  is  omitted,
assumes   your   GitHub  login.  With  -p,  use  private  remote
"git@github.com:USER/REPOSITORY.git".`,
}

func init() {
	CmdRunner.Use(cmdSubmodule)
}

/**
  $ gh submodule add jingweno/gh vendor/gh
  > git submodule add git://github.com/jingweno/gh.git vendor/gh

  $ gh submodule add -p jingweno/gh vendor/gh
  > git submodule add git@github.com:jingweno/gh.git vendor/gh

  $ gh submodule add -b gh --name gh jingweno/gh vendor/gh
  > git submodule add -b gh --name gh git://github.com/jingweno/gh.git vendor/gh
**/

func submodule(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformSubmoduleArgs(args)
	}
}

func transformSubmoduleArgs(args *Args) {
	var idx int
	if idx = args.IndexOfParam("add"); idx == -1 {
		return
	}
	args.RemoveParam(idx)
	transformCloneArgs(args)
	args.InsertParam(idx, "add")
}
