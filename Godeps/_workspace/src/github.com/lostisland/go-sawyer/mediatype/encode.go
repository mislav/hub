package mediatype

import (
	"bytes"
	"fmt"
	"io"
)

var encoders = make(map[string]EncoderFunc)

type EncoderFunc func(w io.Writer) Encoder

type Encoder interface {
	Encode(v interface{}) error
}

func AddEncoder(format string, encfunc EncoderFunc) {
	encoders[format] = encfunc
}

func (m *MediaType) Encoder(w io.Writer) (Encoder, error) {
	if encfunc, ok := encoders[m.Format]; ok {
		return encfunc(w), nil
	}
	return nil, fmt.Errorf("No encoder found for format %s (%s)", m.Format, m.String())
}

func (m *MediaType) Encode(v interface{}) (*bytes.Buffer, error) {
	if v == nil {
		return nil, fmt.Errorf("Nothing to encode")
	}

	buf := new(bytes.Buffer)
	enc, err := m.Encoder(buf)
	if err != nil {
		return buf, err
	}

	return buf, enc.Encode(v)
}
