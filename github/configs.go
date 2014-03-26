package github

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/github/hub/utils"
	"github.com/howeyc/gopass"
)

var (
	defaultConfigsFile = filepath.Join(os.Getenv("HOME"), ".config", "hub")
)

type Host struct {
	Host        string `toml:"host"`
	User        string `toml:"user"`
	AccessToken string `toml:"access_token"`
}

type Configs struct {
	Hosts []Host `toml:"hosts"`
}

func (c *Configs) PromptFor(host string) *Host {
	h := c.find(host)
	if h == nil {
		user := c.PromptForUser()
		pass := c.PromptForPassword(host, user)

		// Create Client with a stub Host
		client := Client{Host: &Host{Host: host}}
		token, err := client.FindOrCreateToken(user, pass, "")
		if err != nil {
			if ce, ok := err.(*ClientError); ok && ce.Is2FAError() {
				code := c.PromptForOTP()
				token, err = client.FindOrCreateToken(user, pass, code)
			}
		}
		utils.Check(err)

		client.Host.AccessToken = token
		currentUser, err := client.CurrentUser()
		utils.Check(err)

		h = &Host{Host: host, User: currentUser.Login, AccessToken: token}
		c.Hosts = append(c.Hosts, *h)
		err = saveTo(configsFile(), c)
		utils.Check(err)
	}

	return h
}

func (c *Configs) PromptForUser() (user string) {
	user = os.Getenv("GITHUB_USER")
	if user != "" {
		return
	}

	fmt.Printf("%s username: ", GitHubHost)
	user = c.scanLine()

	return
}

func (c *Configs) PromptForPassword(host, user string) (pass string) {
	pass = os.Getenv("GITHUB_PASSWORD")
	if pass != "" {
		return
	}

	fmt.Printf("%s password for %s (never stored): ", host, user)
	if isTerminal(os.Stdout.Fd()) {
		pass = string(gopass.GetPasswd())
	} else {
		pass = c.scanLine()
	}

	return
}

func (c *Configs) PromptForOTP() string {
	fmt.Print("two-factor authentication code: ")
	return c.scanLine()
}

func (c *Configs) scanLine() string {
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	utils.Check(scanner.Err())

	return line
}

func (c *Configs) find(host string) *Host {
	for _, h := range c.Hosts {
		if h.Host == host {
			return &h
		}
	}

	return nil
}

func saveTo(filename string, v interface{}) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	return enc.Encode(v)
}

func loadFrom(filename string, c *Configs) (err error) {
	_, err = toml.DecodeFile(filename, c)
	return
}

func configsFile() string {
	configsFile := os.Getenv("GH_CONFIG")
	if configsFile == "" {
		configsFile = defaultConfigsFile
	}

	return configsFile
}

func CurrentConfigs() *Configs {
	c := &Configs{}

	configFile := configsFile()
	err := loadFrom(configFile, c)
	if err != nil {
		// load from YAML
	}

	return c
}

func (c *Configs) DefaultHost() (host *Host) {
	if GitHubHostEnv != "" {
		host = c.PromptFor(GitHubHostEnv)
	} else if len(c.Hosts) > 0 {
		host = c.selectHost()
	} else {
		host = c.PromptFor(DefaultGitHubHost())
	}

	return
}

func (c *Configs) selectHost() *Host {
	options := len(c.Hosts)

	if options == 1 {
		return &c.Hosts[0]
	}

	prompt := "Select host:\n"
	for idx, host := range c.Hosts {
		prompt += fmt.Sprintf(" %d. %s\n", idx+1, host.Host)
	}
	prompt += fmt.Sprint("> ")

	fmt.Printf(prompt)
	index := c.scanLine()
	i, err := strconv.Atoi(index)
	if err != nil || i < 1 || i > options {
		utils.Check(fmt.Errorf("Error: must enter a number [1-%d]", options))
	}

	return &c.Hosts[i-1]
}

func (c *Configs) Save() error {
	return saveTo(configsFile(), c)
}

// Public for testing purpose
func CreateTestConfigs(user, token string) *Configs {
	f, _ := ioutil.TempFile("", "test-config")
	defaultConfigsFile = f.Name()

	host := Host{
		User:        "jingweno",
		AccessToken: "123",
		Host:        GitHubHost,
	}

	c := &Configs{Hosts: []Host{host}}
	saveTo(f.Name(), c)

	return c
}
