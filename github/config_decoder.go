package github

import (
	"io"

	"github.com/BurntSushi/toml"
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
