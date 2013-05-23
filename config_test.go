package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, _ := LoadConfig("./test_support/gh")

	assert.Equal(t, "jingweno", config.User)
	assert.Equal(t, "02a66f3bdde949182bc0d629f1abef0d501e6a53", config.Token)
}
