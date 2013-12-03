package mediatype

import (
	"fmt"
	"io"
)

var decoders = make(map[string]DecoderFunc)

// DecoderFunc is a function that creates a Decoder from an io.Reader.
type DecoderFunc func(r io.Reader) Decoder

// A Decoder will decode the given value to the Decoder's io.Reader.
type Decoder interface {
	Decode(v interface{}) error
}

/*
AddDecoder installs a decoder for a given format.

	AddDecoder("json", func(r io.Reader) Encoder { return json.NewDecoder(r) })
	mt, err := Parse("application/json")
	decoder, err := mt.Decoder(someReader)
*/
func AddDecoder(format string, decfunc DecoderFunc) {
	decoders[format] = decfunc
}

// Decoder finds a decoder based on this MediaType's Format field.  An error is
// returned if a decoder cannot be found.
func (m *MediaType) Decoder(body io.Reader) (Decoder, error) {
	if decfunc, ok := decoders[m.Format]; ok {
		return decfunc(body), nil
	}
	return nil, fmt.Errorf("No decoder found for format %s (%s)", m.Format, m.String())
}

// Encode uses this MediaType's Decoder to decode the io.Reader into the given
// value.
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
