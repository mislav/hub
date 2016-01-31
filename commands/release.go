package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdRelease = &Command{
		Run: listReleases,
		Usage: `
release
release show <TAG>
release create [-dp] [-a <FILE>] [-m <MESSAGE>|-f <FILE>] [-c <COMMIT>] <TAG>
release edit [<options>] <TAG>
`,
		Long: `Manage GitHub releases.

## Commands:

With no arguments, shows a list of existing releases.

With '--include-drafs', include draft releases in the listing.

	* _show_:
		Show GitHub release notes for <TAG>.

		With '--show-downloads', include the "Downloads" section.

	* _create_:
		Create a GitHub release for the specified <TAG> name. If git tag <TAG>
		doesn't exist, it will be created at <COMMIT> (default: HEAD).

	* _edit_:
		Edit the GitHub release for the specified <TAG> name. Accepts the same
		options as _create_ command. Publish a draft with '--draft=false'.

## Options:
	-d, --draft
		Create a draft release.

	-p, --prerelease
		Create a pre-release.

	-a, --asset <FILE>
		Attach a file as an asset for this release.

		If <FILE> is in the "<filename>#<text>" format, the text after the '#'
		character is taken as asset label.

	-m, --message <MESSAGE>
		Use the first line of <MESSAGE> as release title, and the rest as release description.

	-f, --file <FILE>
		Read the release title and description from <FILE>.
	
	-c, --commitish <COMMIT>
		A SHA, tag, or branch name to attach the release to (default: current branch).

	<TAG>
		The git tag name for this release.

## See also:

hub(1), git-tag(1)
	`,
	}

	cmdShowRelease = &Command{
		Key: "show",
		Run: showRelease,
	}

	cmdCreateRelease = &Command{
		Key: "create",
		Run: createRelease,
	}

	cmdEditRelease = &Command{
		Key: "edit",
		Run: editRelease,
	}

	flagReleaseIncludeDrafts,
	flagReleaseShowDownloads,
	flagReleaseDraft,
	flagReleasePrerelease bool

	flagReleaseMessage,
	flagReleaseFile,
	flagReleaseCommitish string

	flagReleaseAssets stringSliceValue
)

func init() {
	cmdRelease.Flag.BoolVarP(&flagReleaseIncludeDrafts, "include-drafts", "d", false, "DRAFTS")

	cmdShowRelease.Flag.BoolVarP(&flagReleaseShowDownloads, "show-downloads", "d", false, "DRAFTS")

	cmdCreateRelease.Flag.BoolVarP(&flagReleaseDraft, "draft", "d", false, "DRAFT")
	cmdCreateRelease.Flag.BoolVarP(&flagReleasePrerelease, "prerelease", "p", false, "PRERELEASE")
	cmdCreateRelease.Flag.VarP(&flagReleaseAssets, "attach", "a", "ATTACH_ASSETS")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseMessage, "message", "m", "", "MESSAGE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseFile, "file", "f", "", "FILE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseCommitish, "commitish", "c", "", "COMMITISH")

	cmdEditRelease.Flag.BoolVarP(&flagReleaseDraft, "draft", "d", false, "DRAFT")
	cmdEditRelease.Flag.BoolVarP(&flagReleasePrerelease, "prerelease", "p", false, "PRERELEASE")
	cmdEditRelease.Flag.VarP(&flagReleaseAssets, "attach", "a", "ATTACH_ASSETS")
	cmdEditRelease.Flag.StringVarP(&flagReleaseMessage, "message", "m", "", "MESSAGE")
	cmdEditRelease.Flag.StringVarP(&flagReleaseFile, "file", "f", "", "FILE")
	cmdEditRelease.Flag.StringVarP(&flagReleaseCommitish, "commitish", "c", "", "COMMITISH")

	cmdRelease.Use(cmdShowRelease)
	cmdRelease.Use(cmdCreateRelease)
	cmdRelease.Use(cmdEditRelease)
	CmdRunner.Use(cmdRelease)
}

func listReleases(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	if args.Noop {
		ui.Printf("Would request list of releases for %s\n", project)
	} else {
		releases, err := gh.FetchReleases(project)
		utils.Check(err)

		for _, release := range releases {
			if !release.Draft || flagReleaseIncludeDrafts {
				ui.Println(release.TagName)
			}
		}
	}

	os.Exit(0)
}

