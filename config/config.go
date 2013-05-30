package config

import (
	"bufio"
	"encoding/json"
	"errors"
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

func Load(user string) (*Config, error) {
	configs, err := loadFrom(DefaultFile)
	if err != nil {
		return nil, err
	}

	for _, c := range configs {
		if c.User == user {
			return c, nil
		}
	}

	return nil, errors.New("There's no matching config for user: " + user)
}

func loadFrom(filename string) ([]*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return doLoadFrom(f)
}

func doLoadFrom(f *os.File) ([]*Config, error) {
	defer f.Close()

	reader := bufio.NewReader(f)
	dec := json.NewDecoder(reader)

	var c []*Config
	err := dec.Decode(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Save(config *Config) error {
	return saveTo(DefaultFile, config)
}

func saveTo(filename string, config *Config) error {
	configs, _ := loadFrom(filename)

	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	var foundConfig *Config
	for _, c := range configs {
		if c.User == config.User {
			foundConfig = c
			break
		}
	}
	if foundConfig == nil {
		configs = append(configs, config)
	} else {
		foundConfig.Token = config.Token
	}

	return doSaveTo(f, configs)
}

func doSaveTo(f *os.File, configs []*Config) error {
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(configs)
}
