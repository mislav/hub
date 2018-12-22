package utils

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/github/hub/ui"
)

func init() {
	initColorCube()
}

var timeNow = time.Now

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

func (c *Color) Distance(other *Color) float64 {
	return math.Sqrt(float64(math.Pow(float64(c.Red - other.Red), 2) + math.Pow(float64(c.Green - other.Green), 2) + math.Pow(float64(c.Blue - other.Blue), 2)))
}

func TimeAgo(t time.Time) string {
	duration := timeNow().Sub(t)
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


var x6colorIndexes = [6]int64{ 0, 95, 135, 175, 215, 255 }
var x6colorCube [216]Color

func initColorCube() {
	i := 0
	for iR := 0; iR < 6; iR++ {
		for iG := 0; iG < 6; iG++ {
			for iB := 0; iB < 6; iB++ {
				x6colorCube[i] = Color{x6colorIndexes[iR], x6colorIndexes[iG], x6colorIndexes[iB]}
				i++
			}
		}
	}
}

func PrintColorCube() {
	for i := 0; i < 216; i++ {
		color := x6colorCube[i];
		intCode := i + 16
		code := fmt.Sprintf("5;%d", intCode)
		ui.Printf("\033[48;%sm %3d %02x %02x %02x \033[0m ",
			code, intCode, color.Red, color.Green, color.Blue)
		if i % 6 == 5 {
			ui.Printf("\n")
		}
	}
}

func ditherTo256ColorCode(color *Color) (code int) {
	iMatch := -1
	minDistance := float64(99999)
	for i := 0; i < 216; i++ {
		distance := color.Distance(&x6colorCube[i])
		if distance < minDistance {
			iMatch = i
			minDistance = distance
		}
	}
	return iMatch + 16
}

var non24bitColorTerms = []string{ "Apple_Terminal" }
var isTerm24bitColorCapableCache bool
var isTerm24bitColorCapableCacheIsInit bool = false
func isTerm24bitColorCapable() (tf bool) {
	if !isTerm24bitColorCapableCacheIsInit {
		isTerm24bitColorCapableCache = true
		myTermProg := os.Getenv("TERM_PROGRAM")
		for i := 0; i < len(non24bitColorTerms); i++ {
			if myTermProg == non24bitColorTerms[i] {
				isTerm24bitColorCapableCache = false
				break
			}
		}
		isTerm24bitColorCapableCacheIsInit = true
	}
	return isTerm24bitColorCapableCache
}

func RgbToTermColorCode(color *Color) (code string) {
	if isTerm24bitColorCapable() {
		code = fmt.Sprintf("2;%d;%d;%d", color.Red, color.Green, color.Blue)
	} else {
		intCode := ditherTo256ColorCode(color)
		code = fmt.Sprintf("5;%d", intCode)
	}
	return
}