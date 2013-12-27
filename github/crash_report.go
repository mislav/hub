package github

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func ReportCrash(err error) {
	if err != nil {
		config := CurrentConfigs()

		buf := make([]byte, 10000)
		runtime.Stack(buf, false)
		stack := formatStack(buf)

		switch config.ReportCrash {
		case "always":
			report(err, stack, config)
		case "never":
			printError(err, stack)
		default:
			printError(err, stack)
			fmt.Print("Would you like to open an issue? ([Y]es/[N]o/[A]lways/N[e]ver): ")
			var confirm string
			fmt.Scan(&confirm)

			always := isOption(confirm, "a", "always")
			if always || isOption(confirm, "y", "yes") {
				report(err, stack, config)
			}

			saveReportConfiguration(config, confirm, always)
		}
		os.Exit(1)
	}
}

func isOption(confirm, short, long string) bool {
	return strings.EqualFold(confirm, short) || strings.EqualFold(confirm, long)
}

func report(err error, stack string, config *Configs) {
	message := "Crash report - %v\nError: %v\nStack:\n```\n%s\n```\n"
	message += `
# Creating crash report:
#
# Now it's the time to be specific. We're not including information about
# the command that you were executing, but knowing a little bit more about it
# would really help us to solve this problem. Feel free to modify the title
# and the description with the error and the stack trace.
`
	message = fmt.Sprintf(message, time.Now(), err, stack)

	GetTitleAndBodyFromEditor("CRASHREPORT", message)
}

func formatStack(buf []byte) string {
	stack := strings.Split(string(buf), "\n")

	stack = append(stack[0:1], stack[3:]...)
	return strings.Join(stack, "\n")
}

func printError(err error, stack string) {
	fmt.Printf("%v\n\n", err)
	fmt.Println(stack)
}

func saveReportConfiguration(config *Configs, confirm string, always bool) {
	if always {
		config.ReportCrash = "always"
		config.Save()
	} else if isOption(confirm, "e", "never") {
		config.ReportCrash = "never"
		config.Save()
	}
}
