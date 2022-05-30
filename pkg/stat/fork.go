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
	Forks []RepositoryEdge

	ForkRepository struct {
		CreatedAt githubv4.DateTime
	}

	RepositoryEdge struct {
		Cursor githubv4.String
		Node   ForkRepository
	}

	RepositoryConnection struct {
		Edges      []RepositoryEdge
		PageInfo   PageInfo
		TotalCount githubv4.Int
	}

	Fork struct {
		List RepositoryConnection `graphql:"forks(first: 100, orderBy: $orderBy, after: $after)"`
	}

	ForkQuery struct {
		Forks Fork `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (f Forks) Chart() Chart {
	now := time.Now()
	var (
		labels  []string
		data    []float64
		dayTime = timex.AllDays(now.Add(-weekDur), now)
	)

	for _, t := range dayTime {
		label := t.Format(labelLayout)
		labels = append(labels, label)
		data = append(data, float64(f.getSpecifiedDate(t)))
	}

	return Chart{Data: data, Labels: labels}
}

func (f Forks) getSpecifiedDate(date time.Time) int {
	zero := timex.Truncate(date)
	var count int
	for _, e := range f {
		if timex.Truncate(e.Node.CreatedAt.Time).Equal(zero) {
			count += 1
		}
	}
	return count
}

func (s Stat) latestWeekForks() Forks {
	var (
		list      Forks
		brk       bool
		after     githubv4.String
		forkQuery ForkQuery
	)

	deadline := time.Now().Add(-7 * 24 * time.Hour)
	arg := map[string]interface{}{
		"after": (*githubv4.String)(nil),
		"owner": githubv4.String(s.owner),
		"name":  githubv4.String(s.repo),
		"orderBy": githubv4.RepositoryOrder{
			Field:     githubv4.RepositoryOrderFieldCreatedAt,
			Direction: githubv4.OrderDirectionDesc,
		},
	}

	for {
		_ = s.graphqlClient.Query(s.ctx, &forkQuery, arg)
		temp := forkQuery.Forks.List.Edges

		for _, e := range temp {
			if e.Node.CreatedAt.Time.Before(deadline) {
				brk = true
				break
			}
			list = append(list, e)
		}
		if brk || !(bool)(forkQuery.Forks.List.PageInfo.HasNextPage) || len(temp) == 0 {
			break
		}

		after = temp[len(temp)-1].Cursor
		arg["after"] = after
	}

	return list
}
