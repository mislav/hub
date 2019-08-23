package github

import (
	"github.com/github/hub/git"
)

func IsHttpsProtocol() bool {
	httpProtocol, _ := git.Config("hub.protocol")
	return httpProtocol == "https"
}
