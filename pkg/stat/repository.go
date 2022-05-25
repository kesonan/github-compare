package stat

import (
	"github.com/shurcooL/githubv4"
)

type (
	Repository struct {
		CreatedAt       githubv4.DateTime
		ForkCount       githubv4.Int
		HomepageUrl     githubv4.URI
		Issues          IssueConnection `graphql:"issues(first: 1, states: $issueStates)"`
		LatestRelease   Release
		LicenseInfo     License
		PrimaryLanguage Language
		NameWithOwner   githubv4.String
		PullRequests    PullRequestConnection `graphql:"pullRequests(first: 1, states: $pullRequestStates)"`
		PushedAt        githubv4.DateTime
		Releases        ReleaseConnection `graphql:"releases(first: 1, orderBy: $orderBy)"`
		StargazerCount  githubv4.Int
		UpdatedAt       githubv4.DateTime
		Watchers        UserConnection `graphql:"watchers(first: 1)"`
	}

	RepositoryQuery struct {
		Repository Repository `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (s Stat) Repository() Repository {
	var repositoryQuery RepositoryQuery
	_ = s.graphqlClient.Query(s.ctx, &repositoryQuery, map[string]interface{}{
		"owner": githubv4.String(s.owner),
		"name":  githubv4.String(s.repo),
		"orderBy": githubv4.ReleaseOrder{
			Field:     githubv4.ReleaseOrderFieldCreatedAt,
			Direction: githubv4.OrderDirectionDesc,
		},
		"issueStates": []githubv4.IssueState{githubv4.IssueStateOpen,
			githubv4.IssueStateClosed},
		"pullRequestStates": []githubv4.PullRequestState{githubv4.PullRequestStateOpen,
			githubv4.PullRequestStateClosed, githubv4.PullRequestStateMerged},
	})
	return repositoryQuery.Repository
}
