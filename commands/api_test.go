package commands

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestMagicValue(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		input    string
		expected interface{}
	}{
		{
			"true",
			true,
		},
		{
			"false",
			false,
		},
		{
			"null",
			nil,
		},
		{
			"50",
			50,
		},
		{
			"@testdata/some-file.txt",
			"this\nis\na\ntest\nfile\n",
		},
		{
			"whatever",
			"whatever",
		},
		{
			"[v1, v2]",
			[]interface{}{"v1", "v2"},
		},
		{
			"[1, true, false, v5]",
			[]interface{}{1, true, false, "v5"},
		},
		{
			"[]",
			[]interface{}{},
		},
		{
			`[v1, v2, "v3,v4"]`,
			[]interface{}{"v1", "v2", "v3,v4"},
		},
		{
			`[v1, v2, v3"]`,
			[]interface{}{"v1", "v2", `v3"`},
		},
		{
			`[v1,    , v3"]`,
			[]interface{}{"v1", "", `v3"`},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			value := magicValue(test.input)
			t.Log(value)
			assert.Equal(t, test.expected, value)
		})
	}
}
