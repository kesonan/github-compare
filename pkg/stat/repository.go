/*
 * MIT License
 *
 * Copyright (c) 2022 anqiansong
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
