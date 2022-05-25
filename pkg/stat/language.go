package stat

import "github.com/shurcooL/githubv4"

type (
	LanguageConnection struct {
		Nodes []Language
	}

	Language struct {
		Color githubv4.String
		Name  githubv4.String
	}
)
