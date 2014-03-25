package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

func TestGetTitleAndBodyFromFlags(t *testing.T) {
	s := "just needs raven\n\nnow it works"
	title, body, err := getTitleAndBodyFromFlags(s, "")

	assert.Equal(t, nil, err)
	assert.Equal(t, "just needs raven", title)
	assert.Equal(t, "now it works", body)

	s = "just needs raven\\n\\nnow it works"
	title, body, err = getTitleAndBodyFromFlags(s, "")

	assert.Equal(t, nil, err)
	assert.Equal(t, "just needs raven", title)
	assert.Equal(t, "now it works", body)
}

func TestDirIsNotEmpty(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)
	ioutil.TempFile(dir, "gh-utils-test-")

	assert.T(t, !isEmptyDir(dir))
}

func TestDirIsEmpty(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	assert.T(t, isEmptyDir(dir))
}

func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "gh-utils-test-")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}
