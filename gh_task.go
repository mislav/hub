// +build gotask

package main

import (
	"github.com/jingweno/gotask/tasking"
	"os"
	"runtime"
)

// Releases gh
//
// Release gh for current operating system. The build artifacts will be in target/VERSION
func TaskRelease(t *tasking.T) {
	t.Log("Updating goxc...")
	err := t.Exec("go get -u github.com/laher/goxc")
	if err != nil {
		t.Errorf("Can't update goxc: %s\n", err)
		return
	}

	t.Log("Removing build target...")
	err = os.RemoveAll("target")
	if err != nil {
		t.Errorf("Can't remove build target: %s\n", err)
		return
	}

	t.Log("Building gh...")
	err = t.Exec("goxc", "-wd=.", "-os="+runtime.GOOS, "-c="+runtime.GOOS)
	if err != nil {
		t.Errorf("Can't build gh: %s\n", err)
		return
	}
}
