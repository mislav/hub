package github

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
	"github.com/github/hub/version"
)

const (
	hubReportCrashConfig = "hub.reportCrash"
	hubProjectOwner      = "github"
	hubProjectName       = "hub"
)

func CaptureCrash() {
	if rec := recover(); rec != nil {
		if err, ok := rec.(error); ok {
			reportCrash(err)
		} else if err, ok := rec.(string); ok {
			reportCrash(errors.New(err))
		}
	}
}

func reportCrash(err error) {
	if err == nil {
		return
	}

	buf := make([]byte, 10000)
	runtime.Stack(buf, false)
	stack := formatStack(buf)

	switch reportCrashConfig() {
	case "always":
		report(err, stack)
	case "never":
		printError(err, stack)
	default:
		printError(err, stack)
		fmt.Print("Would you like to open an issue? ([Y]es/[N]o/[A]lways/N[e]ver): ")
		var confirm string
		fmt.Scan(&confirm)

		always := utils.IsOption(confirm, "a", "always")
		if always || utils.IsOption(confirm, "y", "yes") {
			report(err, stack)
		}

		saveReportConfiguration(confirm, always)
	}
	os.Exit(1)
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

	ui.Println(issue.HtmlUrl)
}

const crashReportTmpl = "Crash report - %v\n\n" +
	"Error (%s): `%v`\n\n" +
	"Stack:\n\n```\n%s\n```\n\n" +
	"Runtime:\n\n```\n%s\n```\n\n" +
	"Version:\n\n```\n%s\n```\n" +
	`
# Creating crash report:
#
# This information will be posted as a new issue under github/hub.
# We're NOT including any information about the command that you were executing,
# but knowing a little bit more about it would really help us to solve this problem.
# Feel free to modify the title and the description for this issue.
`

func reportTitleAndBody(reportedError error, stack string) (title, body string, err error) {
	errType := reflect.TypeOf(reportedError).String()
	fullVersion, _ := version.FullVersion()
	message := fmt.Sprintf(
		crashReportTmpl,
		reportedError,
		errType,
		reportedError,
		stack,
		runtimeInfo(),
		fullVersion,
	)

	editor, err := NewEditor("CRASH_REPORT", "crash report", message)
	if err != nil {
		return "", "", err
	}

	defer editor.DeleteFile()

	return editor.EditTitleAndBody()
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

func printError(err error, stack string) {
	ui.Printf("%v\n\n", err)
	ui.Println(stack)
}

func saveReportConfiguration(confirm string, always bool) {
	if always {
		git.SetGlobalConfig(hubReportCrashConfig, "always")
	} else if utils.IsOption(confirm, "e", "never") {
		git.SetGlobalConfig(hubReportCrashConfig, "never")
	}
}

func reportCrashConfig() (opt string) {
	opt = os.Getenv("HUB_REPORT_CRASH")
	if opt == "" {
		opt, _ = git.GlobalConfig(hubReportCrashConfig)
	}

	return
}
