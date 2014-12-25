package mediatype

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParsesJsonType(t *testing.T) {
	m := Get(t, "application/json")
	assert.Equal(t, "application/json", m.Type)
	assert.Equal(t, "application", m.MainType)
	assert.Equal(t, "json", m.SubType)
	assert.Equal(t, "", m.Suffix)
	assert.Equal(t, "", m.Vendor)
	assert.Equal(t, "json", m.Format)
	assert.Equal(t, false, m.IsVendor())
	assert.Equal(t, 0, len(m.Params))
}

func TestParsesEmptyType(t *testing.T) {
	m := Get(t, "*; q=.2")
	assert.Equal(t, "*", m.Type)
	assert.Equal(t, "*", m.MainType)
	assert.Equal(t, "", m.SubType)
	assert.Equal(t, "", m.Suffix)
	assert.Equal(t, "", m.Vendor)
	assert.Equal(t, "", m.Format)
	assert.Equal(t, false, m.IsVendor())
	assert.Equal(t, 1, len(m.Params))
	assert.Equal(t, ".2", m.Params["q"])
}

func TestSimpleTypeWithParams(t *testing.T) {
	m := Get(t, "text/plain; charset=utf-8")
	assert.Equal(t, "text/plain", m.Type)
	assert.Equal(t, "text", m.MainType)
	assert.Equal(t, "plain", m.SubType)
	assert.Equal(t, "", m.Suffix)
	assert.Equal(t, "", m.Vendor)
	assert.Equal(t, "", m.Format)
	assert.Equal(t, false, m.IsVendor())
	assert.Equal(t, 1, len(m.Params))
	assert.Equal(t, "utf-8", m.Params["charset"])
}

func TestVendorType(t *testing.T) {
	m := Get(t, "application/vnd.json+xml; charset=utf-8")
	assert.Equal(t, "application/vnd.json+xml", m.Type)
	assert.Equal(t, "application", m.MainType)
	assert.Equal(t, "vnd.json", m.SubType)
	assert.Equal(t, "xml", m.Suffix)
	assert.Equal(t, "json", m.Vendor)
	assert.Equal(t, "xml", m.Format)
	assert.Equal(t, true, m.IsVendor())
	assert.Equal(t, 1, len(m.Params))
	assert.Equal(t, "utf-8", m.Params["charset"])
}

func TestSubtypeVersion(t *testing.T) {
	m := Get(t, "application/vnd.abc.v1+xml; version=v2; charset=utf-8")
	assert.Equal(t, "application/vnd.abc.v1+xml", m.Type)
	assert.Equal(t, "application", m.MainType)
	assert.Equal(t, "vnd.abc.v1", m.SubType)
	assert.Equal(t, "xml", m.Suffix)
	assert.Equal(t, "abc", m.Vendor)
	assert.Equal(t, "v1", m.Version)
	assert.Equal(t, "xml", m.Format)
	assert.Equal(t, true, m.IsVendor())
	assert.Equal(t, 2, len(m.Params))
	assert.Equal(t, "utf-8", m.Params["charset"])
	assert.Equal(t, "v2", m.Params["version"])
}

func TestParamVersion(t *testing.T) {
	m := Get(t, "application/vnd.abc+xml; version=v2; charset=utf-8")
	assert.Equal(t, "application/vnd.abc+xml", m.Type)
	assert.Equal(t, "application", m.MainType)
	assert.Equal(t, "vnd.abc", m.SubType)
	assert.Equal(t, "xml", m.Suffix)
	assert.Equal(t, "abc", m.Vendor)
	assert.Equal(t, "v2", m.Version)
	assert.Equal(t, "xml", m.Format)
	assert.Equal(t, true, m.IsVendor())
	assert.Equal(t, 2, len(m.Params))
	assert.Equal(t, "utf-8", m.Params["charset"])
	assert.Equal(t, "v2", m.Params["version"])
}
func Get(t *testing.T, v string) *MediaType {
	m, err := Parse(v)
	if err != nil {
		t.Fatalf("Errors parsing media type %s:\n%s", v, err.Error())
	}
	assert.Equal(t, v, m.String())
	return m
}

// used for encoding/decoding tests
type Person struct {
	Name string
}
