package github

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/github/hub/v2/fixtures"
	"github.com/github/hub/v2/internal/assert"
)

func TestConfigService_TomlLoad(t *testing.T) {
	testConfig := fixtures.SetupTomlTestConfig()
	defer testConfig.TearDown()

	cc := &Config{}
	cs := &configService{
		Encoder: &tomlConfigEncoder{},
		Decoder: &tomlConfigDecoder{},
	}
	err := cs.Load(testConfig.Path, cc)
	assert.Equal(t, nil, err)

	assert.Equal(t, 1, len(cc.Hosts))
	host := cc.Hosts[0]
	assert.Equal(t, "github.com", host.Host)
	assert.Equal(t, "jingweno", host.User)
	assert.Equal(t, "123", host.AccessToken)
	assert.Equal(t, "http", host.Protocol)
}

func TestConfigService_TomlLoad_UnixSocket(t *testing.T) {
	testConfigUnixSocket := fixtures.SetupTomlTestConfigWithUnixSocket()
	defer testConfigUnixSocket.TearDown()

	cc := &Config{}
	cs := &configService{
		Encoder: &tomlConfigEncoder{},
		Decoder: &tomlConfigDecoder{},
	}

	err := cs.Load(testConfigUnixSocket.Path, cc)
	assert.Equal(t, nil, err)

	assert.Equal(t, 1, len(cc.Hosts))
	host := cc.Hosts[0]
	assert.Equal(t, "github.com", host.Host)
	assert.Equal(t, "jingweno", host.User)
	assert.Equal(t, "123", host.AccessToken)
	assert.Equal(t, "http", host.Protocol)
	assert.Equal(t, "/tmp/go.sock", host.UnixSocket)
}

func TestConfigService_YamlLoad(t *testing.T) {
	testConfig := fixtures.SetupTestConfigs()
	defer testConfig.TearDown()

	cc := &Config{}
	cs := &configService{
		Encoder: &yamlConfigEncoder{},
		Decoder: &yamlConfigDecoder{},
	}
	err := cs.Load(testConfig.Path, cc)
	assert.Equal(t, nil, err)

	assert.Equal(t, 1, len(cc.Hosts))
	host := cc.Hosts[0]
	assert.Equal(t, "github.com", host.Host)
	assert.Equal(t, "jingweno", host.User)
	assert.Equal(t, "123", host.AccessToken)
	assert.Equal(t, "http", host.Protocol)
}

func TestConfigService_YamlLoad_Unix_Socket(t *testing.T) {
	testConfigUnixSocket := fixtures.SetupTestConfigsWithUnixSocket()
	defer testConfigUnixSocket.TearDown()

	cc := &Config{}
	cs := &configService{
		Encoder: &yamlConfigEncoder{},
		Decoder: &yamlConfigDecoder{},
	}

	err := cs.Load(testConfigUnixSocket.Path, cc)
	assert.Equal(t, nil, err)

	assert.Equal(t, 1, len(cc.Hosts))
	host := cc.Hosts[0]
	assert.Equal(t, "github.com", host.Host)
	assert.Equal(t, "jingweno", host.User)
	assert.Equal(t, "123", host.AccessToken)
	assert.Equal(t, "http", host.Protocol)
	assert.Equal(t, "/tmp/go.sock", host.UnixSocket)
}

func TestConfigService_YamlLoad_Invalid_HostName(t *testing.T) {
	testConfigInvalidHostName := fixtures.SetupTestConfigsInvalidHostName()
	defer testConfigInvalidHostName.TearDown()

	cc := &Config{}
	cs := &configService{
		Encoder: &yamlConfigEncoder{},
		Decoder: &yamlConfigDecoder{},
	}

	err := cs.Load(testConfigInvalidHostName.Path, cc)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "host name is must be string but got 123", err.Error())
}

