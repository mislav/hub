// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uritemplates

import (
	"encoding/json"
	"os"
	"testing"
)

type spec struct {
	title  string
	values map[string]interface{}
	tests  []specTest
}
type specTest struct {
	template string
	expected []string
}

func loadSpec(t *testing.T, path string) []spec {

	file, err := os.Open(path)
	if err != nil {
		t.Errorf("Failed to load test specification: %s", err)
	}

	stat, _ := file.Stat()
	buffer := make([]byte, stat.Size())
	_, err = file.Read(buffer)
	if err != nil {
		t.Errorf("Failed to load test specification: %s", err)
	}

	var root_ interface{}
	err = json.Unmarshal(buffer, &root_)
	if err != nil {
		t.Errorf("Failed to load test specification: %s", err)
	}

	root := root_.(map[string]interface{})
	results := make([]spec, 1024)
	i := -1
	for title, spec_ := range root {
		i = i + 1
		results[i].title = title
		specMap := spec_.(map[string]interface{})
		results[i].values = specMap["variables"].(map[string]interface{})
		tests := specMap["testcases"].([]interface{})
		results[i].tests = make([]specTest, len(tests))
		for k, test_ := range tests {
			test := test_.([]interface{})
			results[i].tests[k].template = test[0].(string)
			switch typ := test[1].(type) {
			case string:
				results[i].tests[k].expected = make([]string, 1)
				results[i].tests[k].expected[0] = test[1].(string)
			case []interface{}:
				arr := test[1].([]interface{})
				results[i].tests[k].expected = make([]string, len(arr))
				for m, s := range arr {
					results[i].tests[k].expected[m] = s.(string)
				}
			case bool:
				results[i].tests[k].expected = make([]string, 0)
			default:
				t.Errorf("Unrecognized value type %v", typ)
			}
		}
	}
	return results
}

func runSpec(t *testing.T, path string) {
	var spec = loadSpec(t, path)
	for _, group := range spec {
		for _, test := range group.tests {
			template, err := Parse(test.template)
			if err != nil {
				if len(test.expected) > 0 {
					t.Errorf("%s: %s %v", group.title, err, test.template)
				}
				continue
			}
			result, err := template.Expand(group.values)
			if err != nil {
				if len(test.expected) > 0 {
					t.Errorf("%s: %s %v", group.title, err, test.template)
				}
				continue
			} else if len(test.expected) == 0 {
				t.Errorf("%s: should have failed while parsing or expanding %v but got %v", group.title, test.template, result)
				continue
			}
			pass := false
			for _, expected := range test.expected {
				if result == expected {
					pass = true
				}
			}
			if !pass {
				t.Errorf("%s: expected %v, but got %v", group.title, test.expected[0], result)
			}
		}
	}
}

func TestExtended(t *testing.T) {
	runSpec(t, "tests/extended-tests.json")
}

func TestNegative(t *testing.T) {
	runSpec(t, "tests/negative-tests.json")
}

func TestSpecExamples(t *testing.T) {
	runSpec(t, "tests/spec-examples.json")
}

var parse_tests = []struct {
	Template string
	ParseOk  bool
}{
	{
		// Syntax error, too many colons:
		"{opts:1:2}",
		false,
	},
}

func TestParse(t *testing.T) {
	for itest, test := range parse_tests {
		if _, err := Parse(test.Template); err != nil {
			if test.ParseOk {
				t.Errorf("%v", err)
			}
		} else if !test.ParseOk {
			t.Errorf("%d: expected error, got none.", itest)
		}
	}
}

type Location struct {
	Path    []interface{} `uri:"path"`
	Version int           `json:"version"`
	Opts    Options       `opts`
}

type Options struct {
	Format string `uri:"fmt"`
}

var expand_tests = []struct {
	Source   interface{}
	Template string
	Expected string
	ExpandOk bool
}{
	{
		// General struct expansion:
		Location{
			Path:    []interface{}{"main", "quux"},
			Version: 2,
			Opts: Options{
				Format: "pdf",
			},
		},
		"{/path*,Version}{?opts*}",
		"/main/quux/2?fmt=pdf",
		true,
	}, {
		// Pointer to struct:
		&Location{Opts: Options{Format: "pdf"}},
		"{?opts*}",
		"?fmt=pdf",
		true,
	}, {
		// Map expansion cannot be truncated:
		Location{Opts: Options{Format: "pdf"}},
		"{?opts:3}",
		"",
		false,
	}, {
		// Map whose values are not all strings:
		map[string]interface{}{
			"one": map[string]interface{}{
				"two": 42,
			},
		},
		"{?one*}",
		"?two=42",
		true,
	}, {
		// Value of inappropriate type:
		42,
		"{?one*}",
		"",
		false,
	}, {
		// Truncated array whose values are not all strings:
		map[string]interface{}{"one": []interface{}{1234}},
		"{?one:3}",
		"?one=123",
		true,
	},
}

func TestUriTemplate_Expand(t *testing.T) {
	for itest, test := range expand_tests {
		if template, err := Parse(test.Template); err != nil {
			t.Errorf("%d: %v", itest, err)
		} else if expanded, err := template.Expand(test.Source); err != nil {
			if test.ExpandOk {
				t.Errorf("%d: unexpected error: %v", itest, err)
			}
		} else if !test.ExpandOk {
			t.Errorf("%d: expected error, got none.", itest, err)
		} else if expanded != test.Expected {
			t.Errorf("%d: expected %v, got %v", itest, test.Expected, expanded)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("http://localhost:6060{/type,path}{.fmt}{?q*}")
	}
}

func BenchmarkExpand(b *testing.B) {
	templ, _ := Parse("http://localhost:6060{/type,path}{.fmt}{?q*}")
	data := map[string]interface{}{
		"type": "pkg",
		"path": [...]string{"github.com", "jtacoma", "uritemplates"},
		"q": map[string]interface{}{
			"somequery": "x!@#$",
			"other":     "y&*()",
		},
	}
	for i := 0; i < b.N; i++ {
		templ.Expand(data)
	}
}
