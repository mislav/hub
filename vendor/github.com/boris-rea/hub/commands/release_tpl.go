package commands

import (
	"bytes"
	"text/template"
)

const releaseTmpl = `
{{.CS}} {{.Operation}} release {{.TagName}} for {{.ProjectName}}{{if .BranchName}} from {{.BranchName}}{{end}}
{{.CS}}
{{.CS}} Write a message for this release. The first block of
{{.CS}} text is the title and the rest is the description.`

type releaseMsg struct {
	Operation   string
	CS          string
	TagName     string
	ProjectName string
	BranchName  string
}

func renderReleaseTpl(operation, cs, tagName, projectName, branchName string) (string, error) {
	t, err := template.New("releaseTmpl").Parse(releaseTmpl)
	if err != nil {
		return "", err
	}

	msg := &releaseMsg{
		Operation:   operation,
		CS:          cs,
		TagName:     tagName,
		ProjectName: projectName,
		BranchName:  branchName,
	}

	var b bytes.Buffer
	err = t.Execute(&b, msg)

	return b.String(), err
}
