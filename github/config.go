// TODO: remove it in favour of Configs

package github

import (
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jingweno/gh/utils"
	"io/ioutil"
	"regexp"
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

func (c *Config) FetchTwoFactorCode() string {
	var code string
	fmt.Print("two-factor authentication code: ")
	fmt.Scanln(&code)

	return code
}

func (c *Config) FetchCredentials() {
	var changed bool
	if c.User == "" {
		c.FetchUser()
		changed = true
	}

	if c.Token == "" {
		password := c.FetchPassword()
		token, err := findOrCreateToken(c.User, password, "")
		// TODO: return an two factor auth failure error
		if err != nil {
			re := regexp.MustCompile("two-factor authentication OTP code")
			if re.MatchString(fmt.Sprintf("%s", err)) {
				code := c.FetchTwoFactorCode()
				token, err = findOrCreateToken(c.User, password, code)
			}
		}

		utils.Check(err)

		c.Token = token
		changed = true
	}

	if changed {
		err := saveTo(configsFile(), c)
		utils.Check(err)
	}
}

func CurrentConfig() *Config {
	var config Config
	err := loadFrom(configsFile(), &config)
	if err != nil {
		config = Config{}
	}
	config.FetchCredentials()

	return &config
}

// TODO: remove it
func CreateTestConfig(user, token string) *Config {
	f, _ := ioutil.TempFile("", "test-config")
	defaultConfigsFile = f.Name()

	config := Config{User: "jingweno", Token: "123"}
	saveTo(f.Name(), &config)

	return &config
}
