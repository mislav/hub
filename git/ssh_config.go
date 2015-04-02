package git

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	hostReStr = "(?i)^[ \t]*(host|hostname)[ \t]+(.+)$"
)

type SSHConfig map[string]string

func newSSHConfigReader() *SSHConfigReader {
	return &SSHConfigReader{
		Files: []string{
			filepath.Join(os.Getenv("HOME"), ".ssh/config"),
			"/etc/ssh_config",
			"/etc/ssh/ssh_config",
		},
	}
}

type SSHConfigReader struct {
	Files []string
}

func (r *SSHConfigReader) Read() SSHConfig {
	config := make(SSHConfig)
	hostRe := regexp.MustCompile(hostReStr)

	for _, filename := range r.Files {
		r.readFile(config, hostRe, filename)
	}

	return config
}

func (r *SSHConfigReader) readFile(c SSHConfig, re *regexp.Regexp, f string) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()

	hosts := []string{"*"}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		names := strings.Fields(match[2])
		if strings.EqualFold(match[1], "host") {
			hosts = names
		} else {
			for _, host := range hosts {
				for _, name := range names {
					c[host] = name
				}
			}
		}
	}

	return scanner.Err()
}
