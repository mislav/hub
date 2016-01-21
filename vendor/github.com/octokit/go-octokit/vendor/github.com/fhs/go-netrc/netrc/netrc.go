// Copyright © 2010 Fazlul Shahriar <fshahriar@gmail.com>.
// See LICENSE file for license details.

// Package netrc implements a parser for netrc file format.
//
// A netrc file usually resides in $HOME/.netrc and is traditionally used
// by the ftp(1) program to look up login information (username, password,
// etc.) of remote system(s). The file format is (loosely) described in
// this man page: http://linux.die.net/man/5/netrc .
package netrc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode"
	"unicode/utf8"
)

const (
	tkMachine = iota
	tkDefault
	tkLogin
	tkPassword
	tkAccount
	tkMacdef
)

var tokenNames = []string{
	"Machine",
	"Default",
	"Login",
	"Password",
	"Account",
	"Macdef",
}

var keywords = map[string]int{
	"machine":  tkMachine,
	"default":  tkDefault,
	"login":    tkLogin,
	"password": tkPassword,
	"account":  tkAccount,
	"macdef":   tkMacdef,
}

// Machine contains information about a remote machine.
type Machine struct {
	Name     string
	Login    string
	Password string
	Account  string
}

// Macros contains all the macro definitions in a netrc file.
type Macros map[string]string

type token struct {
	kind      int
	macroName string
	value     string
}

type filePos struct {
	name string
	line int
}

// Error represents a netrc file parse error.
type Error struct {
	Filename string
	LineNum  int    // Line number
	Msg      string // Error message
}

// Error returns a string representation of error e.
func (e *Error) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.Filename, e.LineNum, e.Msg)
}

func getWord(b []byte, pos *filePos) (string, []byte) {
	// Skip over leading whitespace
	i := 0
	for i < len(b) {
		r, size := utf8.DecodeRune(b[i:])
		if r == '\n' {
			pos.line++
		}
		if !unicode.IsSpace(r) {
			break
		}
		i += size
	}
	b = b[i:]

	// Find end of word
	i = bytes.IndexFunc(b, unicode.IsSpace)
	if i < 0 {
		i = len(b)
	}
	return string(b[0:i]), b[i:]
}

func getToken(b []byte, pos *filePos) ([]byte, *token, error) {
	word, b := getWord(b, pos)
	if word == "" {
		return b, nil, nil // EOF reached
	}

	t := new(token)
	var ok bool
	t.kind, ok = keywords[word]
	if !ok {
		return b, nil, &Error{pos.name, pos.line, "keyword expected; got " + word}
	}
	if t.kind == tkDefault {
		return b, t, nil
	}

	word, b = getWord(b, pos)
	if word == "" {
		return b, nil, &Error{pos.name, pos.line, "word expected"}
	}
	if t.kind == tkMacdef {
		t.macroName = word

		// Macro value starts on next line. The rest of current line
		// should contain nothing but whitespace
		i := 0
		for i < len(b) {
			r, size := utf8.DecodeRune(b[i:])
			if r == '\n' {
				i += size
				pos.line++
				break
			}
			if !unicode.IsSpace(r) {
				return b, nil, &Error{pos.name, pos.line, "unexpected word"}
			}
			i += size
		}
		b = b[i:]

		// Find end of macro value
		i = bytes.Index(b, []byte("\n\n"))
		if i < 0 { // EOF reached
			i = len(b)
		}
		t.value = string(b[0:i])

		return b[i:], t, nil
	}
	t.value = word
	return b, t, nil
}

func parse(r io.Reader, pos *filePos) ([]*Machine, Macros, error) {
	// TODO(fhs): Clear memory containing password.
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	mach := make([]*Machine, 0, 20)
	mac := make(Macros, 10)
	var defaultSeen bool
	var m *Machine
	var t *token
	for {
		b, t, err = getToken(b, pos)
		if err != nil {
			return nil, nil, err
		}
		if t == nil {
			break
		}
		switch t.kind {
		case tkMacdef:
			mac[t.macroName] = t.value
		case tkDefault:
			if defaultSeen {
				return nil, nil, &Error{pos.name, pos.line, "multiple default token"}
			}
			if m != nil {
				mach, m = append(mach, m), nil
			}
			m = new(Machine)
			m.Name = ""
			defaultSeen = true
		case tkMachine:
			if m != nil {
				mach, m = append(mach, m), nil
			}
			m = new(Machine)
			m.Name = t.value
		case tkLogin:
			if m == nil || m.Login != "" {
				return nil, nil, &Error{pos.name, pos.line, "unexpected token login "}
			}
			m.Login = t.value
		case tkPassword:
			if m == nil || m.Password != "" {
				return nil, nil, &Error{pos.name, pos.line, "unexpected token password"}
			}
			m.Password = t.value
		case tkAccount:
			if m == nil || m.Account != "" {
				return nil, nil, &Error{pos.name, pos.line, "unexpected token account"}
			}
			m.Account = t.value
		}
	}
	if m != nil {
		mach, m = append(mach, m), nil
	}
	return mach, mac, nil
}

// ParseFile parses the netrc file identified by filename and returns the set of
// machine information and macros defined in it. The ``default'' machine,
// which is intended to be used when no machine name matches, is identified
// by an empty machine name. There can be only one ``default'' machine.
//
// If there is a parsing error, an Error is returned.
func ParseFile(filename string) ([]*Machine, Macros, error) {
	// TODO(fhs): Check if file is readable by anyone besides the user if there is password in it.
	fd, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer fd.Close()
	return parse(fd, &filePos{filename, 1})
}

// FindMachine parses the netrc file identified by filename and returns
// the Machine named by name. If no Machine with name name is found, the
// ``default'' machine is returned.
func FindMachine(filename, name string) (*Machine, error) {
	mach, _, err := ParseFile(filename)
	if err != nil {
		return nil, err
	}
	var def *Machine
	for _, m := range mach {
		if m.Name == name {
			return m, nil
		}
		if m.Name == "" {
			def = m
		}
	}
	if def == nil {
		return nil, errors.New("no machine found")
	}
	return def, nil
}
