package github

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestSaveCredentials(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())

	ccreds := Credentials{Host: "github.com", User: "jingweno", AccessToken: "123"}
	c := Configs{Credentials: []Credentials{ccreds}}

	err := saveTo(file.Name(), &c)
	assert.Equal(t, nil, err)

	cc := &Configs{}
	err = loadFrom(file.Name(), cc)
	assert.Equal(t, nil, err)

	creds := cc.Credentials[0]
	assert.Equal(t, "github.com", creds.Host)
	assert.Equal(t, "jingweno", creds.User)
	assert.Equal(t, "123", creds.AccessToken)
}

func TestReadAndSaveDeprecatedConfiguration(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())
	defaultConfigsFile = file.Name()

	file.WriteString(`[{"host":"github.com","user":"jingweno","access_token":"123"}]`)
	file.Close()

	CurrentConfigs()

	expectedConfig := `{"credentials":[{"host":"github.com","user":"jingweno","access_token":"123"}]}
`

	f, _ := os.Open(file.Name())
	content, _ := ioutil.ReadAll(f)
	assert.Equal(t, expectedConfig, string(content))
}
