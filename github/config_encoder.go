package github

import (
	"io"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v1"
)

type configEncoder interface {
	Encode(w io.Writer, v interface{}) error
}

type tomlConfigEncoder struct {
}

func (t *tomlConfigEncoder) Encode(w io.Writer, v interface{}) error {
	enc := toml.NewEncoder(w)
	return enc.Encode(v)
}

type yamlConfigEncoder struct {
}

func (y *yamlConfigEncoder) Encode(w io.Writer, v interface{}) error {
	d, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	n, err := w.Write(d)
	if err == nil && n < len(d) {
		err = io.ErrShortWrite
	}

	return err
}
