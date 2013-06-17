package commands

import (
	"fmt"
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

func removeItem(slice []string, index int) (newSlice []string, item string) {
	if index > len(slice)-1 {
		panic(fmt.Sprintf("Index %d is out of bound", index))
	}

	item = slice[index]
	newSlice = append(slice[:index], slice[index+1:]...)

	return newSlice, item
}
