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
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/kevwan/mapreduce/v2"
	"github.com/shurcooL/githubv4"
)

type (
	Data struct {
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

		Description           string   `json:"description,omitempty"`
		Tags                  []string `json:"tags,omitempty"`
		LatestMonthStargazers Chart    `json:"latestMonthStargazers"`

		LatestWeekForks   Chart `json:"latestWeekForks"`
		LatestWeekCommits Chart `json:"latestWeekCommits"`
		LatestWeekPulls   Chart `json:"latestWeekPulls"`
		LatestWeekIssues  Chart `json:"latestWeekIssues"`
	}

	Chart struct {
		Data   []float64
		Labels []string
	}
)

func Overview(accessToken string, renderColor bool, repos ...string) []Data {
	getDetail := len(repos) == 1
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
			forkWeekChart         Chart
			commitWeekChart       Chart
			pullWeekChart         Chart
			issueWeekChart        Chart
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
		}, func() {
			if getDetail {
				forkWeekChart = s.latestWeekForks().Chart()
			}
		}, func() {
			if getDetail {
				commitWeekChart = s.latestWeekCommits().Chart()
			}
		}, func() {
			if getDetail {
				pullWeekChart = s.latestWeekPRS().Chart()
			}
		}, func() {
			if getDetail {
				issueWeekChart = s.LatestWeekIssues().Chart()
			}
		})

		homePage := ""
		if repo.HomepageUrl.URL != nil {
			homePage = repo.HomepageUrl.URL.String()
		}
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
			FullName:  fmt.Sprintf("%s/%s", s.owner, s.repo),
			StarCount: fmt.Sprintf("%d(%d/d)", totalStarCount, avgStarCount),
			LatestDayStarCount: formatStarTrend(func() (int, int, bool) {
				stars, trend := latestMonthStargazers.LatestDayStars()
				return stars, trend, renderColor
			}()),
			LatestWeekStarCount: formatStarTrend(func() (int, int, bool) {
				stars, trend := latestMonthStargazers.LatestWeekStars()
				return stars, trend, renderColor
			}()),
			LatestMonthStarCount: formatValue(latestMonthStargazers.LatestMonthStars()),
			ForkCount:            fmt.Sprintf("%d(%d/d)", totalForkCount, avgForkCount),
			WatcherCount:         formatValue(repo.Watchers.TotalCount),
			Language: formatLanguage(repo.PrimaryLanguage.Name,
				repo.PrimaryLanguage.Color, renderColor),
			Issue:   fmt.Sprintf("%d/%d", openIssueCount, repo.Issues.TotalCount),
			Pull:    fmt.Sprintf("%d/%d", openPrCount, repo.PullRequests.TotalCount),
			License: formatValue(repo.LicenseInfo.Name),
			Age: formatPeriod(func() time.Duration {
				if repo.CreatedAt.IsZero() {
					return 0
				}
				return time.Since(repo.CreatedAt.Time)
			}()),
			LastPushedAt:     formatDuration(repo.PushedAt.Time),
			LastUpdatedAt:    formatDuration(repo.UpdatedAt.Time),
			LatestReleaseAt:  formatDuration(repo.LatestRelease.PublishedAt.Time),
			ReleaseCount:     formatValue(repo.Releases.TotalCount),
			AvgReleasePeriod: formatPeriod(avgReleasePeriod),
			ContributorCount: formatValue(contributorCount),
			Homepage:         homePage,

			Description:           formatValue(repo.Description),
			Tags:                  repo.RepositoryTopics.List(),
			LatestMonthStargazers: latestMonthStargazers.Chart(),
			LatestWeekForks:       forkWeekChart,
			LatestWeekCommits:     commitWeekChart,
			LatestWeekPulls:       pullWeekChart,
			LatestWeekIssues:      issueWeekChart,
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
	ret := fmt.Sprintf("%v", v)
	if len(ret) == 0 {
		return "N/A"
	}
	return ret
}

func formatLanguage(lang, color githubv4.String, renderColor bool) string {
	if len(lang) == 0 {
		return "N/A"
	}
	if !renderColor {
		return string(lang)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(fmt.Sprintf("%s %s", "◉",
		lang))
}

func formatStarTrend(stars, trend int, renderColor bool) string {
	var trendEmoji, starStr string
	c := color.New()
	starStr = fmt.Sprintf("%d", stars)
	switch {
	case trend < 0:
		if !renderColor {
			trendEmoji = "⇊"
		} else {
			c.Add(color.FgHiRed)
			starStr = c.Sprintf("%d", stars)
			trendEmoji = c.Sprintf("⇊")
		}
	case trend > 0:
		if !renderColor {
			trendEmoji = "⇈"
		} else {
			c.Add(color.FgHiGreen)
			starStr = c.Sprintf("%d", stars)
			trendEmoji = c.Sprintf("⇈")
		}
	default:
		trendEmoji = ""
	}

	return fmt.Sprintf("%s %s", starStr, trendEmoji)
}
