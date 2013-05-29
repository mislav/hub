package main

var cmdCiStatus = &Command{
	Run:   ciStatus,
	Usage: "ci-status [COMMIT]",
	Short: "Show CI status of a commit",
	Long: `Looks up the SHA for COMMIT in GitHub Status API and displays the latest
status. Exits with one of:

success (0), error (1), failure (1), pending (2), no status (3)
`,
}

func ciStatus(cmd *Command, args []string) {
}
