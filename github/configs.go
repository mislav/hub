package github

import (
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jingweno/gh/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var (
	defaultConfigsFile = filepath.Join(os.Getenv("HOME"), ".config", "gh")
)

type Credentials struct {
	Host        string `json:"host"`
	User        string `json:"user"`
	AccessToken string `json:"access_token"`
}

type Configs struct {
	Autoupdate  bool          `json:"autoupdate"`
	Credentials []Credentials `json:"credentials"`
}

func (c *Configs) PromptFor(host string) *Credentials {
	cc := c.find(host)
	if cc == nil {
		user := c.PromptForUser()
		pass := c.PromptForPassword(host, user)

		// Create Client with a stub Credentials
		client := Client{Credentials: &Credentials{Host: host}}
		token, err := client.FindOrCreateToken(user, pass, "")
		if err != nil {
			if ce, ok := err.(*ClientError); ok && ce.Is2FAError() {
				code := c.PromptForOTP()
				token, err = client.FindOrCreateToken(user, pass, code)
			}
		}
		utils.Check(err)

		cc = &Credentials{Host: host, User: user, AccessToken: token}
		c.Credentials = append(c.Credentials, *cc)
		err = saveTo(configsFile(), c)
		utils.Check(err)
	}

	return cc
}

func (c *Configs) PromptForUser() (user string) {
	user = os.Getenv("GITHUB_USER")
	if user != "" {
		return
	}

	fmt.Printf("%s username: ", GitHubHost)
	fmt.Scanln(&user)

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
		fmt.Scanln(&pass)
	}

	return
}

func (c *Configs) PromptForOTP() string {
	var code string
	fmt.Print("two-factor authentication code: ")
	fmt.Scanln(&code)

	return code
}

func (c *Configs) find(host string) *Credentials {
	for _, t := range c.Credentials {
		if t.Host == host {
			return &t
		}
	}

	return nil
}

func saveTo(filename string, v interface{}) error {
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
	return enc.Encode(v)
}

func loadFrom(filename string, c *Configs) error {
	return loadFromFile(filename, c)
}

// Function to load deprecated configuration.
// It's not intended to be used.
func loadFromDeprecated(filename string, c *[]Credentials) error {
	return loadFromFile(filename, c)
}

func loadFromFile(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	for {
		if err := dec.Decode(v); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
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
		// Try deprecated configuration
		var creds []Credentials
		err := loadFromDeprecated(configsFile(), &creds)
		if err != nil {
			creds = make([]Credentials, 0)
		}
		c.Credentials = creds
		saveTo(configFile, c)
	}

	return c
}

func (c *Configs) DefaultCredentials() (credentials *Credentials) {
	if GitHubHostEnv != "" {
		credentials = c.PromptFor(GitHubHostEnv)
	} else if len(c.Credentials) > 0 {
		credentials = c.selectCredentials()
	} else {
		credentials = c.PromptFor(DefaultHost())
	}

	return
}

func (c *Configs) selectCredentials() *Credentials {
	options := len(c.Credentials)

	if options == 1 {
		return &c.Credentials[0]
	}

	prompt := "Select host:\n"
	for idx, creds := range c.Credentials {
		prompt += fmt.Sprintf(" %d. %s\n", idx+1, creds.Host)
	}
	prompt += fmt.Sprint("> ")

	fmt.Printf(prompt)
	var index string
	fmt.Scanln(&index)

	i, err := strconv.Atoi(index)
	if err != nil || i < 1 || i > options {
		utils.Check(fmt.Errorf("Error: must enter a number [1-%d]", options))
	}

	return &c.Credentials[i-1]
}

// Public for testing purpose
func CreateTestConfigs(user, token string) *Configs {
	f, _ := ioutil.TempFile("", "test-config")
	defaultConfigsFile = f.Name()

	creds := []Credentials{
		{User: "jingweno", AccessToken: "123", Host: GitHubHost},
	}

	c := &Configs{Credentials: creds}
	saveTo(f.Name(), c)

	return c
}
