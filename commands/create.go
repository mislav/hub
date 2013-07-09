package commands

var cmdCreate = &Command{
	Run:   create,
	Usage: "create [NAME] [-p] [-d DESCRIPTION] [-h HOMEPAGE]",
	Short: "Create this repository on GitHub and add GitHub as origin",
	Long: `Create a new public GitHub repository from the current git
repository and add remote origin at "git@github.com:USER/REPOSITORY.git";
USER is your GitHub username and REPOSITORY is the current working
directory name. To explicitly name the new repository, pass in NAME,
optionally in ORGANIZATION/NAME form to create under an organization
you're a member of. With -p, create a private repository, and with
-d and -h set the repository's description and homepage URL, respectively.
`,
}

/*
  $ gh create
  ... create repo on github ...
  > git remote add -f origin git@github.com:YOUR_USER/CURRENT_REPO.git

  # with description:
  $ gh create -d 'It shall be mine, all mine!'

  $ gh create recipes
  [ repo created on GitHub ]
  > git remote add origin git@github.com:YOUR_USER/recipes.git

  $ gh create sinatra/recipes
  [ repo created in GitHub organization ]
  > git remote add origin git@github.com:sinatra/recipes.git
*/
func create(command *Command, args *Args) {
}
