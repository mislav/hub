// Copyright © 2010 Fazlul Shahriar <fshahriar@gmail.com>.
// See LICENSE file for license details.

package netrc_test

import (
	. "code.google.com/p/go-netrc/netrc"
	"testing"
)

var expectedMach = []*Machine{
	&Machine{"mail.google.com", "joe@gmail.com", "somethingSecret", "gmail"},
	&Machine{"ray", "demo", "mypassword", ""},
	&Machine{"", "anonymous", "joe@example.com", ""},
}
var expectedMac = Macros{
	"allput": "put src/*",
}

func eqMach(a *Machine, b *Machine) bool {
	return a.Name == b.Name &&
		a.Login == b.Login &&
		a.Password == b.Password &&
		a.Account == b.Account
}

func TestParse(t *testing.T) {
	mach, mac, err := ParseFile("example.netrc")
	if err != nil {
		t.Fatal(err)
	}

	for i, e := range expectedMach {
		if !eqMach(e, mach[i]) {
			t.Errorf("bad machine; expected %v, got %v\n", e, mach[i])
		}
	}

	for k, v := range expectedMac {
		if v != mac[k] {
			t.Errorf("bad macro for %s; expected %s, got %s\n", k, v, mac[k])
		}
	}
}

func TestFindMachine(t *testing.T) {
	m, err := FindMachine("example.netrc", "ray")
	if err != nil {
		t.Fatal(err)
	}
	if !eqMach(m, expectedMach[1]) {
		t.Errorf("bad machine; expected %v, got %v\n", expectedMach[1], m)
	}

	m, err = FindMachine("example.netrc", "non.existent")
	if err != nil {
		t.Fatal(err)
	}
	if !eqMach(m, expectedMach[2]) {
		t.Errorf("bad machine; expected %v, got %v\n", expectedMach[2], m)
	}
}
