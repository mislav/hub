package commands

import (
//"fmt"
//"github.com/jingweno/gh/utils"
//"github.com/jingweno/octokat"
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

func apply(command *Command, args *Args) {
}