func TestConfigService_YamlLoad_Invalid_HostEntry(t *testing.T) {
	testConfigInvalidHostEntry := fixtures.SetupTestConfigsInvalidHostEntry()
	defer testConfigInvalidHostEntry.TearDown()

	cc := &Config{}
	cs := &configService{
		Encoder: &yamlConfigEncoder{},
		Decoder: &yamlConfigDecoder{},
	}

	err := cs.Load(testConfigInvalidHostEntry.Path, cc)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "value of host entry is must be array but got \"hello\"", err.Error())
}

func TestConfigService_YamlLoad_Invalid_PropertyValue(t *testing.T) {
	testConfigInvalidPropertyValue := fixtures.SetupTestConfigsInvalidPropertyValue()
	defer testConfigInvalidPropertyValue.TearDown()

	cc := &Config{}
	cs := &configService{
		Encoder: &yamlConfigEncoder{},
		Decoder: &yamlConfigDecoder{},
	}

	err := cs.Load(testConfigInvalidPropertyValue.Path, cc)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "user is must be string but got <nil>", err.Error())
}

func TestConfigService_TomlSave(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())

	host := &Host{
		Host:        "github.com",
		User:        "jingweno",
		AccessToken: "123",
		Protocol:    "https",
	}
	c := &Config{Hosts: []*Host{host}}

	cs := &configService{
		Encoder: &tomlConfigEncoder{},
		Decoder: &tomlConfigDecoder{},
	}
	err := cs.Save(file.Name(), c)
	assert.Equal(t, nil, err)

	b, _ := ioutil.ReadFile(file.Name())
	content := `[[hosts]]
  host = "github.com"
  user = "jingweno"
  access_token = "123"
  protocol = "https"`
	assert.Equal(t, content, strings.TrimSpace(string(b)))
}

func TestConfigService_TomlSave_UnixSocket(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())

	host := &Host{
		Host:        "github.com",
		User:        "jingweno",
		AccessToken: "123",
		Protocol:    "https",
		UnixSocket:  "/tmp/go.sock",
	}
	c := &Config{Hosts: []*Host{host}}

	cs := &configService{
		Encoder: &tomlConfigEncoder{},
		Decoder: &tomlConfigDecoder{},
	}
	err := cs.Save(file.Name(), c)
	assert.Equal(t, nil, err)

	b, _ := ioutil.ReadFile(file.Name())
	content := `[[hosts]]
  host = "github.com"
  user = "jingweno"
  access_token = "123"
  protocol = "https"
  unix_socket = "/tmp/go.sock"`
	assert.Equal(t, content, strings.TrimSpace(string(b)))
}

func TestConfigService_YamlSave(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())

	host := &Host{
		Host:        "github.com",
		User:        "jingweno",
		AccessToken: "123",
		Protocol:    "https",
	}
	c := &Config{Hosts: []*Host{host}}

	cs := &configService{
		Encoder: &yamlConfigEncoder{},
		Decoder: &yamlConfigDecoder{},
	}
	err := cs.Save(file.Name(), c)
	assert.Equal(t, nil, err)

	b, _ := ioutil.ReadFile(file.Name())
	content := `github.com:
- user: jingweno
  oauth_token: "123"
  protocol: https`
	assert.Equal(t, content, strings.TrimSpace(string(b)))
}

func TestConfigService_YamlSave_UnixSocket(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-gh-config-")
	defer os.RemoveAll(file.Name())

	host := &Host{
		Host:        "github.com",
		User:        "jingweno",
		AccessToken: "123",
		Protocol:    "https",
		UnixSocket:  "/tmp/go.sock",
	}
	c := &Config{Hosts: []*Host{host}}

	cs := &configService{
		Encoder: &yamlConfigEncoder{},
		Decoder: &yamlConfigDecoder{},
	}
	err := cs.Save(file.Name(), c)
	assert.Equal(t, nil, err)

	b, _ := ioutil.ReadFile(file.Name())
	content := `github.com:
- user: jingweno
  oauth_token: "123"
  protocol: https
  unix_socket: /tmp/go.sock`
	assert.Equal(t, content, strings.TrimSpace(string(b)))
}
