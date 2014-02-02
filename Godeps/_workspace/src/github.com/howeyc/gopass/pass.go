// Reads password from terminal.
package gopass

import (
	"fmt"
	"os"
)

// Returns password byte array read from terminal without input being echoed.
// Array of bytes does not include end-of-line characters.
func GetPasswd() []byte {
	pass := make([]byte, 0)
	for v := getch(); ; v = getch() {
		if v == 127 || v == 8 {
			if len(pass) > 0 {
				pass = pass[:len(pass)-1]
			}
		} else if v == 13 || v == 10 {
			break
		} else {
			pass = append(pass, v)
		}
	}
	println()
	return pass
}

// Masking password functionality
// Removed character restrictions
func GetPasswdMasked() []byte {
	secret := make([]byte, 0)
	mask := byte('*')

	pass := make([]byte, 0)
	for v := getch(); ; v = getch() {
		if v == 127 || v == 8 {
			if len(pass) > 0 {
				pass = pass[:len(pass)-1]
			}
			if len(secret) > 0 {
				secret = secret[:len(secret)-1]
				os.Stdout.Write([]byte("\b \b"))
			}
		} else if v == 13 || v == 10 {
			break
		} else {
			secret = append(secret, mask)
			fmt.Print(string(mask))
			pass = append(pass, v)
		}
	}
	println()
	return pass
}
