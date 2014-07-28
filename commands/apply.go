package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/github/hub/github"
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

func init() {
	CmdRunner.Use(cmdApply)
	CmdRunner.Use(cmdAm)
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
func apply(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformApplyArgs(args)
	}
}

func transformApplyArgs(args *Args) {
	urlRegexp := regexp.MustCompile("^https?://(gist\\.)github\\.com/")
	pullRegexp := regexp.MustCompile("^(pull|commit)/([0-9a-f]+)")
	for _, arg := range args.Params {
		var (
			url  string
			gist bool
		)
		projectURL, err := github.ParseURL(arg)
		if err == nil {
			match := pullRegexp.FindStringSubmatch(projectURL.ProjectPath())
			if match != nil {
				url = projectURL.Project.WebURL("", "", match[1]+"/"+match[2])
			}
		} else {
			gist = urlRegexp.MatchString(arg)
			if gist {
				url = arg
			}
		}

		if url == "" {
			continue
		}

		var ext string
		if gist {
			ext = ".txt"
		} else {
			ext = ".patch"
		}

		idx := args.IndexOfParam(arg)
		if filepath.Ext(url) != ext {
			url += ext
		}

		var prefix string
		if gist {
			prefix = "gist-"
		}

		patchFile := filepath.Join(os.TempDir(), prefix+filepath.Base(url))

		args.Before("curl", "-#LA", fmt.Sprintf("gh %s", Version), url, "-o", patchFile)
		args.Params[idx] = patchFile
	}
}
