package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestTransformToTripleDots(t *testing.T) {
	s := "1.0..2.0"
	assert.Equal(t, "1.0...2.0", transformToTripleDots(s))

	s = "1.0...2.0"
	assert.Equal(t, "1.0...2.0", transformToTripleDots(s))
}
