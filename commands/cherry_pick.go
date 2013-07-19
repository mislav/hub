package commands

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

func cherryPick(command *Command, args *Args) {
}
