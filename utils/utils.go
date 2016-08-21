package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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

func TimeAgo(t time.Time) string {
	duration := time.Since(t)
	minutes := duration.Minutes()
	hours := duration.Hours()
	days := hours / 24
	months := days / 30
	years := months / 12

	var val int
	var unit string

	if minutes < 1 {
		return "now"
	} else if hours < 1 {
		val = int(minutes)
		unit = "minute"
	} else if days < 1 {
		val = int(hours)
		unit = "hour"
	} else if months < 1 {
		val = int(days)
		unit = "day"
	} else if years < 1 {
		val = int(months)
		unit = "month"
	} else {
		val = int(years)
		unit = "year"
	}

	var plural string
	if val > 1 {
		plural = "s"
	}
	return fmt.Sprintf("%d %s%s ago", val, unit, plural)
}
