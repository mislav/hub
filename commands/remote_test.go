package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestTransformRemoteArgs(t *testing.T) {
	args := []string{"add", "jingweno"}
	args = transformRemoteArgs(args)

	assert.Equal(t, 3, len(args))
	assert.Equal(t, "add", args[0])
	assert.Equal(t, "jingweno", args[1])
	assert.Equal(t, "git://github.com/jingweno/gh.git", args[2])

	args = []string{"add", "-p", "jingweno"}
	args = transformRemoteArgs(args)

	assert.Equal(t, 3, len(args))
	assert.Equal(t, "add", args[0])
	assert.Equal(t, "jingweno", args[1])
	assert.Equal(t, "git@github.com:jingweno/gh.git", args[2])
}
