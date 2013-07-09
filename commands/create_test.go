package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParseCreateOwnerAndName(t *testing.T) {
	owner, name := parseCreateOwnerAndName("jingweno/gh")

	assert.Equal(t, "jingweno", owner)
	assert.Equal(t, "gh", name)

	owner, name = parseCreateOwnerAndName("gh")

	assert.Equal(t, "", owner)
	assert.Equal(t, "gh", name)

	owner, name = parseCreateOwnerAndName("jingweno/gh/foo")

	assert.Equal(t, "jingweno", owner)
	assert.Equal(t, "gh/foo", name)
}
