// +build !windows

package github

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/github/hub/git"
)

func IsHttpsProtocol() bool {
	httpProcotol, _ := git.Config("hub.protocol")
	if httpProcotol == "https" {
		return true
	}

	httpClone, _ := git.Config("--bool hub.http-clone")
	if httpClone == "true" {
		return true
	}

	return false
}

func isTerminal(fd uintptr) bool {
	return terminal.IsTerminal(int(fd))
}
