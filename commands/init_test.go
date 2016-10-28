package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/github"
)

func setupInitContext() {
	os.Setenv("HUB_PROTOCOL", "git")
	os.Setenv("HUB_CONFIG", "")
	github.CreateTestConfigs("jingweno", "123")
}

func TestEmptyParams(t *testing.T) {
	setupInitContext()

	args := NewArgs([]string{"init"})
	err := transformInitArgs(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, args.IsParamsEmpty())
}

func TestFlagToAddRemote(t *testing.T) {
	setupInitContext()

	args := NewArgs([]string{"init", "-g", "--quiet"})
	err := transformInitArgs(args)
	assert.Equal(t, nil, err)

	commands := args.Commands()
	assert.Equal(t, 2, len(commands))
	assert.Equal(t, "git init --quiet", commands[0].String())

	currentDir, err := os.Getwd()
	assert.Equal(t, nil, err)

	expected := fmt.Sprintf(
		"git --git-dir %s remote add origin git@github.com:jingweno/%s.git",
		filepath.Join(currentDir, ".git"),
		filepath.Base(currentDir),
	)
	assert.Equal(t, expected, commands[1].String())
}

func TestInitInAnotherDir(t *testing.T) {
	setupInitContext()

	args := NewArgs([]string{"init", "-g", "--template", "mytpl", "--shared=umask", "my project"})
	err := transformInitArgs(args)
	assert.Equal(t, nil, err)

	commands := args.Commands()
	assert.Equal(t, 2, len(commands))
	assert.Equal(t, "git init --template mytpl --shared=umask my project", commands[0].String())

	currentDir, err := os.Getwd()
	assert.Equal(t, nil, err)

	expected := fmt.Sprintf(
		"git --git-dir %s remote add origin git@github.com:jingweno/%s.git",
		filepath.Join(currentDir, "my project", ".git"),
		"my-project",
	)
	assert.Equal(t, expected, commands[1].String())
}

func TestSeparateGitDir(t *testing.T) {
	setupInitContext()

	args := NewArgs([]string{"init", "-g", "--separate-git-dir", "/tmp/where-i-play.git", "my/playground"})
	err := transformInitArgs(args)
	assert.Equal(t, nil, err)

	commands := args.Commands()
	assert.Equal(t, 2, len(commands))
	assert.Equal(t, "git init --separate-git-dir /tmp/where-i-play.git my/playground", commands[0].String())

	currentDir, err := os.Getwd()
	assert.Equal(t, nil, err)

	expected := fmt.Sprintf(
		"git --git-dir %s remote add origin git@github.com:jingweno/%s.git",
		filepath.Join(currentDir, "my", "playground", ".git"),
		"playground",
	)
	assert.Equal(t, expected, commands[1].String())
}
