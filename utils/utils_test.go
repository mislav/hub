package utils

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestSearchBrowserLauncher(t *testing.T) {
	browser, args := searchBrowserLauncher("darwin")
	assert.Equal(t, "open", browser)
	assert.Equal(t, "", args)

	browser, args = searchBrowserLauncher("windows")
	assert.Equal(t, "cmd /c start", browser)
	assert.Equal(t, "", args)
}

func TestConcatPaths(t *testing.T) {
	assert.Equal(t, "foo/bar/baz", ConcatPaths("foo", "bar", "baz"))
}
