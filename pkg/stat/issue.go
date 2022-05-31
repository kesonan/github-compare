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
	IssueList []IssueEdge

	IssueEdge struct {
		Cursor githubv4.String
		Node   Issue
	}

	IssueConnection struct {
		Edges      []IssueEdge
		PageInfo   PageInfo
		TotalCount githubv4.Int
	}

	Issue struct {
		CreatedAt githubv4.DateTime
	}

	Issues struct {
		List IssueConnection `graphql:"issues(first: $first, orderBy: $orderBy, after: $after, states: $issueStates)"`
	}

	IssueQuery struct {
		Issue Issues `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (i IssueList) Chart() Chart {
	var (
		labels  []string
		data    []float64
		now     = time.Now()
		dayTime = timex.AllDays(now.Add(-weekDur), now)
	)

	for _, t := range dayTime {
		label := t.Format(labelLayout)
		labels = append(labels, label)
		data = append(data, float64(i.getSpecifiedDate(t)))
	}

	return Chart{Data: data, Labels: labels}
}

func (i IssueList) getSpecifiedDate(date time.Time) int {
	var (
		count int
		zero  = timex.Truncate(date)
	)

	for _, e := range i {
		if timex.Truncate(e.Node.CreatedAt.Time).Equal(zero) {
			count += 1
		}
	}

	return count
}

func (s Stat) OpenIssueCount() githubv4.Int {
	var issueQuery IssueQuery
	_ = s.graphqlClient.Query(s.ctx, &issueQuery, map[string]interface{}{
		"after":       (*githubv4.String)(nil),
		"owner":       githubv4.String(s.owner),
		"name":        githubv4.String(s.repo),
		"first":       1,
		"issueStates": []githubv4.IssueState{githubv4.IssueStateOpen},
		"orderBy": githubv4.IssueOrder{
			Field:     githubv4.IssueOrderFieldCreatedAt,
			Direction: githubv4.OrderDirectionDesc,
		},
	})

	return issueQuery.Issue.List.TotalCount
}

func (s Stat) LatestWeekIssues() IssueList {
	var (
		list       IssueList
		brk        bool
		issueQuery IssueQuery
		after      githubv4.String
		deadline   = time.Now().Add(-timeWeek)
	)

	arg := map[string]interface{}{
		"after":       (*githubv4.String)(nil),
		"owner":       githubv4.String(s.owner),
		"name":        githubv4.String(s.repo),
		"first":       githubv4.Int(100),
		"issueStates": []githubv4.IssueState{githubv4.IssueStateOpen, githubv4.IssueStateClosed},
		"orderBy": githubv4.IssueOrder{
			Field:     githubv4.IssueOrderFieldCreatedAt,
			Direction: githubv4.OrderDirectionDesc,
		},
	}

	for {
		_ = s.graphqlClient.Query(s.ctx, &issueQuery, arg)
		temp := issueQuery.Issue.List.Edges

		for _, e := range temp {
			if e.Node.CreatedAt.Time.Before(deadline) {
				brk = true
				break
			}
			list = append(list, e)
		}
		if brk || !(bool)(issueQuery.Issue.List.PageInfo.HasNextPage) || len(temp) == 0 {
			break
		}

		after = temp[len(temp)-1].Cursor
		arg["after"] = after
	}

	return list
}
