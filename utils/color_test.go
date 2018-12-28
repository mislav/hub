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