func showRelease(cmd *Command, args *Args) {
	tagName := cmd.Arg(0)
	if tagName == "" {
		utils.Check(fmt.Errorf("Missing argument TAG"))
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	if args.Noop {
		ui.Printf("Would display information for `%s' release\n", tagName)
	} else {
		release, err := gh.FetchRelease(project, tagName)
		utils.Check(err)

		body := strings.TrimSpace(release.Body)

		ui.Printf("%s (%s)\n", release.Name, release.TagName)
		if body != "" {
			ui.Printf("\n%s\n", body)
		}
		if flagReleaseShowDownloads {
			ui.Printf("\n## Downloads\n\n")
			for _, asset := range release.Assets {
				ui.Println(asset.DownloadUrl)
			}
			if release.ZipballUrl != "" {
				ui.Println(release.ZipballUrl)
				ui.Println(release.TarballUrl)
			}
		}
	}

	os.Exit(0)
}

func createRelease(cmd *Command, args *Args) {
	tagName := cmd.Arg(0)
	if tagName == "" {
		utils.Check(fmt.Errorf("Missing argument TAG"))
		return
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.CurrentProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	commitish := flagReleaseCommitish
	if commitish == "" {
		currentBranch, err := localRepo.CurrentBranch()
		utils.Check(err)
		commitish = currentBranch.ShortName()
	}

	var title string
	var body string
	var editor *github.Editor

	if cmd.FlagPassed("message") {
		title, body = readMsg(flagReleaseMessage)
	} else if cmd.FlagPassed("file") {
		title, body, err = readMsgFromFile(flagReleaseMessage)
		utils.Check(err)
	} else {
		cs := git.CommentChar()
		message, err := renderReleaseTpl("Creating", cs, tagName, project.String(), commitish)
		utils.Check(err)

		editor, err := github.NewEditor("RELEASE", "release", message)
		utils.Check(err)

		title, body, err = editor.EditTitleAndBody()
		utils.Check(err)
	}

	if title == "" {
		utils.Check(fmt.Errorf("Aborting release due to empty release title"))
	}

	params := &github.Release{
		TagName:         tagName,
		TargetCommitish: commitish,
		Name:            title,
		Body:            body,
		Draft:           flagReleaseDraft,
		Prerelease:      flagReleasePrerelease,
	}

	var release *github.Release

	if args.Noop {
		ui.Printf("Would create release `%s' for %s with tag name `%s'\n", title, project, tagName)
	} else {
		release, err = gh.CreateRelease(project, params)
		utils.Check(err)

		ui.Println(release.HtmlUrl)
	}

	uploadAssets(gh, release, flagReleaseAssets, args)

	if editor != nil {
		editor.DeleteFile()
	}
	os.Exit(0)
}

func editRelease(cmd *Command, args *Args) {
	tagName := cmd.Arg(0)
	if tagName == "" {
		utils.Check(fmt.Errorf("Missing argument TAG"))
		return
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.CurrentProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	release, err := gh.FetchRelease(project, tagName)
	utils.Check(err)

	params := map[string]interface{}{}
	commitish := release.TargetCommitish

	if cmd.FlagPassed("commitish") {
		params["target_commitish"] = flagReleaseCommitish
		commitish = flagReleaseCommitish
	}

	if cmd.FlagPassed("draft") {
		params["draft"] = flagReleaseDraft
	}

	if cmd.FlagPassed("prerelease") {
		params["prerelease"] = flagReleasePrerelease
	}

	var title string
	var body string
	var editor *github.Editor

	if cmd.FlagPassed("message") {
		title, body = readMsg(flagReleaseMessage)
	} else if cmd.FlagPassed("file") {
		title, body, err = readMsgFromFile(flagReleaseMessage)
		utils.Check(err)
	} else {
		cs := git.CommentChar()
		message, err := renderReleaseTpl("Editing", cs, tagName, project.String(), commitish)
		utils.Check(err)

		message = fmt.Sprintf("%s\n\n%s\n%s", release.Name, release.Body, message)
		editor, err := github.NewEditor("RELEASE", "release", message)
		utils.Check(err)

		title, body, err = editor.EditTitleAndBody()
		utils.Check(err)
	}

	if title != "" {
		params["name"] = title
	}
	if body != "" {
		params["body"] = body
	}

	if args.Noop {
		ui.Printf("Would edit release `%s'\n", tagName)
	} else {
		release, err = gh.EditRelease(release, params)
		utils.Check(err)

		if editor != nil {
			editor.DeleteFile()
		}
	}

	uploadAssets(gh, release, flagReleaseAssets, args)
	os.Exit(0)
}

func uploadAssets(gh *github.Client, release *github.Release, assets []string, args *Args) {
	for _, asset := range assets {
		var label string
		parts := strings.SplitN(asset, "#", 2)
		asset = parts[0]
		if len(parts) > 1 {
			label = parts[1]
		}

		if args.Noop {
			if label == "" {
				ui.Errorf("Would attach release asset `%s'\n", asset)
			} else {
				ui.Errorf("Would attach release asset `%s' with label `%s'\n", asset, label)
			}
		} else {
			for _, existingAsset := range release.Assets {
				if existingAsset.Name == filepath.Base(asset) {
					err := gh.DeleteReleaseAsset(&existingAsset)
					utils.Check(err)
					break
				}
			}
			ui.Errorf("Attaching release asset `%s'...\n", asset)
			_, err := gh.UploadReleaseAsset(release, asset, label)
			utils.Check(err)
		}
	}
}
