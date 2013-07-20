package main

import (
	"github.com/jingweno/gh/commands"
	"github.com/jingweno/gh/utils"
	"os"
)

func main() {
  runner := commands.Runner{Args: os.Args[1:]}
	err := runner.Execute()
	utils.Check(err)
}
