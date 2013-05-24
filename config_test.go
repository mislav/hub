package main

import (
	"github.com/bmizerany/assert"
	"os"
	"path/filepath"
	"testing"
)

func _TestSaveConfig(t *testing.T) {
	config := Config{"jingweno", "123"}
	file := "./test_support/test"
	err := SaveConfig(file, config)

	assert.Equal(t, nil, err)

	newConfig, _ := LoadConfig(file)
	assert.Equal(t, "jingweno", newConfig.User)
	assert.Equal(t, "123", newConfig.Token)

	os.RemoveAll(filepath.Dir(file))
}
