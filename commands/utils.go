package commands

import (
	"github.com/jingweno/gh/cmd"
	"github.com/jingweno/gh/utils"
	"os"
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

func isDir(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.IsDir()
}
