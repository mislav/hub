package cmd

import (
	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
	"testing"
)

func TestNew(t *testing.T) {
	execCmd := New("vim --noplugin")
	assert.Equal(t, "vim", execCmd.Name)
	assert.Equal(t, 1, len(execCmd.Args))
	assert.Equal(t, "--noplugin", execCmd.Args[0])
}

func TestWithArg(t *testing.T) {
	execCmd := New("git")
	execCmd.WithArg("log").WithArg("--no-color")
	assert.Equal(t, "git", execCmd.Name)
	assert.Equal(t, 2, len(execCmd.Args))
}
