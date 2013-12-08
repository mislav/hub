package github

import (
	"github.com/bmizerany/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveCredentials(t *testing.T) {
	c := Credentials{Host: "github.com", User: "jingweno", AccessToken: "123"}
	file := "./test_support/test"
	defer os.RemoveAll(filepath.Dir(file))

	err := saveTo(file, &c)
	assert.Equal(t, nil, err)

	var cc Credentials
	err = loadFrom(file, &cc)
	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com", cc.Host)
	assert.Equal(t, "jingweno", cc.User)
	assert.Equal(t, "123", cc.AccessToken)
}
