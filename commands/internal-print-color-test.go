package commands

import "github.com/github/hub/utils"

var cmdInternalPrintColorTest = &Command{
	Run:   internalPrintColorTest,
	Usage: "internal-print-color-test",
	Long: "Print a terminal color test swatch",
}

func init() {
	CmdRunner.Use(cmdInternalPrintColorTest)
}

func internalPrintColorTest(_ *Command, args *Args) {
	args.NoForward()
	utils.PrintColorCube()
}