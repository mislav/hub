package main

import (
	"github.com/bmizerany/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, _ := LoadConfig("./test_support/gh")

	assert.Equal(t, "jingweno", config.User)
	assert.Equal(t, "02a66f3bdde949182bc0d629f1abef0d501e6a53", config.Token)
}

func TestSaveConfig(t *testing.T) {
	config := Config{"jingweno", "123"}
	file := "./test_support/test"
	err := SaveConfig(file, config)

	assert.Equal(t, nil, err)

	newConfig, _ := LoadConfig(file)
	assert.Equal(t, "jingweno", newConfig.User)
	assert.Equal(t, "123", newConfig.Token)

	os.Remove(file)
}
