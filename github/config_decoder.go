package github

import (
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v1"
)

type configDecoder interface {
	Decode(r io.Reader, c *Config) error
}

type tomlConfigDecoder struct {
}

func (t *tomlConfigDecoder) Decode(r io.Reader, c *Config) error {
	_, err := toml.DecodeReader(r, c)
	return err
}

type yamlConfigDecoder struct {
}

func (y *yamlConfigDecoder) Decode(r io.Reader, c *Config) error {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	yc := make(yamlConfig)
	err = yaml.Unmarshal(d, &yc)

	if err != nil {
		return err
	}

	for h, v := range yc {
		host := Host{
			Host:        h,
			User:        v.User,
			AccessToken: v.OAuthToken,
			Protocol:    v.Protocol,
		}
		c.Hosts = append(c.Hosts, host)
	}

	return nil
}
