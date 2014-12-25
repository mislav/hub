package mediatype

import (
	"bytes"
	"github.com/bmizerany/assert"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestAddDecoder(t *testing.T) {
	buf := bytes.NewBufferString("bob")
	mt, err := Parse("application/test+test")
	if err != nil {
		t.Fatalf("Error parsing media type: %s", err.Error())
	}

	person := &Person{}
	err = mt.Decode(person, buf)
	if err != nil {
		t.Fatalf("Error decoding: %s", err.Error())
	}
	assert.Equal(t, "bob", person.Name)
}

func TestRequiresDecoder(t *testing.T) {
	buf := bytes.NewBufferString("bob")
	mt, err := Parse("application/test+whatevs")
	if err != nil {
		t.Fatalf("Error parsing media type: %s", err.Error())
	}

	person := &Person{}
	err = mt.Decode(person, buf)
	if err == nil {
		t.Fatal("No decoding error")
	}

	if !strings.HasPrefix(err.Error(), "No decoder found for format whatevs") {
		t.Fatalf("Bad error: %s", err)
	}
}

func TestSkipsDecoderForNil(t *testing.T) {
	buf := bytes.NewBufferString("bob")
	mt, err := Parse("application/test+whatevs")
	if err != nil {
		t.Fatalf("Error parsing media type: %s", err.Error())
	}

	err = mt.Decode(nil, buf)
	if err != nil {
		t.Fatalf("Decoding error: %s", err.Error())
	}
}

type PersonDecoder struct {
	body io.Reader
}

func (d *PersonDecoder) Decode(v interface{}) error {
	if p, ok := v.(*Person); ok {
		by, err := ioutil.ReadAll(d.body)
		if err != nil {
			return err
		}
		p.Name = string(by)
	}
	return nil
}

func init() {
	AddDecoder("test", func(r io.Reader) Decoder {
		return &PersonDecoder{r}
	})
}
