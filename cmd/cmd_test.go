package cmd

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestWithArg(t *testing.T) {
	execCmd := New("git")
	execCmd.WithArg("log").WithArg("--no-color")
	assert.Equal(t, "git", execCmd.Name)
	assert.Equal(t, 2, len(execCmd.Args))
}
