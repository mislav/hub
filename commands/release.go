package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/github/hub/Godeps/_workspace/src/github.com/octokit/go-octokit/octokit"
	"github.com/github/hub/github"
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
		Usage: "release create [-d] [-p] [-a <ASSETS_DIR>] [-m <MESSAGE>|-f <FILE>] <TAG>",
		Short: "Create a new release in GitHub",
		Long: `Creates a new release in GitHub for the project that the "origin" remote points to.
It requires the name of the tag to release as a first argument.

Specify the assets to include in the release from a directory via "-a".

Without <MESSAGE> or <FILE>, a text editor will open in which title and body
of the release can be entered in the same manner as git commit message.

If "-d" is given, it creates a draft release.

If "-p" is given, it creates a pre-release.
`}

	flagReleaseDraft,
	flagReleasePrerelease bool

	flagReleaseAssets,
	flagReleaseMessage,
	flagReleaseFile string
)

func init() {
	cmdCreateRelease.Flag.BoolVarP(&flagReleaseDraft, "draft", "d", false, "DRAFT")
	cmdCreateRelease.Flag.BoolVarP(&flagReleasePrerelease, "prerelease", "p", false, "PRERELEASE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseAssets, "attach", "a", "", "ATTACH_ASSETS")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseMessage, "message", "m", "", "MESSAGE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseFile, "file", "f", "", "FILE")

	cmdRelease.Use(cmdCreateRelease)
	CmdRunner.Use(cmdRelease)
}

func release(cmd *Command, args *Args) {
	runInLocalRepo(func(localRepo *github.GitHubRepo, project *github.Project, client *github.Client) {
		if args.Noop {
			fmt.Printf("Would request list of releases for %s\n", project)
		} else {
			releases, err := client.Releases(project)
			utils.Check(err)
			var outputs []string
			for _, release := range releases {
				out := fmt.Sprintf("%s (%s)\n%s", release.Name, release.TagName, release.Body)
				outputs = append(outputs, out)
			}

			fmt.Println(strings.Join(outputs, "\n\n"))
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
		currentBranch, err := localRepo.CurrentBranch()
		utils.Check(err)
		branchName := currentBranch.ShortName()

		title, body, err := getTitleAndBodyFromFlags(flagReleaseMessage, flagReleaseFile)
		utils.Check(err)

		if title == "" {
			title, body, err = writeReleaseTitleAndBody(project, tag, branchName)
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

		finalRelease, err := client.CreateRelease(project, params)
		utils.Check(err)

		uploadReleaseAssets(client, finalRelease)

		fmt.Printf("\n\nRelease created: %s", finalRelease.HTMLURL)
	})
}

func writeReleaseTitleAndBody(project *github.Project, tag, currentBranch string) (string, string, error) {
	message := `
# Creating release %s for %s from %s
#
# Write a message for this release. The first block
# of text is the title and the rest is description.
`
	message = fmt.Sprintf(message, tag, project.Name, currentBranch)

	editor, err := github.NewEditor("RELEASE", "release", message)
	if err != nil {
		return "", "", err
	}

	return editor.EditTitleAndBody()
}

func uploadReleaseAssets(client *github.Client, release *octokit.Release) {
	if flagReleaseAssets == "" {
		return
	}

	assetInfo, err := os.Stat(flagReleaseAssets)
	utils.Check(err)

	var wg sync.WaitGroup
	var totalAssets, countAssets uint64

	notifyProgress := func() {
		atomic.AddUint64(&countAssets, uint64(1))
		printUploadProgress(&countAssets, totalAssets)
		wg.Done()
	}

	if assetInfo.IsDir() {
		filepath.Walk(flagReleaseAssets, func(path string, fi os.FileInfo, err error) error {
			if !fi.IsDir() {
				totalAssets += 1
			}
			return nil
		})

		printUploadProgress(&countAssets, totalAssets)

		filepath.Walk(flagReleaseAssets, func(path string, fi os.FileInfo, err error) error {
			if !fi.IsDir() {
				wg.Add(1)
				go uploadAsset(client, release, fi, path, notifyProgress)
			}
			return nil
		})
	} else {
		totalAssets = 1
		printUploadProgress(&countAssets, totalAssets)
		wg.Add(1)
		uploadAsset(client, release, assetInfo, flagReleaseAssets, notifyProgress)
	}

	wg.Wait()
}

func uploadAsset(gh *github.Client, release *octokit.Release, fi os.FileInfo, path string, notifyProgress func()) {
	defer notifyProgress()
	uploadUrl, err := release.UploadURL.Expand(octokit.M{"name": fi.Name()})
	utils.Check(err)

	contentType := detectContentType(path, fi)

	file, err := os.Open(path)
	utils.Check(err)
	defer file.Close()

	err = gh.UploadReleaseAsset(uploadUrl, file, contentType)
	utils.Check(err)
}

func detectContentType(path string, fi os.FileInfo) string {
	file, err := os.Open(path)
	utils.Check(err)
	defer file.Close()

	fileHeader := &bytes.Buffer{}
	headerSize := int64(512)
	if fi.Size() < headerSize {
		headerSize = fi.Size()
	}

	// The content type detection only uses 512 bytes at most.
	// This way we avoid copying the whole content for big files.
	_, err = io.CopyN(fileHeader, file, headerSize)
	utils.Check(err)

	return http.DetectContentType(fileHeader.Bytes())
}

func printUploadProgress(count *uint64, total uint64) {
	out := fmt.Sprintf("Uploading assets (%d/%d)", atomic.LoadUint64(count), total)
	fmt.Print("\r" + out)
}
