package git

import (
	"fmt"
	"github.com/jingweno/gh/cmd"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Version() (string, error) {
	output, err := execGitCmd("version")
	if err != nil {
		return "", fmt.Errorf("Can't load git version")
	}

	return output[0], nil
}

func Dir() (string, error) {
	output, err := execGitCmd("rev-parse", "-q", "--git-dir")
	if err != nil {
		return "", fmt.Errorf("Not a git repository (or any of the parent directories): .git")
	}

	gitDir := output[0]
	gitDir, err = filepath.Abs(gitDir)
	if err != nil {
		return "", err
	}

	return gitDir, nil
}

func HasFile(segments ...string) bool {
	dir, err := Dir()
	if err != nil {
		return false
	}

	s := []string{dir}
	s = append(s, segments...)
	path := filepath.Join(s...)
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func BranchAtRef(refs ...string) (name string, err error) {
	dir, err := Dir()
	if err != nil {
		return
	}

	segments := []string{dir}
	segments = append(segments, refs...)
	path := filepath.Join(segments...)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	n := string(b)
	if strings.HasPrefix(n, "ref: ") {
		name = strings.TrimPrefix(n, "ref: ")
		name = strings.TrimSpace(name)
	} else {
		err = fmt.Errorf("No branch info in %s: %s", path, n)
	}

	return
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
		return "", fmt.Errorf("Can't load git var: GIT_EDITOR")
	}

	return output[0], nil
}

func Head() (string, error) {
	return BranchAtRef("HEAD")
}

func SymbolicFullName(name string) (string, error) {
	output, err := execGitCmd("rev-parse", "--symbolic-full-name", name)
	if err != nil {
		return "", fmt.Errorf("Unknown revision or path not in the working tree: %s", name)
	}

	return output[0], nil
}

func Ref(ref string) (string, error) {
	output, err := execGitCmd("rev-parse", "-q", ref)
	if err != nil {
		return "", fmt.Errorf("Unknown revision or path not in the working tree: %s", ref)
	}

	return output[0], nil
}

func RefList(a, b string) ([]string, error) {
	ref := fmt.Sprintf("%s...%s", a, b)
	output, err := execGitCmd("rev-list", "--cherry-pick", "--right-only", "--no-merges", ref)
	if err != nil {
		return []string{}, fmt.Errorf("Can't load rev-list for %s", ref)
	}

	return output, nil
}

func Show(sha string) (string, error) {
	output, err := execGitCmd("show", "-s", "--format=%w(78,0,0)%s%+b", sha)
	if err != nil {
		return "", fmt.Errorf("Can't show commit for %s", sha)
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
		return "", fmt.Errorf("Can't load git log %s..%s", sha1, sha2)
	}

	return outputs, nil
}

func Remotes() ([]string, error) {
	return execGitCmd("remote", "-v")
}

func Config(name string) (string, error) {
	output, err := execGitCmd("config", name)
	if err != nil {
		return "", fmt.Errorf("Unknown config %s", name)
	}

	return output[0], nil
}

func Alias(name string) (string, error) {
	return Config(fmt.Sprintf("alias.%s", name))
}

func Spawn(command string, args ...string) error {
	cmd := cmd.New("git")
	cmd.WithArg(command)
	for _, a := range args {
		cmd.WithArg(a)
	}

	return cmd.Exec()
}

func execGitCmd(input ...string) (outputs []string, err error) {
	cmd := cmd.New("git")
	for _, i := range input {
		cmd.WithArg(i)
	}

	out, err := cmd.ExecOutput()
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			outputs = append(outputs, string(line))
		}
	}

	return outputs, err
}
