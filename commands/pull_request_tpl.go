package commands

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

const pullRequestTmpl = `{{if .InitMsg}}{{.InitMsg}}
{{end}}
{{.CS}} Requesting a pull to {{.Base}} from {{.Head}}
{{.CS}}
{{.CS}} Write a message for this pull request. The first block
{{.CS}} of text is the title and the rest is the description.{{if .HasCommitLogs}}
{{.CS}}
{{.CS}} Changes:
{{.CS}}
{{.FormattedCommitLogs}}{{end}}`

type pullRequestMsg struct {
	InitMsg    string
	CS         string
	Base       string
	Head       string
	CommitLogs string
}

func (p *pullRequestMsg) HasCommitLogs() bool {
	return len(p.CommitLogs) > 0
}

func (p *pullRequestMsg) FormattedCommitLogs() string {
	startRegexp := regexp.MustCompilePOSIX("^")
	endRegexp := regexp.MustCompilePOSIX(" +$")

	commitLogs := strings.TrimSpace(p.CommitLogs)
	commitLogs = startRegexp.ReplaceAllString(commitLogs, fmt.Sprintf("%s ", p.CS))
	commitLogs = endRegexp.ReplaceAllString(commitLogs, "")

	return commitLogs
}

func renderPullRequestTpl(initMsg, cs, base, head string, commitLogs string) (string, error) {
	t, err := template.New("pullRequestTmpl").Parse(pullRequestTmpl)
	if err != nil {
		return "", err
	}

	msg := &pullRequestMsg{
		InitMsg:    initMsg,
		CS:         cs,
		Base:       base,
		Head:       head,
		CommitLogs: commitLogs,
	}

	var b bytes.Buffer
	err = t.Execute(&b, msg)

	return b.String(), err
}
