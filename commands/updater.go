package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
	goupdate "github.com/inconshreveable/go-update"
)

const ghAutoUpdateConfig = "gh.autoUpdate"

func NewUpdater() *Updater {
	version := os.Getenv("GH_VERSION")
	if version == "" {
		version = Version
	}

	timestampPath := filepath.Join(os.Getenv("HOME"), ".config", "gh-update")
	return &Updater{
		Host:           github.DefaultHost(),
		CurrentVersion: version,
		timestampPath:  timestampPath,
	}
}

type Updater struct {
	Host           string
	CurrentVersion string
	timestampPath  string
}

func (updater *Updater) timeToUpdate() bool {
	if updater.CurrentVersion == "dev" || readTime(updater.timestampPath).After(time.Now()) {
		return false
	}

	// the next update is in about 14 days
	wait := 13*24*time.Hour + randDuration(24*time.Hour)
	return writeTime(updater.timestampPath, time.Now().Add(wait))
}

func (updater *Updater) PromptForUpdate() (err error) {
	if !updater.timeToUpdate() {
		return
	}

	releaseName, version := updater.latestReleaseNameAndVersion()
	if version != "" && version != updater.CurrentVersion {
		switch autoUpdateConfig() {
		case "always":
			err = updater.updateTo(releaseName, version)
		case "never":
			return
		default:
			fmt.Println("There is a newer version of gh available.")
			fmt.Print("Would you like to update? ([Y]es/[N]o/[A]lways/N[e]ver): ")
			var confirm string
			fmt.Scan(&confirm)

			always := utils.IsOption(confirm, "a", "always")
			if always || utils.IsOption(confirm, "y", "yes") {
				err = updater.updateTo(releaseName, version)
			}

			saveAutoUpdateConfiguration(confirm, always)
		}
	}

	return
}

func (updater *Updater) Update() (err error) {
	releaseName, version := updater.latestReleaseNameAndVersion()
	if version == "" {
		fmt.Println("There is no newer version of gh available.")
		return
	}

	if version == updater.CurrentVersion {
		fmt.Printf("You're already on the latest version: %s\n", version)
	} else {
		err = updater.updateTo(releaseName, version)
	}

	return
}

func (updater *Updater) latestReleaseNameAndVersion() (name, version string) {
	// Create Client with a stub Credential
	c := github.Client{Credential: &github.Credential{Host: updater.Host}}
	name, _ = c.GhLatestTagName()
	version = strings.TrimPrefix(name, "v")

	return
}

func (updater *Updater) updateTo(releaseName, version string) (err error) {
	fmt.Printf("Updating gh to %s...\n", version)
	downloadURL := fmt.Sprintf("https://%s/jingweno/gh/releases/download/%s/gh_%s_%s_%s.zip", updater.Host, releaseName, version, runtime.GOOS, runtime.GOARCH)
	path, err := downloadFile(downloadURL)
	if err != nil {
		return
	}

	exec, err := unzipExecutable(path)
	if err != nil {
		return
	}

	err, _ = goupdate.FromFile(exec)
	if err == nil {
		fmt.Println("Done!")
	}

	return
}

func unzipExecutable(path string) (exec string, err error) {
	rc, err := zip.OpenReader(path)
	if err != nil {
		err = fmt.Errorf("Can't open zip file %s: %s", path, err)
		return
	}
	defer rc.Close()

	for _, file := range rc.File {
		if !strings.HasPrefix(file.Name, "gh") {
			continue
		}

		dir := filepath.Dir(path)
		exec, err = unzipFile(file, dir)
		break
	}

	if exec == "" && err == nil {
		err = fmt.Errorf("No gh executable is found in %s", path)
	}

	return
}

func unzipFile(file *zip.File, to string) (exec string, err error) {
	frc, err := file.Open()
	if err != nil {
		err = fmt.Errorf("Can't open zip entry %s when reading: %s", file.Name, err)
		return
	}
	defer frc.Close()

	dest := filepath.Join(to, filepath.Base(file.Name))
	f, err := os.Create(dest)
	if err != nil {
		return
	}
	defer f.Close()

	copied, err := io.Copy(f, frc)
	if err != nil {
		return
	}

	if uint32(copied) != file.UncompressedSize {
		err = fmt.Errorf("Zip entry %s is corrupted", file.Name)
		return
	}

	exec = f.Name()

	return
}

func downloadFile(url string) (path string, err error) {
	dir, err := ioutil.TempDir("", "gh-update")
	if err != nil {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		err = fmt.Errorf("Can't download %s: %d", url, resp.StatusCode)
		return
	}

	file, err := os.Create(filepath.Join(dir, filepath.Base(url)))
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}

	path = file.Name()

	return
}

func randDuration(n time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(n)))
}

func readTime(path string) time.Time {
	p, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return time.Time{}
	}
	if err != nil {
		return time.Now().Add(1000 * time.Hour)
	}

	t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(p)))
	if err != nil {
		return time.Time{}
	}

	return t
}

func writeTime(path string, t time.Time) bool {
	return ioutil.WriteFile(path, []byte(t.Format(time.RFC3339)), 0644) == nil
}

func saveAutoUpdateConfiguration(confirm string, always bool) {
	if always {
		git.SetGlobalConfig(ghAutoUpdateConfig, "always")
	} else if utils.IsOption(confirm, "e", "never") {
		git.SetGlobalConfig(ghAutoUpdateConfig, "never")
	}
}

func autoUpdateConfig() (opt string) {
	opt = os.Getenv("GH_AUTOUPDATE")
	if opt == "" {
		opt, _ = git.GlobalConfig(ghAutoUpdateConfig)
	}

	return
}
