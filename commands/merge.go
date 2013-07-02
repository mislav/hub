package commands

import ()

var cmdMerge = &Command{
	Run:          merge,
	GitExtension: true,
	Usage:        "merge PULLREQ-URL",
	Short:        "Join two or more development histories (branches) together",
	Long: `Merge the pull request with a commit message that includes the pull request
ID and title, similar to the GitHub Merge Button.
`,
}

func merge(command *Command, args *Args) {
}
