package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParseRemoteNames(t *testing.T) {
	args := NewArgs([]string{"fetch", "jingweno,foo"})
	names := parseRemoteNames(args)

	assert.Equal(t, 2, len(names))
	assert.Equal(t, "jingweno", names[0])
	assert.Equal(t, "foo", names[1])
	cmd := args.ToCmd()
	assert.Equal(t, "git fetch --multiple jingweno foo", cmd.String())

	args = NewArgs([]string{"fetch", "--multiple", "jingweno", "foo"})
	names = parseRemoteNames(args)
	assert.Equal(t, 2, len(names))
	assert.Equal(t, "jingweno", names[0])
	assert.Equal(t, "foo", names[1])

	args = NewArgs([]string{"fetch", "--multiple"})
	names = parseRemoteNames(args)
	assert.Equal(t, 0, len(names))
}
