package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
)

var (
	cmdRelease = &Command{
		Run: listReleases,
		Usage: `
release list [--include-drafts] [--exclude-prereleases] [-L <LIMIT>] [-f <FORMAT>]
release show [-f <FORMAT>] <TAG>
release create [-dpoc] [-a <FILE>] [-m <MESSAGE>|-F <FILE>] [-t <TARGET>] <TAG>
release edit [<options>] <TAG>
release download <TAG> [-i <PATTERN>]
release delete <TAG>
`,
		Long: `Manage GitHub Releases for the current repository.

## Commands:

With no arguments, performs ''hub release list''

	* _list_:
		Show a list of existing releases for the current repository.

		With  --include-drafts,  include  draft  releases  in the listing. With
		--prereleases, exclude non-stable releases from the listing.

	* _show_:
		Show GitHub release notes for <TAG>.

		With ''--show-downloads'', include the "Downloads" section.

	* _create_:
		Create a GitHub release for the specified <TAG> name. If git tag <TAG>
		does not exist, it will be created at <TARGET> (default: current branch).

	* _edit_:
		Edit the GitHub release for the specified <TAG> name. Accepts the same
		options as _create_ command. Publish a draft with ''--draft=false''.

		Without ''--message'' or ''--file'', a text editor will open pre-populated with
		the current release title and body. To re-use existing title and body
		unchanged, pass ''-m ""''.

	* _download_:
		Download the assets attached to release for the specified <TAG>.

	* _delete_:
		Delete the release and associated assets for the specified <TAG>. Note that
		this does **not** remove the git tag <TAG>.

## Options:
	-d, --include-drafts
		List drafts together with published releases.

	-p, --exclude-prereleases
		Exclude prereleases from the list.

	-L, --limit
		Display only the first <LIMIT> releases.

	-d, --draft
		Create a draft release.

	-p, --prerelease
		Create a pre-release.

	-a, --attach <FILE>
		Attach a file as an asset for this release.

		If <FILE> is in the "<filename>#<text>" format, the text after the "#"
		character is taken as asset label.

	-m, --message <MESSAGE>
		The text up to the first blank line in <MESSAGE> is treated as the release
		title, and the rest is used as release description in Markdown format.

		When multiple ''--message'' are passed, their values are concatenated with a
		blank line in-between.

		When neither ''--message'' nor ''--file'' were supplied to ''release create'', a
		text editor will open to author the title and description in.

	-F, --file <FILE>
		Read the release title and description from <FILE>. Pass "-" to read from
		standard input instead. See ''--message'' for the formatting rules.

	-e, --edit
		Open the release title and description in a text editor before submitting.
		This can be used in combination with ''--message'' or ''--file''.

	-o, --browse
		Open the new release in a web browser.

	-c, --copy
		Put the URL of the new release to clipboard instead of printing it.

	-t, --commitish <TARGET>
		A commit SHA or branch name to attach the release to, only used if <TAG>
		does not already exist (default: main branch).

	-i, --include <PATTERN>
		Filter the files in the release to those that match the glob <PATTERN>.

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

		%n: newline

		%%: a literal %

	--color[=<WHEN>]
		Enable colored output even if stdout is not a terminal. <WHEN> can be one
		of "always" (default for ''--color''), "never", or "auto" (default).

	<TAG>
		The git tag name for this release.

## See also:

hub(1), git-tag(1)
`,
		KnownFlags: `
		-d, --include-drafts
		-p, --exclude-prereleases
		-L, --limit N
		-f, --format FMT
		--color
`,
	}

	cmdListRelease = &Command{
		Key: "list",
		Run: listReleases,
		KnownFlags: `
		-d, --include-drafts
		-p, --exclude-prereleases
		-L, --limit N
		-f, --format FMT
		--color
`,
	}

	cmdShowRelease = &Command{
		Key: "show",
		Run: showRelease,
		KnownFlags: `
		-d, --show-downloads
		-f, --format FMT
		--color
`,
	}

	cmdCreateRelease = &Command{
		Key: "create",
		Run: createRelease,
		KnownFlags: `
		-e, --edit
		-d, --draft
		-p, --prerelease
		-o, --browse
		-c, --copy
		-a, --attach FILE
		-m, --message MSG
		-F, --file FILE
		-t, --commitish C
`,
	}

	cmdEditRelease = &Command{
		Key: "edit",
		Run: editRelease,
		KnownFlags: `
		-e, --edit
		-d, --draft
		-p, --prerelease
		-a, --attach FILE
		-m, --message MSG
		-F, --file FILE
		-t, --commitish C
`,
	}

	cmdDownloadRelease = &Command{
		Key: "download",
		Run: downloadRelease,
		KnownFlags: `
		-i, --include PATTERN
		`,
	}

	cmdDeleteRelease = &Command{
		Key: "delete",
		Run: deleteRelease,
	}
)

