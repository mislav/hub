package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParsePullRequestId(t *testing.T) {
	url := "https://github.com/jingweno/gh/pull/73"
	assert.Equal(t, "73", parsePullRequestId(url))

	url = "https://github.com/jingweno/gh/pull/"
	assert.Equal(t, "", parsePullRequestId(url))
}

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
