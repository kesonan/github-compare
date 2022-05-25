package stat

import (
	"context"
	"fmt"
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
	if hours < hour {
		return fmt.Sprintf("%v minute(s) ago", minutes)
	}
	if hours < day {
		return fmt.Sprintf("%v hour(s) ago", hours)
	}
	if hours < month {
		return fmt.Sprintf("%v day(s) ago", hours/day)
	}
	if hours < year {
		return fmt.Sprintf("%v month(s) ago", hours/month)
	}
	return at.Format("2006-01-02")
}
