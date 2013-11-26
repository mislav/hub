package cmd

import (
	"fmt"
	"github.com/jingweno/gh/utils"
	"github.com/kballard/go-shellquote"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	Name string
	Args []string
}

func (cmd Cmd) String() string {
	return fmt.Sprintf("%s %s", cmd.Name, strings.Join(cmd.Args, " "))
}

func (cmd *Cmd) WithArg(arg string) *Cmd {
	if arg != "" {
		cmd.Args = append(cmd.Args, arg)
	}

	return cmd
}

func (cmd *Cmd) WithArgs(args ...string) *Cmd {
	for _, arg := range args {
		cmd.WithArg(arg)
	}

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

func New(cmd string) *Cmd {
	cmds, err := shellquote.Split(cmd)
	utils.Check(err)

	name := cmds[0]
	args := make([]string, 0)
	for _, arg := range cmds[1:] {
		args = append(args, arg)
	}
	return &Cmd{Name: name, Args: args}
}

func NewWithArray(cmd []string) *Cmd {
	return &Cmd{Name: cmd[0], Args: cmd[1:]}
}
