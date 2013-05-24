package main

import (
	"os"
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

func (cmd *ExecCmd) ExecOutput() (string, error) {
	output, err := exec.Command(cmd.Name, cmd.Args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (cmd *ExecCmd) Exec() error {
	c := exec.Command(cmd.Name, cmd.Args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

func NewExecCmd(name string) *ExecCmd {
	return &ExecCmd{name, make([]string, 0)}
}
