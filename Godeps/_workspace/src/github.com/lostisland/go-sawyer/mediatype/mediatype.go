package mediatype

import (
	"mime"
	"strings"
)

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

func (m *MediaType) String() string {
	return m.full
}

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
