package github

import (
	"github.com/github/hub/v2/git"
)

func IsHTTPSProtocol() bool {
	httpProtocol, _ := git.Config("hub.protocol")
	return httpProtocol == "https"
}
