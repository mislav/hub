package github

import (
	"encoding/json"

	"github.com/azaky/graphb"
)

func FetchPullRequestsCheckStatusQuery(project *Project, filterParams map[string]interface{}, limit int) (queryObject map[string]interface{}, err error) {
	pullRequestArguments := []graphb.Argument{}
	if filterParams != nil {
		if filterParams["state"] != nil {
			switch filterParams["state"].(string) {
			case "open":
				pullRequestArguments = append(pullRequestArguments, graphb.ArgumentEnum("states", "OPEN"))
			case "closed":
				pullRequestArguments = append(pullRequestArguments, graphb.ArgumentEnum("states", "CLOSED"))
			case "merged":
				pullRequestArguments = append(pullRequestArguments, graphb.ArgumentEnum("states", "MERGED"))
			case "all":
				pullRequestArguments = append(pullRequestArguments, graphb.ArgumentEnumSlice("states", "OPEN", "CLOSED", "MERGED"))
			}
		} else {
			pullRequestArguments = append(pullRequestArguments, graphb.ArgumentEnum("states", "OPEN"))
		}

		if filterParams["head"] != nil {
			pullRequestArguments = append(pullRequestArguments, graphb.ArgumentString("headRefName", filterParams["head"].(string)))
		}

		if filterParams["base"] != nil {
			pullRequestArguments = append(pullRequestArguments, graphb.ArgumentString("baseRefName", filterParams["base"].(string)))
		}

		var sortArgument graphb.Argument
		if filterParams["sort"] != nil {
			switch filterParams["sort"].(string) {
			case "created":
				sortArgument = graphb.ArgumentEnum("field", "CREATED_AT")
			case "updated":
				sortArgument = graphb.ArgumentEnum("field", "UPDATED_AT")
			case "popularity":
				sortArgument = graphb.ArgumentEnum("field", "COMMENTS")
			}
		} else {
			sortArgument = graphb.ArgumentEnum("field", "CREATED_AT")
		}

		var directionArgument graphb.Argument
		if filterParams["direction"] != nil {
			switch filterParams["direction"].(string) {
			case "asc":
				directionArgument = graphb.ArgumentEnum("direction", "ASC")
			case "desc":
				directionArgument = graphb.ArgumentEnum("direction", "DESC")
			}
		}

		issueOrderArgument := graphb.ArgumentCustomType("orderBy", sortArgument, directionArgument)
		pullRequestArguments = append(pullRequestArguments, issueOrderArgument)

		if limit != 0 {
			pullRequestArguments = append(pullRequestArguments, graphb.ArgumentInt("first", limit))
		} else {
			pullRequestArguments = append(pullRequestArguments, graphb.ArgumentInt("first", 100))
		}
	}

	graphbQuery := graphb.MakeQuery(graphb.TypeQuery).
		SetName("").
		SetFields(graphb.MakeField("repository").
			SetArguments(
				graphb.ArgumentString("owner", project.Owner),
				graphb.ArgumentString("name", project.Name),
			).
			SetFields(graphb.MakeField("pullRequests").
				SetArguments(pullRequestArguments...).
				SetFields(graphb.MakeField("nodes").
					SetFields(
						graphb.MakeField("number"),
						graphb.MakeField("title"),
						graphb.MakeField("commits").
							SetArguments(graphb.ArgumentInt("last", 1)).
							SetFields(graphb.MakeField("nodes").
								SetFields(graphb.MakeField("commit").
									SetFields(
										graphb.MakeField("status").
											SetFields(graphb.MakeField("state")),
										graphb.MakeField("checkSuites").SetArguments(graphb.ArgumentInt("first", 50)).
											SetFields(graphb.MakeField("nodes").
												SetFields(graphb.MakeField("conclusion")),
											),
									),
								),
							),
					),
				),
			),
		)

	query, err := graphbQuery.JSON()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(query), &queryObject)
	if err != nil {
		return
	}

	return
}
