package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdRelease = &Command{
		Run: listReleases,
		Usage: `
release [--include-drafts] [--exclude-prereleases] [-L <LIMIT>]
release show <TAG>
release create [-dpoc] [-a <FILE>] [-m <MESSAGE>|-F <FILE>] [-t <TARGET>] <TAG>
release edit [<options>] <TAG>
release download <TAG>
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
	-L, --limit
		Display only the first <LIMIT> releases.

	-d, --draft
		Create a draft release.

	-p, --prerelease
		Create a pre-release.

	-a, --attach <FILE>
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

	-f, --format <FORMAT>
		Pretty print releases using <FORMAT> (default: "%T%n"). See the "PRETTY
		FORMATS" section of git-log(1) for some additional details on how
		placeholders are used in format. The available placeholders for issues are:

		%U: the URL of this release

		%uT: tarball URL

		%uZ: zipball URL

		%uA: asset upload URL

		%S: state (i.e. "draft", "pre-release")

		%sC: set color to yellow or red, depending on state

		%t: release name

		%T: release tag

		%b: body

		%as: the list of assets attached to this release

		%cD: created date-only (no time of day)

		%cr: created date, relative

		%ct: created date, UNIX timestamp

		%cI: created date, ISO 8601 format

		%pD: published date-only (no time of day)

		%pr: published date, relative

		%pt: published date, UNIX timestamp

		%pI: published date, ISO 8601 format

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

	cmdDeleteRelease = &Command{
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
	flagReleaseFormat,
	flagReleaseCommitish string

	flagReleaseAssets stringSliceValue

	flagReleaseLimit int
)

func init() {
	cmdRelease.Flag.BoolVarP(&flagReleaseIncludeDrafts, "include-drafts", "d", false, "DRAFTS")
	cmdRelease.Flag.BoolVarP(&flagReleaseExcludePrereleases, "exclude-prereleases", "p", false, "PRERELEASE")
	cmdRelease.Flag.IntVarP(&flagReleaseLimit, "limit", "L", -1, "LIMIT")
	cmdRelease.Flag.StringVarP(&flagReleaseFormat, "format", "f", "%T%n", "FORMAT")

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
	cmdRelease.Use(cmdDeleteRelease)
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
		releases, err := gh.FetchReleases(project, flagReleaseLimit, func(release *github.Release) bool {
			return (!release.Draft || flagReleaseIncludeDrafts) &&
				(!release.Prerelease || !flagReleaseExcludePrereleases)
		})
		utils.Check(err)

		colorize := ui.IsTerminal(os.Stdout)
		for _, release := range releases {
			ui.Printf(formatRelease(release, flagReleaseFormat, colorize))
		}
	}

	args.NoForward()
}

func formatRelease(release github.Release, format string, colorize bool) string {
	state := ""
	stateColorSwitch := ""
	if release.Draft {
		state = "draft"
		stateColorSwitch = fmt.Sprintf("\033[%dm", 33)
	} else if release.Prerelease {
		state = "pre-release"
		stateColorSwitch = fmt.Sprintf("\033[%dm", 31)
	}

	var createdDate, createdAtISO8601, createdAtUnix, createdAtRelative,
		publishedDate, publishedAtISO8601, publishedAtUnix, publishedAtRelative string
	if !release.CreatedAt.IsZero() {
		createdDate = release.CreatedAt.Format("02 Jan 2006")
		createdAtISO8601 = release.CreatedAt.Format(time.RFC3339)
		createdAtUnix = fmt.Sprintf("%d", release.CreatedAt.Unix())
		createdAtRelative = utils.TimeAgo(release.CreatedAt)
	}
	if !release.PublishedAt.IsZero() {
		publishedDate = release.PublishedAt.Format("02 Jan 2006")
		publishedAtISO8601 = release.PublishedAt.Format(time.RFC3339)
		publishedAtUnix = fmt.Sprintf("%d", release.PublishedAt.Unix())
		publishedAtRelative = utils.TimeAgo(release.PublishedAt)
	}

	assets := make([]string, len(release.Assets))
	for i, asset := range release.Assets {
		assets[i] = fmt.Sprintf("%s\t%s", asset.DownloadUrl, asset.Label)
	}

	placeholders := map[string]string{
		"U":  release.HtmlUrl,
		"uT": release.TarballUrl,
		"uZ": release.ZipballUrl,
		"uA": release.UploadUrl,
		"S":  state,
		"sC": stateColorSwitch,
		"t":  release.Name,
		"T":  release.TagName,
		"b":  release.Body,
		"as": strings.Join(assets, "\n"),
		"cD": createdDate,
		"cI": createdAtISO8601,
		"ct": createdAtUnix,
		"cr": createdAtRelative,
		"pD": publishedDate,
		"pI": publishedAtISO8601,
		"pt": publishedAtUnix,
		"pr": publishedAtRelative,
	}

	return ui.Expand(format, placeholders, colorize)
}

func showRelease(cmd *Command, args *Args) {
	tagName := cmd.Arg(0)
	if tagName == "" {
		utils.Check(fmt.Errorf(cmdRelease.Synopsis()))
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
		utils.Check(fmt.Errorf(cmdRelease.Synopsis()))
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
		utils.Check(fmt.Errorf(cmdRelease.Synopsis()))
		return
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.CurrentProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	messageBuilder := &github.MessageBuilder{
		Filename: "RELEASE_EDITMSG",
		Title:    "release",
	}

	messageBuilder.AddCommentedSection(fmt.Sprintf(`Creating release %s for %s

Write a message for this release. The first block of
text is the title and the rest is the description.`, tagName, project))

	if cmd.FlagPassed("message") {
		messageBuilder.Message = flagReleaseMessage
		messageBuilder.Edit = flagReleaseEdit
	} else if cmd.FlagPassed("file") {
		messageBuilder.Message, err = msgFromFile(flagReleaseFile)
		utils.Check(err)
		messageBuilder.Edit = flagReleaseEdit
	} else {
		messageBuilder.Edit = true
	}

	title, body, err := messageBuilder.Extract()
	utils.Check(err)

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

	messageBuilder.Cleanup()

	uploadAssets(gh, release, flagReleaseAssets, args)
}

func editRelease(cmd *Command, args *Args) {
	tagName := cmd.Arg(0)
	if tagName == "" {
		utils.Check(fmt.Errorf(cmdRelease.Synopsis()))
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

	if cmd.FlagPassed("commitish") {
		params["target_commitish"] = flagReleaseCommitish
	}

	if cmd.FlagPassed("draft") {
		params["draft"] = flagReleaseDraft
	}

	if cmd.FlagPassed("prerelease") {
		params["prerelease"] = flagReleasePrerelease
	}

	messageBuilder := &github.MessageBuilder{
		Filename: "RELEASE_EDITMSG",
		Title:    "release",
	}

	messageBuilder.AddCommentedSection(fmt.Sprintf(`Editing release %s for %s

Write a message for this release. The first block of
text is the title and the rest is the description.`, tagName, project))

	if cmd.FlagPassed("message") {
		messageBuilder.Message = flagReleaseMessage
		messageBuilder.Edit = flagReleaseEdit
	} else if cmd.FlagPassed("file") {
		messageBuilder.Message, err = msgFromFile(flagReleaseFile)
		utils.Check(err)
		messageBuilder.Edit = flagReleaseEdit
	} else {
		messageBuilder.Edit = true
		messageBuilder.Message = fmt.Sprintf("%s\n\n%s", release.Name, release.Body)
	}

	title, body, err := messageBuilder.Extract()
	utils.Check(err)

	if title == "" && !cmd.FlagPassed("message") {
		utils.Check(fmt.Errorf("Aborting editing due to empty release title"))
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

		messageBuilder.Cleanup()
	}

	uploadAssets(gh, release, flagReleaseAssets, args)
	args.NoForward()
}

func deleteRelease(cmd *Command, args *Args) {
	tagName := cmd.Arg(0)
	if tagName == "" {
		utils.Check(fmt.Errorf(cmdRelease.Synopsis()))
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
