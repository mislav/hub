package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type argsFlag struct {
	expectsValue bool
	values       []string
}

func (f *argsFlag) addValue(v string) {
	f.values = append(f.values, v)
}

func (f *argsFlag) lastValue() string {
	l := len(f.values)
	if l > 0 {
		return f.values[l-1]
	} else {
		return ""
	}
}

func (f *argsFlag) reset() {
	if len(f.values) > 0 {
		f.values = []string{}
	}
}

type ArgsParser struct {
	flagMap           map[string]*argsFlag
	flagAliases       map[string]string
	PositionalIndices []int
	terminated        bool
}

func (p *ArgsParser) Parse(args []string) ([]string, error) {
	var flagName string
	var flagValue string
	var hasFlagValue bool
	var i int
	var arg string

	p.terminated = false
	for _, f := range p.flagMap {
		f.reset()
	}
	if len(p.PositionalIndices) > 0 {
		p.PositionalIndices = []int{}
	}

	positional := []string{}
	var parseError error
	logError := func(f string, p ...interface{}) {
		if parseError == nil {
			parseError = fmt.Errorf(f, p...)
		}
	}

	acknowledgeFlag := func() bool {
		canonicalFlagName := flagName
		if n, found := p.flagAliases[flagName]; found {
			canonicalFlagName = n
		}
		f := p.flagMap[canonicalFlagName]
		if f == nil {
			if len(flagName) == 2 {
				logError("unknown shorthand flag: '%s' in %s", flagName[1:], arg)
			} else {
				logError("unknown flag: '%s'", flagName)
			}
			return true
		}
		if f.expectsValue {
			if !hasFlagValue {
				i++
				if i < len(args) {
					flagValue = args[i]
				} else {
					logError("no value given for '%s'", flagName)
					return true
				}
			}
		} else if hasFlagValue && len(flagName) <= 2 {
			flagValue = ""
		}
		f.addValue(flagValue)
		return f.expectsValue
	}

	for i = 0; i < len(args); i++ {
		arg = args[i]

		if p.terminated || arg == "-" {
		} else if arg == "--" {
			if !p.terminated {
				p.terminated = true
				continue
			}
		} else if strings.HasPrefix(arg, "--") {
			flagName = arg
			eq := strings.IndexByte(arg, '=')
			hasFlagValue = eq >= 0
			if hasFlagValue {
				flagName = arg[:eq]
				flagValue = arg[eq+1:]
			}
			acknowledgeFlag()
			continue
		} else if arg[0] == '-' {
			for j := 1; j < len(arg); j++ {
				flagName = "-" + arg[j:j+1]
				flagValue = ""
				hasFlagValue = j+1 < len(arg)
				if hasFlagValue {
					flagValue = arg[j+1:]
				}
				if acknowledgeFlag() {
					break
				}
			}
			continue
		}

		p.PositionalIndices = append(p.PositionalIndices, i)
		positional = append(positional, arg)
	}

	return positional, parseError
}

func (p *ArgsParser) RegisterValue(name string, aliases ...string) {
	f := &argsFlag{expectsValue: true}
	p.flagMap[name] = f
	for _, alias := range aliases {
		p.flagAliases[alias] = name
	}
}

func (p *ArgsParser) RegisterBool(name string, aliases ...string) {
	f := &argsFlag{expectsValue: false}
	p.flagMap[name] = f
	for _, alias := range aliases {
		p.flagAliases[alias] = name
	}
}

func (p *ArgsParser) Value(name string) string {
	if f, found := p.flagMap[name]; found {
		return f.lastValue()
	} else {
		return ""
	}
}

func (p *ArgsParser) AllValues(name string) []string {
	if f, found := p.flagMap[name]; found {
		return f.values
	} else {
		return []string{}
	}
}

func (p *ArgsParser) Bool(name string) bool {
	if f, found := p.flagMap[name]; found {
		return len(f.values) > 0 && f.lastValue() != "false"
	} else {
		return false
	}
}

func (p *ArgsParser) Int(name string) int {
	i, _ := strconv.Atoi(p.Value(name))
	return i
}

func (p *ArgsParser) HasReceived(name string) bool {
	f, found := p.flagMap[name]
	return found && len(f.values) > 0
}

func NewArgsParser() *ArgsParser {
	return &ArgsParser{
		flagMap:     make(map[string]*argsFlag),
		flagAliases: make(map[string]string),
	}
}
