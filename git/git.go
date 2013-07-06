package git

import (
	"errors"
	"fmt"
	"github.com/jingweno/gh/cmd"
	"os/exec"
	"path/filepath"
	"strings"
)

func Version() (string, error) {
	output, err := execGitCmd("version")
	if err != nil {
		return "", errors.New("Can't load git version")
	}

	return output[0], nil
}

func Dir() (string, error) {
	output, err := execGitCmd("rev-parse", "-q", "--git-dir")
	if err != nil {
		return "", errors.New("Not a git repository (or any of the parent directories): .git")
	}

	gitDir := output[0]
	gitDir, err = filepath.Abs(gitDir)
	if err != nil {
		return "", err
	}

	return gitDir, nil
}

func PullReqMsgFile() (string, error) {
	gitDir, err := Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(gitDir, "PULLREQ_EDITMSG"), nil
}

func Editor() (string, error) {
	output, err := execGitCmd("var", "GIT_EDITOR")
	if err != nil {
		return "", errors.New("Can't load git var: GIT_EDITOR")
	}

	return output[0], nil
}

func EditorPath() (string, error) {
	gitEditor, err := Editor()
	if err != nil {
		return "", err
	}

	gitEditorWithParams := strings.Split(gitEditor, " ")
	gitEditor = gitEditorWithParams[0]
	gitEditorParams := gitEditorWithParams[1:]

	editorPath, err := exec.LookPath(gitEditor)
	if err != nil {
		return "", errors.New("Can't locate git editor: " + gitEditor)
	}

	for _, p := range gitEditorParams {
		editorPath = editorPath + " " + p
	}

	return editorPath, nil
}

func Head() (*Branch, error) {
	output, err := execGitCmd("symbolic-ref", "-q", "HEAD")
	if err != nil {
		return nil, errors.New("Can't load git HEAD")
	}

	return &Branch{output[0]}, nil
}

func Ref(ref string) (string, error) {
	output, err := execGitCmd("rev-parse", "-q", ref)
	if err != nil {
		return "", errors.New("Unknown revision or path not in the working tree: " + ref)
	}

	return output[0], nil
}

func Log(sha1, sha2 string) (string, error) {
	execCmd := cmd.New("git")
	execCmd.WithArg("log").WithArg("--no-color")
	execCmd.WithArg("--format=%h (%aN, %ar)%n%w(78,3,3)%s%n%+b")
	execCmd.WithArg("--cherry")
	shaRange := fmt.Sprintf("%s...%s", sha1, sha2)
	execCmd.WithArg(shaRange)

	outputs, err := execCmd.ExecOutput()
	if err != nil {
		return "", errors.New("Can't load git log " + sha1 + ".." + sha2)
	}

	return outputs, nil
}

func Config(name string) (string, error) {
	output, err := execGitCmd("config", name)
	if err != nil {
		return "", fmt.Errorf("Unknown config %s", name)
	}

	return output[0], nil
}

func SysExec(command string, args ...string) error {
	cmd := cmd.New("git")
	cmd.WithArg(command)
	for _, a := range args {
		cmd.WithArg(a)
	}

	return cmd.SysExec()
}

func Spawn(command string, args ...string) error {
	cmd := cmd.New("git")
	cmd.WithArg(command)
	for _, a := range args {
		cmd.WithArg(a)
	}

	out, err := cmd.ExecOutput()
	if err != nil {
		return errors.New(out)
	}

	return nil
}

func execGitCmd(input ...string) (outputs []string, err error) {
	cmd := cmd.New("git")
	for _, i := range input {
		cmd.WithArg(i)
	}

	out, err := cmd.ExecOutput()
	for _, line := range strings.Split(out, "\n") {
		outputs = append(outputs, string(line))
	}

	return outputs, err
}
