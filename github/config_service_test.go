package github

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestConfigService_Load(t *testing.T) {
	testConfig := fixtures.SetupTestConfigs()
	defer testConfig.TearDown()

	cc := &Config{}
	err := newConfigService().Load(testConfig.Path, cc)
	assert.Equal(t, nil, err)

	assert.Equal(t, 1, len(cc.Hosts))
	host := cc.Hosts[0]
	assert.Equal(t, "github.com", host.Host)
	assert.Equal(t, "jingweno", host.User)
	assert.Equal(t, "123", host.AccessToken)
	assert.Equal(t, "http", host.Protocol)
}

func TestConfigService_Save(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())

	host := Host{
		Host:        "github.com",
		User:        "jingweno",
		AccessToken: "123",
		Protocol:    "https",
	}
	c := Config{Hosts: []Host{host}}

	err := newConfigService().Save(file.Name(), &c)
	assert.Equal(t, nil, err)

	b, _ := ioutil.ReadFile(file.Name())
	content := `[[hosts]]
  host = "github.com"
  user = "jingweno"
  access_token = "123"
  protocol = "https"`
	assert.Equal(t, content, strings.TrimSpace(string(b)))
}
