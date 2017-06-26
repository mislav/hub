package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/hub/cmd"
)

var GlobalFlags []string

func Version() (string, error) {
	output, err := gitOutput("version")
	if err == nil {
		return output[0], nil
	} else {
		return "", fmt.Errorf("error running git version: %s", err)
	}
}

var cachedDir string

func Dir() (string, error) {
	if cachedDir != "" {
		return cachedDir, nil
	}

	output, err := gitOutput("rev-parse", "-q", "--git-dir")
	if err != nil {
		return "", fmt.Errorf("Not a git repository (or any of the parent directories): .git")
	}

	var chdir string
	for i, flag := range GlobalFlags {
		if flag == "-C" {
			dir := GlobalFlags[i+1]
			if filepath.IsAbs(dir) {
				chdir = dir
			} else {
				chdir = filepath.Join(chdir, dir)
			}
		}
	}

	gitDir := output[0]

	if !filepath.IsAbs(gitDir) {
		if chdir != "" {
			gitDir = filepath.Join(chdir, gitDir)
		}

		gitDir, err = filepath.Abs(gitDir)
		if err != nil {
			return "", err
		}

		gitDir = filepath.Clean(gitDir)
	}

	cachedDir = gitDir
	return gitDir, nil
}

func WorkdirName() (string, error) {
	output, err := gitOutput("rev-parse", "--show-toplevel")
	if err == nil {
		if len(output) > 0 {
			return output[0], nil
		} else {
			return "", fmt.Errorf("unable to determine git working directory")
		}
	} else {
		return "", err
	}
}

func HasFile(segments ...string) bool {
	// The blessed way to resolve paths within git dir since Git 2.5.0
	output, err := gitOutput("rev-parse", "-q", "--git-path", filepath.Join(segments...))
	if err == nil && output[0] != "--git-path" {
		if _, err := os.Stat(output[0]); err == nil {
			return true
		}
	}

	// Fallback for older git versions
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

func BranchAtRef(paths ...string) (name string, err error) {
	dir, err := Dir()
	if err != nil {
		return
	}

	segments := []string{dir}
	segments = append(segments, paths...)
	path := filepath.Join(segments...)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	n := string(b)
	refPrefix := "ref: "
	if strings.HasPrefix(n, refPrefix) {
		name = strings.TrimPrefix(n, refPrefix)
		name = strings.TrimSpace(name)
	} else {
		err = fmt.Errorf("No branch info in %s: %s", path, n)
	}

	return
}

func Editor() (string, error) {
	output, err := gitOutput("var", "GIT_EDITOR")
	if err != nil {
		return "", fmt.Errorf("Can't load git var: GIT_EDITOR")
	}

	return os.ExpandEnv(output[0]), nil
}

func Head() (string, error) {
	return BranchAtRef("HEAD")
}

func SymbolicFullName(name string) (string, error) {
	output, err := gitOutput("rev-parse", "--symbolic-full-name", name)
	if err != nil {
		return "", fmt.Errorf("Unknown revision or path not in the working tree: %s", name)
	}

	return output[0], nil
}

func Ref(ref string) (string, error) {
	output, err := gitOutput("rev-parse", "-q", ref)
	if err != nil {
		return "", fmt.Errorf("Unknown revision or path not in the working tree: %s", ref)
	}

	return output[0], nil
}

func RefList(a, b string) ([]string, error) {
	ref := fmt.Sprintf("%s...%s", a, b)
	output, err := gitOutput("rev-list", "--cherry-pick", "--right-only", "--no-merges", ref)
	if err != nil {
		return []string{}, fmt.Errorf("Can't load rev-list for %s", ref)
	}

	return output, nil
}

func NewRange(a, b string) (*Range, error) {
	output, err := gitOutput("rev-parse", "-q", a, b)
	if err != nil {
		return nil, err
	}

	return &Range{output[0], output[1]}, nil
}

type Range struct {
	A string
	B string
}

func (r *Range) IsIdentical() bool {
	return strings.EqualFold(r.A, r.B)
}

func (r *Range) IsAncestor() bool {
	cmd := gitCmd("merge-base", "--is-ancestor", r.A, r.B)
	return cmd.Success()
}

func CommentChar() string {
	char, err := Config("core.commentchar")
	if err != nil {
		char = "#"
	}

	return char
}

func Show(sha string) (string, error) {
	cmd := cmd.New("git")
	cmd.WithArg("show").WithArg("-s").WithArg("--format=%s%n%+b").WithArg(sha)

	output, err := cmd.CombinedOutput()
	output = strings.TrimSpace(output)

	return output, err
}

func Log(sha1, sha2 string) (string, error) {
	execCmd := cmd.New("git")
	execCmd.WithArg("-c").WithArg("log.showSignature=false").WithArg("log").WithArg("--no-color")
	execCmd.WithArg("--format=%h (%aN, %ar)%n%w(78,3,3)%s%n%+b")
	execCmd.WithArg("--cherry")
	shaRange := fmt.Sprintf("%s...%s", sha1, sha2)
	execCmd.WithArg(shaRange)

	outputs, err := execCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Can't load git log %s..%s", sha1, sha2)
	}

	return outputs, nil
}

