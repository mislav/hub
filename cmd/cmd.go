package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
	"github.com/kballard/go-shellquote"
)

type Cmd struct {
	Name   string
	Args   []string
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
}

func (cmd Cmd) String() string {
	return fmt.Sprintf("%s %s", cmd.Name, strings.Join(cmd.Args, " "))
}

func (cmd *Cmd) WithArg(arg string) *Cmd {
	cmd.Args = append(cmd.Args, arg)

	return cmd
}

func (cmd *Cmd) WithArgs(args ...string) *Cmd {
	for _, arg := range args {
		cmd.WithArg(arg)
	}

	return cmd
}

func (cmd *Cmd) CombinedOutput() (string, error) {
	verboseLog(cmd)
	output, err := exec.Command(cmd.Name, cmd.Args...).CombinedOutput()

	return string(output), err
}

func (cmd *Cmd) Success() bool {
	verboseLog(cmd)
	err := exec.Command(cmd.Name, cmd.Args...).Run()
	return err == nil
}

// Run runs command with `Exec` on platforms except Windows
// which only supports `Spawn`
func (cmd *Cmd) Run() error {
	if runtime.GOOS == "windows" {
		return cmd.Spawn()
	} else {
		return cmd.Exec()
	}
}

// Spawn runs command with spawn(3)
func (cmd *Cmd) Spawn() error {
	verboseLog(cmd)
	c := exec.Command(cmd.Name, cmd.Args...)
	c.Stdin = cmd.Stdin
	c.Stdout = cmd.Stdout
	c.Stderr = cmd.Stderr

	return c.Run()
}

// Exec runs command with exec(3)
// Note that Windows doesn't support exec(3): http://golang.org/src/pkg/syscall/exec_windows.go#L339
func (cmd *Cmd) Exec() error {
	verboseLog(cmd)

	binary, err := exec.LookPath(cmd.Name)
	if err != nil {
		return &exec.Error{
			Name: cmd.Name,
			Err:  fmt.Errorf("command not found"),
		}
	}

	args := []string{binary}
	args = append(args, cmd.Args...)

	return syscall.Exec(binary, args, os.Environ())
}

func New(cmd string) *Cmd {
	cmds, err := shellquote.Split(cmd)
	utils.Check(err)

	name := cmds[0]
	args := make([]string, 0)
	for _, arg := range cmds[1:] {
		args = append(args, arg)
	}
	return &Cmd{Name: name, Args: args, Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
}

func NewWithArray(cmd []string) *Cmd {
	return &Cmd{Name: cmd[0], Args: cmd[1:], Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
}

func verboseLog(cmd *Cmd) {
	if os.Getenv("HUB_VERBOSE") != "" {
		msg := fmt.Sprintf("$ %s %s", cmd.Name, strings.Join(cmd.Args, " "))
		if ui.IsTerminal(os.Stderr) {
			msg = fmt.Sprintf("\033[35m%s\033[0m", msg)
		}
		ui.Errorln(msg)
	}
}
