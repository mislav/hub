package commands

import (
	"github.com/bmizerany/assert"
	"github.com/jingweno/gh/github"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestTransformRemoteArgs(t *testing.T) {
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

	github.DefaultConfigFile = "./test_support/gh"
	config := github.Config{User: "jingweno", Token: "123"}
	github.SaveConfig(&config)
	defer os.RemoveAll(filepath.Dir(github.DefaultConfigFile))

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
}
