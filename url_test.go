package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestMustMatchGitUrl(t *testing.T) {
	url, _ := mustMatchGitUrl("git://github.com/jingweno/gh.git")
	assert.Equal(t, "git://github.com/jingweno/gh.git", url[0])

	url, _ = mustMatchGitUrl("git@github.com:jingweno/gh.git")
	assert.Equal(t, "git@github.com:jingweno/gh.git", url[0])

	url, _ = mustMatchGitUrl("https://github.com/jingweno/gh.git")
	assert.Equal(t, "https://github.com/jingweno/gh.git", url[0])
}
