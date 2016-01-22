package cmd

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestNew(t *testing.T) {
	execCmd := New("vim --noplugin")
	assert.Equal(t, "vim", execCmd.Name)
	assert.Equal(t, 1, len(execCmd.Args))
	assert.Equal(t, "--noplugin", execCmd.Args[0])
}

func TestWithArg(t *testing.T) {
	execCmd := New("git")
	execCmd.WithArg("command").WithArg("--amend").WithArg("-m").WithArg(`""`)
	assert.Equal(t, "git", execCmd.Name)
	assert.Equal(t, 4, len(execCmd.Args))
}
