package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Cmd struct {
	Name string
	Args []string
}

func (cmd Cmd) String() string {
	return fmt.Sprintf("%s %s", cmd.Name, strings.Join(cmd.Args, " "))
}

func (cmd *Cmd) WithArg(arg string) *Cmd {
	cmd.Args = append(cmd.Args, arg)

	return cmd
}

func (cmd *Cmd) WithArgs(args ...string) *Cmd {
	cmd.Args = append(cmd.Args, args...)

	return cmd
}

func (cmd *Cmd) ExecOutput() (string, error) {
	output, err := exec.Command(cmd.Name, cmd.Args...).CombinedOutput()

	return string(output), err
}

func (cmd *Cmd) Exec() error {
	binary, lookErr := exec.LookPath(cmd.Name)
	if lookErr != nil {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}

	c := exec.Command(binary, cmd.Args...)
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
