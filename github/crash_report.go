package github

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
	"github.com/github/hub/v2/version"
)

const (
	hubReportCrashConfig = "hub.reportCrash"
	hubProjectOwner      = "github"
	hubProjectName       = "hub"
)

func CaptureCrash() {
	if rec := recover(); rec != nil {
		switch err := rec.(type) {
		case error:
			reportCrash(err)
		case string:
			reportCrash(errors.New(err))
		default:
			return
		}
		os.Exit(1)
	}
}

func reportCrash(err error) {
	buf := make([]byte, 10000)
	runtime.Stack(buf, false)
	stack := formatStack(buf)

	ui.Errorf("%v\n\n", err)
	ui.Errorln(stack)

	isTerm := ui.IsTerminal(os.Stdin) && ui.IsTerminal(os.Stdout)
	if !isTerm || reportCrashConfig() == "never" {
		return
	}

	ui.Print("Would you like to open an issue? ([y]es / [N]o / n[e]ver): ")
	var confirm string
	prompt := bufio.NewScanner(os.Stdin)
	if prompt.Scan() {
		confirm = prompt.Text()
	}
	if prompt.Err() != nil {
		return
	}

	if isOption(confirm, "y", "yes") {
		report(err, stack)
	} else if isOption(confirm, "e", "never") {
		git.SetGlobalConfig(hubReportCrashConfig, "never")
	}
}

func isOption(confirm, short, long string) bool {
	return strings.EqualFold(confirm, short) || strings.EqualFold(confirm, long)
}

func report(reportedError error, stack string) {
	title, body, err := reportTitleAndBody(reportedError, stack)
	utils.Check(err)

	project := NewProject(hubProjectOwner, hubProjectName, GitHubHost)

	gh := NewClient(project.Host)

	params := map[string]interface{}{
		"title":  title,
		"body":   body,
		"labels": []string{"Crash Report"},
	}

	issue, err := gh.CreateIssue(project, params)
	utils.Check(err)

	ui.Println(issue.HTMLURL)
}

const crashReportTmpl = "Crash report - %v\n\n" +
	"Error (%s): `%v`\n\n" +
	"Stack:\n\n```\n%s\n```\n\n" +
	"Runtime:\n\n```\n%s\n```\n\n" +
	"Version:\n\n```\n%s\nhub version %s\n```\n"

func reportTitleAndBody(reportedError error, stack string) (title, body string, err error) {
	errType := reflect.TypeOf(reportedError).String()
	gitVersion, gitErr := git.Version()
	if gitErr != nil {
		gitVersion = "git unavailable!"
	}
	message := fmt.Sprintf(
		crashReportTmpl,
		reportedError,
		errType,
		reportedError,
		stack,
		runtimeInfo(),
		gitVersion,
		version.Version,
	)

	messageBuilder := &MessageBuilder{
		Filename: "CRASH_REPORT",
		Title:    "crash report",
		Message:  message,
		Edit:     true,
	}
	messageBuilder.AddCommentedSection(`Creating crash report:

This information will be posted as a new issue under github/hub.
We're NOT including any information about the command that you were executing,
but knowing a little bit more about it would really help us to solve this problem.
Feel free to modify the title and the description for this issue.`)

	title, body, err = messageBuilder.Extract()
	if err != nil {
		return
	}
	defer messageBuilder.Cleanup()

	return
}

func runtimeInfo() string {
	return fmt.Sprintf("GOOS: %s\nGOARCH: %s", runtime.GOOS, runtime.GOARCH)
}

func formatStack(buf []byte) string {
	buf = bytes.Trim(buf, "\x00")

	stack := strings.Split(string(buf), "\n")
	stack = append(stack[0:1], stack[5:]...)

	return strings.Join(stack, "\n")
}

func reportCrashConfig() (opt string) {
	opt = os.Getenv("HUB_REPORT_CRASH")
	if opt == "" {
		opt, _ = git.GlobalConfig(hubReportCrashConfig)
	}

	return
}