func Remotes() ([]string, error) {
	return gitOutput("remote", "-v")
}

func Config(name string) (string, error) {
	return gitGetConfig(name)
}

func ConfigAll(name string) ([]string, error) {
	mode := "--get-all"
	if strings.Contains(name, "*") {
		mode = "--get-regexp"
	}

	lines, err := gitOutput(gitConfigCommand([]string{mode, name})...)
	if err != nil {
		err = fmt.Errorf("Unknown config %s", name)
	}
	return lines, err
}

func GlobalConfig(name string) (string, error) {
	return gitGetConfig("--global", name)
}

func SetGlobalConfig(name, value string) error {
	_, err := gitConfig("--global", name, value)
	return err
}

func gitGetConfig(args ...string) (string, error) {
	output, err := gitOutput(gitConfigCommand(args)...)
	if err != nil {
		return "", fmt.Errorf("Unknown config %s", args[len(args)-1])
	}

	if len(output) == 0 {
		return "", nil
	}

	return output[0], nil
}

func gitConfig(args ...string) ([]string, error) {
	return gitOutput(gitConfigCommand(args)...)
}

func gitConfigCommand(args []string) []string {
	cmd := []string{"config"}
	return append(cmd, args...)
}

func Alias(name string) (string, error) {
	return Config(fmt.Sprintf("alias.%s", name))
}

func Run(args ...string) error {
	cmd := gitCmd(args...)
	return cmd.Run()
}

func Spawn(args ...string) error {
	cmd := gitCmd(args...)
	return cmd.Spawn()
}

func Quiet(args ...string) bool {
	cmd := gitCmd(args...)
	return cmd.Success()
}

func IsGitDir(dir string) bool {
	cmd := cmd.New("git")
	cmd.WithArgs("--git-dir="+dir, "rev-parse", "--git-dir")
	return cmd.Success()
}

func LocalBranches() ([]string, error) {
	lines, err := gitOutput("branch", "--list")
	if err == nil {
		for i, line := range lines {
			lines[i] = strings.TrimPrefix(line, "* ")
			lines[i] = strings.TrimPrefix(lines[i], "  ")
		}
	}
	return lines, err
}

func gitOutput(input ...string) (outputs []string, err error) {
	cmd := gitCmd(input...)

	out, err := cmd.CombinedOutput()
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) != "" {
			outputs = append(outputs, string(line))
		}
	}

	return outputs, err
}

func gitCmd(args ...string) *cmd.Cmd {
	cmd := cmd.New("git")

	for _, v := range GlobalFlags {
		cmd.WithArg(v)
	}

	for _, a := range args {
		cmd.WithArg(a)
	}

	return cmd
}

func IsBuiltInGitCommand(command string) bool {
	helpCommandOutput, err := gitOutput("help", "-a")
	if err != nil {
		return false
	}
	for _, helpCommandOutputLine := range helpCommandOutput {
		if strings.HasPrefix(helpCommandOutputLine, "  ") {
			for _, gitCommand := range strings.Split(helpCommandOutputLine, " ") {
				if gitCommand == command {
					return true
				}
			}
		}
	}
	return false
}
