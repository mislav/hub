// +build  !windows
package github

import (
	"syscall"
)

func resetConsole() {
	fd, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0660)
	if err == nil {
		syscall.Dup2(fd, syscall.Stdin)
	}
}
