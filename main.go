package main

import (
	"github.com/jingweno/gh/commands"
	"os"
)

func main() {
	err := commands.CmdRunner.Execute()
	os.Exit(err.ExitCode)
}
