package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestTransformRemoteArgs(t *testing.T) {
	args := NewArgs([]string{"add", "jingweno"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.Size())
	assert.Equal(t, "add", args.First())
	assert.Equal(t, "jingweno", args.Get(1))
	assert.Equal(t, "git://github.com/jingweno/gh.git", args.Get(2))

	args = NewArgs([]string{"add", "-p", "jingweno"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.Size())
	assert.Equal(t, "add", args.First())
	assert.Equal(t, "jingweno", args.Get(1))
	assert.Equal(t, "git@github.com:jingweno/gh.git", args.Get(2))
}
