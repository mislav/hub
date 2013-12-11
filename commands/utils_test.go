package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParseRepoNameOwner(t *testing.T) {
	owner, repo, match := parseRepoNameOwner("jingweno")

	assert.T(t, match)
	assert.Equal(t, "jingweno", owner)
	assert.Equal(t, "", repo)

	owner, repo, match = parseRepoNameOwner("jingweno/gh")

	assert.T(t, match)
	assert.Equal(t, "jingweno", owner)
	assert.Equal(t, "gh", repo)
}
