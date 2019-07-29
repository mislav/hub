package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/github/hub/ui"
)

type Cmd struct {
	Name   string
	Args   []string
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
	Shell  bool
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

func (cmd *Cmd) Output() (string, error) {
	verboseLog(cmd)
	c := exec.Command(cmd.Name, cmd.Args...)
	c.Stderr = cmd.Stderr
	output, err := c.Output()

	return string(output), err
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
	if isWindows() {
		return cmd.Spawn()
	} else {
		return cmd.Exec()
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows" || detectWSL()
}

var detectedWSL bool
var detectedWSLContents string

// https://github.com/Microsoft/WSL/issues/423#issuecomment-221627364
func detectWSL() bool {
	if !detectedWSL {
		b := make([]byte, 1024)
		f, err := os.Open("/proc/version")
		if err == nil {
			f.Read(b)
			f.Close()
			detectedWSLContents = string(b)
		}
		detectedWSL = true
	}
	return strings.Contains(detectedWSLContents, "Microsoft")
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

func New(name string) *Cmd {
	return &Cmd{
		Name:   name,
		Args:   []string{},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func NewWithArray(cmd []string) *Cmd {
	return &Cmd{Name: cmd[0], Args: cmd[1:], Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
}

func NewWithShell(cmd []string) *Cmd {
	args := []string{"-c", fmt.Sprintf(`%s "$@"`, cmd[0])}
	for _, arg := range cmd {
		args = append(args, arg)
	}
	return &Cmd{
		Name:   findShellPath(),
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Shell:  true,
	}
}

// findShellPath returns the location of the `sh` binary, or `sh` if it is
// likely to be found in the PATH.
func findShellPath() string {
	// If this is not Windows, we presume that it's a Unix system and the
	// user has sh in the PATH already.
	if runtime.GOOS != "windows" {
		return "sh"
	}

	path := os.Getenv("PATH")
	dirs := strings.Split(path, ";")

	// Assume this is Cygwin or some Unix environment and our path is
	// already set appropriately.
	if len(dirs) == 1 {
		return "sh"
	}

	for _, dir := range dirs {
		// Git for Windows and PortableGit put the Git binary in the
		// `cmd` directory, which is in the PATH, but sh.exe is in the
		// `bin` directory next to it. If we found git.exe, put the
		// directory with sh at the end of PATH so we can invoke it.
		shPath := filepath.Join(dir, "..", "bin", "sh.exe")
		if _, err := os.Stat(shPath); err == nil {
			return shPath
		}
	}
	return "sh"
}

func verboseLog(cmd *Cmd) {
	if os.Getenv("HUB_VERBOSE") != "" {
		msg := fmt.Sprintf("$ %s %s", cmd.Name, strings.Join(cmd.Args, " "))
		if cmd.Shell {
			msg = fmt.Sprintf("$ %s", strings.Join(cmd.Args[2:], " "))
		}
		if ui.IsTerminal(os.Stderr) {
			msg = fmt.Sprintf("\033[35m%s\033[0m", msg)
		}
		ui.Errorln(msg)
	}
}
