package commands

import (
	"fmt"
	"io"
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
release [--include-drafts] [--exclude-prereleases]
release show <TAG>
release create [-dpoc] [-a <FILE>] [-m <MESSAGE>|-F <FILE>] [-t <TARGET>] <TAG>
release edit [<options>] <TAG>
release delete <TAG>
`,
		Long: `Manage GitHub releases.

## Commands:

With no arguments, shows a list of existing releases.

With '--include-drafts', include draft releases in the listing.
With '--exclude-prereleases', exclude non-stable releases from the listing.

	* _show_:
		Show GitHub release notes for <TAG>.

		With '--show-downloads', include the "Downloads" section.

	* _create_:
		Create a GitHub release for the specified <TAG> name. If git tag <TAG>
		doesn't exist, it will be created at <TARGET> (default: current branch).

	* _edit_:
		Edit the GitHub release for the specified <TAG> name. Accepts the same
		options as _create_ command. Publish a draft with '--draft=false'.

		When <MESSAGE> or <FILE> are not specified, a text editor will open
		pre-populated with current release title and body. To re-use existing title
		and body unchanged, pass '-m ""'.

	* _download_:
	  Download the assets attached to release for the specified <TAG>.

	* _delete_:
	  Delete the release and associated assets for the specified <TAG>.

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

	-F, --file <FILE>
		Read the release title and description from <FILE>.

	-e, --edit
		Further edit the contents of <FILE> in a text editor before submitting.

	-o, --browse
		Open the new release in a web browser.

	-c, --copy
		Put the URL of the new release to clipboard instead of printing it.

	-t, --commitish <TARGET>
		A commit SHA or branch name to attach the release to, only used if <TAG>
		doesn't already exist (default: main branch).

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

	cmdDownloadRelease = &Command{
		Key: "download",
		Run: downloadRelease,
	}

	cmdDeleteRelase = &Command{
		Key: "delete",
		Run: deleteRelease,
	}

	flagReleaseIncludeDrafts,
	flagReleaseExcludePrereleases,
	flagReleaseShowDownloads,
	flagReleaseDraft,
	flagReleaseEdit,
	flagReleaseBrowse,
	flagReleaseCopy,
	flagReleasePrerelease bool

	flagReleaseMessage,
	flagReleaseFile,
	flagReleaseCommitish string

	flagReleaseAssets stringSliceValue
)

func init() {
	cmdRelease.Flag.BoolVarP(&flagReleaseIncludeDrafts, "include-drafts", "d", false, "DRAFTS")
	cmdRelease.Flag.BoolVarP(&flagReleaseExcludePrereleases, "exclude-prereleases", "p", false, "PRERELEASE")

	cmdShowRelease.Flag.BoolVarP(&flagReleaseShowDownloads, "show-downloads", "d", false, "DRAFTS")

	cmdCreateRelease.Flag.BoolVarP(&flagReleaseEdit, "edit", "e", false, "EDIT")
	cmdCreateRelease.Flag.BoolVarP(&flagReleaseDraft, "draft", "d", false, "DRAFT")
	cmdCreateRelease.Flag.BoolVarP(&flagReleasePrerelease, "prerelease", "p", false, "PRERELEASE")
	cmdCreateRelease.Flag.BoolVarP(&flagReleaseBrowse, "browse", "o", false, "BROWSE")
	cmdCreateRelease.Flag.BoolVarP(&flagReleaseCopy, "copy", "c", false, "COPY")
	cmdCreateRelease.Flag.VarP(&flagReleaseAssets, "attach", "a", "ATTACH_ASSETS")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseMessage, "message", "m", "", "MESSAGE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseFile, "file", "F", "", "FILE")
	cmdCreateRelease.Flag.StringVarP(&flagReleaseCommitish, "commitish", "t", "", "COMMITISH")

	cmdEditRelease.Flag.BoolVarP(&flagReleaseEdit, "edit", "e", false, "EDIT")
	cmdEditRelease.Flag.BoolVarP(&flagReleaseDraft, "draft", "d", false, "DRAFT")
	cmdEditRelease.Flag.BoolVarP(&flagReleasePrerelease, "prerelease", "p", false, "PRERELEASE")
	cmdEditRelease.Flag.VarP(&flagReleaseAssets, "attach", "a", "ATTACH_ASSETS")
	cmdEditRelease.Flag.StringVarP(&flagReleaseMessage, "message", "m", "", "MESSAGE")
	cmdEditRelease.Flag.StringVarP(&flagReleaseFile, "file", "F", "", "FILE")
	cmdEditRelease.Flag.StringVarP(&flagReleaseCommitish, "commitish", "t", "", "COMMITISH")

	cmdRelease.Use(cmdShowRelease)
	cmdRelease.Use(cmdCreateRelease)
	cmdRelease.Use(cmdEditRelease)
	cmdRelease.Use(cmdDownloadRelease)
	cmdRelease.Use(cmdDeleteRelase)
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
			if (!release.Draft || flagReleaseIncludeDrafts) &&
				(!release.Prerelease || !flagReleaseExcludePrereleases) {
				ui.Println(release.TagName)
			}
		}
	}

	args.NoForward()
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

		ui.Println(release.Name)
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

	args.NoForward()
}

