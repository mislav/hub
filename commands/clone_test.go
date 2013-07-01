package commands

import (
	"github.com/bmizerany/assert"
	"github.com/jingweno/gh/github"
	"testing"
)

func TestTransformCloneArgs(t *testing.T) {
	args := NewArgs([]string{"jingweno/gh"})
	transformCloneArgs(args)

	assert.Equal(t, 1, args.Size())
	assert.Equal(t, "git://github.com/jingweno/gh.git", args.First())

	args = NewArgs([]string{"-p", "jingweno/gh"})
	transformCloneArgs(args)

	assert.Equal(t, 1, args.Size())
	assert.Equal(t, "git@github.com:jingweno/gh.git", args.First())

	args = NewArgs([]string{"-p", "jekyll_and_hyde"})
	config := github.Config{User: "jingweno", Token: "123"}
	github.SaveConfig(&config)
	transformCloneArgs(args)

	assert.Equal(t, 1, args.Size())
	assert.Equal(t, "git@github.com:jingweno/jekyll_and_hyde.git", args.First())
}

func TestParseCloneNameAndOwner(t *testing.T) {
	arg := "jekyll_and_hyde"
	config := github.Config{User: "jingweno", Token: "123"}
	github.SaveConfig(&config)

	name, owner := parseCloneNameAndOwner(arg)
	assert.Equal(t, "jekyll_and_hyde", name)
	assert.Equal(t, "jingweno", owner)
}
