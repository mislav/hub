package github

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/github/hub/v2/internal/assert"
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
