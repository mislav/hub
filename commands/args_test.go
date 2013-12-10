package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewArgs(t *testing.T) {
	args := NewArgs([]string{})
	assert.Equal(t, "", args.Command)
	assert.Equal(t, 0, args.ParamsSize())

	args = NewArgs([]string{"command"})
	assert.Equal(t, "command", args.Command)
	assert.Equal(t, 0, args.ParamsSize())

	args = NewArgs([]string{"command", "args"})
	assert.Equal(t, "command", args.Command)
	assert.Equal(t, 1, args.ParamsSize())
}

func TestArgs_Words(t *testing.T) {
	args := NewArgs([]string{"--no-ff", "master"})
	a := args.Words()

	assert.Equal(t, 1, len(a))
	assert.Equal(t, "master", a[0])
}

func TestArgs_Insert(t *testing.T) {
	args := NewArgs([]string{"command", "1", "2", "3", "4"})
	args.InsertParam(0, "foo")

	assert.Equal(t, 5, args.ParamsSize())
	assert.Equal(t, "foo", args.FirstParam())

	args = NewArgs([]string{"command", "1", "2", "3", "4"})
	args.InsertParam(3, "foo")

	assert.Equal(t, 5, args.ParamsSize())
	assert.Equal(t, "foo", args.Params[3])
}

func TestArgs_Remove(t *testing.T) {
	args := NewArgs([]string{"1", "2", "3", "4"})

	item := args.RemoveParam(1)
	assert.Equal(t, "3", item)
	assert.Equal(t, 2, args.ParamsSize())
	assert.Equal(t, "2", args.FirstParam())
	assert.Equal(t, "4", args.GetParam(1))
}
