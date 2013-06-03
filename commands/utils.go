package commands

import (
	"github.com/jingweno/gh/cmd"
	"github.com/jingweno/gh/utils"
)

func browserCommand(url string) error {
	launcher, err := utils.BrowserLauncher()
	if err != nil {
		return err
	}

	launcher = append(launcher, url)
	c := cmd.NewWithArray(launcher)
	return c.Exec()
}
