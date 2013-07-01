package github

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jingweno/gh/utils"
	"os"
	"path/filepath"
)

type Config struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

func (c *Config) FetchUser() string {
	if c.User == "" {
		var user string
		msg := fmt.Sprintf("%s username: ", GitHubHost)
		fmt.Print(msg)
		fmt.Scanln(&user)
		c.User = user
	}

	return c.User
}

func (c *Config) FetchPassword() string {
	msg := fmt.Sprintf("%s password for %s (never stored): ", GitHubHost, c.User)
	fmt.Print(msg)

	pass := gopass.GetPasswd()
	if len(pass) == 0 {
		utils.Check(errors.New("Password cannot be empty"))
	}

	return string(pass)
}

func (c *Config) FetchCredentials() {
	var changed bool
	if c.User == "" {
		c.FetchUser()
		changed = true
	}

	if c.Token == "" {
		password := c.FetchPassword()
		token, err := findOrCreateToken(c.User, password)
		utils.Check(err)

		c.Token = token
		changed = true
	}

	if changed {
		err := SaveConfig(c)
		utils.Check(err)
	}
}

var (
	DefaultConfigFile = filepath.Join(os.Getenv("HOME"), ".config", "gh")
)

func CurrentConfig() *Config {
	config, err := loadConfig()
	if err != nil {
		config = Config{}
	}
	config.FetchCredentials()

	return &config
}

func loadConfig() (Config, error) {
	return loadFrom(DefaultConfigFile)
}

func loadFrom(filename string) (Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}

	return doLoadFrom(f)
}

func doLoadFrom(f *os.File) (Config, error) {
	defer f.Close()

	reader := bufio.NewReader(f)
	dec := json.NewDecoder(reader)

	var c Config
	err := dec.Decode(&c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func SaveConfig(config *Config) error {
	return saveTo(DefaultConfigFile, config)
}

func saveTo(filename string, config *Config) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	return doSaveTo(f, config)
}

func doSaveTo(f *os.File, config *Config) error {
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(config)
}
