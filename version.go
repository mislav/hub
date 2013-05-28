package main

import (
	"fmt"
)

const Version = "0.0.1"

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "Show gh version",
	Long:  `Version shows the gh client version showstring.`,
}

func runVersion(cmd *Command, args []string) {
	gitVersion, err := git.Version()
	check(err)

	ghVersion := fmt.Sprintf("gh version %s", Version)

	fmt.Println(gitVersion)
	fmt.Println(ghVersion)
}
