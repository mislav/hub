package git

import (
	"errors"
	"fmt"
	"github.com/jingweno/gh/cmd"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Version() (string, error) {
	output, err := execGitCmd([]string{"version"})
	if err != nil {
		return "", errors.New("Can't load git version")
	}

	return output[0], nil
}

func Dir() (string, error) {
	output, err := execGitCmd([]string{"rev-parse", "-q", "--git-dir"})
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
	output, err := execGitCmd([]string{"var", "GIT_EDITOR"})
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

func Head() (string, error) {
	output, err := execGitCmd([]string{"symbolic-ref", "-q", "--short", "HEAD"})
	if err != nil {
		return "master", errors.New("Can't load git HEAD")
	}

	return output[0], nil
}

func Ref(ref string) (string, error) {
	output, err := execGitCmd([]string{"rev-parse", "-q", ref})
	if err != nil {
		return "", errors.New("Unknown revision or path not in the working tree: " + ref)
	}

	return output[0], nil
}

// FIXME: only care about origin push remote now
func Remote() (string, error) {
	r := regexp.MustCompile("origin\t(.+github.com.+) \\(push\\)")
	output, err := execGitCmd([]string{"remote", "-v"})
	if err != nil {
		return "", errors.New("Can't load git remote")
	}

	for _, o := range output {
		if r.MatchString(o) {
			return r.FindStringSubmatch(o)[1], nil
		}
	}

	return "", errors.New("Can't find git remote (push)")
}

func AddRemote(name, url string) error {
	_, err := execGitCmd([]string{"remote", "add", "-f", name, url})

	return err
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

func execGitCmd(input []string) (outputs []string, err error) {
	cmd := cmd.New("git")
	for _, i := range input {
		cmd.WithArg(i)
	}

	out, err := cmd.ExecOutput()
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(out, "\n") {
		outputs = append(outputs, string(line))
	}

	return outputs, nil
}
