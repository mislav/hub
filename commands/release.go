package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"github.com/jingweno/go-octokit/octokit"
	"os"
	"path/filepath"
	"strings"
)

var (
	cmdReleases = &Command{
		Run:   releases,
		Usage: "releases",
		Short: "Retrieve releases from GitHub",
		Long:  `Retrieve releases from GitHub for the project that the "origin" remote points to.`}

	cmdRelease = &Command{
		Run:   release,
		Usage: "release TAG [-d] [-p] [-a <ASSETS_DIR>] [-m <MESSAGE>|-f <FILE>]",
		Short: "Create a new release in GitHub",
		Long: `Create a new release in GitHub for the project that the "origin" remote points to.
- It requires the name of the tag to release as a first argument.
- The assets to include in the release are taken from releases/TAG or from the directory specified by -a.
- Use the flag -d to create a draft.
- Use the flag -p to create a prerelease.
`}

	flagReleaseDraft,
	flagReleasePrerelease bool

	flagReleaseAssetsDir,
	flagReleaseMessage,
	flagReleaseFile string
)

func init() {
	cmdRelease.Flag.BoolVar(&flagReleaseDraft, "d", false, "DRAFT")
	cmdRelease.Flag.BoolVar(&flagReleasePrerelease, "p", false, "PRERELEASE")
	cmdRelease.Flag.StringVar(&flagReleaseAssetsDir, "a", "", "ASSETS_DIR")
	cmdRelease.Flag.StringVar(&flagReleaseMessage, "m", "", "MESSAGE")
	cmdRelease.Flag.StringVar(&flagReleaseFile, "f", "", "FILE")
}

func releases(cmd *Command, args *Args) {
	runInLocalRepo(func(localRepo *github.GitHubRepo, project *github.Project, gh *github.Client) {
		if args.Noop {
			fmt.Printf("Would request list of releases for %s\n", project)
		} else {
			releases, err := gh.Releases(project)
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

func release(cmd *Command, args *Args) {
	tag := args.LastParam()

	assetsDir, err := getAssetsDirectory(flagReleaseAssetsDir, tag)
	utils.Check(err)

	runInLocalRepo(func(localRepo *github.GitHubRepo, project *github.Project, gh *github.Client) {
		currentBranch, err := localRepo.CurrentBranch()
		utils.Check(err)

		title, body, err := github.GetTitleAndBodyFromFlags(flagReleaseMessage, flagReleaseFile)
		utils.Check(err)

		if title == "" {
			title, body, err = github.GetTitleAndBodyFromEditor(nil)
			utils.Check(err)
		}

		params := octokit.ReleaseParams{
			TagName:         tag,
			TargetCommitish: currentBranch.ShortName(),
			Name:            title,
			Body:            body,
			Draft:           flagReleaseDraft,
			Prerelease:      flagReleasePrerelease}

		finalRelease, err := gh.CreateRelease(project, params)
		utils.Check(err)

		uploadReleaseAssets(gh, finalRelease, assetsDir)
	})
}

func runInLocalRepo(fn func(localRepo *github.GitHubRepo, project *github.Project, client *github.Client)) {
	localRepo := github.LocalRepo()
	project, err := localRepo.CurrentProject()
	utils.Check(err)

	client := github.NewClient(project.Host)
	fn(localRepo, project, client)

	os.Exit(0)
}

func getAssetsDirectory(assetsDir, tag string) (string, error) {
	if assetsDir == "" {
		pwd, err := os.Getwd()
		utils.Check(err)

		assetsDir = filepath.Join(pwd, "releases", tag)
	}

	if !utils.IsDir(assetsDir) {
		return "", fmt.Errorf("The assets directory doesn't exist: %s", assetsDir)
	}

	if utils.IsEmptyDir(assetsDir) {
		return "", fmt.Errorf("The assets directory is empty: %s", assetsDir)
	}

	return assetsDir, nil
}

func uploadReleaseAssets(gh *github.Client, release *octokit.Release, assetsDir string) {
	filepath.Walk(assetsDir, func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			fmt.Printf("- Uploading asset %s\n", fi.Name())

			file, err := os.Open(path)
			utils.Check(err)
			defer file.Close()

			err = gh.UploadReleaseAsset(release, file, fi)
			utils.Check(err)
		}

		return nil
	})
}