func downloadRelease(cmd *Command, args *Args) {
	tagName := cmd.Arg(0)
	if tagName == "" {
		utils.Check(fmt.Errorf("Missing argument TAG"))
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	release, err := gh.FetchRelease(project, tagName)
	utils.Check(err)

	for _, asset := range release.Assets {
		ui.Printf("Downloading %s ...\n", asset.Name)
		err := downloadReleaseAsset(asset, gh)
		utils.Check(err)
	}

	args.NoForward()
}

func downloadReleaseAsset(asset github.ReleaseAsset, gh *github.Client) (err error) {
	assetReader, err := gh.DownloadReleaseAsset(asset.ApiUrl)
	if err != nil {
		return
	}
	defer assetReader.Close()

	assetFile, err := os.OpenFile(asset.Name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return
	}
	defer assetFile.Close()

	_, err = io.Copy(assetFile, assetReader)
	if err != nil {
		return
	}
	return
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

	var title string
	var body string
	var editor *github.Editor

	if cmd.FlagPassed("message") {
		title, body = readMsg(flagReleaseMessage)
	} else if cmd.FlagPassed("file") {
		title, body, editor, err = readMsgFromFile(flagReleaseFile, flagReleaseEdit, "RELEASE", "release")
		utils.Check(err)
	} else {
		cs := git.CommentChar()
		message, err := renderReleaseTpl("Creating", cs, tagName, project.String(), flagReleaseCommitish)
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
		TargetCommitish: flagReleaseCommitish,
		Name:            title,
		Body:            body,
		Draft:           flagReleaseDraft,
		Prerelease:      flagReleasePrerelease,
	}

	var release *github.Release

	args.NoForward()
	if args.Noop {
		ui.Printf("Would create release `%s' for %s with tag name `%s'\n", title, project, tagName)
	} else {
		release, err = gh.CreateRelease(project, params)
		utils.Check(err)

		printBrowseOrCopy(args, release.HtmlUrl, flagReleaseBrowse, flagReleaseCopy)
	}

	if editor != nil {
		editor.DeleteFile()
	}

	uploadAssets(gh, release, flagReleaseAssets, args)
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
		title, body, editor, err = readMsgFromFile(flagReleaseFile, flagReleaseEdit, "RELEASE", "release")
		utils.Check(err)

		if title == "" {
			utils.Check(fmt.Errorf("Aborting editing due to empty release title"))
		}
	} else {
		cs := git.CommentChar()
		message, err := renderReleaseTpl("Editing", cs, tagName, project.String(), commitish)
		utils.Check(err)

		message = fmt.Sprintf("%s\n\n%s\n%s", release.Name, release.Body, message)
		editor, err := github.NewEditor("RELEASE", "release", message)
		utils.Check(err)

		title, body, err = editor.EditTitleAndBody()
		utils.Check(err)

		if title == "" {
			utils.Check(fmt.Errorf("Aborting editing due to empty release title"))
		}
	}

	if title != "" {
		params["name"] = title
	}
	if body != "" {
		params["body"] = body
	}

	if len(params) > 0 {
		if args.Noop {
			ui.Printf("Would edit release `%s'\n", tagName)
		} else {
			release, err = gh.EditRelease(release, params)
			utils.Check(err)
		}

		if editor != nil {
			editor.DeleteFile()
		}
	}

	uploadAssets(gh, release, flagReleaseAssets, args)
	args.NoForward()
}

func deleteRelease(cmd *Command, args *Args) {
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

	if args.Noop {
		message := fmt.Sprintf("Deleting release related to %s...", tagName)
		ui.Println(message)
	} else {
		err = gh.DeleteRelease(release)
		utils.Check(err)
	}

	args.NoForward()
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
