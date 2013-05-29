package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParseNameAndOwner(t *testing.T) {
	owner, name := parseOwnerAndName()
	assert.Equal(t, "gh", name)
	assert.Equal(t, "jingweno", owner)
}

func TestMustMatchGitUrl(t *testing.T) {
	url, _ := mustMatchGitUrl("git://github.com/jingweno/gh.git")
	assert.Equal(t, "git://github.com/jingweno/gh.git", url[0])

	url, _ = mustMatchGitUrl("git@github.com:jingweno/gh.git")
	assert.Equal(t, "git@github.com:jingweno/gh.git", url[0])

	url, _ = mustMatchGitUrl("https://github.com/jingweno/gh.git")
	assert.Equal(t, "https://github.com/jingweno/gh.git", url[0])
}
