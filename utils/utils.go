package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/github/hub/ui"
)

func Check(err error) {
	if err != nil {
		ui.Errorln(err)
		os.Exit(1)
	}
}

func ConcatPaths(paths ...string) string {
	return strings.Join(paths, "/")
}

func BrowserLauncher() ([]string, error) {
	browser := os.Getenv("BROWSER")
	if browser == "" {
		browser = searchBrowserLauncher(runtime.GOOS)
	}

	if browser == "" {
		return nil, errors.New("Please set $BROWSER to a web launcher")
	}

	return strings.Split(browser, " "), nil
}

func searchBrowserLauncher(goos string) (browser string) {
	switch goos {
	case "darwin":
		browser = "open"
	case "windows":
		browser = "cmd /c start"
	default:
		candidates := []string{"xdg-open", "cygstart", "x-www-browser", "firefox",
			"opera", "mozilla", "netscape"}
		for _, b := range candidates {
			path, err := exec.LookPath(b)
			if err == nil {
				browser = path
				break
			}
		}
	}

	return browser
}

func DirName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	name := filepath.Base(dir)
	name = strings.Replace(name, " ", "-", -1)
	return name, nil
}

func IsOption(confirm, short, long string) bool {
	return strings.EqualFold(confirm, short) || strings.EqualFold(confirm, long)
}
