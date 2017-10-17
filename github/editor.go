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

	cs := git.CommentChar()

	editor = &Editor{
		Program:    program,
		Topic:      topic,
		File:       messageFile,
		Message:    message,
		CS:         cs,
		openEditor: openTextEditor,
	}

	return
}

type Editor struct {
	Program    string
	Topic      string
	File       string
	Message    string
	CS         string
	openEditor func(program, file string) error
}

func (e *Editor) DeleteFile() error {
	return os.Remove(e.File)
}

func (e *Editor) EditTitleAndBody() (title, body string, err error) {
	content, err := e.openAndEdit()
	if err != nil {
		return
	}

	content = bytes.TrimSpace(content)
	reader := bytes.NewReader(content)
	title, body, err = readTitleAndBody(reader, e.CS)

	if err != nil || title == "" {
		defer e.DeleteFile()
	}

	return
}

func (e *Editor) openAndEdit() (content []byte, err error) {
	err = e.writeContent()
	if err != nil {
		return
	}

	err = e.openEditor(e.Program, e.File)
	if err != nil {
		err = fmt.Errorf("error using text editor for %s message", e.Topic)
		defer e.DeleteFile()
		return
	}

	content, err = e.readContent()

	return
}

func (e *Editor) writeContent() (err error) {
	// only write message if file doesn't exist
	if !e.isFileExist() && e.Message != "" {
		err = ioutil.WriteFile(e.File, []byte(e.Message), 0644)
		if err != nil {
			return
		}
	}

	return
}

func (e *Editor) isFileExist() bool {
	_, err := os.Stat(e.File)
	return err == nil || !os.IsNotExist(err)
}

func (e *Editor) readContent() (content []byte, err error) {
	return ioutil.ReadFile(e.File)
}

func openTextEditor(program, file string) error {
	editCmd := cmd.New(program)
	r := regexp.MustCompile(`\b(?:[gm]?vim|vi)(?:\.exe)?$`)
	if r.MatchString(editCmd.Name) {
		editCmd.WithArg("--cmd")
		editCmd.WithArg("set ft=gitcommit tw=0 wrap lbr")
	}
	editCmd.WithArg(file)
	// Reattach stdin to the console before opening the editor
	setConsole(editCmd)

	return editCmd.Spawn()
}

func readTitleAndBody(reader io.Reader, cs string) (title, body string, err error) {
	var titleParts, bodyParts []string

	r := regexp.MustCompile("\\S")
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, cs) {
			continue
		}

		if len(bodyParts) == 0 && r.MatchString(line) {
			titleParts = append(titleParts, line)
		} else {
			bodyParts = append(bodyParts, line)
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	title = strings.Join(titleParts, " ")
	title = strings.TrimSpace(title)

	body = strings.Join(bodyParts, "\n")
	body = strings.TrimSpace(body)

	return
}

func getMessageFile(about string) (string, error) {
	gitDir, err := git.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(gitDir, fmt.Sprintf("%s_EDITMSG", about)), nil
}
