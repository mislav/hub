package utils

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func init() {
	initColorCube()
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
	return (0.299*float32(c.Red) +
		0.587*float32(c.Green) +
		0.114*float32(c.Blue)) / 255
}

func (c *Color) Distance(other *Color) float64 {
	return math.Sqrt(float64(math.Pow(float64(c.Red-other.Red), 2) +
		math.Pow(float64(c.Green-other.Green), 2) +
		math.Pow(float64(c.Blue-other.Blue), 2)))
}

var x6colorIndexes = [6]int64{0, 95, 135, 175, 215, 255}
var x6colorCube [216]Color

func initColorCube() {
	i := 0
	for iR := 0; iR < 6; iR++ {
		for iG := 0; iG < 6; iG++ {
			for iB := 0; iB < 6; iB++ {
				x6colorCube[i] = Color{
					x6colorIndexes[iR],
					x6colorIndexes[iG],
					x6colorIndexes[iB],
				}
				i++
			}
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

var non24bitColorTerms = []string{
	"Apple_Terminal",
}
var isTerm24bitColorCapableCache bool
var isTerm24bitColorCapableCacheIsInit bool = false

func isTerm24bitColorCapable() bool {
	if !isTerm24bitColorCapableCacheIsInit {
		isTerm24bitColorCapableCache = true
		myTermProg := os.Getenv("TERM_PROGRAM")
		for _, brokenTerm := range non24bitColorTerms {
			if myTermProg == brokenTerm {
				isTerm24bitColorCapableCache = false
				break
			}
		}
		isTerm24bitColorCapableCacheIsInit = true
	}
	return isTerm24bitColorCapableCache
}

func RgbToTermColorCode(color *Color) string {
	if isTerm24bitColorCapable() {
		return fmt.Sprintf("2;%d;%d;%d", color.Red, color.Green, color.Blue)
	} else {
		intCode := ditherTo256ColorCode(color)
		return fmt.Sprintf("5;%d", intCode)
	}
}
