package stat

import "github.com/google/go-github/v44/github"

func (s Stat) ContributorCount() int {
	listOpt := &github.ListContributorsOptions{
		Anon:        "true",
		ListOptions: github.ListOptions{Page: 1, PerPage: 1},
	}
	_, resp, _ := s.restClient.Repositories.ListContributors(s.ctx, s.owner, s.repo, listOpt)
	return s.GetTotal(resp)
}
