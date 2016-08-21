package ui

import (
	"testing"
)

type expanderTest struct {
	name     string
	format   string
	values   map[string]string
	colorize bool
	expect   string
}

func testExpander(t *testing.T, tests []expanderTest) {
	for _, test := range tests {
		if got := Expand(test.format, test.values, test.colorize); got != test.expect {
			t.Errorf("%s: Expand(%q, ...) = %q, want %q", test.name, test.format, got, test.expect)
		}
	}
}

func TestExpand(t *testing.T) {
	testExpander(t, []expanderTest{
		{
			name:   "Simple example",
			format: "The author of %h was %an, %ar%nThe title was >>%s<<%n",
			values: map[string]string{
				"h":  "fe6e0ee",
				"an": "Junio C Hamano",
				"ar": "23 hours ago",
				"s":  "t4119: test autocomputing -p<n> for traditional diff input.",
			},
			expect: "The author of fe6e0ee was Junio C Hamano, 23 hours ago\nThe title was >>t4119: test autocomputing -p<n> for traditional diff input.<<\n",
		},
		{
			name:   "Percent sign, middle and trailing",
			format: "%%a %%b %",
			values: map[string]string{"a": "A variable that should not be used."},
			expect: "%a %b %",
		},
		{
			name:     "Colors",
			format:   "%Cred%r %Cgreen%g %Cblue%b%Creset normal",
			values:   map[string]string{"r": "RED", "g": "GREEN", "b": "BLUE"},
			colorize: true,
			expect:   "\033[31mRED \033[32mGREEN \033[34mBLUE\033[m normal",
		},
		{
			name:     "Colors not colorized",
			format:   "%Cred%r %Cgreen%g %Cblue%b%Creset normal",
			values:   map[string]string{"r": "RED", "g": "GREEN", "b": "BLUE"},
			colorize: false,
			expect:   "RED GREEN BLUE normal",
		},
		{
			name:   "Byte from hex code",
			format: "%x00 %x3712%x61 %x%x1%xga",
			expect: "\x00 \x3712a %x%x1%xga",
		},
	})
}

func TestExpand_Modifiers(t *testing.T) {
	testExpander(t, []expanderTest{
		{
			name:   "plus modifier, conditional line",
			format: "line1%+a line2%+b line3",
			values: map[string]string{"a": "A", "b": ""},
			expect: "line1\nA line2 line3",
		},
		{
			name:   "blank modifier, conditional blank",
			format: "word1% a word2% b word3",
			values: map[string]string{"a": "A", "b": ""},
			expect: "word1 A word2 word3",
		},
		{
			name:   "minus modifier, crush preceding line-feeds",
			format: "word1%n%n%-a",
			values: map[string]string{"a": ""},
			expect: "word1",
		},
	})
}

func TestExpand_Padding(t *testing.T) {
	testExpander(t, []expanderTest{
		{
			name:   "padding",
			format: "%<(10)%a",
			values: map[string]string{"a": "012"},
			expect: "012       ",
		},
		{
			name:   "padding, wrong number",
			format: "%<(1a)%a",
			values: map[string]string{"a": "012"},
			expect: "%<(1a)012",
		},
		{
			name:   "padding left",
			format: "%>(10)%a",
			values: map[string]string{"a": "012"},
			expect: "       012",
		},
		{
			name:   "padding middle",
			format: "%><(10)%a",
			values: map[string]string{"a": "0123"},
			expect: "   0123   ",
		},
		{
			name:   "padding middle (odd # of blanks)",
			format: "%><(10)%a",
			values: map[string]string{"a": "012"},
			expect: "   012    ",
		},
		{
			name:   "padding uses extra blank on the left",
			format: "%>>(5)|    %a",
			values: map[string]string{"a": "0123456"},
			expect: "|  0123456",
		},
		{
			name:   "padding until column N",
			format: "%>|(10)abcdef%a",
			values: map[string]string{"a": "012"},
			expect: "abcdef 012",
		},
	})
}

func TestExpand_Truncing(t *testing.T) {
	testExpander(t, []expanderTest{
		{
			name:   "truncing",
			format: "%>(5,trunc)%a",
			values: map[string]string{"a": "0123456"},
			expect: "012..",
		},
		{
			name:   "truncing on the right",
			format: "%>(5,rtrunc)%a",
			values: map[string]string{"a": "0123456"},
			expect: "..456",
		},
		{
			name:   "truncing in the middle",
			format: "%>(6,mtrunc)%a",
			values: map[string]string{"a": "0123456"},
			expect: "01..56",
		},
		{
			name:   "truncing in the middle (odd # of chars)",
			format: "%>(5,mtrunc)%a",
			values: map[string]string{"a": "0123456"},
			expect: "0..56",
		},
		{
			name:   "truncing not enough space",
			format: "%>(1,trunc)%a",
			values: map[string]string{"a": "0123456"},
			expect: "..",
		},
		{
			name:   "truncing but use extra blanks on the left",
			format: "%>>(3,trunc)|   %a",
			values: map[string]string{"a": "0123456"},
			expect: "|0123..",
		},
	})
}
