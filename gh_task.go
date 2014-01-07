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
		"github.com/laher/goxc",
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

	path := fmt.Sprintf("%s%c%s", filepath.Join(gopath, "bin"), os.PathListSeparator, os.Getenv("PATH"))
	os.Setenv("PATH", path)
	t.Logf("PATH=%s\n", path)

	t.Logf("Packaging for %s...\n", runtime.GOOS)

	t.Log("Installing dependencies...")
	TaskInstallDeps(t)

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	ghPath := filepath.Join(gopath, "src", "github.com", "jingweno", "gh")
	t.Logf("Copying source from %s to %s\n", pwd, ghPath)
	err = copyDir(pwd, ghPath)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir(ghPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Cross-compiling...")
	godepPath := filepath.Join(ghPath, "Godeps", "_workspace")
	gopath = fmt.Sprintf("%s%c%s", gopath, os.PathListSeparator, godepPath)
	os.Setenv("GOPATH", gopath)
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
	err = copyBuildArtifacts(source, target)
	if err != nil {
		t.Fatal(err)
	}
}

// NAME
//    cross-compile - cross-compiles gh for current platform.
//
// DESCRIPTION
//    Cross-compiles gh for current platform. Build artifacts will be in target/VERSION
func TaskCrossCompile(t *tasking.T) {
	t.Logf("Cross-compiling gh for %s...\n", runtime.GOOS)
	t.Logf("GOPATH=%s\n", os.Getenv("GOPATH"))
	err := t.Exec("goxc", "-wd=.", "-os="+runtime.GOOS, "-c="+runtime.GOOS)
	if err != nil {
		t.Fatalf("Can't cross-compile gh: %s\n", err)
	}
}

func copyBuildArtifacts(srcDir, destDir string) error {
	artifacts := findBuildArtifacts(srcDir)
	for _, artifact := range artifacts {
		target := filepath.Join(destDir, filepath.Base(artifact))
		fmt.Printf("Copying %s to %s\n", artifact, target)
		err := copyFile(artifact, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(source, dest string) error {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)

	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}
	}

	return err
}

func copyDir(source, dest string) (err error) {
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("Source is not a directory")
	}

	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return fmt.Errorf("Destination already exists")
	}

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)
	for _, entry := range entries {
		sfp := filepath.Join(source, entry.Name())
		dfp := filepath.Join(dest, entry.Name())
		if entry.IsDir() {
			err = copyDir(sfp, dfp)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(sfp, dfp)
			if err != nil {
				return err
			}
		}
	}

	return
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
