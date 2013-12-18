// +build gotask

package main

import (
	"fmt"
	"github.com/jingweno/gotask/tasking"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// NAME
//    install-deps - install dependencies with go get
//
// DESCRIPTION
//    Install dependencies with go get.
func TaskInstallDeps(t *tasking.T) {
	deps := []string{
		"github.com/kr/godep",
		"github.com/laher/goxc",
		"github.com/jingweno/gh",
	}

	for _, dep := range deps {
		t.Logf("Installing %s\n", dep)
		err := t.Exec("go get", dep)
		if err != nil {
			t.Fatalf("Can't download dependency %s", err)
		}
	}
}

// NAME
//    package - cross compile gh and package it
//
// DESCRIPTION
//    Cross compile gh and package it into PWD/target
func TaskPackage(t *tasking.T) {
	gopath, err := ioutil.TempDir("", "gh-build")
	os.Setenv("GOPATH", gopath)
	t.Logf("GOPATH=%s\n", gopath)

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Packaging for %s...\n", runtime.GOOS)

	t.Log("Installing dependencies...")
	TaskInstallDeps(t)

	ghPath := filepath.Join(gopath, "src", "github.com", "jingweno", "gh")
	err = os.Chdir(ghPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Cross-compiling...")
	godepPath := filepath.Join(ghPath, "Godeps", "_workspace")
	os.Setenv("GOPATH", fmt.Sprintf("%s:%s", godepPath, gopath))
	TaskCrossCompile(t)

	source := filepath.Join(ghPath, "target")
	target := filepath.Join(pwd, "target")
	t.Logf("Copying build artifacts from %s to %s...\n", source, target)
	_, err = os.Stat(target)
	if err != nil {
		err = os.Mkdir(target, 0777)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Args = append(t.Args, source, target)
	TaskCopyBuildArtifacts(t)
}

// NAME
//    cross-compile - cross-compiles gh for current platform.
//
// DESCRIPTION
//    Cross-compiles gh for current platform. Build artifacts will be in target/VERSION
func TaskCrossCompile(t *tasking.T) {
	t.Logf("Cross-compiling gh for %s...\n", runtime.GOOS)
	err := t.Exec("goxc", "-wd=.", "-os="+runtime.GOOS, "-c="+runtime.GOOS)
	if err != nil {
		t.Fatalf("Can't cross-compile gh: %s\n", err)
	}
}

// NAME
//    copy-build-artifacts - copy build artifacts from source to target
//
// DESCRIPTION
//    Copy build artifacts from source to target.
//    For example, `gotask copy-build-artifacts src dest`
func TaskCopyBuildArtifacts(t *tasking.T) {
	if len(t.Args) < 2 {
		t.Fatal("Missing source or target")
	}

	srcDir := t.Args[0]
	destDir := t.Args[1]

	artifacts := findBuildArtifacts(srcDir)
	for _, artifact := range artifacts {
		target := filepath.Join(destDir, filepath.Base(artifact))
		t.Logf("Copying %s to %s\n", artifact, target)
		err := copyFile(artifact, target)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}

func findBuildArtifacts(root string) (artifacts []string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext == ".deb" || ext == ".zip" || ext == ".gz" {
			artifacts = append(artifacts, path)
		}

		return nil
	})

	return
}
