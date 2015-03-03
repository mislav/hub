package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/hub/Godeps/_workspace/src/github.com/octokit/go-octokit/octokit"
	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdRelease = &Command{
		Run:   release,
		Usage: "release",
		Short: "Retrieve releases from GitHub",
		Long:  `Retrieves releases from GitHub for the project that the "origin" remote points to.`}

	cmdCreateRelease = &Command{
		Key:   "create",
		Run:   createRelease,
		Usage: "release create [-d] [-p] [-a <ASSETS_FILE>] [-m <MESSAGE>|-f <FILE>] <TAG>",
		Short: "Create a new release in GitHub",
		Long: `Creates a new release in GitHub for the project that the "origin" remote points to.
It requires the name of the tag to release as a first argument.

Specify the assets to include in the release via "-a".

Without <MESSAGE> or <FILE>, a text editor will open in which title and body
of the release can be entered in the same manner as git commit message.

If "-d" is given, it creates a draft release.

If "-p" is given, it creates a pre-release.
`}

	flagReleaseDraft,
	flagReleasePrerelease bool

	flagReleaseMessage,
	flagReleaseFile string

	flagReleaseAssets stringSliceValue
)

func init() {
	cmdCreateRelease.Flag.BoolVarP(&flagReleaseDraft, "draft", "d", false, "DRAFT")
	cmdCreateRelease.Flag.BoolVarP(&flagReleasePrerelease, "prerelease", "p", false, "PRERELEASE")
	cmdCreateRelease.Flag.VarP(&flagReleaseAssets, "attach", "a", "ATTACH_ASSETS")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseMessage, "message", "m", "", "MESSAGE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseFile, "file", "f", "", "FILE")

	cmdRelease.Use(cmdCreateRelease)
	CmdRunner.Use(cmdRelease)
}

func release(cmd *Command, args *Args) {
	runInLocalRepo(func(localRepo *github.GitHubRepo, project *github.Project, client *github.Client) {
		if args.Noop {
			ui.Printf("Would request list of releases for %s\n", project)
		} else {
			releases, err := client.Releases(project)
			utils.Check(err)
			var outputs []string
			for _, release := range releases {
				out := fmt.Sprintf("%s (%s)\n%s", release.Name, release.TagName, release.Body)
				outputs = append(outputs, out)
			}

			ui.Println(strings.Join(outputs, "\n\n"))
		}
	})
}

func createRelease(cmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		utils.Check(fmt.Errorf("Missed argument TAG"))
		return
	}

	tag := args.LastParam()
	runInLocalRepo(func(localRepo *github.GitHubRepo, project *github.Project, client *github.Client) {
		release, err := client.Release(project, tag)
		utils.Check(err)

		if release == nil {
			currentBranch, err := localRepo.CurrentBranch()
			utils.Check(err)
			branchName := currentBranch.ShortName()

			title, body, err := getTitleAndBodyFromFlags(flagReleaseMessage, flagReleaseFile)
			utils.Check(err)

			var editor *github.Editor
			if title == "" {
				cs := git.CommentChar()
				message, err := renderReleaseTpl(cs, tag, project.Name, branchName)
				utils.Check(err)

				editor, err = github.NewEditor("RELEASE", "release", message)
				utils.Check(err)

				title, body, err = editor.EditTitleAndBody()
				utils.Check(err)
			}

			params := octokit.ReleaseParams{
				TagName:         tag,
				TargetCommitish: branchName,
				Name:            title,
				Body:            body,
				Draft:           flagReleaseDraft,
				Prerelease:      flagReleasePrerelease,
			}
			release, err = client.CreateRelease(project, params)
			utils.Check(err)

			if editor != nil {
				defer editor.DeleteFile()
			}
		}

		if len(flagReleaseAssets) > 0 {
			paths := make([]string, 0)
			for _, asset := range flagReleaseAssets {
				finder := assetFinder{}
				p, err := finder.Find(asset)
				utils.Check(err)

				paths = append(paths, p...)
			}

			uploader := assetUploader{
				Client:  client,
				Release: release,
			}
			err = uploader.UploadAll(paths)
			if err != nil {
				ui.Println("")
				utils.Check(err)
			}
		}

		ui.Printf("\n%s\n", release.HTMLURL)
	})
}

type assetUploader struct {
	Client  *github.Client
	Release *octokit.Release
}

func (a *assetUploader) UploadAll(paths []string) error {
	errUploadChan := make(chan string)
	successChan := make(chan bool)
	total := len(paths)
	count := 0

	for _, path := range paths {
		go a.uploadAsync(path, successChan, errUploadChan)
	}

	a.printUploadProgress(count, total)

	errUploads := make([]string, 0)
	for {
		select {
		case _ = <-successChan:
			count++
			a.printUploadProgress(count, total)
		case errUpload := <-errUploadChan:
			errUploads = append(errUploads, errUpload)
			count++
			a.printUploadProgress(count, total)
		}

		if count == total {
			break
		}
	}

	var err error
	if len(errUploads) > 0 {
		err = fmt.Errorf("Error uploading %s", strings.Join(errUploads, ", "))
	}

	return err
}

func (a *assetUploader) uploadAsync(path string, successChan chan bool, errUploadChan chan string) {
	err := a.Upload(path)
	if err == nil {
		successChan <- true
	} else {
		errUploadChan <- path
	}
}

func (a *assetUploader) Upload(path string) error {
	contentType, err := a.detectContentType(path)
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	uploadUrl, err := a.Release.UploadURL.Expand(octokit.M{"name": filepath.Base(path)})
	if err != nil {
		return err
	}

	return a.Client.UploadReleaseAsset(uploadUrl, f, contentType)
}

func (a *assetUploader) detectContentType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileHeader := &bytes.Buffer{}
	headerSize := int64(512)
	if fi.Size() < headerSize {
		headerSize = fi.Size()
	}

	// The content type detection only uses 512 bytes at most.
	// This way we avoid copying the whole content for big files.
	_, err = io.CopyN(fileHeader, file, headerSize)
	if err != nil {
		return "", err
	}

	t := http.DetectContentType(fileHeader.Bytes())

	return strings.Split(t, ";")[0], nil
}

func (a *assetUploader) printUploadProgress(count int, total int) {
	out := fmt.Sprintf("Uploading assets (%d/%d)", count, total)
	fmt.Print("\r" + out)
}

type assetFinder struct {
}

func (a *assetFinder) Find(path string) ([]string, error) {
	result := make([]string, 0)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			result = append(result, path)
		}

		return nil
	})

	return result, err
}
