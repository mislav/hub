package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"os"
	"strings"
)

var cmdRelease = &Command{
	Run:   release,
	Usage: "release",
	Short: "Manipulate releases on GitHub",
	Long: `Manipulates releases on GitHub for the project that the "origin" remote
points to.
`,
}

func release(cmd *Command, args *Args) {
	gh := github.New()
	if args.Noop {
		fmt.Printf("Would request list of releases for %s\n", gh.Project)
	} else {
		releases, err := gh.Releases()
		utils.Check(err)
		var outputs []string
		for _, release := range releases {
			out := fmt.Sprintf("%s (%s)\n%s", release.Name, release.TagName, release.Body)
			outputs = append(outputs, out)
		}

		fmt.Println(strings.Join(outputs, "\n\n"))
	}

	os.Exit(0)
}
