package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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

func CommandPath(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		cmd = cmd + ".exe"
	}

	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", err
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(path)
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

type Color struct {
	Red   int64
	Green int64
	Blue  int64
}

func NewColor(hex string) (*Color, error) {
	red, err := strconv.ParseInt(hex[0:2], 16, 16)
	if err != nil {
		return nil, err
	}
	green, err := strconv.ParseInt(hex[2:4], 16, 16)
	if err != nil {
		return nil, err
	}
	blue, err := strconv.ParseInt(hex[4:6], 16, 16)
	if err != nil {
		return nil, err
	}

	return &Color{
		Red:   red,
		Green: green,
		Blue:  blue,
	}, nil
}

func (c *Color) Brightness() float32 {
	return (0.299*float32(c.Red) + 0.587*float32(c.Green) + 0.114*float32(c.Blue)) / 255
}
