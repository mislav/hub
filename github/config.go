package github

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh/terminal"
)

var defaultConfigsFile string

func init() {
	homeDir, err := homedir.Dir()
	utils.Check(err)

	defaultConfigsFile = filepath.Join(homeDir, ".config", "hub")
}

type yamlHost struct {
	User       string `yaml:"user"`
	OAuthToken string `yaml:"oauth_token"`
	Protocol   string `yaml:"protocol"`
}

type yamlConfig map[string][]yamlHost

type Host struct {
	Host        string `toml:"host"`
	User        string `toml:"user"`
	AccessToken string `toml:"access_token"`
	Protocol    string `toml:"protocol"`
}

type Config struct {
	Hosts []*Host `toml:"hosts"`
}

func (c *Config) PromptForHost(host string) (h *Host, err error) {
	token := c.DetectToken()
	tokenFromEnv := token != ""

	h = c.Find(host)
	if h != nil {
		if h.User == "" {
			utils.Check(CheckWriteable(configsFile()))
			// User is missing from the config: this is a broken config probably
			// because it was created with an old (broken) version of hub. Let's fix
			// it now. See issue #1007 for details.
			user := c.PromptForUser(host)
			if user == "" {
				utils.Check(fmt.Errorf("missing user"))
			}
			h.User = user
			err := newConfigService().Save(configsFile(), c)
			utils.Check(err)
		}
		if tokenFromEnv {
			h.AccessToken = token
		} else {
			return
		}
	} else {
		h = &Host{
			Host:        host,
			AccessToken: token,
			Protocol:    "https",
		}
		c.Hosts = append(c.Hosts, h)
	}

	client := NewClientWithHost(h)

	if !tokenFromEnv {
		utils.Check(CheckWriteable(configsFile()))
		err = c.authorizeClient(client, host)
		if err != nil {
			return
		}
	}

	currentUser, err := client.CurrentUser()
	if err != nil {
		return
	}
	h.User = currentUser.Login

	if !tokenFromEnv {
		err = newConfigService().Save(configsFile(), c)
	}

	return
}

func (c *Config) authorizeClient(client *Client, host string) (err error) {
	user := c.PromptForUser(host)
	pass := c.PromptForPassword(host, user)

	var code, token string
	for {
		token, err = client.FindOrCreateToken(user, pass, code)
		if err == nil {
			break
		}

		if ae, ok := err.(*AuthError); ok && ae.IsRequired2FACodeError() {
			if code != "" {
				ui.Errorln("warning: invalid two-factor code")
			}
			code = c.PromptForOTP()
		} else {
			break
		}
	}

	if err == nil {
		client.Host.AccessToken = token
	}

	return
}

func (c *Config) DetectToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

func (c *Config) PromptForUser(host string) (user string) {
	user = os.Getenv("GITHUB_USER")
	if user != "" {
		return
	}

	ui.Printf("%s username: ", host)
	user = c.scanLine()

	return
}

func (c *Config) PromptForPassword(host, user string) (pass string) {
	pass = os.Getenv("GITHUB_PASSWORD")
	if pass != "" {
		return
	}

	ui.Printf("%s password for %s (never stored): ", host, user)
	if ui.IsTerminal(os.Stdin) {
		if password, err := getPassword(); err == nil {
			pass = password
		}
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

func getPassword() (string, error) {
	stdin := int(syscall.Stdin)
	initialTermState, err := terminal.GetState(stdin)
	if err != nil {
		return "", err
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		terminal.Restore(stdin, initialTermState)
		switch sig := s.(type) {
		case syscall.Signal:
			if int(sig) == 2 {
				fmt.Println("^C")
			}
		}
		os.Exit(1)
	}()

	passBytes, err := terminal.ReadPassword(stdin)
	if err != nil {
		return "", err
	}

	signal.Stop(c)
	fmt.Print("\n")
	return string(passBytes), nil
}

func (c *Config) Find(host string) *Host {
	for _, h := range c.Hosts {
		if h.Host == host {
			return h
		}
	}

	return nil
}

func (c *Config) selectHost() *Host {
	options := len(c.Hosts)

	if options == 1 {
		return c.Hosts[0]
	}

	prompt := "Select host:\n"
	for idx, host := range c.Hosts {
		prompt += fmt.Sprintf(" %d. %s\n", idx+1, host.Host)
	}
	prompt += fmt.Sprint("> ")

	ui.Printf(prompt)
	index := c.scanLine()
	i, err := strconv.Atoi(index)
	if err != nil || i < 1 || i > options {
		utils.Check(fmt.Errorf("Error: must enter a number [1-%d]", options))
	}

	return c.Hosts[i-1]
}

func configsFile() string {
	configsFile := os.Getenv("HUB_CONFIG")
	if configsFile == "" {
		configsFile = defaultConfigsFile
	}

	return configsFile
}

var currentConfig *Config
var configLoadedFrom = ""

func CurrentConfig() *Config {
	filename := configsFile()
	if configLoadedFrom != filename {
		currentConfig = &Config{}
		newConfigService().Load(filename, currentConfig)
		configLoadedFrom = filename
	}

	return currentConfig
}

func (c *Config) DefaultHost() (host *Host, err error) {
	if GitHubHostEnv != "" {
		host, err = c.PromptForHost(GitHubHostEnv)
	} else if len(c.Hosts) > 0 {
		host = c.selectHost()
		// HACK: forces host to inherit GITHUB_TOKEN if applicable
		host, err = c.PromptForHost(host.Host)
	} else {
		host, err = c.PromptForHost(DefaultGitHubHost())
	}

	return
}

// CheckWriteable checks if config file is writeable. This should
// be called before asking for credentials and only if current
// operation needs to update the file. See issue #1314 for details.
func CheckWriteable(filename string) error {
	// Check if file exists already. if it doesn't, we will delete it after
	// checking for writeabilty
	fileExistsAlready := false

	if _, err := os.Stat(filename); err == nil {
		fileExistsAlready = true
	}

	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	w.Close()

	if !fileExistsAlready {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// Public for testing purpose
func CreateTestConfigs(user, token string) *Config {
	f, _ := ioutil.TempFile("", "test-config")
	defaultConfigsFile = f.Name()

	host := &Host{
		User:        "jingweno",
		AccessToken: "123",
		Host:        GitHubHost,
	}

	c := &Config{Hosts: []*Host{host}}
	err := newConfigService().Save(f.Name(), c)
	if err != nil {
		panic(err)
	}

	return c
}
