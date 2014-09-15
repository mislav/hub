package github

import (
	"io"

	"github.com/BurntSushi/toml"
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
