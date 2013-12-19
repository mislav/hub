package commands

import (
	"archive/zip"
	"fmt"
	goupdate "github.com/inconshreveable/go-update"
	"github.com/jingweno/gh/github"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	updateTimestampPath = filepath.Join(os.Getenv("HOME"), ".config", "gh-update")
)

func NewUpdater() *Updater {
	return &Updater{Host: github.GitHubHost, CurrentVersion: Version}
}

type Updater struct {
	Host           string
	CurrentVersion string
}

func (update *Updater) WantUpdate() bool {
	if update.CurrentVersion == "dev" || readTime(updateTimestampPath).After(time.Now()) {
		return false
	}

	wait := 12*time.Hour + randDuration(8*time.Hour)
	return writeTime(updateTimestampPath, time.Now().Add(wait))
}

func (updater *Updater) Update() (err error) {
	client := github.NewClient(updater.Host)
	releases, err := client.Releases(github.NewProject("jingweno", "gh", updater.Host))
	if err != nil {
		err = fmt.Errorf("Error fetching releases: %s", err)
		return
	}

	latestRelease := releases[0]
	tagName := latestRelease.TagName
	version := strings.TrimPrefix(tagName, "v")
	if version == updater.CurrentVersion {
		fmt.Printf("You're already on the latest version: %s\n", updater.CurrentVersion)
		return
	}

	fmt.Printf("Updating gh to %s...\n", version)
	downloadURL := fmt.Sprintf("https://%s/jingweno/gh/releases/download/%s/gh_%s_%s_%s.zip", updater.Host, tagName, version, runtime.GOOS, runtime.GOARCH)
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
	t, err := time.Parse(time.RFC3339, string(p))
	if err != nil {
		return time.Now().Add(1000 * time.Hour)
	}
	return t
}

func writeTime(path string, t time.Time) bool {
	return ioutil.WriteFile(path, []byte(t.Format(time.RFC3339)), 0644) == nil
}
