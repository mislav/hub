package commands

import (
	"github.com/bmizerany/assert"
	"github.com/github/hub/github"
	"os"
	"regexp"
	"testing"
)

func TestParseRepoNameOwner(t *testing.T) {
	owner, repo := parseRepoNameOwner("jingweno")
	assert.Equal(t, "jingweno", owner)
	assert.Equal(t, "", repo)

	owner, repo = parseRepoNameOwner("jingweno/gh")
	assert.Equal(t, "jingweno", owner)
	assert.Equal(t, "gh", repo)
}

func TestTransformRemoteArgs(t *testing.T) {
	os.Setenv("HUB_PROTOCOL", "git")
	github.CreateTestConfigs("jingweno", "123")

	args := NewArgs([]string{"remote", "add", "jingweno"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "jingweno", args.GetParam(1))
	reg := regexp.MustCompile("^git@github\\.com:jingweno/.+\\.git$")
	assert.T(t, reg.MatchString(args.GetParam(2)))

	args = NewArgs([]string{"remote", "add", "mislav"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "mislav", args.GetParam(1))
	reg = regexp.MustCompile("^git://github\\.com/mislav/.+\\.git$")
	assert.T(t, reg.MatchString(args.GetParam(2)))

	args = NewArgs([]string{"remote", "add", "-p", "mislav"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "mislav", args.GetParam(1))
	reg = regexp.MustCompile("^git@github\\.com:mislav/.+\\.git$")
	assert.T(t, reg.MatchString(args.GetParam(2)))

	args = NewArgs([]string{"remote", "add", "origin"})
	transformRemoteArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "add", args.FirstParam())
	assert.Equal(t, "origin", args.GetParam(1))
	reg = regexp.MustCompile("^git@github\\.com:jingweno/.+\\.git$")
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
