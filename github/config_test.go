package github

import (
	"github.com/bmizerany/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveConfig(t *testing.T) {
	config := Config{"jingweno", "123"}
	file := "./test_support/test"
	defer os.RemoveAll(filepath.Dir(file))

	err := saveTo(file, &config)
	assert.Equal(t, nil, err)

	config, err = loadFrom(file)
	assert.Equal(t, nil, err)
	assert.Equal(t, "jingweno", config.User)
	assert.Equal(t, "123", config.Token)

	newConfig := Config{"foo", "456"}
	err = saveTo(file, &newConfig)
	assert.Equal(t, nil, err)

	config, err = loadFrom(file)
	assert.Equal(t, "foo", config.User)
	assert.Equal(t, "456", config.Token)
}
