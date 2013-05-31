package cmd

import (
	"os"
	"os/exec"
)

type Cmd struct {
	Name string
	Args []string
}

func (cmd *Cmd) WithArg(arg string) *Cmd {
	cmd.Args = append(cmd.Args, arg)

	return cmd
}

func (cmd *Cmd) ExecOutput() (string, error) {
	output, err := exec.Command(cmd.Name, cmd.Args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (cmd *Cmd) Exec() error {
	c := exec.Command(cmd.Name, cmd.Args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

func New(name string) *Cmd {
	return &Cmd{name, make([]string, 0)}
}

func NewWithArray(cmd []string) *Cmd {
	return &Cmd{cmd[0], cmd[1:]}
}
