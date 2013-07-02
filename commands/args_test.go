package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewArgs(t *testing.T) {
	args := NewArgs([]string{})
	assert.Equal(t, "", args.Command)
	assert.Equal(t, 0, args.Size())

	args = NewArgs([]string{"command"})
	assert.Equal(t, "command", args.Command)
	assert.Equal(t, 0, args.Size())

	args = NewArgs([]string{"command", "args"})
	assert.Equal(t, "command", args.Command)
	assert.Equal(t, 1, args.Size())
}

func TestRemove(t *testing.T) {
	args := NewArgs([]string{"1", "2", "3", "4"})

	item := args.Remove(1)
	assert.Equal(t, "3", item)
	assert.Equal(t, 2, args.Size())
	assert.Equal(t, "2", args.First())
	assert.Equal(t, "4", args.Get(1))
}
