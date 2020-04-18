package commands

import (
	"github.com/github/hub/v2/internal/assert"
	"testing"
)

func TestParseRange(t *testing.T) {
	s := "1.0..2.0"
	assert.Equal(t, "1.0...2.0", parseCompareRange(s))

	s = "1.0...2.0"
	assert.Equal(t, "1.0...2.0", parseCompareRange(s))
}
