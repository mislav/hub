package shellquote

import (
	"reflect"
	"testing"
)

func TestSimpleSplit(t *testing.T) {
	for _, elem := range simpleSplitTest {
		output, err := Split(elem.input)
		if err != nil {
			t.Errorf("Input %q, got error %#v", elem.input, err)
		} else if !reflect.DeepEqual(output, elem.output) {
			t.Errorf("Input %q, got %q, expected %q", elem.input, output, elem.output)
		}
	}
}

func TestErrorSplit(t *testing.T) {
	for _, elem := range errorSplitTest {
		_, err := Split(elem.input)
		if err != elem.error {
			t.Errorf("Input %q, got error %#v, expected error %#v", elem.input, err, elem.error)
		}
	}
}

var simpleSplitTest = []struct {
	input  string
	output []string
}{
	{"hello", []string{"hello"}},
	{"hello goodbye", []string{"hello", "goodbye"}},
	{"hello   goodbye", []string{"hello", "goodbye"}},
	{"glob* test?", []string{"glob*", "test?"}},
	{"don\\'t you know the dewey decimal system\\?", []string{"don't", "you", "know", "the", "dewey", "decimal", "system?"}},
	{"'don'\\''t you know the dewey decimal system?'", []string{"don't you know the dewey decimal system?"}},
	{"one '' two", []string{"one", "", "two"}},
	{"text with\\\na backslash-escaped newline", []string{"text", "witha", "backslash-escaped", "newline"}},
	{"text \"with\na\" quoted newline", []string{"text", "with\na", "quoted", "newline"}},
	{"\"quoted\\d\\\\\\\" text with\\\na backslash-escaped newline\"", []string{"quoted\\d\\\" text witha backslash-escaped newline"}},
	{"foo\"bar\"baz", []string{"foobarbaz"}},
}

var errorSplitTest = []struct {
	input string
	error error
}{
	{"don't worry", UnterminatedSingleQuoteError},
	{"'test'\\''ing", UnterminatedSingleQuoteError},
	{"\"foo'bar", UnterminatedDoubleQuoteError},
	{"foo\\", UnterminatedEscapeError},
}
