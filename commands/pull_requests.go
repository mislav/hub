package commands

import (
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
	"github.com/github/hub/ui"
	"strings"
	"strconv"
	"time"
	"math"
	"fmt"
)

var cmdListPullRequest = &Command{
	Run: listPullRequest,
	Usage: `
pull-requests [-a] [-o <ORG>]
`,
	Long: `List GitHub pull requests.

## Options:
	-a, --all
		To list all pull requests for all repositories in the current organisation

	-o, --org <ORG>
		To list pull requests for all repositories in the given github organisation

## See also:

hub(1), hub-pull-request(1), hub-merge(1), hub-checkout(1)
`,
}

var (
	flagPullRequestOrganisation string
	flagPullRequestAll bool
)

func init() {
	cmdListPullRequest.Flag.StringVarP(&flagPullRequestOrganisation, "org", "o", "", "ORG")
	cmdListPullRequest.Flag.BoolVarP(&flagPullRequestAll, "all", "a", false, "ALL")

	CmdRunner.Use(cmdListPullRequest)
}

func listPullRequest(cmd *Command, args *Args) {
	runInMainOrCurrentProject(func(localRepo *github.GitHubRepo, project *github.Project, gh *github.Client) {
		organisation := flagPullRequestOrganisation
		if flagPullRequestAll && len(organisation) == 0 {
			organisation = project.Owner
		}

		if len(organisation) == 0 {
			pullRequests, err := gh.PullRequests(project)
			utils.Check(err)
			for _, pr := range pullRequests {
				url := pr.HTMLURL
				if flagIssueAssignee == "" ||
					strings.EqualFold(pr.Assignee.Login, flagIssueAssignee) {
					ui.Printf("% 7d] %s ( %s ) age: %s\n", pr.Number, pr.Title, url, durationText(pr.CreatedAt))
				}
			}
		} else {
			ui.Printf("Pull requests for organisation: %s\n", organisation)
			repos, err := gh.OrgRepositories(organisation)
			utils.Check(err)
			maxLen := 1
			for _, repo := range repos {
				l := len(repo.Name)
				if l > maxLen {
					maxLen = l
				}
			}
			for _, repo := range repos {
				name := repo.Name
				repoProject := &github.Project{
					Owner: organisation,
					Name: name,
					Protocol: project.Protocol,
					Host: project.Host,
				}
				pullRequests, err := gh.PullRequests(repoProject)
				utils.Check(err)
				for _, pr := range pullRequests {
					url := pr.HTMLURL
					if flagIssueAssignee == "" ||
						strings.EqualFold(pr.Assignee.Login, flagIssueAssignee) {
						ui.Printf("%" + strconv.Itoa(maxLen) + "s | % 7d] %s ( %s ) age: %s\n", name, pr.Number, pr.Title, url, durationText(pr.CreatedAt))
					}
				}
			}
		}
	})
}

func durationText(t time.Time) string {
	s := time.Since(t)
	hours := math.Floor(s.Hours())
	mins := s.Minutes() - (hours * 60)
	days := math.Floor(hours / 24)
	hours = hours - (days * 24)
	return fmt.Sprintf("%0.0f days %0.0f:%02.1f hours", days, hours, mins)
}
