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
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/kevwan/mapreduce/v2"
	"github.com/shurcooL/githubv4"
)

type Data struct {
	Age                  string `json:"age"`
	AvgReleasePeriod     string `json:"avgReleasePeriod,omitempty"`
	ContributorCount     string `json:"contributorCount,omitempty"`
	ForkCount            string `json:"forkCount,omitempty"`
	FullName             string `json:"fullName,omitempty"`
	Homepage             string `json:"homepage,omitempty"`
	Issue                string `json:"issue"`
	Language             string `json:"language,omitempty"`
	LastPushedAt         string `json:"lastPushedAt"`
	LatestReleaseAt      string `json:"latestReleaseAt"`
	LastUpdatedAt        string `json:"lastUpdatedAt"`
	LatestDayStarCount   string `json:"latestDayStarCount"`
	LatestMonthStarCount string `json:"latestMonthStarCount"`
	LatestWeekStarCount  string `json:"latestWeekStarCount"`
	License              string `json:"license,omitempty"`
	Pull                 string `json:"pull"`
	ReleaseCount         string `json:"releaseCount,omitempty"`
	StarCount            string `json:"starCount,omitempty"`
	WatcherCount         string `json:"watcherCount,omitempty"`
}

func Overview(accessToken string, repos ...string) []Data {
	reduce, _ := mapreduce.MapReduce(func(source chan<- *Stat) {
		for _, r := range repos {
			s := NewStat(r, accessToken)
			source <- s
		}
	}, func(s *Stat, writer mapreduce.Writer[[]Data], cancel func(error)) {
		var (
			repo                  Repository
			openIssueCount        githubv4.Int
			openPrCount           githubv4.Int
			contributorCount      int
			latestMonthStargazers StargazerEdges
			list                  []Data
		)
		mapreduce.FinishVoid(func() {
			repo = s.Repository()
		}, func() {
			openIssueCount = s.OpenIssueCount()
		}, func() {
			openPrCount = s.OpenPullRequestCount()
		}, func() {
			contributorCount = s.ContributorCount()
		}, func() {
			latestMonthStargazers = s.latestMonthStargazers()
		})

		releaseCount := repo.Releases.TotalCount
		totalStarCount := int(repo.StargazerCount)
		avgStarCount := totalStarCount
		totalForkCount := int(repo.ForkCount)
		avgForkCount := totalForkCount
		avgReleasePeriod := time.Duration(0)
		ageDuration := time.Since(repo.CreatedAt.Time)
		ageDays := int(ageDuration.Hours() / 24)
		if releaseCount > 0 {
			avgReleasePeriod = ageDuration / time.Duration(releaseCount)
		}
		if ageDays > 1 {
			avgStarCount = totalStarCount / ageDays
			avgForkCount = totalForkCount / ageDays
		}

		list = append(list, Data{
			FullName:             formatValue(repo.NameWithOwner),
			StarCount:            fmt.Sprintf("%d(%d/d)", totalStarCount, avgStarCount),
			LatestDayStarCount:   formatStarTrend(latestMonthStargazers.LatestDayStars()),
			LatestWeekStarCount:  formatStarTrend(latestMonthStargazers.LatestWeekStars()),
			LatestMonthStarCount: formatValue(latestMonthStargazers.LatestMonthStars()),
			ForkCount:            fmt.Sprintf("%d(%d/d)", totalForkCount, avgForkCount),
			WatcherCount:         formatValue(repo.Watchers.TotalCount),
			Language:             formatValue(repo.PrimaryLanguage.Name),
			Issue:                fmt.Sprintf("%d/%d", openIssueCount, repo.Issues.TotalCount),
			Pull:                 fmt.Sprintf("%d/%d", openPrCount, repo.PullRequests.TotalCount),
			License:              formatValue(repo.LicenseInfo.Name),
			Age:                  formatPeriod(time.Since(repo.CreatedAt.Time)),
			LastPushedAt:         formatDuration(repo.PushedAt.Time),
			LastUpdatedAt:        formatDuration(repo.UpdatedAt.Time),
			LatestReleaseAt:      formatDuration(repo.LatestRelease.PublishedAt.Time),
			ReleaseCount:         formatValue(repo.Releases.TotalCount),
			AvgReleasePeriod:     formatPeriod(avgReleasePeriod),
			ContributorCount:     formatValue(contributorCount),
			Homepage:             repo.HomepageUrl.String(),
		})
		writer.Write(list)
	}, func(pipe <-chan []Data, writer mapreduce.Writer[[]Data], cancel func(error)) {
		var list []Data
		for p := range pipe {
			list = append(list, p...)
		}
		writer.Write(list)
	}, mapreduce.WithWorkers(len(repos)))

	m := make(map[string]Data, len(reduce))
	for _, e := range reduce {
		m[e.FullName] = e
	}

	var list []Data
	for _, r := range repos {
		if data, ok := m[r]; ok {
			list = append(list, data)
		}
	}
	return list
}

func formatValue(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func formatStarTrend(stars, trend int) string {
	var trendEmoji string
	c := color.New(color.FgHiWhite)
	switch {
	case trend < 0:
		c.Add(color.BgHiRed)
		trendEmoji = c.Sprintf("(down)")
	case trend > 0:
		c.Add(color.BgHiGreen)
		trendEmoji = c.Sprintf("(up)")
	default:
		trendEmoji = ""
	}
	return fmt.Sprintf("%d %s", stars, trendEmoji)
}
