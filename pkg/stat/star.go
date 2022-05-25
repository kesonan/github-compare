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
	"time"

	"github.com/shurcooL/githubv4"
)

type (
	StargazerEdges []StargazerEdge

	StargazerEdge struct {
		Cursor    githubv4.String
		StarredAt githubv4.DateTime
	}

	StargazerConnection struct {
		Edges      []StargazerEdge
		PageInfo   PageInfo
		TotalCount githubv4.Int
	}

	Stargazer struct {
		Stargazers StargazerConnection `graphql:"stargazers(first: 100, orderBy: $orderBy, after: $after)"`
	}

	StargazerQuery struct {
		Stargazer Stargazer `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (s StargazerEdges) LatestDayStars() (int, int) {
	y, m, d := time.Now().Date()
	deadlineOfToday := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	deadlineOfYesterday := deadlineOfToday.AddDate(0, 0, -1)
	var starsOfToday, starsOfYesterday int
	for _, e := range s {
		if e.StarredAt.Time.After(deadlineOfToday) {
			starsOfToday += 1
		}
		if e.StarredAt.Time.Before(deadlineOfToday) && e.StarredAt.Time.After(deadlineOfYesterday) {
			starsOfYesterday += 1
		}

	}
	return starsOfToday, starsOfToday - starsOfYesterday
}

func (s StargazerEdges) LatestWeekStars() (int, int) {
	deadlineOfLatest7Days := time.Now().AddDate(0, 0, -7)
	deadlineOfPre7Days := deadlineOfLatest7Days.AddDate(0, 0, -7)
	var starsOfLatest7Days, starsOfPre7Days int
	for _, e := range s {
		if e.StarredAt.Time.After(deadlineOfLatest7Days) {
			starsOfLatest7Days += 1
		}
		if e.StarredAt.Time.Before(deadlineOfLatest7Days) && e.StarredAt.Time.After(deadlineOfPre7Days) {
			starsOfPre7Days += 1
		}
	}

	return starsOfLatest7Days, starsOfLatest7Days - starsOfPre7Days
}

func (s StargazerEdges) LatestMonthStars() int {
	return len(s)
}

func (s Stat) latestMonthStargazers() StargazerEdges {
	var (
		list           []StargazerEdge
		brk            bool
		after          githubv4.String
		stargazerQuery StargazerQuery
	)

	deadline := time.Now().AddDate(0, -1, 0)
	arg := map[string]interface{}{
		"after": (*githubv4.String)(nil),
		"owner": githubv4.String(s.owner),
		"name":  githubv4.String(s.repo),
		"orderBy": githubv4.StarOrder{
			Field:     githubv4.StarOrderFieldStarredAt,
			Direction: githubv4.OrderDirectionDesc,
		},
	}

	for {
		_ = s.graphqlClient.Query(s.ctx, &stargazerQuery, arg)
		temp := stargazerQuery.Stargazer.Stargazers.Edges

		for _, e := range temp {
			if e.StarredAt.Time.Before(deadline) {
				brk = true
				break
			}
			list = append(list, e)
		}
		if brk || !(bool)(stargazerQuery.Stargazer.Stargazers.PageInfo.HasNextPage) || len(temp) == 0 {
			break
		}

		after = temp[len(temp)-1].Cursor
		arg["after"] = after
	}

	return list
}
