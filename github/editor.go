package github

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/github/hub/v2/cmd"
	"github.com/github/hub/v2/git"
	"github.com/kballard/go-shellquote"
)

const Scissors = "------------------------ >8 ------------------------"

func NewEditor(filename, topic, message string) (editor *Editor, err error) {
	gitDir, err := git.Dir()
	if err != nil {
		return
	}
	messageFile := filepath.Join(gitDir, filename)

	program, err := git.Editor()
	if err != nil {
		return
	}

	cs, err := git.CommentChar(message)
	if err != nil {
		return
	}

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
	Program           string
	Topic             string
	File              string
	Message           string
	CS                string
	addedFirstComment bool
	openEditor        func(program, file string) error
}

func (e *Editor) AddCommentedSection(text string) {
	if !e.addedFirstComment {
		scissors := e.CS + " " + Scissors + "\n"
		scissors += e.CS + " Do not modify or remove the line above.\n"
		scissors += e.CS + " Everything below it will be ignored.\n"
		e.Message = e.Message + "\n" + scissors
		e.addedFirstComment = true
	}

	e.Message = e.Message + "\n" + text
}

func (e *Editor) DeleteFile() error {
	return os.Remove(e.File)
}

func (e *Editor) EditContent() (content string, err error) {
	b, err := e.openAndEdit()
	if err != nil {
		return
	}

	b = bytes.TrimSpace(b)
	reader := bytes.NewReader(b)
	scanner := bufio.NewScanner(reader)
	unquotedLines := []string{}

	scissorsLine := e.CS + " " + Scissors
	for scanner.Scan() {
		line := scanner.Text()
		if line == scissorsLine {
			break
		}
		unquotedLines = append(unquotedLines, line)
	}
	if err = scanner.Err(); err != nil {
		return
	}

	content = strings.Join(unquotedLines, "\n")
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
	if !e.isFileExist() {
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
	programArgs, err := shellquote.Split(program)
	if err != nil {
		return err
	}
	editCmd := cmd.NewWithArray(programArgs)
	r := regexp.MustCompile(`\b(?:[gm]?vim)(?:\.exe)?$`)
	if r.MatchString(editCmd.Name) {
		editCmd.WithArg("--cmd")
		editCmd.WithArg("set ft=gitcommit tw=0 wrap lbr")
	}
	editCmd.WithArg(file)
	// Reattach stdin to the console before opening the editor
	setConsole(editCmd)

	return editCmd.Spawn()
}
