package commands

import (
	"testing"
	"time"

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
	format := "%sC%>(8)%i%Creset  %t%  l%n"
	testFormatIssue(t, []formatIssueTest{
		{
			name: "standard usage",
			issue: github.Issue{
				Number:    42,
				Title:     "Just an Issue",
				State:     "open",
				User:      &github.User{Login: "pcorpet"},
				Body:      "Body of the\nissue",
				Assignees: []github.User{{Login: "mislav"}},
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
				User:   &github.User{Login: "octocat"},
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
				User:   &github.User{Login: "octocat"},
				Labels: []github.IssueLabel{
					{Name: "bug", Color: "800000"},
					{Name: "reproduced", Color: "55ff55"},
				},
			},
			format:   format,
			colorize: true,
			expect:   "\033[32m     #42\033[m  An issue with labels  \033[38;2;255;255;255;48;2;128;0;0m bug \033[m \033[38;2;0;0;0;48;2;85;255;85m reproduced \033[m\n",
		},
		{
			name: "not colorized",
			issue: github.Issue{
				Number: 42,
				Title:  "Just an Issue",
				State:  "open",
				User:   &github.User{Login: "octocat"},
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
				User:   &github.User{Login: "octocat"},
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
	createdAt, err := time.Parse(time.RFC822Z, "16 Mar 15 12:34 +0000")
	if err != nil {
		t.Fatal(err)
	}
	updatedAt, err := time.Parse(time.RFC822Z, "17 Mar 15 12:34 +0900")
	if err != nil {
		t.Fatal(err)
	}

	issue := github.Issue{
		Number: 42,
		Title:  "Just an Issue",
		State:  "open",
		User:   &github.User{Login: "pcorpet"},
		Body:   "Body of the\nissue",
		Assignees: []github.User{
			{Login: "mislav"},
			{Login: "josh"},
		},
		Labels: []github.IssueLabel{
			{Name: "bug", Color: "880000"},
			{Name: "feature", Color: "008800"},
		},
		HtmlUrl:  "the://url",
		Comments: 12,
		Milestone: &github.Milestone{
			Number: 31,
			Title:  "2.2-stable",
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	testFormatIssue(t, []formatIssueTest{
		{
			name:     "number",
			issue:    issue,
			format:   "%I",
			colorize: true,
			expect:   "42",
		},
		{
			name:     "hashed number",
			issue:    issue,
			format:   "%i",
			colorize: true,
			expect:   "#42",
		},
		{
			name:     "state as text",
			issue:    issue,
			format:   "%S",
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
			expect:   "\033[38;2;255;255;255;48;2;136;0;0m bug \033[m \033[38;2;255;255;255;48;2;0;136;0m feature \033[m",
		},
		{
			name:     "label not colorized",
			issue:    issue,
			format:   "%l",
			colorize: false,
			expect:   " bug   feature ",
		},
		{
			name:     "raw labels",
			issue:    issue,
			format:   "%L",
			colorize: true,
			expect:   "bug, feature",
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
			format:   "%au",
			colorize: true,
			expect:   "pcorpet",
		},
		{
			name:     "assignee login",
			issue:    issue,
			format:   "%as",
			colorize: true,
			expect:   "mislav, josh",
		},
		{
			name: "assignee login but not assigned",
			issue: github.Issue{
				State: "open",
				User:  &github.User{Login: "pcorpet"},
			},
			format:   "%as",
			colorize: true,
			expect:   "",
		},
		{
			name:     "milestone number",
			issue:    issue,
			format:   "%Mn",
			colorize: true,
			expect:   "31",
		},
		{
			name:     "milestone title",
			issue:    issue,
			format:   "%Mt",
			colorize: true,
			expect:   "2.2-stable",
		},
		{
			name:     "comments number",
			issue:    issue,
			format:   "%Nc",
			colorize: true,
			expect:   "(12)",
		},
		{
			name:     "raw comments number",
			issue:    issue,
			format:   "%NC",
			colorize: true,
			expect:   "12",
		},
		{
			name:     "issue URL",
			issue:    issue,
			format:   "%U",
			colorize: true,
			expect:   "the://url",
		},
		{
			name:     "created date",
			issue:    issue,
			format:   "%cD",
			colorize: true,
			expect:   "16 Mar 2015",
		},
		{
			name:     "created time ISO 8601",
			issue:    issue,
			format:   "%cI",
			colorize: true,
			expect:   "2015-03-16T12:34:00Z",
		},
		{
			name:     "created time Unix",
			issue:    issue,
			format:   "%ct",
			colorize: true,
			expect:   "1426509240",
		},
		{
			name:     "updated date",
			issue:    issue,
			format:   "%uD",
			colorize: true,
			expect:   "17 Mar 2015",
		},
		{
			name:     "updated time ISO 8601",
			issue:    issue,
			format:   "%uI",
			colorize: true,
			expect:   "2015-03-17T12:34:00+09:00",
		},
		{
			name:     "updated time Unix",
			issue:    issue,
			format:   "%ut",
			colorize: true,
			expect:   "1426563240",
		},
	})
}
