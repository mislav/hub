// +build windows

package github

func isTerminal(fd uintptr) bool {
	return true
}
