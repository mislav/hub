package commands

import (
	"github.com/bmizerany/assert"
	"github.com/jingweno/gh/github"
	"os"
	"regexp"
	"testing"
)

func TestTransformRemoteArgs(t *testing.T) {
	os.Setenv("GH_PROTOCOL", "git")
	args := NewArgs([]string{"remote", "add", "jingweno"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "jingweno", args.GetParam(1))
	reg := regexp.MustCompile("^git://github.com/jingweno/.+\\.git$")
	assert.T(t, reg.MatchString(args.GetParam(2)))

	args = NewArgs([]string{"remote", "add", "-p", "jingweno"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "jingweno", args.GetParam(1))
	reg = regexp.MustCompile("^git@github.com:jingweno/.+\\.git$")
	assert.T(t, reg.MatchString(args.GetParam(2)))

	github.CreateTestConfig("jingweno", "123")

	args = NewArgs([]string{"remote", "add", "origin"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "origin", args.GetParam(1))
	reg = regexp.MustCompile("^git://github.com/.+/.+\\.git$")
	assert.T(t, reg.MatchString(args.GetParam(2)))

	args = NewArgs([]string{"remote", "add", "jingweno", "git@github.com:jingweno/gh.git"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "jingweno", args.GetParam(1))
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "git@github.com:jingweno/gh.git", args.GetParam(2))

	args = NewArgs([]string{"remote", "add", "-p", "origin", "org/foo"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "origin", args.GetParam(1))
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "git@github.com:org/foo.git", args.GetParam(2))
}
