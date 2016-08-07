package commands

import (
	"testing"

	"github.com/github/hub/github"
)

type formatIssueTest struct {
	name     string
	issue    github.Issue
	format   string
	colorize bool
	expect   string
}

func testFormatIssue(t *testing.T, tests []formatIssueTest) {
	for _, test := range tests {
		if got := formatIssue(test.issue, test.format, test.colorize); got != test.expect {
			t.Errorf("%s: formatIssue(..., %q, %t) = %q, want %q", test.name, test.format, test.colorize, got, test.expect)
		}
	}
}

func TestFormatIssue(t *testing.T) {
	format := "%sC%>(8)%ih%Creset  %t%  l%n"
	testFormatIssue(t, []formatIssueTest{
		{
			name: "standard usage",
			issue: github.Issue{
				Number:   42,
				Title:    "Just an Issue",
				State:    "open",
				User:     &github.User{Login: "pcorpet"},
				Body:     "Body of the\nissue",
				Assignee: &github.User{Login: "mislav"},
			},
			format:   format,
			colorize: true,
			expect:   "\033[32m     #42\033[m  Just an Issue\n",
		},
		{
			name: "closed issue colored differently",
			issue: github.Issue{
				Number: 42,
				Title:  "Just an Issue",
				State:  "closed",
			},
			format:   format,
			colorize: true,
			expect:   "\033[31m     #42\033[m  Just an Issue\n",
		},
		{
			name: "labels",
			issue: github.Issue{
				Number: 42,
				Title:  "An issue with labels",
				State:  "open",
				Labels: []github.IssueLabel{
					{Name: "bug", Color: "800000"},
					{Name: "reproduced", Color: "55ff55"},
				},
			},
			format:   format,
			colorize: true,
			expect:   "\033[32m     #42\033[m  An issue with labels  \033[38;5;15;48;2;128;0;0m bug \033[m \033[38;5;16;48;2;85;255;85m reproduced \033[m\n",
		},
		{
			name: "not colorized",
			issue: github.Issue{
				Number: 42,
				Title:  "Just an Issue",
				State:  "open",
			},
			format:   format,
			colorize: false,
			expect:   "     #42  Just an Issue\n",
		},
		{
			name: "labels not colorized",
			issue: github.Issue{
				Number: 42,
				Title:  "An issue with labels",
				State:  "open",
				Labels: []github.IssueLabel{
					{Name: "bug", Color: "880000"},
					{Name: "reproduced", Color: "55ff55"},
				},
			},
			format:   format,
			colorize: false,
			expect:   "     #42  An issue with labels   bug   reproduced \n",
		},
	})
}

func TestFormatIssue_customFormatString(t *testing.T) {
	issue := github.Issue{
		Number:   42,
		Title:    "Just an Issue",
		State:    "open",
		User:     &github.User{Login: "pcorpet"},
		Body:     "Body of the\nissue",
		Assignee: &github.User{Login: "mislav"},
		Labels: []github.IssueLabel{
			{Name: "bug", Color: "880000"},
		},
	}

	testFormatIssue(t, []formatIssueTest{
		{
			name:     "number",
			issue:    issue,
			format:   "%in",
			colorize: true,
			expect:   "42",
		},
		{
			name:     "hashed number",
			issue:    issue,
			format:   "%ih",
			colorize: true,
			expect:   "#42",
		},
		{
			name:     "state as text",
			issue:    issue,
			format:   "%st",
			colorize: true,
			expect:   "open",
		},
		{
			name:     "state as color switch",
			issue:    issue,
			format:   "%sC",
			colorize: true,
			expect:   "\033[32m",
		},
		{
			name:     "state as color switch non colorized",
			issue:    issue,
			format:   "%sC",
			colorize: false,
			expect:   "",
		},
		{
			name:     "title",
			issue:    issue,
			format:   "%t",
			colorize: true,
			expect:   "Just an Issue",
		},
		{
			name:     "label colorized",
			issue:    issue,
			format:   "%l",
			colorize: true,
			expect:   "\033[38;5;15;48;2;136;0;0m bug \033[m",
		},
		{
			name:     "label not colorized",
			issue:    issue,
			format:   "%l",
			colorize: false,
			expect:   " bug ",
		},
		{
			name:     "body",
			issue:    issue,
			format:   "%b",
			colorize: true,
			expect:   "Body of the\nissue",
		},
		{
			name:     "user login",
			issue:    issue,
			format:   "%u",
			colorize: true,
			expect:   "pcorpet",
		},
		{
			name:     "assignee login",
			issue:    issue,
			format:   "%a",
			colorize: true,
			expect:   "mislav",
		},
		{
			name: "assignee login but not assigned",
			issue: github.Issue{
				State: "open",
				User:  &github.User{Login: "pcorpet"},
			},
			format:   "%a",
			colorize: true,
			expect:   "",
		},
	})
}
