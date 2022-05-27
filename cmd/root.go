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

	rootCmd = &cobra.Command{
		Use:   "github-compare",
		Short: "A cli tool to compare two github repositories",
		Args:  cobra.RangeArgs(1, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
			s.Suffix = " Loading..."
			s.Start() // Start the spinner

			data, err := checkAndGet(s, true, args...)
			if err != nil {
				return err
			}
			return render(s, data...)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&githubAccessToken, "token", "t", "",
		"github access token")
	rootCmd.AddCommand(exportCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(codeFailure)
	}
}

func checkAndGet(s *spinner.Spinner, renderColor bool, args ...string) ([]stat.Data, error) {
	if err := validateGithubRepo(args...); err != nil {
		return nil, err
	}

	data := stat.Overview(githubAccessToken, renderColor, args...)
	s.Stop()
	return data, nil
}
