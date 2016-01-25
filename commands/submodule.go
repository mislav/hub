package commands

var cmdSubmodule = &Command{
	Run:          submodule,
	GitExtension: true,
	Usage:        "submodule add [-p] [<OPTIONS>] [<USER>/]<REPOSITORY> <DESTINATION>",
	Long: `Add a git submodule for a GitHub repository.

## Examples:
		$ hub submodule add jingweno/gh vendor/gh
		> git submodule add git://github.com/jingweno/gh.git vendor/gh

## See also:

hub-remote(1), hub(1), git-submodule(1)
`,
}

func init() {
	CmdRunner.Use(cmdSubmodule)
}

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
