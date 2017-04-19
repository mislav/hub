package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func build(path string) error {
	goDir, err := ioutil.TempDir("", "hub-build")
	if err != nil {
		return err
	}

	os.Setenv("GOPATH", goDir)

	hubDir := filepath.Join(goDir, "src", "github.com", "github", "hub")
	err = copyDir(path, hubDir)
	if err != nil {
		return err
	}

	version, err := hubVersion(path)
	if err != nil {
		return err
	}

	err = os.Chdir(path)
	if err != nil {
		return err
	}

	versionFlag := fmt.Sprintf("-X github.com/github/hub/commands.Version %s", version)
	return runCmd("go", "build", "-ldflags", versionFlag)
}
