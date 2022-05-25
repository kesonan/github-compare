package cmd

import (
	"os"
	"time"

	"github.com/anqiansong/github-compare/pkg/stat"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

const codeFailure = 1

var (
	githubAccessToken string
	rootCmd           = &cobra.Command{
		Use:   "github-compare",
		Short: "A cli tool to compare two github repositories",
		Args:  cobra.RangeArgs(1, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateGithubRepo(args...); err != nil {
				return err
			}
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
			s.Suffix = " Loading..."
			s.Start() // Start the spinner
			data := stat.Overview(githubAccessToken, args...)
			return render(s, data...)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&githubAccessToken, "token", "t", "",
		"github access token")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(codeFailure)
	}
}
