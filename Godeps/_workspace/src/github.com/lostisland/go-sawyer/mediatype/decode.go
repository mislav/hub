package mediatype

import (
	"fmt"
	"io"
)

var decoders = make(map[string]DecoderFunc)

type DecoderFunc func(r io.Reader) Decoder

type Decoder interface {
	Decode(v interface{}) error
}

func AddDecoder(format string, decfunc DecoderFunc) {
	decoders[format] = decfunc
}

func (m *MediaType) Decoder(body io.Reader) (Decoder, error) {
	if decfunc, ok := decoders[m.Format]; ok {
		return decfunc(body), nil
	}
	return nil, fmt.Errorf("No decoder found for format %s (%s)", m.Format, m.String())
}

func (m *MediaType) Decode(v interface{}, body io.Reader) error {
	if v == nil {
		return nil
	}

	dec, err := m.Decoder(body)
	if err != nil {
		return err
	}

	return dec.Decode(v)
}
