package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func FetchGitDir() (string, error) {
	output, err := execGitCmd([]string{"rev-parse", "-q", "--git-dir"})
	if err != nil {
		return "", err
	}

	gitDir := output[0]
	gitDir, err = filepath.Abs(gitDir)
	if err != nil {
		return "", err
	}

	return gitDir, nil
}

func FetchGitEditor() (string, error) {
	output, err := execGitCmd([]string{"var", "GIT_EDITOR"})
	if err != nil {
		return "", err
	}

	return output[0], nil
}

func FetchGitOwner() (string, error) {
	remote, err := FetchGitRemote()
	if err != nil {
		return "", err
	}

	return mustMatchGitUrl(remote)[1], nil
}

func FetchGitProject() (string, error) {
	remote, err := FetchGitRemote()
	if err != nil {
		return "", err
	}

	return mustMatchGitUrl(remote)[2], nil
}

func FetchGitHead() (string, error) {
	output, err := execGitCmd([]string{"symbolic-ref", "-q", "--short", "HEAD"})
	if err != nil {
		return "master", err
	}

	return output[0], nil
}

// FIXME: only care about origin push remote now
func FetchGitRemote() (string, error) {
	r := regexp.MustCompile("origin\t(.+github.com.+) \\(push\\)")
	output, err := execGitCmd([]string{"remote", "-v"})
	if err != nil {
		return "", err
	}

	for _, o := range output {
		if r.MatchString(o) {
			return r.FindStringSubmatch(o)[1], nil
		}
	}

	return "", errors.New("Can't find remote")
}

func FetchGitCommitLogs(sha1, sha2 string) (string, error) {
	execCmd := NewExecCmd("git")
	execCmd.WithArg("log").WithArg("--no-color")
	execCmd.WithArg("--format=%h (%aN, %ar)%n%w(78,3,3)%s%n%+b")
	execCmd.WithArg("--cherry")
	shaRange := fmt.Sprintf("%s...%s", sha1, sha2)
	execCmd.WithArg(shaRange)

	outputs, err := execCmd.ExecOutput()
	if err != nil {
		return "", err
	}

	return outputs, nil
}

func execGitCmd(input []string) (outputs []string, err error) {
	cmd := NewExecCmd("git")
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
