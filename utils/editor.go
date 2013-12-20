package utils

import (
	"github.com/jingweno/gh/git"
	"io/ioutil"
	"os"
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

func GetTitleAndBodyFromEditor(fn func(messageFile string) error) (title, body string, err error) {
	messageFile, err := git.PullReqMsgFile()
	if err != nil {
		return
	}
	defer os.Remove(messageFile)

	if fn != nil {
		err = fn(messageFile)
		if err != nil {
			return
		}
	}

	editor, err := git.Editor()
	if err != nil {
		return
	}

	err = editTitleAndBody(editor, messageFile)
	if err != nil {
		err = fmt.Errorf("error using text editor for release message")
		return
	}

	title, body, err = readTitleAndBody(messageFile)
	if err != nil {
		return
	}

	return
}

func editTitleAndBody(editor, messageFile string) error {
	editCmd := cmd.New(editor)
	r := regexp.MustCompile("[mg]?vi[m]$")
	if r.MatchString(editor) {
		editCmd.WithArg("-c")
		editCmd.WithArg("set ft=gitcommit tw=0 wrap lbr")
	}
	editCmd.WithArg(messageFile)

	return editCmd.Exec()
}

func readTitleAndBody(messageFile string) (title, body string, err error) {
	f, err := os.Open(messageFile)
	defer f.Close()
	if err != nil {
		return "", "", err
	}

	reader := bufio.NewReader(f)

	return readTitleAndBodyFrom(reader)
}

func readTitleAndBodyFrom(reader *bufio.Reader) (title, body string, err error) {
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
