package stat

import (
	"github.com/shurcooL/githubv4"
)

type (
	IssueConnection struct {
		TotalCount githubv4.Int
	}

	Issue struct {
		Issues IssueConnection `graphql:"issues(first: 1, states: $issueStates)"`
	}

	IssueQuery struct {
		Issue Issue `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (s Stat) OpenIssueCount() githubv4.Int {
	var issueQuery IssueQuery
	_ = s.graphqlClient.Query(s.ctx, &issueQuery, map[string]interface{}{
		"owner":       githubv4.String(s.owner),
		"name":        githubv4.String(s.repo),
		"issueStates": []githubv4.IssueState{githubv4.IssueStateOpen},
	})
	return issueQuery.Issue.Issues.TotalCount
}
