package cmd

import (
	"os"
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

func TestNew(t *testing.T) {
	execCmd := New("vim --noplugin")
	assert.Equal(t, "vim --noplugin", execCmd.Name)
	assert.Equal(t, 0, len(execCmd.Args))
}

func TestWithArg(t *testing.T) {
	execCmd := New("git")
	execCmd.WithArg("command").WithArg("--amend").WithArg("-m").WithArg(`""`)
	assert.Equal(t, "git", execCmd.Name)
	assert.Equal(t, 4, len(execCmd.Args))
}

func TestInvokingShell(t *testing.T) {
	sh := NewWithShell([]string{"$FOO", "hello"})
	sh.WithArg("happy world")
	defer func() {
		os.Unsetenv("FOO")
	}()

	os.Setenv("FOO", "echo")

	output, err := sh.Output()
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello happy world\n", output)
}

func TestInvokingShellComplex(t *testing.T) {
	sh := NewWithShell([]string{"$FOO hey", "hello"})
	sh.WithArg("happy world")
	defer func() {
		os.Unsetenv("FOO")
	}()

	os.Setenv("FOO", "echo")

	output, err := sh.Output()
	assert.Equal(t, nil, err)
	assert.Equal(t, "hey hello happy world\n", output)
}
