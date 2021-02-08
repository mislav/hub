package utils

import (
	"errors"
	"reflect"
	"testing"
)

func equal(t *testing.T, expected, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected: %#v, got: %#v", expected, got)
	}
}

func TestArgsParser(t *testing.T) {
	p := NewArgsParser()
	p.RegisterValue("--hello", "-e")
	p.RegisterValue("--origin", "-o")
	args := []string{"--hello", "world", "one", "--", "--two"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{"one", "--two"}, rest)
	equal(t, "world", p.Value("--hello"))
	equal(t, true, p.HasReceived("--hello"))
	equal(t, "", p.Value("-e"))
	equal(t, false, p.HasReceived("-e"))
	equal(t, "", p.Value("--origin"))
	equal(t, false, p.HasReceived("--origin"))
	equal(t, []int{2, 4}, p.PositionalIndices)
}

func TestArgsParser_RepeatedInvocation(t *testing.T) {
	p := NewArgsParser()
	p.RegisterValue("--hello", "-e")
	p.RegisterValue("--origin", "-o")

	rest, err := p.Parse([]string{"--hello", "world", "--", "one"})
	equal(t, nil, err)
	equal(t, []string{"one"}, rest)
	equal(t, []int{3}, p.PositionalIndices)
	equal(t, true, p.HasReceived("--hello"))
	equal(t, "world", p.Value("--hello"))
	equal(t, false, p.HasReceived("--origin"))
	equal(t, true, p.HasTerminated)

	rest, err = p.Parse([]string{"two", "-oupstream"})
	equal(t, nil, err)
	equal(t, []string{"two"}, rest)
	equal(t, []int{0}, p.PositionalIndices)
	equal(t, false, p.HasReceived("--hello"))
	equal(t, true, p.HasReceived("--origin"))
	equal(t, "upstream", p.Value("--origin"))
	equal(t, false, p.HasTerminated)
}

func TestArgsParser_UnknownFlag(t *testing.T) {
	p := NewArgsParser()
	p.RegisterValue("--hello")
	p.RegisterBool("--yes", "-y")

	args := []string{"--hello", "world", "--nonexist", "one", "--", "--two"}
	rest, err := p.Parse(args)
	equal(t, errors.New("unknown flag: '--nonexist'"), err)
	equal(t, []string{"one", "--two"}, rest)

	rest, err = p.Parse([]string{"one", "-yelp"})
	equal(t, errors.New("unknown shorthand flag: 'e' in -yelp"), err)
	equal(t, []string{"one"}, rest)
	equal(t, true, p.Bool("--yes"))
}

func TestArgsParser_BlankArgs(t *testing.T) {
	p := NewArgsParser()
	rest, err := p.Parse([]string{"", ""})
	equal(t, nil, err)
	equal(t, []string{"", ""}, rest)
	equal(t, []int{0, 1}, p.PositionalIndices)
}

func TestArgsParser_Values(t *testing.T) {
	p := NewArgsParser()
	p.RegisterValue("--origin", "-o")
	args := []string{"--origin=a=b", "--origin=", "--origin", "c", "-o"}
	rest, err := p.Parse(args)
	equal(t, errors.New("no value given for '-o'"), err)
	equal(t, []string{}, rest)
	equal(t, []string{"a=b", "", "c"}, p.AllValues("--origin"))
}

func TestArgsParser_Bool(t *testing.T) {
	p := NewArgsParser()
	p.RegisterBool("--noop")
	p.RegisterBool("--color")
	p.RegisterBool("--draft", "-d")
	args := []string{"-d", "--draft=false", "--color=auto"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{}, rest)
	equal(t, false, p.Bool("--draft"))
	equal(t, true, p.HasReceived("--draft"))
	equal(t, false, p.HasReceived("-d"))
	equal(t, false, p.HasReceived("--noop"))
	equal(t, false, p.Bool("--noop"))
	equal(t, true, p.HasReceived("--color"))
	equal(t, "auto", p.Value("--color"))
}

func TestArgsParser_BoolValue(t *testing.T) {
	p := NewArgsParser()
	p.RegisterBool("--draft")
	args := []string{"--draft=yes pls"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{}, rest)
	equal(t, true, p.HasReceived("--draft"))
	equal(t, true, p.Bool("--draft"))
	equal(t, "yes pls", p.Value("--draft"))
}

func TestArgsParser_BoolValue_multiple(t *testing.T) {
	p := NewArgsParser()
	p.RegisterBool("--draft")
	p.RegisterBool("--prerelease")
	args := []string{"--draft=false", "--prerelease"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{}, rest)
	equal(t, false, p.Bool("--draft"))
	equal(t, true, p.Bool("--prerelease"))
}

func TestArgsParser_Shorthand(t *testing.T) {
	p := NewArgsParser()
	p.RegisterValue("--origin", "-o")
	p.RegisterBool("--draft", "-d")
	p.RegisterBool("--copy", "-c")
	args := []string{"-co", "one", "-dotwo"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{}, rest)
	equal(t, []string{"one", "two"}, p.AllValues("--origin"))
	equal(t, true, p.Bool("--draft"))
	equal(t, true, p.Bool("--copy"))
}

func TestArgsParser_ShorthandEdgeCase(t *testing.T) {
	p := NewArgsParser()
	p.RegisterBool("--draft", "-d")
	p.RegisterBool("-f")
	p.RegisterBool("-a")
	p.RegisterBool("-l")
	p.RegisterBool("-s")
	p.RegisterBool("-e")
	args := []string{"-dfalse"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{}, rest)
	equal(t, true, p.Bool("--draft"))
}

func TestArgsParser_Dashes(t *testing.T) {
	p := NewArgsParser()
	p.RegisterValue("--file", "-F")
	args := []string{"-F-", "-", "--", "-F", "--"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{"-", "-F", "--"}, rest)
	equal(t, "-", p.Value("--file"))
}

func TestArgsParser_RepeatedArg(t *testing.T) {
	p := NewArgsParser()
	p.RegisterValue("--msg", "-m")
	args := []string{"--msg=hello", "-m", "world", "--msg", "how", "-mare you?"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{}, rest)
	equal(t, "are you?", p.Value("--msg"))
	equal(t, []string{"hello", "world", "how", "are you?"}, p.AllValues("--msg"))
}

func TestArgsParser_Int(t *testing.T) {
	p := NewArgsParser()
	p.RegisterValue("--limit", "-L")
	p.RegisterValue("--depth", "-d")
	args := []string{"-L24", "-d", "-3"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{}, rest)
	equal(t, true, p.HasReceived("--limit"))
	equal(t, 24, p.Int("--limit"))
	equal(t, true, p.HasReceived("--depth"))
	equal(t, -3, p.Int("--depth"))
}

func TestArgsParser_WithUsage(t *testing.T) {
	p := NewArgsParserWithUsage(`
		-L, --limit N
			retrieve at most N records
		-d, --draft
			save as draft
		--message=<msg>, -m <msg>
			set message body
	`)
	args := []string{"-L24", "-d", "-mhello"}
	rest, err := p.Parse(args)
	equal(t, nil, err)
	equal(t, []string{}, rest)
	equal(t, "24", p.Value("--limit"))
	equal(t, true, p.Bool("--draft"))
	equal(t, "hello", p.Value("--message"))
}
