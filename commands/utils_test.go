package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

func TestReadMsg(t *testing.T) {
	title, body := readMsg("")
	assert.Equal(t, "", title)
	assert.Equal(t, "", body)

	title, body = readMsg("my pull title")
	assert.Equal(t, "my pull title", title)
	assert.Equal(t, "", body)

	title, body = readMsg("my pull title\n\nmy description\n\nanother line")
	assert.Equal(t, "my pull title", title)
	assert.Equal(t, "my description\n\nanother line", body)

	title, body = readMsg("my pull\ntitle\n\nmy description\n\nanother line")
	assert.Equal(t, "my pull title", title)
	assert.Equal(t, "my description\n\nanother line", body)
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
