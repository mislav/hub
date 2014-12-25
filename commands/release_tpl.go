package commands

import (
	"bytes"
	"text/template"
)

const releaseTmpl = `{{.CS}} Creating release {{.TagName}} for {{.ProjectName}} from {{.BranchName}}
{{.CS}}
{{.CS}} Write a message for this release. The first block
{{.CS}} of text is the title and the rest is description.`

type releaseMsg struct {
	CS          string
	TagName     string
	ProjectName string
	BranchName  string
}

func renderReleaseTpl(cs, tagName, projectName, branchName string) (string, error) {
	t, err := template.New("releaseTmpl").Parse(releaseTmpl)
	if err != nil {
		return "", err
	}

	msg := &releaseMsg{
		CS:          cs,
		TagName:     tagName,
		ProjectName: projectName,
		BranchName:  branchName,
	}

	var b bytes.Buffer
	err = t.Execute(&b, msg)

	return b.String(), err
}
