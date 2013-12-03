package mediatype

import (
	"bytes"
	"fmt"
	"io"
)

var encoders = make(map[string]EncoderFunc)

// EncoderFunc is a function that creates an Encoder from an io.Writer.
type EncoderFunc func(w io.Writer) Encoder

// An Encoder will encode the given value to the Encoder's io.Writer.
type Encoder interface {
	Encode(v interface{}) error
}

/*
AddEncoder installs an encoder for a given format.

  AddEncoder("json", func(w io.Writer) Encoder { return json.NewEncoder(w) })
	mt, err := Parse("application/json")
	encoder, err := mt.Encoder(someWriter)
*/
func AddEncoder(format string, encfunc EncoderFunc) {
	encoders[format] = encfunc
}

// Encoder finds an encoder based on this MediaType's Format field.  An error is
// returned if an encoder cannot be found.
func (m *MediaType) Encoder(w io.Writer) (Encoder, error) {
	if encfunc, ok := encoders[m.Format]; ok {
		return encfunc(w), nil
	}
	return nil, fmt.Errorf("No encoder found for format %s (%s)", m.Format, m.String())
}

// Encode uses this MediaType's Encoder to encode the given value into a
// bytes.Buffer.
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
