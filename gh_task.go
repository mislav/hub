// +build gotask

package main

import (
	"github.com/jingweno/gotask/tasking"
	"os"
	"runtime"
)

// NAME
//    cross-compile-all - cross-compiles gh for all supported platforms.
//
// DESCRIPTION
//    Cross-compiles gh for all supported platforms. Build artifacts will be in target/VERSION.
//    This only works on darwin with Vagrant setup.
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
	t.Log("Downloading gh...")
	err = t.Exec("vagrant ssh -c 'rm -rf ~/gocode && go get github.com/jingweno/gh'")
	if err != nil {
		t.Errorf("Can't download gh on linux: %s\n", err)
		return
	}

	t.Log("Cross-compiling gh...")
	err = t.Exec("vagrant ssh -c 'cd ~/gocode/src/github.com/jingweno/gh && ./script/bootstrap && GOPATH=`godep path`:$GOPATH gotask cross-compile'")
	if err != nil {
		t.Errorf("Can't cross-compile gh on linux: %s\n", err)
		return
	}

	t.Log("Moving build artifacts...")
	err = t.Exec("vagrant ssh -c 'cp -R ~/gocode/src/github.com/jingweno/gh/target/* ~/target/'")
	if err != nil {
		t.Errorf("Can't cross-compile gh on linux: %s\n", err)
		return
	}
}

// NAME
//    cross-compile - cross-compiles gh for current platform.
//
// DESCRIPTION
//    Cross-compiles gh for current platform. Build artifacts will be in target/VERSION
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
