package main

import (
	"os/exec"
)

type ExecCmd struct {
	Name string
	Args []string
}

func (cmd *ExecCmd) WithArg(arg string) *ExecCmd {
	cmd.Args = append(cmd.Args, arg)

	return cmd
}

func (cmd *ExecCmd) Exec() (out string, err error) {
	output, err := exec.Command(cmd.Name, cmd.Args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func NewExecCmd(name string) *ExecCmd {
	return &ExecCmd{name, make([]string, 0)}
}
