package commands

import (
	"regexp"
	"testing"

	"github.com/bmizerany/assert"
)

func TestRunHelp(t *testing.T) {
	args := NewArgs([]string{"help", "clone"})
	runHelp(cmdHelp, args)

	// fallthrough to git's own help
	assert.Equal(t, "help", args.Command)
	assert.Equal(t, []string{"clone"}, args.Params)

	args = NewArgs([]string{"help", "-a"})
	runHelp(cmdHelp, args)

	cmds := args.Commands()
	assert.Equal(t, 3, len(cmds))
	assert.Equal(t, "git help -a", cmds[0].String())
	assert.Equal(t, "echo \nhub custom commands\n", cmds[1].String())
	// print out in the format of echo COMMANDS
	// the length of the split strings should be len(COMMANDS) + 1, including "echo"
	split := regexp.MustCompile("\\s+").Split(cmds[2].String(), -1)
	assert.Equal(t, len(customCommands())+1, len(split))
}
