package commands

import (
	"github.com/github/hub/v2/internal/assert"
	"testing"
)

func TestTransformPushArgs(t *testing.T) {
	args := NewArgs([]string{"push", "origin,staging,qa", "bert_timeout"})
	transformPushArgs(args)
	cmds := args.Commands()

	assert.Equal(t, 3, len(cmds))
	assert.Equal(t, "git push origin bert_timeout", cmds[0].String())
	assert.Equal(t, "git push staging bert_timeout", cmds[1].String())

	// TODO: travis-ci doesn't have HEAD
	//args = NewArgs([]string{"push", "origin"})
	//transformPushArgs(args)
	//cmds = args.Commands()

	//assert.Equal(t, 1, len(cmds))
	//pushRegexp := regexp.MustCompile("git push origin .+")
	//assert.T(t, pushRegexp.MatchString(cmds[0].String()))
}
