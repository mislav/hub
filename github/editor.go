package github

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jingweno/gh/cmd"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetTitleAndBodyFromFlags(messageFlag, fileFlag string) (title, body string, err error) {
	if messageFlag != "" {
		title, body = readMsg(messageFlag)
	} else if fileFlag != "" {
		var (
			content []byte
			err     error
		)

		if fileFlag == "-" {
			content, err = ioutil.ReadAll(os.Stdin)
		} else {
			content, err = ioutil.ReadFile(fileFlag)
		}
		utils.Check(err)
		title, body = readMsg(string(content))
	}

	return
}

func NewEditor(topic, message string) (editor *Editor, err error) {
	messageFile, err := getMessageFile(topic)
	if err != nil {
		return
	}

	program, err := git.Editor()
	if err != nil {
		return
	}

	editor = &Editor{
		Program: program,
		File:    messageFile,
		Message: message,
		doEdit:  doTextEditorEdit,
	}

	return
}

type Editor struct {
	Program string
	File    string
	Message string
	doEdit  func(program, file string) error
}

func (e *Editor) Edit() (content []byte, err error) {
	if e.Message != "" {
		err = ioutil.WriteFile(e.File, []byte(e.Message), 0644)
		if err != nil {
			return
		}
	}
	defer os.Remove(e.File)

	err = e.doEdit(e.Program, e.File)
	if err != nil {
		err = fmt.Errorf("error using text editor for editing message")
		return
	}

	content, err = ioutil.ReadFile(e.File)

	return
}

func (e *Editor) EditTitleAndBody() (title, body string, err error) {
	content, err := e.Edit()
	if err != nil {
		return
	}

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

func readMsg(msg string) (title, body string) {
	split := strings.SplitN(msg, "\n\n", 2)
	title = strings.TrimSpace(split[0])
	if len(split) > 1 {
		body = strings.TrimSpace(split[1])
	}

	return
}

func getMessageFile(about string) (string, error) {
	gitDir, err := git.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(gitDir, fmt.Sprintf("%s_EDITMSG", about)), nil
}
