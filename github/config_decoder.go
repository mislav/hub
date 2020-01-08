package github

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
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

	yc := yaml.MapSlice{}
	err = yaml.Unmarshal(d, &yc)

	if err != nil {
		return err
	}

	for _, hostEntry := range yc {
		v, ok := hostEntry.Value.([]interface{})
		if !ok {
			return fmt.Errorf("value of host entry is must be array but got %#v", hostEntry.Value)
		}
		if len(v) < 1 {
			continue
		}
		hostName, ok := hostEntry.Key.(string)
		if !ok {
			return fmt.Errorf("host name is must be string but got %#v", hostEntry.Key)
		}
		host := &Host{Host: hostName}
		for _, prop := range v[0].(yaml.MapSlice) {
			propName, ok := prop.Key.(string)
			if !ok {
				return fmt.Errorf("property name is must be string but got %#v", prop.Key)
			}
			switch propName {
			case "user":
				host.User, ok = prop.Value.(string)
			case "oauth_token":
				host.AccessToken, ok = prop.Value.(string)
			case "protocol":
				host.Protocol, ok = prop.Value.(string)
			case "unix_socket":
				host.UnixSocket, ok = prop.Value.(string)
			}
			if !ok {
				return fmt.Errorf("%s is must be string but got %#v", propName, prop.Value)
			}
		}
		c.Hosts = append(c.Hosts, host)
	}

	return nil
}
