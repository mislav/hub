package github

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParseURL(t *testing.T) {
	url, err :=
		ParseURL("https://github.com/jingweno/gh/pulls/21")
	assert.Equal(t, nil, err)
	assert.Equal(t, "jingweno", url.Owner)
	assert.Equal(t, "gh", url.Name)
	assert.Equal(t, "pulls/21", url.ProjectPath())

	url, err =
		ParseURL("https://github.com/jingweno/gh")
	assert.Equal(t, nil, err)
	assert.Equal(t, "jingweno", url.Owner)
	assert.Equal(t, "gh", url.Name)
	assert.Equal(t, "", url.ProjectPath())

	url, err =
		ParseURL("https://github.com/jingweno/gh/")
	assert.Equal(t, nil, err)
	assert.Equal(t, "jingweno", url.Owner)
	assert.Equal(t, "gh", url.Name)
	assert.Equal(t, "", url.ProjectPath())
}
