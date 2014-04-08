package commands

import (
	"github.com/bmizerany/assert"
	"os"
	"testing"
)

func TestParseCherryPickProjectAndSha(t *testing.T) {
	ref := "https://github.com/jingweno/gh/commit/a319d88#comments"
	project, sha := parseCherryPickProjectAndSha(ref)

	assert.Equal(t, "jingweno", project.Owner)
	assert.Equal(t, "gh", project.Name)
	assert.Equal(t, "a319d88", sha)

	ref = "https://github.com/jingweno/gh/commit/a319d88#comments"
	project, sha = parseCherryPickProjectAndSha(ref)

	assert.Equal(t, "jingweno", project.Owner)
	assert.Equal(t, "gh", project.Name)
	assert.Equal(t, "a319d88", sha)
}

func TestTransformCherryPickArgs(t *testing.T) {
	os.Setenv("HUB_PROTOCOL", "git")
	args := NewArgs([]string{"cherry-pick", "https://github.com/jingweno/gh/commit/a319d88#comments"})
	transformCherryPickArgs(args)

	cmds := args.Commands()
	assert.Equal(t, 2, len(cmds))
	assert.Equal(t, "git remote add -f jingweno git://github.com/jingweno/gh.git", cmds[0].String())
	assert.Equal(t, "git cherry-pick a319d88", cmds[1].String())
}
