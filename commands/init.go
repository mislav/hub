package commands

var cmdInit = &Command{
	Run:          gitInit,
	GitExtension: true,
	Usage:        "init -g",
	Short:        "Create an empty git repository or reinitialize an existing one",
	Long: `Create a git repository as with git-init(1) and add remote origin at
"git@github.com:USER/REPOSITORY.git"; USER is your GitHub username and
REPOSITORY is the current working directory's basename.
`,
}

/*
  $ gh init -g
  > git init
  > git remote add origin git@github.com:USER/REPO.git
*/
func gitInit(command *Command, args *Args) {
}
