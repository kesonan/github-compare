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
	"path/filepath"
	"strings"
	"time"

	"github.com/anqiansong/github-compare/pkg/stat"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

const codeFailure = 1

var (
	githubAccessToken string
	jsonStyle         bool
	tableStyle        bool
	yamlStyle         bool

	rootCmd = &cobra.Command{
		Use:   "github-compare",
		Short: "A cli tool to compare two github repositories",
		Args:  cobra.RangeArgs(1, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateGithubRepo(args...); err != nil {
				return err
			}

			printStyle := styleTable
			if jsonStyle {
				printStyle = styleJSON
			} else if yamlStyle {
				printStyle = styleYAML
			}

			data, err := getData(printStyle == styleTable && len(outputFile) == 0, args...)
			if err != nil {
				return err
			}

			if len(outputFile) > 0 {
				tp := getExportType(outputFile, printStyle)
				return export(data, tp)
			}

			return render(printStyle, data...)
		},
	}
)

func getExportType(outputFile string, printStyle style) string {
	ext := strings.TrimPrefix(filepath.Ext(outputFile), ".")
	switch strings.ToLower(ext) {
	case "json":
		return exportTPJSON
	case "yaml", "yml":
		return exportTPYAML
	case "csv":
		return exportTPCSV
	default:
		if printStyle != styleTable {
			return string(printStyle)
		}
		return exportTPJSON
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&githubAccessToken, "token", "t", "",
		"github access token")
	rootCmd.PersistentFlags().BoolVar(&tableStyle, "table", true,
		"print with table style(default)")
	rootCmd.PersistentFlags().BoolVar(&jsonStyle, "json", false, "print with json style")
	rootCmd.PersistentFlags().BoolVar(&yamlStyle, "yaml", false, "print with yaml style")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "file", "f", "", "output to a specified file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(codeFailure)
	}
}

func getData(renderColor bool, args ...string) ([]stat.Data, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Suffix = " Loading..."
	s.Start() // Start the spinner
	data := stat.Overview(githubAccessToken, renderColor, args...)
	s.Stop()
	return data, nil
}
