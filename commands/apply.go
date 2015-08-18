package commands

import (
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdApply = &Command{
	Run:          apply,
	GitExtension: true,
	Usage:        "apply GITHUB-URL",
	Short:        "Apply a patch to files and/or to the index",
	Long: `Downloads the patch file for the pull request or commit at the URL and
applies that patch from disk with git am or git apply. Similar to
cherry-pick, but doesn't add new remotes. git am creates commits while
preserving authorship info while <code>apply</code> only applies the
patch to the working copy.
`,
}

var cmdAm = &Command{
	Run:          apply,
	GitExtension: true,
	Usage:        "am GITHUB-URL",
	Short:        "Apply a patch to files and/or to the index",
	Long: `Downloads the patch file for the pull request or commit at the URL and
applies that patch from disk with git am or git apply. Similar to
cherry-pick, but doesn't add new remotes. git am creates commits while
preserving authorship info while <code>apply</code> only applies the
patch to the working copy.
`,
}

var cmdApr = &Command{
	Run:          applyPullRequest,
	GitExtension: true,
	Usage:        "apr PR-NUMBER",
	Short:        "Apply a patch to files and/or to the index",
	Long: `Downloads the patch file for the pull request and
applies that patch from disk with git am or git apply. Similar to
cherry-pick, but doesn't add new remotes. git am creates commits while
preserving authorship info while <code>apply</code> only applies the
patch to the working copy.
`,
}

func init() {
	CmdRunner.Use(cmdApply)
	CmdRunner.Use(cmdAm)
	CmdRunner.Use(cmdApr)
}

/*
  $ gh apply https://github.com/jingweno/gh/pull/55
  > curl https://github.com/jingweno/gh/pull/55.patch -o /tmp/55.patch
  > git apply /tmp/55.patch

  $ git apply --ignore-whitespace https://github.com/jingweno/gh/commit/fdb9921
  > curl https://github.com/jingweno/gh/commit/fdb9921.patch -o /tmp/fdb9921.patch
  > git apply --ignore-whitespace /tmp/fdb9921.patch

  $ git apply https://gist.github.com/8da7fb575debd88c54cf
  > curl https://gist.github.com/8da7fb575debd88c54cf.txt -o /tmp/gist-8da7fb575debd88c54cf.txt
  > git apply /tmp/gist-8da7fb575debd88c54cf.txt
*/
func applyPullRequest(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformAprArgs(args)
		transformApplyArgs(args)
	}
}

func apply(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformApplyArgs(args)
	}
}

func transformAprArgs(args *Args) {
	length := args.ParamsSize()
	if length > 0 {
		remote, err := git.Config("remote.origin.url")
		utils.Check(err)
		urlBase := strings.TrimSuffix(strings.TrimSuffix(remote, ".git"), "/")
		args.ReplaceParam(length-1, urlBase+"/pull/"+strings.TrimPrefix(args.GetParam(length-1), "#"))
		args.InsertParam(length-1, "-3")
		args.Command = "am"
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
