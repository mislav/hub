package commands

import (
	"testing"

	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
	"github.com/github/hub/cmd"
)

func TestTransformAprArgs(t *testing.T) {
	SetRemoteOriginUrl("https://github.com/hub.git")

	args := NewArgs([]string{"apr", "33"})
	transformAprArgs(args)
	cmd := args.ToCmd()
	assert.Equal(t, "git am -3 https://github.com/hub/pull/33", cmd.String())

	args = NewArgs([]string{"apr", "#33"})
	transformAprArgs(args)
	cmd = args.ToCmd()
	assert.Equal(t, "git am -3 https://github.com/hub/pull/33", cmd.String())

	args = NewArgs([]string{"apr", "-q", "33"})
	transformAprArgs(args)
	cmd = args.ToCmd()
	assert.Equal(t, "git am -q -3 https://github.com/hub/pull/33", cmd.String())

	SetRemoteOriginUrl("https://github.com/hub/")

	args = NewArgs([]string{"apr", "33"})
	transformAprArgs(args)
	cmd = args.ToCmd()
	assert.Equal(t, "git am -3 https://github.com/hub/pull/33", cmd.String())

	SetRemoteOriginUrl("https://github.com/hub")

	args = NewArgs([]string{"apr", "33"})
	transformAprArgs(args)
	cmd = args.ToCmd()
	assert.Equal(t, "git am -3 https://github.com/hub/pull/33", cmd.String())
}

func SetRemoteOriginUrl(url string) {
	cmd1 := cmd.New("git config --unset remote.origin.url")
	cmd2 := cmd.New("git config --add remote.origin.url " + url)
	cmd1.Run()
	cmd2.Run()
}
