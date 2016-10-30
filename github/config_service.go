package github

import (
	"os"
	"path/filepath"
)

func newConfigService() *configService {
	return &configService{
		Encoder: &yamlConfigEncoder{},
		Decoder: &yamlConfigDecoder{},
	}
}

type configService struct {
	Encoder configEncoder
	Decoder configDecoder
}

func (s *configService) Save(filename string, c *Config) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer w.Close()

	return s.Encoder.Encode(w, c)
}

func (s *configService) Load(filename string, c *Config) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	return s.Decoder.Decode(r, c)
}

// CheckWriteable checks if config file is writeable. This should
// be called before asking for credentials and only if current
// operation needs to update the file. See issue #1314 for details.
func (s *configService) CheckWriteable(filename string) error {
    err := os.MkdirAll(filepath.Dir(filename), 0771)
    if err != nil {
        return err
    }

    w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
    if err != nil {
        return err
    }
    w.Close()
    return nil
}
