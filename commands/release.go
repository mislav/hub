package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
	"github.com/octokit/go-octokit/octokit"
)

var (
	cmdRelease = &Command{
		Run:   release,
		Usage: "release",
		Long:  "Retrieve releases from GitHub",
	}

	cmdCreateRelease = &Command{
		Key:   "create",
		Run:   createRelease,
		Usage: "release create [-d] [-p] [-a <FILE>] [-m <MESSAGE>|-f <FILE>] [-c <COMMIT>] <TAG>",
		Long: `Create a GitHub release.

## Options:
	-d, --draft
		Create a draft release.

	-p, --prerelease
		Create a pre-release.

	-a, --asset <FILE>
		Attach a file as an asset for this release.

	-m, --message <MESSAGE>
		Use the first line of <MESSAGE> as release title, and the rest as release description.

	-f, --file <FILE>
		Read the release title and description from <FILE>.
	
	-c, --commitish <COMMIT>
		A SHA, tag, or branch name to attach the release to (default: current branch).

	<TAG>
		The git tag name for this release.
`}

	flagReleaseDraft,
	flagReleasePrerelease bool

	flagReleaseMessage,
	flagReleaseFile,
	flagReleaseCommitish string

	flagReleaseAssets stringSliceValue
)

func init() {
	cmdCreateRelease.Flag.BoolVarP(&flagReleaseDraft, "draft", "d", false, "DRAFT")
	cmdCreateRelease.Flag.BoolVarP(&flagReleasePrerelease, "prerelease", "p", false, "PRERELEASE")
	cmdCreateRelease.Flag.VarP(&flagReleaseAssets, "attach", "a", "ATTACH_ASSETS")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseMessage, "message", "m", "", "MESSAGE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseFile, "file", "f", "", "FILE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseCommitish, "commitish", "c", "", "COMMITISH")

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
			commitish := flagReleaseCommitish
			if commitish == "" {
				currentBranch, err := localRepo.CurrentBranch()
				utils.Check(err)
				commitish = currentBranch.ShortName()
			}

			title, body, err := getTitleAndBodyFromFlags(flagReleaseMessage, flagReleaseFile)
			utils.Check(err)

			var editor *github.Editor
			if title == "" {
				cs := git.CommentChar()
				message, err := renderReleaseTpl(cs, tag, project.Name, commitish)
				utils.Check(err)

				editor, err = github.NewEditor("RELEASE", "release", message)
				utils.Check(err)

				title, body, err = editor.EditTitleAndBody()
				utils.Check(err)
			}

			params := octokit.ReleaseParams{
				TagName:         tag,
				TargetCommitish: commitish,
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
