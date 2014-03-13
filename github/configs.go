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

type Credential struct {
	Host        string `toml:"host"`
	User        string `toml:"user"`
	AccessToken string `toml:"access_token"`
}

type Configs struct {
	Credentials []Credential `toml:"credentials"`
}

func (c *Configs) allCredentials() (creds []Credential) {
	for _, cred := range c.Credentials {
		creds = append(creds, cred)
	}

	return
}

func (c *Configs) PromptFor(host string) *Credential {
	cd := c.find(host)
	if cd == nil {
		user := c.PromptForUser()
		pass := c.PromptForPassword(host, user)

		// Create Client with a stub Credential
		client := Client{Credential: &Credential{Host: host}}
		token, err := client.FindOrCreateToken(user, pass, "")
		if err != nil {
			if ce, ok := err.(*ClientError); ok && ce.Is2FAError() {
				code := c.PromptForOTP()
				token, err = client.FindOrCreateToken(user, pass, code)
			}
		}
		utils.Check(err)

		client.Credential.AccessToken = token
		currentUser, err := client.CurrentUser()
		utils.Check(err)

		cd = &Credential{Host: host, User: currentUser.Login, AccessToken: token}
		c.Credentials = append(c.Credentials, *cd)
		err = saveTo(configsFile(), c)
		utils.Check(err)
	}

	return cd
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
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			pass = scanner.Text()
		}
	}

	return
}

func (c *Configs) PromptForOTP() string {
	var code string
	fmt.Print("two-factor authentication code: ")
	fmt.Scanln(&code)

	return code
}

func (c *Configs) find(host string) *Credential {
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

func (c *Configs) DefaultCredential() (credential *Credential) {
	if GitHubHostEnv != "" {
		credential = c.PromptFor(GitHubHostEnv)
	} else if len(c.Credentials) > 0 {
		credential = c.selectCredential()
	} else {
		credential = c.PromptFor(DefaultHost())
	}

	return
}

func (c *Configs) selectCredential() *Credential {
	creds := c.allCredentials()
	options := len(creds)

	if options == 1 {
		return &creds[0]
	}

	prompt := "Select host:\n"
	for idx, cred := range creds {
		prompt += fmt.Sprintf(" %d. %s\n", idx+1, cred.Host)
	}
	prompt += fmt.Sprint("> ")

	fmt.Printf(prompt)
	var index string
	fmt.Scanln(&index)

	i, err := strconv.Atoi(index)
	if err != nil || i < 1 || i > options {
		utils.Check(fmt.Errorf("Error: must enter a number [1-%d]", options))
	}

	return &creds[i-1]
}

func (c *Configs) Save() error {
	return saveTo(configsFile(), c)
}

// Public for testing purpose
func CreateTestConfigs(user, token string) *Configs {
	f, _ := ioutil.TempFile("", "test-config")
	defaultConfigsFile = f.Name()

	cred := Credential{
		User:        "jingweno",
		AccessToken: "123",
		Host:        GitHubHost,
	}

	c := &Configs{Credentials: []Credential{cred}}
	saveTo(f.Name(), c)

	return c
}
