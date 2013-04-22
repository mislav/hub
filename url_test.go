package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestMustMatchGitUrl(t *testing.T) {
	assert.Equal(t, "git://github.com/jingweno/gh.git", mustMatchGitUrl("git://github.com/jingweno/gh.git")[0])
	assert.Equal(t, "git@github.com:jingweno/gh.git", mustMatchGitUrl("git@github.com:jingweno/gh.git")[0])
	assert.Equal(t, "https://github.com/jingweno/gh.git", mustMatchGitUrl("https://github.com/jingweno/gh.git")[0])
}
