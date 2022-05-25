package cmd

import (
	"fmt"
	"regexp"
)

const repoRegex = `(?m)^[\w-]+\/[\w-]+`

func validateGithubRepo(name ...string) error {
	re := regexp.MustCompile(repoRegex)
	for _, e := range name {
		all := re.FindAllString(e, -1)
		if len(all) > 0 && all[0] == e {
			continue
		}
		return fmt.Errorf("invalid github repo name: %s", name)
	}
	return nil
}
