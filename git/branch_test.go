package git

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestShortName(t *testing.T) {
	b := Branch{"refs/heads/master"}
	assert.Equal(t, "master", b.ShortName())
}
