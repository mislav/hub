package github

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
)

func TestEditor_openAndEdit_deleteFileWhenOpeningEditorFails(t *testing.T) {
	tempFile, _ := ioutil.TempFile("", "editor-test")
	tempFile.Close()

	ioutil.WriteFile(tempFile.Name(), []byte("hello"), 0644)
	editor := Editor{
		Program: "memory",
		File:    tempFile.Name(),
		Topic:   "test",
		openEditor: func(program string, file string) error {
			assert.Equal(t, "memory", program)
			assert.Equal(t, tempFile.Name(), file)
			return fmt.Errorf("error")
		},
	}

	_, err := os.Stat(tempFile.Name())
	assert.Equal(t, nil, err)

	_, err = editor.openAndEdit()
	assert.Equal(t, "error using text editor for test message", fmt.Sprintf("%s", err))

	// file is removed if there's error
	_, err = os.Stat(tempFile.Name())
	assert.T(t, os.IsNotExist(err))
}

func TestEditor_openAndEdit_readFileIfExist(t *testing.T) {
	tempFile, _ := ioutil.TempFile("", "editor-test")
	tempFile.Close()

	ioutil.WriteFile(tempFile.Name(), []byte("hello"), 0644)
	editor := Editor{
		Program: "memory",
		File:    tempFile.Name(),
		openEditor: func(program string, file string) error {
			assert.Equal(t, "memory", program)
			assert.Equal(t, tempFile.Name(), file)

			return nil
		},
	}

	content, err := editor.openAndEdit()
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello", string(content))
}

func TestEditor_openAndEdit_writeFileIfNotExist(t *testing.T) {
	tempFile, _ := ioutil.TempFile("", "PULLREQ")
	tempFile.Close()

	editor := Editor{
		Program: "memory",
		File:    tempFile.Name(),
		openEditor: func(program string, file string) error {
			assert.Equal(t, "memory", program)
			assert.Equal(t, tempFile.Name(), file)

			return ioutil.WriteFile(file, []byte("hello"), 0644)
		},
	}

	content, err := editor.openAndEdit()
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello", string(content))
}

func TestEditor_EditTitleAndBodyEmptyTitle(t *testing.T) {
	tempFile, _ := ioutil.TempFile("", "PULLREQ")
	tempFile.Close()

	editor := Editor{
		Program: "memory",
		File:    tempFile.Name(),
		CS:      "#",
		openEditor: func(program string, file string) error {
			assert.Equal(t, "memory", program)
			assert.Equal(t, tempFile.Name(), file)
			return ioutil.WriteFile(file, []byte(""), 0644)
		},
	}

	title, body, err := editor.EditTitleAndBody()
	assert.Equal(t, nil, err)
	assert.Equal(t, "", title)
	assert.Equal(t, "", body)

	_, err = os.Stat(tempFile.Name())
	assert.T(t, os.IsNotExist(err))
}

func TestEditor_EditTitleAndBody(t *testing.T) {
	tempFile, _ := ioutil.TempFile("", "PULLREQ")
	tempFile.Close()

	editor := Editor{
		Program: "memory",
		File:    tempFile.Name(),
		CS:      "#",
		openEditor: func(program string, file string) error {
			assert.Equal(t, "memory", program)
			assert.Equal(t, tempFile.Name(), file)

			message := `A title
A title continues

A body
A body continues
# comment
`
			return ioutil.WriteFile(file, []byte(message), 0644)
		},
	}

	title, body, err := editor.EditTitleAndBody()
	assert.Equal(t, nil, err)
	assert.Equal(t, "A title A title continues", title)
	assert.Equal(t, "A body\nA body continues", body)
}

func TestReadTitleAndBody(t *testing.T) {
	message := `A title
A title continues

A body
A body continues
# comment
`
	r := strings.NewReader(message)
	reader := bufio.NewReader(r)
	title, body, err := readTitleAndBody(reader, "#")
	assert.Equal(t, nil, err)
	assert.Equal(t, "A title A title continues", title)
	assert.Equal(t, "A body\nA body continues", body)

	message = `# Dat title

/ This line is commented out.

Dem body.
`
	r = strings.NewReader(message)
	reader = bufio.NewReader(r)
	title, body, err = readTitleAndBody(reader, "/")
	assert.Equal(t, nil, err)
	assert.Equal(t, "# Dat title", title)
	assert.Equal(t, "Dem body.", body)
}

func TestGetMessageFile(t *testing.T) {
	gitPullReqMsgFile, _ := getMessageFile("PULLREQ")
	assert.T(t, strings.Contains(gitPullReqMsgFile, "PULLREQ_EDITMSG"))
}
