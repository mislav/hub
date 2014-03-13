package github

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

func TestConfigs_loadFrom(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())

	content := `[[credentials]]
  host = "https://github.com"
  user = "jingweno"
  access_token = "123"`
	ioutil.WriteFile(file.Name(), []byte(content), os.ModePerm)

	cc := &Configs{}
	err := loadFrom(file.Name(), cc)
	assert.Equal(t, nil, err)

	assert.Equal(t, 1, len(cc.Credentials))
	cred := cc.allCredentials()[0]
	assert.Equal(t, "https://github.com", cred.Host)
	assert.Equal(t, "jingweno", cred.User)
	assert.Equal(t, "123", cred.AccessToken)
}

func TestConfigs_saveTo(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())

	cred := Credential{Host: "https://github.com", User: "jingweno", AccessToken: "123"}
	c := Configs{Credentials: []Credential{cred}}

	err := saveTo(file.Name(), &c)
	assert.Equal(t, nil, err)

	b, _ := ioutil.ReadFile(file.Name())
	content := `[[credentials]]
  host = "https://github.com"
  user = "jingweno"
  access_token = "123"`
	assert.Equal(t, content, string(b))
}
