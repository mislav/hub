package commands

import (
	"os"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestParseCherryPickProjectAndSha(t *testing.T) {
	testConfigs := fixtures.SetupTestConfigs()
	defer testConfigs.TearDown()

	ref := "https://github.com/jingweno/gh/commit/a319d88#comments"
	project, sha := parseCherryPickProjectAndSha(ref)

	assert.Equal(t, "jingweno", project.Owner)
	assert.Equal(t, "gh", project.Name)
	assert.Equal(t, "github.com", project.Host)
	assert.Equal(t, "https", project.Protocol)
	assert.Equal(t, "a319d88", sha)

	ref = "https://github.com/jingweno/gh/commit/a319d88#comments"
	project, sha = parseCherryPickProjectAndSha(ref)

	assert.Equal(t, "jingweno", project.Owner)
	assert.Equal(t, "gh", project.Name)
	assert.Equal(t, "a319d88", sha)
}

func TestTransformCherryPickArgs(t *testing.T) {
	testConfigs := fixtures.SetupTestConfigs()
	defer testConfigs.TearDown()

	args := NewArgs([]string{})
	transformCherryPickArgs(args)
	cmds := args.Commands()
	assert.Equal(t, 1, len(cmds))

	os.Setenv("HUB_PROTOCOL", "git")
	defer os.Setenv("HUB_PROTOCOL", "")
	args = NewArgs([]string{"cherry-pick", "https://github.com/jingweno/gh/commit/a319d88#comments"})
	transformCherryPickArgs(args)

	cmds = args.Commands()
	assert.Equal(t, 2, len(cmds))
	assert.Equal(t, "git remote add -f --no-tags jingweno git://github.com/jingweno/gh.git", cmds[0].String())
	assert.Equal(t, "git cherry-pick a319d88", cmds[1].String())
}
