package main

import (
	"fmt"
)

const Version = "dev"

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "Show gh version",
	Long:  `Version shows the gh client version showstring.`,
}

func runVersion(cmd *Command, args []string) {
	fmt.Println(Version)
}
