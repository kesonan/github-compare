package stat

import (
	"github.com/shurcooL/githubv4"
)

type (
	PullRequestConnection struct {
		TotalCount githubv4.Int
	}

	PullRequest struct {
		PullRequests PullRequestConnection `graphql:"pullRequests(first: 1, states: $pullRequestStates)"`
	}

	PRQuery struct {
		PullRequest PullRequest `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (s Stat) OpenPullRequestCount() githubv4.Int {
	var prQuery PRQuery
	_ = s.graphqlClient.Query(s.ctx, &prQuery, map[string]interface{}{
		"owner":             githubv4.String(s.owner),
		"name":              githubv4.String(s.repo),
		"pullRequestStates": []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
	})
	return prQuery.PullRequest.PullRequests.TotalCount
}
