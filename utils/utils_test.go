package utils

import (
	"os"

	"github.com/bmizerany/assert"
	"testing"
)

func TestSearchBrowserLauncher(t *testing.T) {
	browser := searchBrowserLauncher("darwin")
	assert.Equal(t, "open", browser)

	browser = searchBrowserLauncher("windows")
	assert.Equal(t, "cmd /c start", browser)
}

func TestConcatPaths(t *testing.T) {
	assert.Equal(t, "foo/bar/baz", ConcatPaths("foo", "bar", "baz"))
}

func TestCleanDirName(t *testing.T) {
	if err := os.MkdirAll("test-clean-dir-name", 0700); err != nil {
		t.Fatalf("Impossible to create temp dir for testing: %v", err)
	}
	if err := os.Chdir("test-clean-dir-name"); err != nil {
		t.Fatalf("Impossible to switch to temp dir for testing: %v", err)
	}
	dirs := []string{"foo", "with multiple spaces", "with-hyphen"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatalf("Impossible to create temp dir %q for testing: %v", dir, err)
		}
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"", "test-clean-dir-name"},
		{"foo", "foo"},
		{"with multiple spaces", "with-multiple-spaces"},
		{"with-hyphen", "with-hyphen"},
	}

	for _, test := range tests {
		if got, err := CleanDirName(test.input); err != nil {
			t.Errorf("CleanDirName(%q) raised %q", test.input, err)
		} else if want := test.expected; got != want {
			t.Errorf("CleanDirName(%q) = %q, want %q", test.input, got, want)
		}
	}

	if err := os.Chdir(".."); err != nil {
		t.Fatalf("Impossible to switch back from temp dir after testing: %v", err)
	}
	if err := os.RemoveAll("test-clean-dir-nam"); err != nil {
		t.Fatalf("Impossible to clean up after testing: %v", err)
	}
}
