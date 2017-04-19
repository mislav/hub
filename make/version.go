package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func hubVersion(path string) (string, error) {
	output, err := runCmdOutput("git", "describe", "--tags", "HEAD")
	if err == nil {
		v := strings.TrimPrefix(output[0], "v")
		return strings.TrimSpace(v), nil
	}

	vf := filepath.Join(path, "commands", "version.go")
	content, err := ioutil.ReadFile(vf)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`var Version = "(.+)"`)
	if !r.Match(content) {
		return "", fmt.Errorf("Can't find version in %s", vf)
	}

	version := string(r.FindSubmatch(content)[1])
	output, err = runCmdOutput("git", "rev-parse", "--short", "HEAD")
	if err == nil {
		headSha := strings.TrimSpace(output[0])
		version = fmt.Sprintf("%s-q%s", version, headSha)
	}

	return version, nil
}

func version(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	version, err := hubVersion(path)
	if err != nil {
		return err
	}

	fmt.Println(version)
	return nil
}
