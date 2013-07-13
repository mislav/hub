package commands

import (
	"github.com/bmizerany/assert"
	"github.com/jingweno/gh/github"
	"os"
	"path/filepath"
	"testing"
)

func TestTransformSubmoduleArgs(t *testing.T) {
	github.DefaultConfigFile = "./test_support/clone_gh"
	config := github.Config{User: "jingweno", Token: "123"}
	github.SaveConfig(&config)
	defer os.RemoveAll(filepath.Dir(github.DefaultConfigFile))

	args := NewArgs([]string{"submodule", "add", "foo/gh", "foo/gh"})
	transformSubmoduleArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "git://github.com/foo/gh.git", args.GetParam(1))

	args = NewArgs([]string{"submodule", "add", "-p", "foo/gh", "foo/gh"})
	transformSubmoduleArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "git@github.com:foo/gh.git", args.GetParam(1))

	args = NewArgs([]string{"submodule", "add", "jingweno/gh", "jingweno/gh"})
	transformSubmoduleArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "git@github.com:jingweno/gh.git", args.GetParam(1))

	args = NewArgs([]string{"submodule", "add", "-p", "acl-services/devise-acl", "foo/devise-acl"})
	transformSubmoduleArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "git@github.com:acl-services/devise-acl.git", args.GetParam(1))

	args = NewArgs([]string{"submodule", "add", "jekyll_and_hyde", "foo/jekyll_and_hyde"})
	transformSubmoduleArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "git://github.com/jingweno/jekyll_and_hyde.git", args.GetParam(1))

	args = NewArgs([]string{"submodule", "add", "-p", "jekyll_and_hyde", "foo/jekyll_and_hyde"})
	transformSubmoduleArgs(args)

	assert.Equal(t, 3, args.ParamsSize())
	assert.Equal(t, "git@github.com:jingweno/jekyll_and_hyde.git", args.GetParam(1))
}
