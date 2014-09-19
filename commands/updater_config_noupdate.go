// +build noupdate

package commands

import "os"

func init() {
	os.Setenv("HUB_AUTOUPDATE", "never")
}
