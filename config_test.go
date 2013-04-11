package main

import (
	"github.com/bmizerany/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	home := os.Getenv("HOME")
	config := loadConfig(home + "/.config/gh")

	assert.Equal(t, "jingweno", config.User)
}
