package stat

import "github.com/shurcooL/githubv4"

type (
	Release struct {
		CreatedAt   githubv4.DateTime
		PublishedAt githubv4.DateTime
	}

	ReleaseConnection struct {
		TotalCount githubv4.Int
	}
)
