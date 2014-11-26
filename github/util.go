package github

import (
	"github.com/github/hub/Godeps/_workspace/src/github.com/mattn/go-isatty"
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
	return isatty.IsTerminal(fd)
}
