package config

import (
	"github.com/bmizerany/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestSave(t *testing.T) {
	config := Config{"jingweno", "123"}
	file := "./test_support/test"
	defer os.RemoveAll(filepath.Dir(file))

	err := saveTo(file, &config)
	assert.Equal(t, nil, err)

	configs, err := loadFrom(file)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(configs))
	assert.Equal(t, "jingweno", configs[0].User)
	assert.Equal(t, "123", configs[0].Token)

	newConfig := Config{"foo", "456"}
	err = saveTo(file, &newConfig)
	assert.Equal(t, nil, err)

	configs, err = loadFrom(file)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(configs))
	assert.Equal(t, "jingweno", configs[0].User)
	assert.Equal(t, "123", configs[0].Token)
	assert.Equal(t, "foo", configs[1].User)
	assert.Equal(t, "456", configs[1].Token)
}
