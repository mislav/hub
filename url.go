package main

import (
	"regexp"
)

func mustMatchGitUrl(url string) []string {
	httpRegex := regexp.MustCompile("https://.+/(.+)/(.+).git")
	if httpRegex.MatchString(url) {
		return httpRegex.FindStringSubmatch(url)
	}

	readOnlyRegex := regexp.MustCompile("git://.+/(.+)/(.+).git")
	if readOnlyRegex.MatchString(url) {
		return readOnlyRegex.FindStringSubmatch(url)
	}

	sshRegex := regexp.MustCompile(".+:(.+)/(.+).git")
	if sshRegex.MatchString(url) {
		return sshRegex.FindStringSubmatch(url)
	}

	panic("Can't find owner")
}