func init() {
	cmdRelease.Use(cmdListRelease)
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

	flagReleaseLimit := args.Flag.Int("--limit")
	flagReleaseIncludeDrafts := args.Flag.Bool("--include-drafts")
	flagReleaseExcludePrereleases := args.Flag.Bool("--exclude-prereleases")

	if args.Noop {
		ui.Printf("Would request list of releases for %s\n", project)
	} else {
		releases, err := gh.FetchReleases(project, flagReleaseLimit, func(release *github.Release) bool {
			return (!release.Draft || flagReleaseIncludeDrafts) &&
				(!release.Prerelease || !flagReleaseExcludePrereleases)
		})
		utils.Check(err)

		colorize := colorizeOutput(args.Flag.HasReceived("--color"), args.Flag.Value("--color"))
		for _, release := range releases {
			flagReleaseFormat := "%T%n"
			if args.Flag.HasReceived("--format") {
				flagReleaseFormat = args.Flag.Value("--format")
			}
			ui.Print(formatRelease(release, flagReleaseFormat, colorize))
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
		assets[i] = fmt.Sprintf("%s\t%s", asset.DownloadURL, asset.Label)
	}

	placeholders := map[string]string{
		"U":  release.HTMLURL,
		"uT": release.TarballURL,
		"uZ": release.ZipballURL,
		"uA": release.UploadURL,
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
	tagName := ""
	if args.ParamsSize() > 0 {
		tagName = args.GetParam(0)
	}
	if tagName == "" {
		utils.Check(cmd.UsageError(""))
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	args.NoForward()

	if args.Noop {
		ui.Printf("Would display information for `%s' release\n", tagName)
	} else {
		release, err := gh.FetchRelease(project, tagName)
		utils.Check(err)

		body := strings.TrimSpace(release.Body)

		colorize := colorizeOutput(args.Flag.HasReceived("--color"), args.Flag.Value("--color"))
		if flagShowReleaseFormat := args.Flag.Value("--format"); flagShowReleaseFormat != "" {
			ui.Print(formatRelease(*release, flagShowReleaseFormat, colorize))
			return
		}

		ui.Println(release.Name)
		if body != "" {
			ui.Printf("\n%s\n", body)
		}
		if args.Flag.Bool("--show-downloads") {
			ui.Printf("\n## Downloads\n\n")
			for _, asset := range release.Assets {
				ui.Println(asset.DownloadURL)
			}
			if release.ZipballURL != "" {
				ui.Println(release.ZipballURL)
				ui.Println(release.TarballURL)
			}
		}
	}
}

func downloadRelease(cmd *Command, args *Args) {
	tagName := ""
	if args.ParamsSize() > 0 {
		tagName = args.GetParam(0)
	}
	if tagName == "" {
		utils.Check(cmd.UsageError(""))
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	release, err := gh.FetchRelease(project, tagName)
	utils.Check(err)

	hasPattern := args.Flag.HasReceived("--include")
	found := false
	for _, asset := range release.Assets {
		if hasPattern {
			isMatch, err := filepath.Match(args.Flag.Value("--include"), asset.Name)
			utils.Check(err)
			if !isMatch {
				continue
			}
		}

		found = true
		ui.Printf("Downloading %s ...\n", asset.Name)
		err := downloadReleaseAsset(asset, gh)
		utils.Check(err)
	}

	if !found && hasPattern {
		names := []string{}
		for _, asset := range release.Assets {
			names = append(names, asset.Name)
		}
		utils.Check(fmt.Errorf("the `--include` pattern did not match any available assets:\n%s", strings.Join(names, "\n")))
	}

	args.NoForward()
}

func downloadReleaseAsset(asset github.ReleaseAsset, gh *github.Client) (err error) {
	assetReader, err := gh.DownloadReleaseAsset(asset.APIURL)
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
	tagName := ""
	if args.ParamsSize() > 0 {
		tagName = args.GetParam(0)
	}
	if tagName == "" {
		utils.Check(cmd.UsageError(""))
		return
	}

	assetsToUpload, close, err := openAssetFiles(args.Flag.AllValues("--attach"))
	utils.Check(err)
	defer close()

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	messageBuilder := &github.MessageBuilder{
		Filename: "RELEASE_EDITMSG",
		Title:    "release",
	}

	messageBuilder.AddCommentedSection(fmt.Sprintf(`Creating release %s for %s

Write a message for this release. The first block of
text is the title and the rest is the description.`, tagName, project))

	flagReleaseMessage := args.Flag.AllValues("--message")
	if len(flagReleaseMessage) > 0 {
		messageBuilder.Message = strings.Join(flagReleaseMessage, "\n\n")
		messageBuilder.Edit = args.Flag.Bool("--edit")
	} else if args.Flag.HasReceived("--file") {
		messageBuilder.Message, err = msgFromFile(args.Flag.Value("--file"))
		utils.Check(err)
		messageBuilder.Edit = args.Flag.Bool("--edit")
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
		TargetCommitish: args.Flag.Value("--commitish"),
		Name:            title,
		Body:            body,
		Draft:           args.Flag.Bool("--draft"),
		Prerelease:      args.Flag.Bool("--prerelease"),
	}

	var release *github.Release

	args.NoForward()
	if args.Noop {
		ui.Printf("Would create release `%s' for %s with tag name `%s'\n", title, project, tagName)
	} else {
		release, err = gh.CreateRelease(project, params)
		utils.Check(err)

		flagReleaseBrowse := args.Flag.Bool("--browse")
		flagReleaseCopy := args.Flag.Bool("--copy")
		printBrowseOrCopy(args, release.HTMLURL, flagReleaseBrowse, flagReleaseCopy)
	}

	messageBuilder.Cleanup()

	numAssets := len(assetsToUpload)
	if numAssets == 0 {
		return
	}
	if args.Noop {
		ui.Printf("Would attach %d %s\n", numAssets, pluralize(numAssets, "asset"))
	} else {
		ui.Errorf("Attaching %d %s...\n", numAssets, pluralize(numAssets, "asset"))
		uploaded, err := gh.UploadReleaseAssets(release, assetsToUpload)
		if err != nil {
			failed := []string{}
			for _, a := range assetsToUpload[len(uploaded):] {
				failed = append(failed, fmt.Sprintf("-a %s", a.Name))
			}
			ui.Errorf("The release was created, but attaching %d %s failed. ", len(failed), pluralize(len(failed), "asset"))
			ui.Errorf("You can retry with:\n%s release edit %s -m '' %s\n\n", "hub", release.TagName, strings.Join(failed, " "))
			utils.Check(err)
		}
	}
}

func editRelease(cmd *Command, args *Args) {
	tagName := ""
	if args.ParamsSize() > 0 {
		tagName = args.GetParam(0)
	}
	if tagName == "" {
		utils.Check(cmd.UsageError(""))
		return
	}

	assetsToUpload, close, err := openAssetFiles(args.Flag.AllValues("--attach"))
	utils.Check(err)
	defer close()

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	release, err := gh.FetchRelease(project, tagName)
	utils.Check(err)

	params := map[string]interface{}{}
	if args.Flag.HasReceived("--commitish") {
		params["target_commitish"] = args.Flag.Value("--commitish")
	}
	if args.Flag.HasReceived("--draft") {
		params["draft"] = args.Flag.Bool("--draft")
	}
	if args.Flag.HasReceived("--prerelease") {
		params["prerelease"] = args.Flag.Bool("--prerelease")
	}

	messageBuilder := &github.MessageBuilder{
		Filename: "RELEASE_EDITMSG",
		Title:    "release",
	}

	messageBuilder.AddCommentedSection(fmt.Sprintf(`Editing release %s for %s

Write a message for this release. The first block of
text is the title and the rest is the description.`, tagName, project))

	flagReleaseMessage := args.Flag.AllValues("--message")
	if len(flagReleaseMessage) > 0 {
		messageBuilder.Message = strings.Join(flagReleaseMessage, "\n\n")
		messageBuilder.Edit = args.Flag.Bool("--edit")
	} else if args.Flag.HasReceived("--file") {
		messageBuilder.Message, err = msgFromFile(args.Flag.Value("--file"))
		utils.Check(err)
		messageBuilder.Edit = args.Flag.Bool("--edit")
	} else {
		messageBuilder.Edit = true
		messageBuilder.Message = strings.Replace(fmt.Sprintf("%s\n\n%s", release.Name, release.Body), "\r\n", "\n", -1)
	}

	title, body, err := messageBuilder.Extract()
	utils.Check(err)

	if title == "" && len(flagReleaseMessage) == 0 {
		utils.Check(fmt.Errorf("Aborting editing due to empty release title"))
	}

	if title != "" {
		params["name"] = title
	}
	if body != "" {
		params["body"] = body
	}

	args.NoForward()
	if len(params) > 0 {
		if args.Noop {
			ui.Printf("Would edit release `%s'\n", tagName)
		} else {
			release, err = gh.EditRelease(release, params)
			utils.Check(err)
		}

		messageBuilder.Cleanup()
	}

	numAssets := len(assetsToUpload)
	if numAssets == 0 {
		return
	}
	if args.Noop {
		ui.Printf("Would attach %d %s\n", numAssets, pluralize(numAssets, "asset"))
	} else {
		ui.Errorf("Attaching %d %s...\n", numAssets, pluralize(numAssets, "asset"))
		uploaded, err := gh.UploadReleaseAssets(release, assetsToUpload)
		if err != nil {
			failed := []string{}
			for _, a := range assetsToUpload[len(uploaded):] {
				failed = append(failed, a.Name)
			}
			ui.Errorf("Attaching these assets failed:\n%s\n\n", strings.Join(failed, "\n"))
			utils.Check(err)
		}
	}
}

func deleteRelease(cmd *Command, args *Args) {
	tagName := ""
	if args.ParamsSize() > 0 {
		tagName = args.GetParam(0)
	}
	if tagName == "" {
		utils.Check(cmd.UsageError(""))
		return
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
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

func openAssetFiles(args []string) ([]github.LocalAsset, func(), error) {
	assets := []github.LocalAsset{}
	files := []*os.File{}

	for _, arg := range args {
		var label string
		parts := strings.SplitN(arg, "#", 2)
		path := parts[0]
		if len(parts) > 1 {
			label = parts[1]
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, nil, err
		}
		stat, err := file.Stat()
		if err != nil {
			return nil, nil, err
		}
		files = append(files, file)

		assets = append(assets, github.LocalAsset{
			Name:     path,
			Label:    label,
			Size:     stat.Size(),
			Contents: file,
		})
	}

	close := func() {
		for _, f := range files {
			f.Close()
		}
	}

	return assets, close, nil
}

func pluralize(count int, label string) string {
	if count == 1 {
		return label
	}
	return fmt.Sprintf("%ss", label)
}
