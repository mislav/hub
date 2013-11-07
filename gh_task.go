// +build gotask

package main

import (
	"github.com/jingweno/gotask/tasking"
	"os"
	"runtime"
)

// Cross-compiles gh for all supported platforms.
//
// Cross-compiles gh for all supported platforms. The build artifacts
// will be in target/VERSION. This only works on darwin with Vagrant setup.
func TaskCrossCompileAll(t *tasking.T) {
	t.Log("Removing build target...")
	err := os.RemoveAll("target")
	if err != nil {
		t.Errorf("Can't remove build target: %s\n", err)
		return
	}

	// for current
	t.Logf("Compiling for %s...\n", runtime.GOOS)
	TaskCrossCompile(t)
	if t.Failed() {
		return
	}

	// for linux
	t.Log("Compiling for linux...")
	err = t.Exec("vagrant ssh -c 'cd ~/src/github.com/jingweno/gh && git pull origin master && gotask cross-compile'")
	if err != nil {
		t.Errorf("Can't compile on linux: %s\n", err)
		return
	}
}

// Cross-compiles gh for current operating system.
//
// Cross-compiles gh for current operating system. The build artifacts will be in target/VERSION
func TaskCrossCompile(t *tasking.T) {
	t.Log("Updating goxc...")
	err := t.Exec("go get -u github.com/laher/goxc")
	if err != nil {
		t.Errorf("Can't update goxc: %s\n", err)
		return
	}

	t.Logf("Cross-compiling gh for %s...\n", runtime.GOOS)
	err = t.Exec("goxc", "-wd=.", "-os="+runtime.GOOS, "-c="+runtime.GOOS)
	if err != nil {
		t.Errorf("Can't cross-compile gh: %s\n", err)
		return
	}
}
