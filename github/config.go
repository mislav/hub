package github

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/github/hub/utils"
	"github.com/howeyc/gopass"
)

var (
	defaultConfigsFile = filepath.Join(os.Getenv("HOME"), ".config", "hub")
)

type yamlHost struct {
	User       string `yaml:"user"`
	OAuthToken string `yaml:"oauth_token"`
	Protocol   string `yaml:"protocol"`
}

type yamlConfig map[string]yamlHost

type Host struct {
	Host        string `toml:"host"`
	User        string `toml:"user"`
	AccessToken string `toml:"access_token"`
	Protocol    string `toml:"protocol"`
}

type Config struct {
	Hosts []Host `toml:"hosts"`
}

func (c *Config) PromptForHost(host string) (h *Host, err error) {
	h = c.Find(host)
	if h != nil {
		return
	}

	user := c.PromptForUser()
	pass := c.PromptForPassword(host, user)

	client := NewClient(host)
	token, e := client.FindOrCreateToken(user, pass, "")
	if e != nil {
		if ae, ok := e.(*AuthError); ok && ae.Is2FAError() {
			code := c.PromptForOTP()
			token, err = client.FindOrCreateToken(user, pass, code)
		} else {
			err = e
		}
	}

	if err != nil {
		return
	}

	client.Host.AccessToken = token
	currentUser, err := client.CurrentUser()
	if err != nil {
		return
	}

	h = &Host{
		Host:        host,
		User:        currentUser.Login,
		AccessToken: token,
		Protocol:    "https",
	}
	c.Hosts = append(c.Hosts, *h)
	err = newConfigService().Save(configsFile(), c)

	return
}

func (c *Config) PromptForUser() (user string) {
	user = os.Getenv("GITHUB_USER")
	if user != "" {
		return
	}

	fmt.Printf("%s username: ", GitHubHost)
	user = c.scanLine()

	return
}

func (c *Config) PromptForPassword(host, user string) (pass string) {
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

func (c *Config) PromptForOTP() string {
	fmt.Print("two-factor authentication code: ")
	return c.scanLine()
}

func (c *Config) scanLine() string {
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	utils.Check(scanner.Err())

	return line
}

func (c *Config) Find(host string) *Host {
	for _, h := range c.Hosts {
		if h.Host == host {
			return &h
		}
	}

	return nil
}

func (c *Config) selectHost() *Host {
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

func configsFile() string {
	configsFile := os.Getenv("GH_CONFIG")
	if configsFile == "" {
		configsFile = defaultConfigsFile
	}

	return configsFile
}

func CurrentConfig() *Config {
	c := &Config{}

	err := newConfigService().Load(configsFile(), c)
	if err != nil {
		panic(err)
	}

	return c
}

func (c *Config) DefaultHost() (host *Host, err error) {
	if GitHubHostEnv != "" {
		host, err = c.PromptForHost(GitHubHostEnv)
	} else if len(c.Hosts) > 0 {
		host = c.selectHost()
	} else {
		host, err = c.PromptForHost(DefaultGitHubHost())
	}

	return
}

// Public for testing purpose
func CreateTestConfigs(user, token string) *Config {
	f, _ := ioutil.TempFile("", "test-config")
	defaultConfigsFile = f.Name()

	host := Host{
		User:        "jingweno",
		AccessToken: "123",
		Host:        GitHubHost,
	}

	c := &Config{Hosts: []Host{host}}
	err := newConfigService().Save(f.Name(), c)
	if err != nil {
		panic(err)
	}

	return c
}
