package github

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func ReportCrash(err error) {
	if err != nil {
		config := CurrentConfigs()

		buf := make([]byte, 10000)
		runtime.Stack(buf, false)
		stack := string(buf)

		switch config.ReportCrash {
		case "always":
			report(err, stack, config)
		case "never":
			fmt.Printf("fatal: %v\n", err)
			fmt.Println(stack)
		default:
			fmt.Printf("fatal: %v\n", err)
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
