package github

import (
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v1"
)

type configDecoder interface {
	Decode(r io.Reader, v interface{}) error
}

type tomlConfigDecoder struct {
}

func (t *tomlConfigDecoder) Decode(r io.Reader, v interface{}) error {
	_, err := toml.DecodeReader(r, v)
	return err
}

type yamlConfigDecoder struct {
}

func (y *yamlConfigDecoder) Decode(r io.Reader, v interface{}) error {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(d, v)
}
