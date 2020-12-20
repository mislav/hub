package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

func TestDirIsNotEmpty(t *testing.T) {
	dir, err := ioutil.TempDir("", "gh-utils-test-")
	if err != nil {
		t.Fatal(err)
	}
	f, err := ioutil.TempFile(dir, "gh-file-test-")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	defer os.RemoveAll(dir)

	assert.T(t, !isEmptyDir(dir))
}

func TestDirIsEmpty(t *testing.T) {
	dir, err := ioutil.TempDir("", "gh-utils-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	assert.T(t, isEmptyDir(dir))
}
