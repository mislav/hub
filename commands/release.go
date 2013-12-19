package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"os"
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
		Usage: "release TAG [-d] [-p] [-a <ASSETS_DIR>]",
		Short: "Create a new release in GitHub",
		Long: `Create a new release in GitHub for the project that the "origin" remote points to.
- It requires the name of the tag to release as a first argument.
- The assets to include in the release are taken from releases/TAG or from the directory specified by -a.
- Use the flag -d to create a draft.
- Use the flag -p to create a prerelease.
`}

	flagReleaseDraft,
	flagReleasePrerelease bool
	flagReleaseAssetsDir string
)

func init() {
	cmdRelease.Flag.BoolVar(&flagReleaseDraft, "d", false, "DRAFT")
	cmdRelease.Flag.BoolVar(&flagReleasePrerelease, "p", false, "PRERELEASE")
	cmdRelease.Flag.StringVar(&flagReleaseAssetsDir, "a", "", "ASSETS_DIR")
}

func releases(cmd *Command, args *Args) {
	runInLocalRepo(func(project *github.Project, gh *github.Client) {
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
}

func runInLocalRepo(fn func(project *github.Project, client *github.Client)) {
	localRepo := github.LocalRepo()
	project, err := localRepo.CurrentProject()
	utils.Check(err)

	client := github.NewClient(project.Host)
	fn(project, client)

	os.Exit(0)
}
