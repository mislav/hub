package github

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/github/hub/cmd"
	"github.com/github/hub/git"
)

func NewEditor(filePrefix, topic, message string) (editor *Editor, err error) {
	messageFile, err := getMessageFile(filePrefix)
	if err != nil {
		return
	}

	program, err := git.Editor()
	if err != nil {
		return
	}

	editor = &Editor{
		Program:    program,
		Topic:      topic,
		File:       messageFile,
		Message:    message,
		openEditor: doTextEditorEdit,
	}

	return
}

type Editor struct {
	Program    string
	Topic      string
	File       string
	Message    string
	openEditor func(program, file string) error
}

func (e *Editor) isFileExist() bool {
	_, err := os.Stat(e.File)
	return err == nil || !os.IsNotExist(err)
}

func (e *Editor) WriteContent() (err error) {
	// only write message if file doesn't exist
	if !e.isFileExist() && e.Message != "" {
		err = ioutil.WriteFile(e.File, []byte(e.Message), 0644)
		if err != nil {
			return
		}
	}

	return
}

func (e *Editor) readContent() (content []byte, err error) {
	return ioutil.ReadFile(e.File)
}

func (e *Editor) openAndEdit() (content []byte, err error) {
	err = e.WriteContent()
	if err != nil {
		return
	}

	// always remove file after editing is done
	defer os.Remove(e.File)

	err = e.openEditor(e.Program, e.File)
	if err != nil {
		err = fmt.Errorf("error using text editor for %s message", e.Topic)
		return
	}

	content, err = e.readContent()

	return
}

func (e *Editor) EditTitleAndBody() (title, body string, err error) {
	content, err := e.openAndEdit()
	if err != nil {
		return
	}

	content = bytes.TrimSpace(content)
	reader := bufio.NewReader(bytes.NewReader(content))
	title, body, err = readTitleAndBody(reader)

	return
}

func doTextEditorEdit(program, file string) error {
	editCmd := cmd.New(program)
	r := regexp.MustCompile("[mg]?vi[m]$")
	if r.MatchString(program) {
		editCmd.WithArg("-c")
		editCmd.WithArg("set ft=gitcommit tw=0 wrap lbr")
	}
	editCmd.WithArg(file)

	return editCmd.Exec()
}

func readTitleAndBody(reader *bufio.Reader) (title, body string, err error) {
	r := regexp.MustCompile("\\S")
	var titleParts, bodyParts []string

	line, err := readLine(reader)
	for err == nil {
		if strings.HasPrefix(line, "#") {
			break
		}

		if len(bodyParts) == 0 && r.MatchString(line) {
			titleParts = append(titleParts, line)
		} else {
			bodyParts = append(bodyParts, line)
		}

		line, err = readLine(reader)
	}

	if err == io.EOF {
		err = nil
	}

	title = strings.Join(titleParts, " ")
	title = strings.TrimSpace(title)

	body = strings.Join(bodyParts, "\n")
	body = strings.TrimSpace(body)

	return
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}

func getMessageFile(about string) (string, error) {
	gitDir, err := git.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(gitDir, fmt.Sprintf("%s_EDITMSG", about)), nil
}
