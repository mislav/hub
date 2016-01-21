// Package xdgbasedir contains helper functions for data/config/cache directory and file lookup
package xdgbasedir

// Pulled from http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html

import (
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	defaultDataHomeDirectory = "/.local/share"
	defaultConfigDirectory   = "/.config"
	defaultCacheDirectory    = "/.cache"
)

var (
	// Overridable for testing
	osGetEnv     = os.Getenv
	userCurrent  = user.Current
	osStat       = os.Stat
	osIsNotExist = os.IsNotExist
)

// DataHomeDirectory is a single base directory relative to which user-specific data files should be written
func DataHomeDirectory() (string, error) {
	return getInEnvOrJoinWithHome("XDG_DATA_HOME", defaultDataHomeDirectory)
}

// GetDataFileLocation returns the location of a data file
func GetDataFileLocation(filename string) (string, error) {
	return fileLocationRetrievalHelper(filename, dataDirectories(), DataHomeDirectory)
}

func dataDirectories() []string {
	return uniqueDirsOnVariable("XDG_DATA_DIRS", "/usr/local/share/:/usr/share/")
}

// ConfigHomeDirectory is a single base directory relative to which user-specific config files should be written
func ConfigHomeDirectory() (string, error) {
	return getInEnvOrJoinWithHome("XDG_CONFIG_HOME", defaultConfigDirectory)
}

// GetConfigFileLocation returns the location of a config file
func GetConfigFileLocation(filename string) (string, error) {
	return fileLocationRetrievalHelper(filename, configDirectories(), ConfigHomeDirectory)
}

func configDirectories() []string {
	return uniqueDirsOnVariable("XDG_CONFIG_DIRS", "/etc/xdg")
}

// CacheDirectory is a single base directory relative to which user specific non-essential data files should be stored.
func CacheDirectory() (string, error) {
	return getInEnvOrJoinWithHome("XDG_CACHE_HOME", defaultCacheDirectory)
}

// GetCacheFileLocation returns the location of the cache file
func GetCacheFileLocation(filename string) (string, error) {
	return fileLocationRetrievalHelper(filename, []string{}, CacheDirectory)
}

func fileLocationRetrievalHelper(filename string, dirs []string, defaultDirectoryFunc func() (string, error)) (string, error) {
	// The default location should be checked first
	defaultDir, defaultDirectoryError := defaultDirectoryFunc()
	if defaultDirectoryError != nil {
		dirs = append([]string{defaultDir}, dirs...)
	}

	for _, dir := range dirs {
		fileLoc := path.Join(dir, filename)
		fileInfo, err := osStat(fileLoc)
		if err != nil {
			continue
		}

		return fileInfo.Name(), nil
	}
	return path.Join(defaultDir, filename), defaultDirectoryError
}

func uniqueDirsOnVariable(envVar string, defaultVal string) []string {
	dataDirs := osGetEnv(envVar)
	if dataDirs == "" {
		dataDirs = defaultVal
	}
	return splitAndReturnUnique(dataDirs, ":")
}

func splitAndReturnUnique(str string, sep string) []string {
	parts := strings.Split(str, sep)
	var ret []string
	usedMap := make(map[string]bool)
	for _, p := range parts {
		_, found := usedMap[p]
		if !found {
			usedMap[p] = true
			ret = append(ret, p)
		}
	}
	return ret
}

func getInEnvOrJoinWithHome(envName string, directory string) (string, error) {
	configHome := osGetEnv(envName)
	if configHome != "" {
		return configHome, nil
	}
	return joinWithHome(directory)
}

func joinWithHome(dir string) (string, error) {
        homeDir := os.Getenv("HOME")

	if homeDir == "" {
		usr, err := userCurrent()
		if err != nil {
		   return "", err
		}
		homeDir = usr.HomeDir
	}

	return path.Join(homeDir, dir), nil
}
