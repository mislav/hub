package commands

var cmdBrowse = &Command{
	Run:   browse,
	Usage: "browse [-u [USER/]REPOSITORY] [-p SUBPAGE]",
	Short: "Open a GitHub page in the default browser",
	Long: `Open repository's GitHub page in the system's default web browser using
open(1) or the BROWSER env variable. If the repository isn't specified,
browse opens the page of the repository found in the current directory.
If SUBPAGE is specified, the browser will open on the specified subpage:
one of "wiki", "commits", "issues" or other (the default is "tree").
`,
}

func browse(cmd *Command, args []string) {
}
