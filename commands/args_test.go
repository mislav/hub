package commands

import (
	"testing"

	"github.com/bmizerany/assert"
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

	args = NewArgs([]string{"--version"})
	assert.Equal(t, "--version", args.Command)
	assert.Equal(t, 0, args.ParamsSize())

	args = NewArgs([]string{"--help"})
	assert.Equal(t, "--help", args.Command)
	assert.Equal(t, 0, args.ParamsSize())
}

func TestArgs_Words(t *testing.T) {
	args := NewArgs([]string{"merge", "--no-ff", "master", "-m", "message"})
	a := args.Words()

	assert.Equal(t, 2, len(a))
	assert.Equal(t, "master", a[0])
	assert.Equal(t, "message", a[1])
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

	args = NewArgs([]string{"checkout", "-b"})
	args.InsertParam(1, "foo")

	assert.Equal(t, 2, args.ParamsSize())
	assert.Equal(t, "-b", args.Params[0])
	assert.Equal(t, "foo", args.Params[1])

	args = NewArgs([]string{"checkout"})
	args.InsertParam(1, "foo")

	assert.Equal(t, 1, args.ParamsSize())
	assert.Equal(t, "foo", args.Params[0])
}

func TestArgs_Remove(t *testing.T) {
	args := NewArgs([]string{"1", "2", "3", "4"})

	item := args.RemoveParam(1)
	assert.Equal(t, "3", item)
	assert.Equal(t, 2, args.ParamsSize())
	assert.Equal(t, "2", args.FirstParam())
	assert.Equal(t, "4", args.GetParam(1))
}

func TestArgs_GlobalFlags(t *testing.T) {
	args := NewArgs([]string{"-c", "key=value", "status", "-s", "-b"})
	assert.Equal(t, "status", args.Command)
	assert.Equal(t, []string{"-c", "key=value"}, args.GlobalFlags)
	assert.Equal(t, []string{"-s", "-b"}, args.Params)
	assert.Equal(t, false, args.Noop)
}

func TestArgs_GlobalFlags_Noop(t *testing.T) {
	args := NewArgs([]string{"-c", "key=value", "--noop", "--literal-pathspecs", "status", "-s", "-b"})
	assert.Equal(t, "status", args.Command)
	assert.Equal(t, []string{"-c", "key=value", "--literal-pathspecs"}, args.GlobalFlags)
	assert.Equal(t, []string{"-s", "-b"}, args.Params)
	assert.Equal(t, true, args.Noop)
}

func TestArgs_GlobalFlags_NoopTwice(t *testing.T) {
	args := NewArgs([]string{"--noop", "--bare", "--noop", "status"})
	assert.Equal(t, "status", args.Command)
	assert.Equal(t, []string{"--bare"}, args.GlobalFlags)
	assert.Equal(t, 0, len(args.Params))
	assert.Equal(t, true, args.Noop)
}

func TestArgs_GlobalFlags_Repeated(t *testing.T) {
	args := NewArgs([]string{"-C", "mydir", "-c", "a=b", "--bare", "-c", "c=d", "-c", "e=f", "status"})
	assert.Equal(t, "status", args.Command)
	assert.Equal(t, []string{"-C", "mydir", "-c", "a=b", "--bare", "-c", "c=d", "-c", "e=f"}, args.GlobalFlags)
	assert.Equal(t, 0, len(args.Params))
	assert.Equal(t, false, args.Noop)
}

func TestArgs_GlobalFlags_Propagate(t *testing.T) {
	args := NewArgs([]string{"-c", "key=value", "status"})
	cmd := args.ToCmd()
	assert.Equal(t, []string{"-c", "key=value", "status"}, cmd.Args)
}

func TestArgs_GlobalFlags_Replaced(t *testing.T) {
	args := NewArgs([]string{"-c", "key=value", "status"})
	args.Replace("open", "", "-a", "http://example.com")
	cmd := args.ToCmd()
	assert.Equal(t, "open", cmd.Name)
	assert.Equal(t, []string{"-a", "http://example.com"}, cmd.Args)
}
