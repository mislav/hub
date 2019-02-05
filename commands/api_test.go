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
			[]string{"v1", "v2"},
		},
		{
			"[]",
			[]string{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			value := magicValue(test.input)
			assert.Equal(t, test.expected, value)
		})
	}
}
