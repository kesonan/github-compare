// MIT License
//
// Copyright (c) 2022 anqiansong
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package stat

import (
	"time"

	"github.com/anqiansong/github-compare/pkg/timex"
	"github.com/shurcooL/githubv4"
)

type (
	PullRequestList []PullRequestEdge

	PullRequestNode struct {
		CreatedAt githubv4.DateTime
	}

	PullRequestEdge struct {
		Cursor githubv4.String
		Node   PullRequestNode
	}

	PullRequestConnection struct {
		Edges      []PullRequestEdge
		PageInfo   PageInfo
		TotalCount githubv4.Int
	}

	PullRequest struct {
		List PullRequestConnection `graphql:"pullRequests(first: $first, orderBy: $orderBy, states: $pullRequestStates)"`
	}

	PRQuery struct {
		PullRequest PullRequest `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (s Stat) OpenPullRequestCount() githubv4.Int {
	var prQuery PRQuery
	_ = s.graphqlClient.Query(s.ctx, &prQuery, map[string]interface{}{
		"after":             (*githubv4.String)(nil),
		"owner":             githubv4.String(s.owner),
		"name":              githubv4.String(s.repo),
		"first":             1,
		"pullRequestStates": []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
		"orderBy": githubv4.PullRequestOrder{
			Field:     githubv4.PullRequestOrderFieldCreatedAt,
			Direction: githubv4.OrderDirectionDesc,
		},
	})

	return prQuery.PullRequest.List.TotalCount
}

func (p PullRequestList) Chart() Chart {
	now := time.Now()
	var (
		dayCount = make(map[string]int)
		labels   []string
		data     []float64
		dayTime  = timex.AllDays(now, now.Add(7*24*time.Hour))
	)

	for _, t := range dayTime {
		label := t.Format("02/01")
		labels = append(labels, label)
		dayCount[label] += p.getSpecifiedDate(t)
	}

	return Chart{Data: data, Labels: labels}
}

func (p PullRequestList) getSpecifiedDate(date time.Time) int {
	zero := timex.Truncate(date)
	var count int
	for _, e := range p {
		if timex.Truncate(e.Node.CreatedAt.Time).Equal(zero) {
			count += 1
		}
	}
	return count
}

func (s Stat) latestWeekPRS() PullRequestList {
	var (
		list    PullRequestList
		brk     bool
		after   githubv4.String
		prQuery PRQuery
	)

	deadline := time.Now().Add(-7 * 24 * time.Hour)
	arg := map[string]interface{}{
		"after": (*githubv4.String)(nil),
		"owner": githubv4.String(s.owner),
		"name":  githubv4.String(s.repo),
		"first": 100,
		"pullRequestStates": []githubv4.PullRequestState{githubv4.PullRequestStateOpen,
			githubv4.PullRequestStateClosed, githubv4.PullRequestStateMerged},
		"orderBy": githubv4.PullRequestOrder{
			Field:     githubv4.PullRequestOrderFieldCreatedAt,
			Direction: githubv4.OrderDirectionDesc,
		},
	}

	for {
		_ = s.graphqlClient.Query(s.ctx, &prQuery, arg)
		temp := prQuery.PullRequest.List.Edges

		for _, e := range temp {
			if e.Node.CreatedAt.Time.Before(deadline) {
				brk = true
				break
			}
			list = append(list, e)
		}
		if brk || !(bool)(prQuery.PullRequest.List.PageInfo.HasNextPage) || len(temp) == 0 {
			break
		}

		after = temp[len(temp)-1].Cursor
		arg["after"] = after
	}

	return list
}
