package git

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ProtocolRe     = regexp.MustCompile("^[a-zA-Z_-]+://")
	SshConfigFiles = []string{
		filepath.Join(os.Getenv("HOME"), ".ssh/config"),
		"/etc/ssh_config",
		"/etc/ssh/ssh_config",
	}
	SshConfig map[string]string
)

func ParseURL(rawurl string) (u *url.URL, err error) {
	if !ProtocolRe.MatchString(rawurl) && strings.Contains(rawurl, ":") {
		rawurl = "ssh://" + strings.Replace(rawurl, ":", "/", 1)
	}

	u, err = url.Parse(rawurl)
	if err == nil {
		if SshConfig == nil {
			SshConfig = readSshConfig()
		}
		if SshConfig[u.Host] != "" {
			u.Host = SshConfig[u.Host]
		}
	}
	return
}

func readSshConfig() map[string]string {
	config := make(map[string]string)
	hostRe := regexp.MustCompile("^[ \t]*(Host|HostName)[ \t]+(.+)$")

	for _, filename := range SshConfigFiles {
		file, err := os.Open(filename)
		if err != nil {
			continue
		}
		hosts := []string{"*"}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			match := hostRe.FindStringSubmatch(line)
			if match == nil {
				continue
			}
			names := strings.Fields(match[2])
			if match[1] == "Host" {
				hosts = names
			} else {
				for _, host := range hosts {
					for _, name := range names {
						config[host] = name
					}
				}
			}
		}
		file.Close()
	}

	return config
}
