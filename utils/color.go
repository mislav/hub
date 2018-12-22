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
	Red   uint8
	Green uint8
	Blue  uint8
}

func NewColor(hex string) (*Color, error) {
	red, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return nil, err
	}
	green, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return nil, err
	}
	blue, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return nil, err
	}

	return &Color{
		Red:   uint8(red),
		Green: uint8(green),
		Blue:  uint8(blue),
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

func rgbComponentToBoldValue(component uint8) float64 {
	srgb := float64(component) / 255
	if srgb <= 0.03928 {
		return srgb / 12.92
	} else {
		return math.Pow(((srgb + 0.055) / 1.055), 2.4)
	}
}

func (c *Color) Luminance() float64 {
	return 0.2126*rgbComponentToBoldValue(c.Red) +
		0.7152*rgbComponentToBoldValue(c.Green) +
		0.0722*rgbComponentToBoldValue(c.Blue)
}

func (c *Color) ContrastRatio(other *Color) float64 {
	L := c.Luminance()
	otherL := other.Luminance()
	var L1, L2 float64
	if L > otherL {
		L1, L2 = L, otherL
	} else {
		L1, L2 = otherL, L
	}
	ratio := (L1 + 0.05) / (L2 + 0.05)
	return ratio
}

type HslColor struct {
	Hue        float64
	Saturation float64
	Lightness  float64
}

func (c *Color) ToHsl() *HslColor {
	rPrime := float64(c.Red) / 255
	gPrime := float64(c.Green) / 255
	bPrime := float64(c.Blue) / 255
	cMax := math.Max(rPrime, math.Max(gPrime, bPrime))
	cMin := math.Min(rPrime, math.Min(gPrime, bPrime))
	delta := cMax - cMin

	var H float64
	if delta == 0 {
		H = 0
	} else if cMax == rPrime {
		H = 60 * math.Mod((gPrime-bPrime)/delta, 6)
	} else if cMax == gPrime {
		H = 60 * (((bPrime - rPrime) / delta) + 2)
	} else {
		H = 60 * (((rPrime - gPrime) / delta) + 4)
	}
	var L float64
	L = (cMax + cMin) / 2
	var S float64
	if delta == 0 {
		S = 0
	} else {
		S = delta / (1 - math.Abs((2*L)-1))
	}

	return &HslColor{H, S, L}
}

func (c *HslColor) ToRgb() *Color {
	C := (1 - math.Abs(2*c.Lightness-1)) * c.Saturation
	X := C * (1 - math.Abs(math.Mod(c.Hue/60, 2)-1))
	m := c.Lightness - (C / 2)
	var rPrime, gPrime, bPrime float64
	switch {
	case c.Hue < 60.0:
		rPrime, gPrime, bPrime = C, X, 0
	case c.Hue < 120.0:
		rPrime, gPrime, bPrime = X, C, 0
	case c.Hue < 180.0:
		rPrime, gPrime, bPrime = 0, C, X
	case c.Hue < 240.0:
		rPrime, gPrime, bPrime = 0, X, C
	case c.Hue < 300.0:
		rPrime, gPrime, bPrime = X, 0, C
	case c.Hue < 360.0:
		rPrime, gPrime, bPrime = C, 0, X
	}
	R := uint8((rPrime + m) * 255)
	G := uint8((gPrime + m) * 255)
	B := uint8((bPrime + m) * 255)

	return &Color{R, G, B}
}

func (c *HslColor) ScaleLightness(x float64) *HslColor {
	newLightness := math.Round(math.Min(math.Max(0, c.Lightness*(1.0+x)), 1.0))
	return &HslColor{c.Hue, c.Saturation, newLightness}
}

func (c *HslColor) ScaleSaturation(x float64) *HslColor {
	newSaturation := math.Round(math.Min(math.Max(0, c.Lightness*(1.0+x)), 1.0))
	return &HslColor{c.Hue, newSaturation, c.Lightness}
}

var x6colorIndexes = [6]uint8{0, 95, 135, 175, 215, 255}
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
