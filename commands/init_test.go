package commands

import (
	"github.com/bmizerany/assert"
	"github.com/jingweno/gh/github"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestTransformInitArgs(t *testing.T) {
	github.DefaultConfigFile = "./test_support/gh"
	config := github.Config{User: "jingweno", Token: "123"}
	github.SaveConfig(&config)
	defer os.RemoveAll(filepath.Dir(github.DefaultConfigFile))

	args := NewArgs([]string{"init"})
	err := transformInitArgs(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, args.IsParamsEmpty())

	args = NewArgs([]string{"init", "-g"})
	err = transformInitArgs(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, args.IsParamsEmpty())

	commands := args.Commands()
	assert.Equal(t, 2, len(commands))
	assert.Equal(t, "git init", commands[0].String())
	reg := regexp.MustCompile("git remote add origin git@github.com:jingweno/.+\\.git")
	assert.T(t, reg.MatchString(commands[1].String()))
}
