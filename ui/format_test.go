package ui

import (
	"testing"
)

func TestExpand(t *testing.T) {
	tests := []struct {
		name   string
		format string
		values map[string]string
		expect string
	}{
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
			name:   "Colors",
			format: "%Cred%r %Cgreen%g %Cblue%b%Creset normal",
			values: map[string]string{"r": "RED", "g": "GREEN", "b": "BLUE"},
			expect: "\033[31mRED \033[32mGREEN \033[34mBLUE\033[m normal",
		},
		{
			name:   "Byte from hex code",
			format: "%x00 %x3712%x61 %x%x1%xga",
			expect: "\x00 \x3712a %x%x1%xga",
		},
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
	}

	for _, test := range tests {
		if got := Expand(test.format, test.values); got != test.expect {
			t.Errorf("%s: Expand(%q, ...) = %q, want %q", test.name, test.format, got, test.expect)
		}
	}
}
