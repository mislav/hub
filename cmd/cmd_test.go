package cmd

import (
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

func Test_String(t *testing.T) {
	c := Cmd{
		Name: "echo",
		Args: []string{"hi", "hello world", "don't", `"fake news"`},
	}
	assert.Equal(t, `echo hi "hello world" "don't" '"fake news"'`, c.String())
}
