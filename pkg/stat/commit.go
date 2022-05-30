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
	CommitList []CommitCommentEdge

	CommitComment struct {
		CreatedAt githubv4.DateTime
	}

	CommitCommentEdge struct {
		Cursor githubv4.String
		Node   ForkRepository
	}

	CommitCommentConnection struct {
		Edges      []CommitCommentEdge
		PageInfo   PageInfo
		TotalCount githubv4.Int
	}

	Commit struct {
		List CommitCommentConnection `graphql:"commitComments(first: 100, after: $after)"`
	}

	CommitQuery struct {
		Commits Commit `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (c CommitList) Chart() Chart {
	now := time.Now()
	var (
		dayCount = make(map[string]int)
		labels   []string
		data     []float64
		dayTime  = timex.AllDays(now.Add(-7*24*time.Hour), now)
	)

	for _, t := range dayTime {
		label := t.Format("02/01")
		labels = append(labels, label)
		dayCount[label] += c.getSpecifiedDate(t)
	}

	return Chart{Data: data, Labels: labels}
}

func (c CommitList) getSpecifiedDate(date time.Time) int {
	zero := timex.Truncate(date)
	var count int
	for _, e := range c {
		if timex.Truncate(e.Node.CreatedAt.Time).Equal(zero) {
			count += 1
		}
	}
	return count
}

func (s Stat) latestWeekCommits() CommitList {
	var (
		list        CommitList
		brk         bool
		after       githubv4.String
		commitQuery CommitQuery
	)

	deadline := time.Now().Add(-7 * 24 * time.Hour)
	arg := map[string]interface{}{
		"after": (*githubv4.String)(nil),
		"owner": githubv4.String(s.owner),
		"name":  githubv4.String(s.repo),
	}

	for {
		_ = s.graphqlClient.Query(s.ctx, &commitQuery, arg)
		temp := commitQuery.Commits.List.Edges

		for _, e := range temp {
			if e.Node.CreatedAt.Time.Before(deadline) {
				brk = true
				break
			}
			list = append(list, e)
		}
		if brk || !(bool)(commitQuery.Commits.List.PageInfo.HasNextPage) || len(temp) == 0 {
			break
		}

		after = temp[len(temp)-1].Cursor
		arg["after"] = after
	}

	return list
}
