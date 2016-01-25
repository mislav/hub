package commands

import (
	"io"
	"io/ioutil"
	"regexp"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdApply = &Command{
	Run:          apply,
	GitExtension: true,
	Usage:        "apply <GITHUB-URL>",
	Long: `Download a patch from GitHub and apply it locally.

## Options:
	<GITHUB-URL>
		A URL to a pull request or commit on GitHub.

## Examples:
		$ hub apply https://github.com/jingweno/gh/pull/55
		> curl https://github.com/jingweno/gh/pull/55.patch -o /tmp/55.patch
		> git apply /tmp/55.patch

## See also:

hub-am(1), hub(1), git-apply(1)
`,
}

var cmdAm = &Command{
	Run:          apply,
	GitExtension: true,
	Usage:        "am [-3] <GITHUB-URL>",
	Long: `Replicate commits from a GitHub pull request locally.

## Options:
	-3
		(Recommended) See git-am(1).

	<GITHUB-URL>
		A URL to a pull request or commit on GitHub.

## Examples:
		$ hub am -3 https://github.com/jingweno/gh/pull/55
		> curl https://github.com/jingweno/gh/pull/55.patch -o /tmp/55.patch
		> git am -3 /tmp/55.patch

## See also:

hub-apply(1), hub-cherry-pick(1), hub(1), git-am(1)
`,
}

func init() {
	CmdRunner.Use(cmdApply)
	CmdRunner.Use(cmdAm)
}

func apply(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformApplyArgs(args)
	}
}

func transformApplyArgs(args *Args) {
	gistRegexp := regexp.MustCompile("^https?://gist\\.github\\.com/([\\w.-]+/)?([a-f0-9]+)")
	pullRegexp := regexp.MustCompile("^(pull|commit)/([0-9a-f]+)")
	for _, arg := range args.Params {
		var (
			patch    io.ReadCloser
			apiError error
		)
		projectURL, err := github.ParseURL(arg)
		if err == nil {
			gh := github.NewClient(projectURL.Project.Host)
			match := pullRegexp.FindStringSubmatch(projectURL.ProjectPath())
			if match != nil {
				if match[1] == "pull" {
					patch, apiError = gh.PullRequestPatch(projectURL.Project, match[2])
				} else {
					patch, apiError = gh.CommitPatch(projectURL.Project, match[2])
				}
			}
		} else {
			match := gistRegexp.FindStringSubmatch(arg)
			if match != nil {
				// TODO: support Enterprise gist
				gh := github.NewClient(github.GitHubHost)
				patch, apiError = gh.GistPatch(match[2])
			}
		}

		utils.Check(apiError)
		if patch == nil {
			continue
		}

		idx := args.IndexOfParam(arg)
		patchFile, err := ioutil.TempFile("", "hub")
		utils.Check(err)

		_, err = io.Copy(patchFile, patch)
		utils.Check(err)

		patchFile.Close()
		patch.Close()

		args.Params[idx] = patchFile.Name()
	}
}
