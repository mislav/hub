package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func _TestWithArg(t *testing.T) {
	execCmd := NewExecCmd("git")
	execCmd.WithArg("log").WithArg("--no-color")
	assert.Equal(t, "git", execCmd.Name)
	assert.Equal(t, 2, len(execCmd.Args))
}
