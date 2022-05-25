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
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v44/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	hour  = 1
	day   = 24 * hour
	month = 30 * day
	year  = 12 * month
)

type (
	Stat struct {
		owner         string
		repo          string
		graphqlClient *githubv4.Client
		restClient    *github.Client
		ctx           context.Context
	}

	PageInfo struct {
		EndCursor       githubv4.String
		HasNextPage     githubv4.Boolean
		HasPreviousPage githubv4.Boolean
		StartCursor     githubv4.String
	}
)

func NewStat(repo string, accessToken ...string) *Stat {
	token := getAccessToken(accessToken...)
	if len(token) == 0 {
		log.Fatalln("missing access token")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, ts)
	splits := strings.Split(repo, "/")
	graphqlClient := githubv4.NewClient(httpClient)
	restClient := github.NewClient(httpClient)

	return &Stat{owner: splits[0], repo: splits[1], graphqlClient: graphqlClient,
		restClient: restClient, ctx: context.Background()}
}

func (s Stat) GetTotal(resp *github.Response) int {
	if resp == nil {
		return 0
	}
	return resp.LastPage
}

func getAccessToken(accessToken ...string) string {
	for _, e := range accessToken {
		if len(e) > 0 {
			return e
		}
	}

	return os.Getenv("GITHUB_ACCESS_TOKEN")
}

func formatPeriod(duration time.Duration) string {
	if duration == 0 {
		return "N/A"
	}
	hours := duration.Hours()

	return fmt.Sprintf("%d days", int(hours/float64(24)))
}

func formatDuration(at time.Time) string {
	if at.IsZero() {
		return "N/A"
	}

	duration := time.Since(at)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes())

	if hours == 0 && minutes == 0 && duration.Seconds() < float64(60) {
		return fmt.Sprintf("%v seconds(s) ago", int(duration.Seconds()))
	}
	switch {
	case hours < hour:
		return fmt.Sprintf("%v minute(s) ago", minutes)
	case hours < day:
		return fmt.Sprintf("%v hour(s) ago", hours)
	case hours < month:
		return fmt.Sprintf("%v day(s) ago", hours/day)
	case hours < year:
		return fmt.Sprintf("%v month(s) ago", hours/month)
	default:
		return at.Format("2006-01-02")
	}
}
