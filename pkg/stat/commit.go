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
	"github.com/google/go-github/v44/github"
)

type CommitList []*github.RepositoryCommit

func (c CommitList) Chart() Chart {
	var (
		labels  []string
		data    []float64
		now     = time.Now()
		dayTime = timex.AllDays(now.Add(-weekDur), now)
	)

	for _, t := range dayTime {
		label := t.Format(labelLayout)
		labels = append(labels, label)
		data = append(data, float64(c.getSpecifiedDate(t)))
	}

	return Chart{Data: data, Labels: labels}
}

func (c CommitList) getSpecifiedDate(date time.Time) int {
	var (
		count int
		zero  = timex.Truncate(date)
	)

	for _, e := range c {
		commit := e.Commit
		if commit == nil {
			continue
		}

		committer := commit.Author
		if committer == nil {
			continue
		}

		if timex.Truncate(committer.GetDate()).Equal(zero) {
			count += 1
		}
	}

	return count
}

func (s Stat) latestWeekCommits() CommitList {
	var (
		page  = 1
		list  CommitList
		until = time.Now()
		since = time.Now().Add(-timeWeek)
	)

	for {
		ret, resp, err := s.restClient.Repositories.ListCommits(s.ctx, s.owner, s.repo,
			&github.CommitsListOptions{
				Since: since,
				Until: until,
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: 100,
				},
			})
		if err != nil {
			return list
		}

		list = append(list, ret...)
		if page >= resp.LastPage {
			return list
		}

		page = resp.NextPage
	}
}
