package github

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"os"
	"reflect"
	"runtime"
	"strings"
)

const (
	ghReportCrashConfig = "gh.reportCrash"
	ghProjectOwner      = "jingweno"
	ghProjectName       = "gh"
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

		always := isOption(confirm, "a", "always")
		if always || isOption(confirm, "y", "yes") {
			report(err, stack)
		}

		saveReportConfiguration(confirm, always)
	}
	os.Exit(1)
}

func isOption(confirm, short, long string) bool {
	return strings.EqualFold(confirm, short) || strings.EqualFold(confirm, long)
}

func report(reportedError error, stack string) {
	title, body, err := reportTitleAndBody(reportedError, stack)
	utils.Check(err)

	project := NewProject(ghProjectOwner, ghProjectName, GitHubHost)

	gh := NewClient(project.Host)

	issue, err := gh.CreateIssue(project, title, body, []string{"Crash Report"})
	utils.Check(err)

	fmt.Println(issue.HTMLURL)
}

func reportTitleAndBody(reportedError error, stack string) (title, body string, err error) {
	message := "Crash report - %v\n\nError (%s): `%v`\n\nStack:\n\n```\n%s\n```\n\nRuntime:\n\n```\n%s\n```\n\n"
	message += `
# Creating crash report:
#
# This information will be posted as a new issue under jingweno/gh.
# We're NOT including any information about the command that you were executing,
# but knowing a little bit more about it would really help us to solve this problem.
# Feel free to modify the title and the description for this issue.
`

	errType := reflect.TypeOf(reportedError).String()
	message = fmt.Sprintf(message, reportedError, errType, reportedError, stack, runtimeInfo())

	editor, err := NewEditor("CRASH_REPORT", message)
	if err != nil {
		return "", "", err
	}

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
	fmt.Printf("%v\n\n", err)
	fmt.Println(stack)
}

func saveReportConfiguration(confirm string, always bool) {
	if always {
		git.SetGlobalConfig(ghReportCrashConfig, "always")
	} else if isOption(confirm, "e", "never") {
		git.SetGlobalConfig(ghReportCrashConfig, "never")
	}
}

func reportCrashConfig() (opt string) {
	opt = os.Getenv("GH_REPORT_CRASH")
	if opt == "" {
		opt, _ = git.GlobalConfig(ghReportCrashConfig)
	}

	return
}
