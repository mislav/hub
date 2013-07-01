package commands

import (
	"github.com/bmizerany/assert"
	"regexp"
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

	args = NewArgs([]string{"jekyll_and_hyde"})
	transformCloneArgs(args)

	assert.Equal(t, 1, args.Size())
	reg := regexp.MustCompile("^git://github.com/.+/jekyll_and_hyde.git$")
	assert.T(t, reg.MatchString(args.First()))

	args = NewArgs([]string{"-p", "jekyll_and_hyde"})
	transformCloneArgs(args)

	assert.Equal(t, 1, args.Size())
	reg = regexp.MustCompile("^git@github.com:.+/jekyll_and_hyde.git$")
	assert.T(t, reg.MatchString(args.First()))
}
