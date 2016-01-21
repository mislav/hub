package mediatype

import (
	"github.com/bmizerany/assert"
	"io"
	"strings"
	"testing"
)

func TestAddEncoder(t *testing.T) {
	mt, err := Parse("application/test+test")
	if err != nil {
		t.Fatalf("Error parsing media type: %s", err.Error())
	}

	person := &Person{"bob"}
	buf, err := mt.Encode(person)
	if err != nil {
		t.Fatalf("Error encoding: %s", err.Error())
	}

	assert.Equal(t, "bob", buf.String())
}

func TetRequiresEncoder(t *testing.T) {
	mt, err := Parse("application/test+whatevs")
	if err != nil {
		t.Fatalf("Error parsing media type: %s", err.Error())
	}

	person := &Person{"bob"}
	_, err = mt.Encode(person)
	if err == nil {
		t.Fatal("No encoding error")
	}

	if !strings.HasPrefix(err.Error(), "No encoder found for format whatevs") {
		t.Fatalf("Bad error: %s", err)
	}
}

func TetRequiresEncodedResource(t *testing.T) {
	mt, err := Parse("application/test+test")
	if err != nil {
		t.Fatalf("Error parsing media type: %s", err.Error())
	}

	_, err = mt.Encode(nil)
	if err == nil {
		t.Fatal("No encoding error")
	}

	assert.Equal(t, "Nothing to encode", err.Error())
}

type PersonEncoder struct {
	body io.Writer
}

func (d *PersonEncoder) Encode(v interface{}) error {
	if p, ok := v.(*Person); ok {
		d.body.Write([]byte(p.Name))
	}
	return nil
}

func init() {
	AddEncoder("test", func(w io.Writer) Encoder {
		return &PersonEncoder{w}
	})
}
