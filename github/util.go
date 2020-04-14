package github

import (
	"github.com/github/hub/git"
)

func IsHTTPSProtocol() bool {
	httpProtocol, _ := git.Config("hub.protocol")
	return httpProtocol == "https"
}
