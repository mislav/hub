package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
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

func (cmd *Cmd) SysExec() error {
	binary, lookErr := exec.LookPath(cmd.Name)
	if lookErr != nil {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}

	args := []string{cmd.Name}
	args = append(args, cmd.Args...)

	env := os.Environ()

	return syscall.Exec(binary, args, env)
}

func New(name string) *Cmd {
	return &Cmd{name, make([]string, 0)}
}

func NewWithArray(cmd []string) *Cmd {
	return &Cmd{cmd[0], cmd[1:]}
}
