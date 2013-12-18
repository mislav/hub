package commands

import (
	"archive/zip"
	"fmt"
	updater "github.com/inconshreveable/go-update"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var cmdUpdate = &Command{
	Run:   update,
	Usage: "update",
	Short: "Update gh",
	Long: `Update gh to the latest version.

Examples:
  git update
`,
}

func update(cmd *Command, args *Args) {
	err := doUpdate()
	utils.Check(err)
	os.Exit(0)
}

func doUpdate() (err error) {
	client := github.NewClient(github.GitHubHost)
	releases, err := client.Releases(github.NewProject("jingweno", "gh", github.GitHubHost))
	if err != nil {
		err = fmt.Errorf("Error fetching releases: %s", err)
		return
	}

	latestRelease := releases[0]
	tagName := latestRelease.TagName
	version := strings.TrimPrefix(tagName, "v")
	if version == Version {
		err = fmt.Errorf("You're already on the latest version: %s", Version)
		return
	}

	fmt.Printf("Updating gh to release %s...\n", version)
	downloadURL := fmt.Sprintf("https://github.com/jingweno/gh/releases/download/%s/gh_%s-snapshot_%s_%s.zip", tagName, version, runtime.GOOS, runtime.GOARCH)
	path, err := downloadFile(downloadURL)
	if err != nil {
		err = fmt.Errorf("Can't download update from %s to %s", downloadURL, path)
		return
	}

	exec, err := unzipExecutable(path)
	if err != nil {
		err = fmt.Errorf("Can't unzip gh executable: %s", err)
		return
	}

	err, _ = updater.FromFile(exec)
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

	file, err := os.Create(filepath.Join(dir, filepath.Base(url)))
	if err != nil {
		return
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}

	path = file.Name()

	return
}
