package config

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

var DefaultFile string

func init() {
	DefaultFile = filepath.Join(os.Getenv("HOME"), ".config", "gh")
}

func Load() (*Config, error) {
	return loadFrom(DefaultFile)
}

func loadFrom(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	dec := json.NewDecoder(reader)

	var c Config
	err = dec.Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func Save(config Config) error {
	return saveTo(DefaultFile, config)
}

func saveTo(filename string, config Config) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.Encode(config)

	return nil
}
