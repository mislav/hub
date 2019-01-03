package utils

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestColorBrightness(t *testing.T) {
	c, err := NewColor("880000")
	assert.Equal(t, nil, err)
	actual := c.Brightness()
	assert.Equal(t, float32(0.15946665406227112), actual)
}

func TestRoundDown(t *testing.T) {
	assert.Equal(t, 3.0, round(3.0))
	assert.Equal(t, 3.0, round(3.01))
	assert.Equal(t, 3.0, round(3.49))
}

func TestRoundUp(t *testing.T) {
	assert.Equal(t, 3.0, round(2.5))
	assert.Equal(t, 3.0, round(2.51))
	assert.Equal(t, 3.0, round(2.99))
}

func TestRoundNegative(t *testing.T) {
	assert.Equal(t, -2.0, round(-2.49))
	assert.Equal(t, -3.0, round(-2.5))
	assert.Equal(t, -3.0, round(-2.51))
}
