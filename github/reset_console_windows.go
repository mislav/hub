//go:build windows
// +build windows

package github

import "github.com/github/hub/v2/cmd"

// This does nothing on windows
func setConsole(cmd *cmd.Cmd) {
}
