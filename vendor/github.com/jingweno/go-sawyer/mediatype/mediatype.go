// Package mediatype contains helpers for parsing media type strings.  Uses
// RFC4288 as a guide.
package mediatype

import (
	"mime"
	"strings"
)

/*
A MediaType is a parsed representation of a media type string.

  application/vnd.github.raw+json; version=3; charset=utf-8

This gets broken up into the various fields:

- Type: application/vnd.github.raw+json
- MainType: application
- SubType: vnd.github.raw
- Suffix: json
- Vendor: github
- Version: raw
- Format: json
- Params:
    version: 3
    charset: utf-8

There are a few special behaviors that prioritize custom media types for APIs:

If an API identifies with an "application/vnd" type, the Vendor and Version
fields are parsed from the remainder.  The Version's semantic meaning depends on
the application.

If it's not an "application/vnd" type, the Version field is taken from the
"version" parameter.

The Format is taken from the Suffix by default.  If not available, it is guessed
by looking for common strings anywhere in the media type.  For instance,
"application/json" will identify as the "json" Format.

The Format is used to get an Encoder and a Decoder.
*/
type MediaType struct {
	full     string
	Type     string
	MainType string
	SubType  string
	Suffix   string
	Vendor   string
	Version  string
	Format   string
	Params   map[string]string
}

// Parse builds a *MediaType from a given media type string.
func Parse(v string) (*MediaType, error) {
	mt, params, err := mime.ParseMediaType(v)
	if err != nil {
		return nil, err
	}

	return parse(&MediaType{
		full:   v,
		Type:   mt,
		Params: params,
	})
}

// String returns the full string representation of the MediaType.
func (m *MediaType) String() string {
	return m.full
}

// IsVendor determines if this MediaType is associated with commercially
// available products.
func (m *MediaType) IsVendor() bool {
	return len(m.Vendor) > 0
}

func parse(m *MediaType) (*MediaType, error) {
	pieces := strings.Split(m.Type, typeSplit)
	m.MainType = pieces[0]
	if len(pieces) > 1 {
		subpieces := strings.Split(pieces[1], suffixSplit)
		m.SubType = subpieces[0]
		if len(subpieces) > 1 {
			m.Suffix = subpieces[1]
		}
	}

	if strings.HasPrefix(m.SubType, vndPrefix) {
		if vnd := m.SubType[vndLen:]; len(vnd) > 0 {
			args := strings.SplitN(vnd, vndSplit, 2)
			m.Vendor = args[0]
			if len(args) > 1 {
				m.Version = args[1]
			}
		}
	}

	if len(m.Version) == 0 {
		if v, ok := m.Params[versionKey]; ok {
			m.Version = v
		}
	}

	if len(m.Suffix) > 0 {
		m.Format = m.Suffix
	} else {
		guessFormat(m)
	}

	return m, nil
}

func guessFormat(m *MediaType) {
	for _, fmt := range guessableTypes {
		if strings.Contains(m.Type, fmt) {
			m.Format = fmt
			return
		}
	}
}

const (
	typeSplit   = "/"
	suffixSplit = "+"
	versionKey  = "version"
	vndPrefix   = "vnd."
	vndLen      = 4
	vndSplit    = "."
)

var guessableTypes = []string{"json", "xml"}
